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

	// Admin group.
	// These all access the api by the admin backdoor so are exempt from authorization.
	display.Router.GET("/admin/reset", display.AdminReset)
	display.Router.GET("/admin/play-as/:username", display.SelectUser)
	display.Router.GET("/admin/dashboard", display.AdminDashboard)
	display.Router.GET("/data/", display.DataHandler)

	// Everything below is specific to a user and protected

	// Routes that fetch and display tables
	display.Router.GET("/commodities", display.FindPlayer(), display.ShowCommodities)
	display.Router.GET("/industries", display.FindPlayer(), display.ShowIndustries)
	display.Router.GET("/classes", display.FindPlayer(), display.ShowClasses)
	display.Router.GET("/industry_stocks", display.FindPlayer(), display.ShowIndustryStocks)
	display.Router.GET("/class_stocks", display.FindPlayer(), display.ShowClassStocks)
	display.Router.GET("/trace", display.FindPlayer(), display.ShowTrace)
	display.Router.GET("/industry/:id", display.FindPlayer(), display.ShowIndustry)
	display.Router.GET("/commodity/:id", display.FindPlayer(), display.ShowCommodity)
	display.Router.GET("/class/:id", display.FindPlayer(), display.ShowClass)
	display.Router.GET("/", display.FindPlayer(), display.ShowIndexPage)

	// Routes that look at tables that are already in client memory
	display.Router.GET("/back", display.FindPlayer(), display.Back)
	display.Router.GET("/forward", display.FindPlayer(), display.Forward)

	//Routes that do things
	display.Router.GET("/quit/", display.FindPlayer(), display.Quit)
	display.Router.GET("/action/:action", display.FindPlayer(), display.ActionHandler)
	display.Router.GET("/user/create/:id", display.FindPlayer(), display.CreateSimulation)
	display.Router.GET("/user/switch/:id", display.FindPlayer(), display.SwitchSimulation)
	display.Router.GET("/user/delete/:id", display.FindPlayer(), display.DeleteSimulation)
	display.Router.GET("/user/restart/:id", display.FindPlayer(), display.RestartSimulation)
	display.Router.GET("/user/dashboard", display.FindPlayer(), display.UserDashboard)

	// Grab user data from the server at startup. Currently, this is fixed.
	// Note that user data on the server are changing as users come and go,
	// so we have to resynchronise whenever we make use of these data.
	// This is the purpose of the player handlers.
	fetch.InitialiseTemplates()
	fetch.InitialiseUsers()

	// Uncomment in extremis for very verbose diagnostic. As a first resort use the /Data endpoint when simulation is running.
	// display.ListData()

	display.Router.Run() // Run the server
}
