package errs

import (
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrBadRequest           = CustomError{400, "Bad Request"}
	ErrUnauthorized         = CustomError{401, "Unauthorized Access"}
	ErrInvalidCredentials   = CustomError{401, "Credentials is incorrect"}
	ErrInvalidToken         = CustomError{404, "Request can not be served"}
	ErrEncounteredOnRequest = CustomError{401, "Some errors encountered on request, try again later"}

	ErrBvnRequired      = errors.New("BVN is required")
	ErrDobRequired      = errors.New("Date of birth (DOB) is required")
	ErrInvalidBvnFormat = errors.New("BVN must be 11 digits")
	ErrInvalidDobFormat = errors.New("Invalid DOB format. Expected format: 03/07/2023")
	ErrDobMismatch      = errors.New("Date of birth mismatch, Kindly provide us with the dob used in your bvn")
	ErrBvnFailed        = errors.New("BVN validation failed")
	ErrBaseFailed       = errors.New("Base url failed")
	ErrBvnNotValid      = errors.New("BVN provided is not valid")

	ErrInvalidAccountNumber = errors.New("Accoun number provided is not valid, is not up to 10 digits")
)

func MapThirdPartyError(httpStatusCode int) CustomError {
	switch httpStatusCode {
	case 400:
		return ErrBadRequest
	case 401:
		return ErrUnauthorized
	case 412:
		return ErrInvalidCredentials
	default:
		return CustomError{http.StatusNotFound, "Operation failed, please try again later"}
	}
}

func MapThirdPartyErrorWithMessage(httpStatusCode int, message string) CustomError {
	switch {
	case httpStatusCode == 400 && message == "":
		return ErrBadRequest
	case httpStatusCode == 401 && message == "":
		return ErrUnauthorized
	case httpStatusCode == 412 && message == "":
		return ErrInvalidCredentials
	default:
		return CustomError{http.StatusNotFound, "Operation failed, please try again later"}
	}
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

/*
func thirdPartyAPICall() (response *http.Response, customErr CustomError) {
	// Make the third-party API request and get the response
	// ...

	if response.StatusCode >= 400 {
		customErr = MapThirdPartyError(response.StatusCode, "")
	}

	return response, customErr
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	_, customErr := thirdPartyAPICall()
	if customErr.Code != 0 {
		// Handle the error, returning customErr as the response
		jsonResponse, _ := json.Marshal(customErr)
		w.WriteHeader(customErr.Code)
		w.Write(jsonResponse)
		return
	}

	// Process the successful response from the third-party API
	// ...
}
*/
