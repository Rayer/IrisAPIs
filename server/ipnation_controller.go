package main

import (
	"IrisAPIs"
	"github.com/gin-gonic/gin"
	"github.com/moogar0880/problems"
	"net/http"
)

type IpNationCountries struct {
	IrisAPIs.IpNationCountries
}

type IpNationCountriesBulk struct {
	IpAddresses []string `json:"ip_addr_list"`
}

type IpNationCountriesBulkResponse struct {
	IpAddressResult map[string]string `json:"ip_addr_result"`
}

type IpNationMyIPResponse struct {
	IpAddr        string
	Country       string
	CountrySymbol string
	Lat           float32
	Lon           float32
}

// IpToNation godoc
// @Summary IP to Nation
// @Description Look up in database, find which nation belongs to an IP
// @Tags Ip2Nation
// @Param ip query string true "IP address"
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} main.IpNationCountries
// @Failure 400 {object} problems.DefaultProblem
// @Router /ip2nation [get]
func (c *Controller) IpToNation(ctx *gin.Context) {
	ipAddr := ctx.Query("ip")
	if ipAddr == "" {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, "No query parameter : ip")
		ctx.JSON(400, err400)
		return
	}
	res, err := c.IpNationService.GetIPNation(ipAddr)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctx.JSON(500, err500)
		return
	}
	ctx.JSON(200, res)
}

// IpToNationBulk godoc
// @Summary IP to Nation in bulk
// @Description Look up in database, find which nation belongs to an IP
// @Tags Ip2Nation
// @Produce json
// @param input body IpNationCountriesBulk true "IP Addresses"
// @Security ApiKeyAuth
// @Success 200 {object} IpNationCountriesBulkResponse
// @Failure 400 {object} problems.DefaultProblem
// @Router /ip2nation/bulk [post]
func (c *Controller) IpToNationBulk(ctx *gin.Context) {
	bulkInput := IpNationCountriesBulk{}
	err := ctx.BindJSON(&bulkInput)
	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, err.Error())
		ctx.JSON(400, err400)
		return
	}

	ret := make(map[string]string)
	for _, ip := range bulkInput.IpAddresses {
		if _, ok := ret[ip]; !ok {
			res, err := c.IpNationService.GetIPNation(ip)
			if err != nil {
				ret[ip] = "ERROR : " + err.Error()
				continue
			}
			ret[ip] = res.Country
		}
	}

	ctx.JSON(http.StatusOK, IpNationCountriesBulkResponse{IpAddressResult: ret})
}

// IpToNationMyIP godoc
// @Summary Lookup my IP information
// @Description Detect client IP address and look up information
// @Tags Ip2Nation
// @Produce json
// @Success 200 {object} IpNationMyIPResponse
// @Failure 500 {object} problems.DefaultProblem
// @Router /ip2nation/myip [get]
func (c *Controller) IpToNationMyIP(ctx *gin.Context) {
	ipAddr := ctx.GetHeader("X-Forwarded-For")
	if ipAddr == "" {
		ipAddr = ctx.ClientIP()
	}

	i, err := c.IpNationService.GetIPNation(ipAddr)
	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, err400)
		return
	}

	ctx.JSON(http.StatusOK, IpNationMyIPResponse{
		IpAddr:        ipAddr,
		Country:       i.Country,
		CountrySymbol: i.IsoCode_3,
		Lat:           i.Lat,
		Lon:           i.Lon,
	})

}
