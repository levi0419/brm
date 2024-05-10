package responder

import (
	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func JsonResponse(ctx *gin.Context, statusCode int, status bool, message string, data any) {
	var response ResponseBody
	response.Success = status
	response.Message = message
	response.Message = message
	response.Data = data

	/*toBeSent, _ := json.Marshal(response)

	res, err := encryption.EncryptResponseData(toBeSent)
	if err != nil {
		fmt.Println("error while encryption response payload ", err)
	}
	*/
	ctx.JSON(statusCode, response)
}

func ResponseData(status bool, message string, data any) ResponseBody {
	var response ResponseBody
	response.Success = status
	response.Message = message
	response.Message = message
	response.Data = data

	/*toBeSent, _ := json.Marshal(response)

	res, err := encryption.EncryptResponseData(toBeSent)
	if err != nil {
		fmt.Println("error while encryption response payload ", err)
	}
	*/
	return response
}
