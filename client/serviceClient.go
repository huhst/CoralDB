package client

import (
	"CoralDB/server"
	"github.com/gin-gonic/gin"
)

var Cserver *server.Server

type Aircondition struct {
	Temperature       string `json:"temperature" xml:"temperature" form:"temperature" query:"temperature"`
	TargetTemperature string `json:"target" xml:"target" form:"target" query:"target"`
}
type Refrigerator struct {
	ColdTemperature string `json:"cold" xml:"cold" form:"cold" query:"cold"`
	IceTemperature  string `json:"ice" xml:"ice" form:"ice" query:"ice"`
}

func getAData(c *gin.Context) {
	var a Aircondition
	a.Temperature = Cserver.Get("123"+"-QAQ-"+"123456", "Temperature", "data", false)
	a.TargetTemperature = Cserver.Get("123"+"-QAQ-"+"123456", "TargetTemperature", "data", false)
	c.JSON(200, a)
}

func getRData(c *gin.Context) {
	var r Refrigerator
	r.ColdTemperature = Cserver.Get("123"+"-QAQ-"+"123456", "ColdTemperature", "data", false)
	r.IceTemperature = Cserver.Get("123"+"-QAQ-"+"123456", "IceTemperature", "data", false)
	c.JSON(200, r)
}

func StartClient() {
	r := gin.Default()
	r.GET("/aircondition", getAData)
	r.GET("/refrigerator", getRData)
	r.Run(":8080")
}
