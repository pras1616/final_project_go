package main

import (
	"final_project/controllers"
	"final_project/database"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"final_project/docs"
	"final_project/helpers"
)

func init() {
	database.ConnectDB()
}

func main() {
	r := gin.Default()

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8888"
	docs.SwaggerInfo.Schemes = []string{"http"}

	ctrl := controllers.NewCarsController(database.GetDB_Comment(), database.GetDB_User(), database.GetDB_Photo(), database.GetDB_SocialMedia())
	userRouter := r.Group("/users")
	{
		userRouter.PUT("/:userId", helpers.Authentication(), ctrl.UpdateUser)
		userRouter.DELETE("/:userId", helpers.Authentication(), ctrl.DeleteUser)
		userRouter.Use(helpers.BasicAuth)
		userRouter.POST("/register", ctrl.Register)
		userRouter.POST("/login", ctrl.Login)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8888")
}
