package image

import (
	"github.com/gin-gonic/gin"
	"githup.com/htl/tcclienttest/internal/pkg/core"
	v1 "githup.com/htl/tcclienttest/internal/pkg/meta/v1"
	"githup.com/htl/tcclienttest/pkg/log"
	"strconv"
)

// Get get an user by the user identifier.
func (i *ImageController) Get(c *gin.Context) {
	log.Infof("get image function called.")
	imageid64, _ := strconv.ParseInt(c.Param("imageid"), 0, 10)
	image, err := i.srv.Images().Get(c, imageid64, v1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, image)
}
