module IrisAPIsServer

go 1.13

replace IrisAPIs => ./../

require (
	IrisAPIs v0.0.0-00010101000000-000000000000
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/moogar0880/problems v0.1.1
	github.com/sirupsen/logrus v1.6.0
)
