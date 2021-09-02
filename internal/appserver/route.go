package appserver

import (
	"github.com/gin-gonic/gin"
	"githup.com/htl/tcclienttest/internal/appserver/controller/v1/image"
	"githup.com/htl/tcclienttest/internal/appserver/controller/v1/terminal"
	"githup.com/htl/tcclienttest/internal/appserver/store/mysql"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	// Middlewares.

	// v1 handlers, requiring authentication
	storeIns, _ := mysql.GetMySQLFactoryOr(nil)
	v1 := g.Group("/v1")
	{
		// user RESTful resource
		imagev1 := v1.Group("/images")
		{
			imageController := image.NewImageController(storeIns)

			//userv1.POST("", userController)
			//userv1.Use(auto.AuthFunc(), middleware.Validation())
			//// v1.PUT("/find_password", userController.FindPassword)
			//userv1.DELETE("", userController.DeleteCollection) // admin api
			//userv1.DELETE(":name", userController.Delete)      // admin api
			//userv1.PUT(":name/change-password", userController.ChangePassword)
			//userv1.PUT(":name", userController.Update)
			imagev1.GET("", imageController.Get)
			//imagev1.GET(":name", userController.Get) // admin api
		}

		//v1.Use(auto.AuthFunc())

		// policy RESTful resource
		terminalv1 := v1.Group("/terminals")
		{
			terminalController := terminal.NewTerminalController(storeIns)

			terminalv1.GET("", terminalController.Get)
		}

		// secret RESTful resource
		//secretv1 := v1.Group("/secrets", middleware.Publish())
		//{
		//	secretController := secret.NewSecretController(storeIns)
		//
		//	secretv1.POST("", secretController.Create)
		//	secretv1.DELETE(":name", secretController.Delete)
		//	secretv1.PUT(":name", secretController.Update)
		//	secretv1.GET("", secretController.List)
		//	secretv1.GET(":name", secretController.Get)
		//}
	}

	return g
}
