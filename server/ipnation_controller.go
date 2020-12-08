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
	ctxCp := ctx.Copy()
	ipAddr := ctxCp.Query("ip")
	if ipAddr == "" {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, "No query parameter : ip")
		ctxCp.JSON(400, err400)
		return
	}
	res, err := c.IpNationService.GetIPNation(ipAddr)
	if err != nil {
		err500 := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())
		ctxCp.JSON(500, err500)
		return
	}
	ctxCp.JSON(200, res)
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
	ctxCp := ctx.Copy()
	bulkInput := IpNationCountriesBulk{}
	err := ctxCp.BindJSON(&bulkInput)
	if err != nil {
		err400 := problems.NewDetailedProblem(http.StatusBadRequest, err.Error())
		ctxCp.JSON(400, err400)
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

	ctxCp.JSON(http.StatusOK, IpNationCountriesBulkResponse{IpAddressResult: ret})

}
