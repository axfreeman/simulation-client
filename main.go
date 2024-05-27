package main

import (
	"capfront/display"
	"capfront/fetch"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	display.Router.Use(gin.Recovery())

	// load the templates
	display.Router.LoadHTMLGlob("./templates/**/*")
	fmt.Println("The Rosy Dawn of Capitalism has begun")

	// Admin group. These all have access to the api by
	// the admin backdoor so they are exempt from authorization.
	display.Router.GET("/admin/reset", display.AdminReset)
	display.Router.GET("/admin/choose-players", display.Lock)
	display.Router.GET("/admin/play-as/:username", display.SelectUser)
	display.Router.GET("/admin/dashboard", display.AdminDashboard)

	// User group. These access the api via SynchWithServer()
	// which requires user to be authorized.
	// TODO couldn't get grouping to work. Pretty sure it did not work as per spec.

	display.Router.GET("/action/:action", display.SynchWithServer(), display.ActionHandler)
	display.Router.GET("/commodities", display.SynchWithServer(), display.ShowCommodities)
	display.Router.GET("/industries", display.SynchWithServer(), display.ShowIndustries)
	display.Router.GET("/classes", display.SynchWithServer(), display.ShowClasses)
	display.Router.GET("/industry_stocks", display.SynchWithServer(), display.ShowIndustryStocks)
	display.Router.GET("/class_stocks", display.SynchWithServer(), display.ShowClassStocks)
	display.Router.GET("/trace", display.SynchWithServer(), display.ShowTrace)
	display.Router.GET("/industry/:id", display.SynchWithServer(), display.ShowIndustry)
	display.Router.GET("/commodity/:id", display.SynchWithServer(), display.ShowCommodity)
	display.Router.GET("/class/:id", display.SynchWithServer(), display.ShowClass)
	display.Router.GET("/user/create/:id", display.SynchWithServer(), display.CreateSimulation)
	display.Router.GET("/user/switch/:id", display.SynchWithServer(), display.SwitchSimulation)
	display.Router.GET("/user/delete/:id", display.SynchWithServer(), display.DeleteSimulation)
	display.Router.GET("/user/restart/:id", display.SynchWithServer(), display.RestartSimulation)
	display.Router.GET("/", display.SynchWithServer(), display.ShowIndexPage)
	display.Router.GET("/data/", display.SynchWithServer(), display.DataHandler)
	display.Router.GET("/user/dashboard", display.SynchWithServer(), display.UserDashboard)
	display.Router.GET("/back", display.SynchWithServer(), display.Back)
	display.Router.GET("/forward", display.SynchWithServer(), display.Forward)
	// api.Use(display.SynchWithServer())
	fetch.Initialise()
	display.Router.Run("localhost:8080") // Run the server

}
