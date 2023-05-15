package services

import (
	"wind-guide/protobuf_data"
	"wind-guide/response"

	"github.com/Wind-318/wind-chimes/logger"
	"github.com/lesismal/arpc"
	"google.golang.org/protobuf/proto"
)

// SendErrorResponse is a function for sending error response.
func SendErrorResponse(ctx *arpc.Context, errorCode string, errorMessage string, err error) {
	// Record log.
	logger.Logger.Warn().Err(err).Msg(errorMessage)
	// Send response.
	res, err := proto.Marshal(&protobuf_data.RegisterResponse{
		Code:    errorCode,
		Message: errorMessage,
	})
	if err != nil {
		logger.Logger.Warn().Err(err).Msg(response.MarshalFailedCode)
		ctx.Error(err.Error())
		return
	}
	ctx.Write(res)
}
