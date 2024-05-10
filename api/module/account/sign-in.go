package account

import (
	_ "context"
	"database/sql"
	_ "encoding/json"
	"errors"
	"fmt"

	"net/http"
	_ "strconv"
	"time"

	"github.com/brm/api/shared"
	db "github.com/brm/db/sqlc"
	_ "github.com/brm/errs"
	"github.com/brm/logger"
	"github.com/brm/responder"
	"github.com/brm/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
	@Route("/login", "POST")
	@Tags("Onboarding")
	@Summary("Login a user")
	@Request(loginUserRequest)
	@Response(200, loginUserResponse)
	@Response(400, errorResponse bad request)
	@Response(500, errorResponse internal server error)
	l
*/

type loginUserRequest struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required" min:"6" max:"20"`
}


func loginUser(ctx *gin.Context, server shared.IServer) {
    // Validate the JSON request first
    var req loginUserRequest

    logger.Info("Login user request", zap.String("email", req.Email))

    if err := ctx.ShouldBindJSON(&req); err != nil {
        logger.Error("failed in binding json request", zap.Error(err))
        responder.JsonResponse(ctx, http.StatusBadRequest, false, err.Error(), nil)
        return
    }

    // Call the loginUserHandler function
    response, err := loginUserHandler(ctx, server, req.Email, req.Password)
    if err != nil {
        logger.Error("login failed", zap.Error(err))
        responder.JsonResponse(ctx, http.StatusUnprocessableEntity, false, err.Error(), nil)
        
    }

    logger.Info("Login user response", zap.Any("response", response))

    responder.JsonResponse(ctx, http.StatusOK, true, "Login Successfully", response)
}

type LoginResponse struct {
    SessionID       		string               `json:"session_id"`
    Token           		string               `json:"token"`
	RefreshToken    		string				 `json:"refresh_token"`
    TokenExpiresAt  		string               `json:"token_expires_at"`
	RefreshTokenExpireAt	string				 `json:"resfresh_token_expires_at"`
    User            		userResponse         `json:"user"`
	Discount                DiscountResponse     `json:"discount"`
}

