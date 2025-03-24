package utils

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type ErrResponse struct {
	Success  bool   `json:"success"`
	Code     int    `json:"code"`
	ErrorMsg string `json:"error_msg"`
}

// SendErrResponse  pack error response
func SendErrResponse(ctx context.Context, c *app.RequestContext, code int, err error) {
	// todo edit custom code
	c.JSON(code, ErrResponse{Success: false, Code: code, ErrorMsg: err.Error()})
}

// SendSuccessResponse  pack success response
func SendSuccessResponse(ctx context.Context, c *app.RequestContext, code int, data interface{}) {
	// todo edit custom code
	c.JSON(code, data)
}
