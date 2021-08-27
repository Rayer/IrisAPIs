package main

import (
	"IrisAPIs"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetaLogger(t *testing.T) {
	tests := []struct {
		name              string
		url               string
		header            map[string][]string
		wantIpAddr        string
		wantCorrelationId string
		wantSid           string
		wantBidObjId      string
	}{
		{
			name: "SuccessCase",
			url:  "https://aa.bb.cc?bidobjid=aabbc&sid=ddeef",
			header: http.Header{
				"X-Forwarded-For":  {"1.2.3.4"},
				"X-Correlation-ID": {"CORRELATION"},
			},
			wantIpAddr:        "1.2.3.4",
			wantCorrelationId: "CORRELATION",
			wantSid:           "ddeef",
			wantBidObjId:      "aabbc",
		},
		{
			name: "SuccessCase2",
			url:  "https://aa.bb.cc?bidobjid=aabbc&sid=ddeef",
			header: http.Header{
				"X-Forwarded-For": {"1.2.3.4"},
				"X-Request-ID":    {"CORRELATION"},
			},
			wantIpAddr:        "1.2.3.4",
			wantCorrelationId: "CORRELATION",
			wantSid:           "ddeef",
			wantBidObjId:      "aabbc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := InjectLoggerMiddleware(logrus.New())
			g, _ := createGinTestItems()
			g.Request, _ = http.NewRequest("GET", tt.url, nil)
			g.Request.Header = tt.header
			middleware(g)
			meta := IrisAPIs.GetMeta(g)
			assert.Equal(t, tt.wantIpAddr, meta.IpAddress)
			assert.Equal(t, tt.wantCorrelationId, meta.CorrelationId)

		})
	}
}

//Only port this, we only care about Panic should not leak outside
func TestPanicClean(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := gin.New()
	password := "my-super-secret-password"
	router.Use(RecoveryWithLogger(gLogger))
	router.GET("/recovery", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusBadRequest)
		panic("Oupps, Houston, we have a problem")
	})
	// RUN

	req := httptest.NewRequest("GET", "/recovery", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", password))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check the buffer does not have the secret key
	assert.NotContains(t, buffer.String(), password)
}
