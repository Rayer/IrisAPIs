package main

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strconv"
)

func createGinTestItems() (g *gin.Context, r *httptest.ResponseRecorder) {
	r = httptest.NewRecorder()
	g, _ = gin.CreateTestContext(r)
	return
}

func stringAsIntAddBy(num string, addBy int) string {
	val, err := strconv.Atoi(num)
	if err != nil {
		return ""
	}
	return strconv.Itoa(val + addBy)
}
