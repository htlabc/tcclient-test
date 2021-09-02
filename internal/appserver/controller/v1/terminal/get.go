package terminal

import (
	"github.com/gin-gonic/gin"
	"githup.com/htl/tcclienttest/internal/pkg/core"
	v1 "githup.com/htl/tcclienttest/internal/pkg/meta/v1"
	"githup.com/htl/tcclienttest/pkg/log"
	"strconv"
)

func (i *TerminalController) Get(c *gin.Context) {
	log.Infof("get user function called.")
	imageid64, _ := strconv.ParseInt(c.Param("deviceid"), 0, 10)
	image, err := i.srv.Images().Get(c, imageid64, v1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, image)
}