func loginUserHandler(ctx *gin.Context, server shared.IServer, email, password string) (LoginResponse, error) {
	var response LoginResponse
	var err error

	// Check if the user exists in the database
	user, err := server.GetDbStore().GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			logger.Error("user record not found", zap.Error(err))
			return response, errors.New("user not found")
		}

		logger.Error("failed to fetch user by email", zap.Error(err))
		return response, errors.New("failed to fetch user by email")
	}

	// Compare the password and the hash
	if err := utils.CheckHashedText(password, user.Password); err != nil {
		logger.Error("invalid credentials", zap.Error(err))
		return response, errors.New("invalid credentials")
	}

	// Get country and state information from the user's IP address
	// ipAddress := ctx.GetHeader("X-Forwarded-For")
	// if ipAddress == "" {
	// 	ipAddress = ctx.Request.RemoteAddr
	// }

	ipAddress := "82.24.36.163"

	country, _, err := utils.GetCountryAndStateFromIP(ipAddress)
	if err != nil {
		logger.Error("failed to get country and state from IP", zap.Error(err))
		return response, errors.New("failed to get country and state from IP")
	}


	// Check if the user is from a third-world country
	isThirdWorld := utils.IsThirdWorldCountry(country)

	fmt.Println("world country", isThirdWorld)

	if isThirdWorld {
		fmt.Println("User is from a third-world country")
	
		// Check if the user already has a discount
		existingDiscount, _ := server.GetDbStore().GetUserDiscountByUserID(ctx, user.ID)
		// Check if the user does not have an existing discount
		if existingDiscount == (db.Discount{}) {
			expirationTime := time.Now().Add(7 * 24 * time.Hour)

			// Create the user discount
			userDiscount := db.AddDiscountParams{
				UserID:         user.ID,
				Label:          "discount",
				ExpirationTime: expirationTime,
				Code:           utils.Generate(), 
				Status:         true,
				CreatedAt:      time.Now(),
			}

			// If the current time exceeds the expiration time, set the status to false
			if time.Now().After(expirationTime) {
				userDiscount.Status = false
			}

			// Save the user discount to the database
			if _, err := server.GetDbStore().AddDiscount(ctx, userDiscount); err != nil {
				logger.Error("failed to add user discount", zap.Error(err))
				return response, errors.New("failed to add user discount")
			}

			fmt.Println("Discount created for user from a third-world country")
		} else {
			fmt.Println("User already has a discount")
		}
	}

	discount, _ := server.GetDbStore().GetUserDiscountByUserID(ctx, user.ID)

	// Check if the user has an existing session
	existSession, _ := server.GetDbStore().GetUserSessionByUserID(ctx, user.ID)
	if existSession == (db.UserSession{}) {
		// No existing session found, proceed to create a new session
		token, payload, err := server.GetTokenMaker().CreateToken(
			user.Uuid,
			server.GetConfig().AccessTokenDuration,
		)
		if err != nil {
			logger.Error("failed to create access token", zap.Error(err))
			return response, errors.New("failed to create access token")
		}

		refreshToken, refreshPayload, err := server.GetTokenMaker().CreateToken(
			user.Uuid,
			server.GetConfig().RefreshTokenDuration,
		)
		if err != nil {
			logger.Error("failed to create refresh token", zap.Error(err))
			return response, errors.New("failed to create refresh token")
		}

		session, err := server.GetDbStore().CreateUserSession(ctx, db.CreateUserSessionParams{
			SessionID:    utils.NewUuid(),
			UserID:       user.ID,
			Token:        token,
			RefreshToken: refreshToken,
			UserAgent:    ctx.Request.UserAgent(),
			Ip:           ctx.Request.RemoteAddr,
			Channel:      "mobile",
			ExpiresAt:    payload.ExpiredAt,
		})

		if err != nil {
			logger.Error("failed to create user session", zap.Error(err))
			return response, errors.New("failed to create user session")
		}

		// Populate the response struct with data from the database for the new session
		response = LoginResponse{
			SessionID:             session.SessionID.String(),
			Token:                 token,
			TokenExpiresAt:        payload.ExpiredAt.Format(time.RFC3339),
			RefreshToken:          refreshToken,
			RefreshTokenExpireAt:  refreshPayload.ExpiredAt.Format(time.RFC3339),
			User: userResponse{
				Uuid:          user.Uuid.String(),
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				PhoneNumber:   user.PhoneNumber,
				Email:         user.Email,
				IsEmailVerified:  user.Email,
				CreatedAt:     user.CreatedAt,
			},
		}
		if isThirdWorld {
			if discount == (db.Discount{}) {
				response.Discount.Message = "You are eligible for a discount as a user from a third-world country"
			} else {
				response.Discount = DiscountResponse{
					Label:          discount.Label,
					Code:           discount.Code,
					ExpirationTime: discount.ExpirationTime,
					CreatedAt:      discount.CreatedAt,
				}
			}
		} else {
			// User is not from a third-world country
			response.Discount.Message = "You are not eligible for a discount"
		}
	
	} else {
		// Create tokens
		updatedToken, updatedPayload, err := server.GetTokenMaker().CreateToken(
			user.Uuid,
			server.GetConfig().AccessTokenDuration,
		)
		if err != nil {
			logger.Error("failed to create access token", zap.Error(err))
			return response, errors.New("failed to create access token")
		}

		updatedRefreshToken, updatedRefreshPayload, err := server.GetTokenMaker().CreateToken(
			user.Uuid,
			server.GetConfig().RefreshTokenDuration,
		)
		if err != nil {
			logger.Error("failed to create refresh token", zap.Error(err))
			return response, errors.New("failed to create refresh token")
		}

		sessionResp, err := server.GetDbStore().UpdateUserSession(ctx, db.UpdateUserSessionParams{
			UserID:         user.ID,
			SessionID:      existSession.SessionID,
			Token:          updatedToken,
			RefreshToken:   updatedRefreshToken,
			UserAgent:      ctx.Request.UserAgent(),
			Ip:             ctx.Request.RemoteAddr,
			Channel:        "web",
			ExpiresAt:      updatedPayload.ExpiredAt,
		})

		if err != nil {
			logger.Error("failed to update user session", zap.Error(err))
			return response, errors.New("failed to update user session")
		}


		// Populate the response struct with data from the database
		response = LoginResponse{
			SessionID:             sessionResp.SessionID.String(),
			Token:                 updatedToken,
			TokenExpiresAt:        updatedPayload.ExpiredAt.Format(time.RFC3339),
			RefreshToken:          updatedRefreshToken,
			RefreshTokenExpireAt:  updatedRefreshPayload.ExpiredAt.Format(time.RFC3339),
			User: userResponse{
				Uuid:          user.Uuid.String(),
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				PhoneNumber:   user.PhoneNumber,
				Email:         user.Email,
				IsEmailVerified:  user.IsEmailVerified,
				CreatedAt:     user.CreatedAt,
			},
		}
		if isThirdWorld {
			if discount == (db.Discount{}) {
				response.Discount.Message = "You are eligible for a discount as a user from a third-world country"
			} else {
				response.Discount = DiscountResponse{
					Label:          discount.Label,
					Code:           discount.Code,
					ExpirationTime: discount.ExpirationTime,
					CreatedAt:      discount.CreatedAt,
				}
			}
		} else {
			// User is not from a third-world country
			response.Discount.Message = "You are not eligible for a discount"
		}

	}

	now := time.Now()
	nullTime := sql.NullTime{Time: now, Valid: true}

	if _, err := server.GetDbStore().UpdateUserLastLogin(ctx, db.UpdateUserLastLoginParams{
		Uuid:     user.Uuid,
		LastLogin: nullTime,
	}); err != nil {
		logger.Error("failed to update last login time", zap.Error(err))
		return response, errors.New("failed to update last login time")
	}

	return response, nil
}
