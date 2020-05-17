package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.NoRoute(func(c *gin.Context) {
		err404 := problems.NewStatusProblem(404)
		err404.Detail = "No such route!"
		c.JSON(404, err404)
	})

	r.NoMethod(func(c *gin.Context) {
		err404 := problems.NewStatusProblem(404)
		err404.Detail = "No such method!"
		c.JSON(404, err404)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})

	r.GET("/currency", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"anchor": "This is a pen!",
		})
	})

	r.Run()
}
