package core

import (
	"github.com/gin-gonic/gin"
	"githup.com/htl/tcclienttest/pkg/log"
	"net/http"
)

type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

// WriteResponse write an error or the response data into http response body.
// It use errors.ParseCoder to parse any error into errors.Coder
// errors.Coder contains error code, user-safe error message and http status code.
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%#+v", err)

		//coder := errors.ParseCoder(err)
		//c.JSON(coder.HTTPStatus(), ErrResponse{
		//	Code:      coder.Code,
		//	Message:   coder.String(),
		//	Reference: coder.Reference(),
		//}
		//)

		return
	}

	c.JSON(http.StatusOK, data)
}