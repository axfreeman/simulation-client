package main

import (
	"capfront/api"
	"capfront/display"
	"capfront/models"
	"capfront/utils"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// Runs once at startup.
// Retrieve users and templates from the server database.
func Initialise() {
	// Retrieve users on the server
	if !api.FetchAdminObject(utils.APISOURCE+`admin/users`, `users`) {
		log.Fatal("Could not retrieve user information from the server. Stopping")
	}
	for _, item := range models.AdminUserList {
		user := models.User{UserName: item.UserName, CurrentSimulationID: item.CurrentSimulationID, ApiKey: item.ApiKey}
		models.Users[item.UserName] = &user
	}

	// Retrieve the templates on the server
	if !api.FetchAdminObject(utils.APISOURCE+`templates/templates`, `templates`) {
		log.Fatal("Could not retrieve templates information from the server. Stopping")
	}
}

// Helper function to list out users and templates
func ListData() {
	fmt.Printf("\nTemplateList has %d elements which are:\n", len(models.TemplateList))
	for i := 0; i < len(models.TemplateList); i++ {
		fmt.Println(models.TemplateList[i])
	}
	m, _ := json.MarshalIndent(models.Users, " ", " ")
	fmt.Println(string(m))
}

func main() {
	r := gin.New()
	r.LoadHTMLGlob("./templates/**/*") // load the templates

	fmt.Println("Welcome to capitalism")

	r.GET("/action/:action", display.ActionHandler)

	r.GET("/commodities", display.ShowCommodities)
	r.GET("/industries", display.ShowIndustries)
	r.GET("/classes", display.ShowClasses)
	r.GET("/industry_stocks", display.ShowIndustryStocks)
	r.GET("/class_stocks", display.ShowClassStocks)
	r.GET("/trace", display.ShowTrace)

	r.GET("/industry/:id", display.ShowIndustry)
	r.GET("/commodity/:id", display.ShowCommodity)
	r.GET("/class/:id", display.ShowClass)

	r.GET("/user/create/:id", display.CreateSimulation)
	r.GET("/user/switch/:id", display.SwitchSimulation)
	r.GET("/user/delete/:id", display.DeleteSimulation)
	r.GET("/user/restart/:id", display.RestartSimulation)

	r.GET("/", display.ShowIndexPage)
	r.GET("/data/", display.DataHandler)
	r.GET("/user/dashboard", display.UserDashboard)
	r.GET("/admin/dashboard", display.AdminDashboard)
	r.GET("/admin/reset", display.AdminReset)

	Initialise()
	ListData()

	r.Run() // Run the server

}
