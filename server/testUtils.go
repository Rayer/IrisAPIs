package main

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func createGinTestItems() (g *gin.Context, r *httptest.ResponseRecorder) {
	r = httptest.NewRecorder()
	g, _ = gin.CreateTestContext(r)
	return
}
