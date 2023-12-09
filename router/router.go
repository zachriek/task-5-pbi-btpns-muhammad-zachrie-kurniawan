package router

import (
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/controllers"
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.Default()

	public := router.Group("")
	{
		public.POST("/users/register", controllers.Register)
		public.POST("/users/login", controllers.Login)
	}

	protected := router.Group("")
	{
		protected.Use(middlewares.Authenticate())
		{
			protected.GET("/users/:id", controllers.GetUserByID)
			protected.PUT("/users/:id", controllers.UpdateUserByID)
			protected.DELETE("/users/:id", controllers.DeleteUserByID)

			protected.GET("/photos", controllers.GetAllPhotos)
			protected.GET("/photos/:id", controllers.GetPhotoByID)
			protected.POST("/photos", controllers.CreatePhoto)
			protected.PUT("/photos/:id", controllers.UpdatePhotoByID)
			protected.DELETE("/photos/:id", controllers.DeletePhotoByID)
		}
	}

	return router
}
