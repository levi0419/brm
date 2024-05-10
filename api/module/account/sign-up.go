package account

import (
	"net/http"
	"strings"
	"time"

	"github.com/brm/api/shared"
	db "github.com/brm/db/sqlc"
	"github.com/brm/logger"
	"github.com/brm/responder"
	"github.com/brm/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Uuid:              user.Uuid.String(),
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		PhoneNumber:       user.PhoneNumber,
		IsEmailVerified:   user.IsEmailVerified,
		CreatedAt:         user.CreatedAt,
	}
}

func newCreateUserResponse(user db.UserSignUpTxResult) userResponse {
	
	return userResponse{
		Uuid:                  user.User.Uuid.String(),
		FirstName:             user.User.FirstName,
		LastName:              user.User.LastName,
		Email:                 user.User.Email,
		PhoneNumber:           user.User.PhoneNumber,
		CreatedAt:             user.User.CreatedAt,
		IsEmailVerified:       user.User.IsEmailVerified,
		
	}
}

type userResponse struct {
	Uuid                  string    `json:"uuid"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	PhoneNumber           string    `json:"phone_number"`
	Email                 string    `json:"email"`
	IsEmailVerified       string      `json:"is_email_verified"`
	CreatedAt             time.Time `json:"created_at"`
}

type DiscountResponse struct {
    Label          string    `json:"label"`
    ExpirationTime time.Time `json:"expiration_time"`
    Code           string    `json:"code"`
    Status         bool      `json:"status"`
    CreatedAt      time.Time `json:"created_at"`
	Message        string    `json:"message"`
}


type createUserRequest struct {
	FirstName   string `json:"first_name" binding:"required,alphanum"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	DeviceID    string `json:"device_id" binding:"required,min=6"`
}

func createUser(ctx *gin.Context, server shared.IServer) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Error("sign up validation error: ", zap.Error(err))
		responder.JsonResponse(ctx, http.StatusBadRequest, false, err.Error(), nil)
		return
	}

	hashedPassword, err := utils.HashText(req.Password)
	if err != nil {
		logger.Error("sign up password hashing error: ", zap.Error(err))
		responder.JsonResponse(ctx, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	createUserArg := db.CreateUserParams{
		Uuid:        utils.NewUuid(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       strings.ToLower(strings.TrimSpace(req.Email)),
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
		DeviceID:    req.DeviceID,
	}

	userDeviceArg := db.AddUserDeviceParams{
		DeviceID:      req.DeviceID,
		IpAddress:     ctx.ClientIP(),
		ClientDetails: ctx.Request.UserAgent(),
	}

	user, err := server.GetDbStore().UserSignUpTx(ctx, createUserArg, userDeviceArg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok { // try converting the error to pq.Error type
			switch pqErr.Code.Name() {
			case "unique_violation":
				logger.Error("sign up unique_violation error: ", zap.Error(err))
				responder.JsonResponse(ctx, http.StatusForbidden, false, err.Error(), nil)
				return
			}
		}
		logger.Error("sign up tx error: ", zap.Error(err))
		responder.JsonResponse(ctx, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}


	rsp, err := loginUserHandler(ctx, server, user.User.Email, req.Password)
	if err != nil {
		logger.Error("login user func error: ", zap.Error(err))
		responder.JsonResponse(ctx, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	//rsp := newCreateUserResponse(user)
	logger.Info("sign up response: ", zap.Any("data", rsp))
	responder.JsonResponse(ctx, http.StatusCreated, true, "user created successfully", rsp)
}
