package main

import (
	"capfront/api"
	"capfront/display"
	"fmt"
)

func main() {

	api.Router.LoadHTMLGlob("./templates/**/*") // load the templates

	fmt.Println("Welcome to capitalism")

	api.Router.GET("/action/:action", display.ActionHandler)

	api.Router.GET("/commodities", display.ShowCommodities)
	api.Router.GET("/industries", display.ShowIndustries)
	api.Router.GET("/classes", display.ShowClasses)
	api.Router.GET("/industry_stocks", display.ShowIndustryStocks)
	api.Router.GET("/class_stocks", display.ShowClassStocks)
	api.Router.GET("/trace", display.ShowTrace)

	api.Router.GET("/industry/:id", display.ShowIndustry)
	api.Router.GET("/commodity/:id", display.ShowCommodity)
	api.Router.GET("/class/:id", display.ShowClass)

	api.Router.GET("/user/create/:id", display.CreateSimulation)
	api.Router.GET("/user/switch/:id", display.SwitchSimulation)
	api.Router.GET("/user/delete/:id", display.DeleteSimulation)
	api.Router.GET("/user/restart/:id", display.RestartSimulation)

	api.Router.GET("/", display.ShowIndexPage)
	api.Router.GET("/data/", display.DataHandler)
	api.Router.GET("/user/dashboard", display.UserDashboard)
	api.Router.GET("/admin/dashboard", display.AdminDashboard)
	api.Router.GET("/admin/reset", display.AdminReset)

	api.Initialise()
	api.ListData()

	api.Router.Run() // Run the server

}
