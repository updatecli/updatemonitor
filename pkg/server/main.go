package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/updatecli/updatemonitor/pkg/dashboard"
	"github.com/updatecli/updatemonitor/pkg/version"
)

type Options struct {
	ReadOnly bool
}

type Server struct {
	Options Options
}

func Create(c *gin.Context) {

	var d dashboard.Dashboard
	err := d.Init()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	if err := c.BindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	err = d.Save()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully"})
}

func Delete(c *gin.Context) {
	ID := c.Param("id")

	result, err := dashboard.DeletebyID(ID)

	res := map[string]interface{}{"data": result}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	if result.DeletedCount < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No data to delete"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Dashboard deleted successfully", "data": res})
}

func FindAll(c *gin.Context) {
	dashboards, err := dashboard.SearchAll()
	if err != nil {
		c.JSON(501, gin.H{"err": err.Error()})
	}

	type data struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	var dataset []data

	for _, dashboard := range dashboards {
		dataset = append(dataset, data{
			ID:   dashboard.ID.Hex(),
			Name: dashboard.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": dataset})
}

func FindByID(c *gin.Context) {
	var err error
	var d dashboard.Dashboard

	ID := c.Param("id")

	d, err = dashboard.SearchbyID(ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success!", "data": d})
}

func Landing(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hey what's up?"})
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func About(c *gin.Context) {

	v := struct {
		Golang    string
		Api       string
		BuildTime string
	}{
		Golang:    version.GoVersion,
		Api:       version.Version,
		BuildTime: version.BuildTime,
	}

	c.JSON(200, gin.H{
		"version": v,
	})
}

func Update(c *gin.Context) {

	ID := c.Param("id")
	var dashboard dashboard.Dashboard

	if err := c.BindJSON(&dashboard); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	result, err := dashboard.SaveByID(ID)

	res := map[string]interface{}{"data": result}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	if result.MatchedCount < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Data doesn't exist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "data updated successfully!", "data": res})
}

func (s *Server) Run() {
	r := gin.Default()
	r.GET("/", Landing)
	r.GET("/ping", Ping)
	r.GET("/about", About)
	r.GET("/dashboards", FindAll)
	r.GET("/dashboards/:id", FindByID)

	if !s.Options.ReadOnly {
		r.POST("/dashboards", Create)
		r.PUT("/dashboards/:id", Update)
		r.DELETE("/dashboards/:id", Delete)
	}

	// listen and server on 0.0.0.0:8080
	if err := r.Run(); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
}
