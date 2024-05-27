// fetch.go

// separate package for getting data from the api.
// seems to be necessary to keep this separate from the api package to
// avoid circular imports (culprits are the models and api packages).
//
// Arrived at pragmatically - haven't really thought it through.

package fetch

import (
	"capfront/api"
	"capfront/models"
	"capfront/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// Iterates through ApiList to refresh all user objects for one user
//
//	Returns: false if any table fails.
//	Returns: true if all tables succeed.
func FetchUserObjects(ctx *gin.Context, username string) bool {
	user := models.Users[username]
	if !user.Sim.Fetch() {
		utils.Trace(utils.Red, "Sim did not fetch\n")
	}
	// Reminder: a dataset is a repository for all objects at one stage of the simulation.
	dataSet := *user.Datasets[user.TimeStamp]

	for key, value := range dataSet {
		if !value.Fetch() {
			log.Output(1, fmt.Sprintf("Could not retrieve server data for the new dataset with key %s\n", key))
		}
	}
	utils.Trace(utils.Cyan, "Refresh complete\n")
	return true
}

// Runs once at startup.
// Retrieve users and templates from the server database.
func Initialise() {
	// Retrieve the templates on the server
	if !api.FetchGlobalObject(utils.APISOURCE+`templates/templates`, &models.TemplateList) {
		log.Fatal("Could not retrieve templates information from the server. Stopping")
	}

	// Retrieve users on the server
	if !api.FetchGlobalObject(utils.APISOURCE+`admin/users`, &models.AdminUserList) {
		log.Fatal("Could not retrieve user information from the server. Stopping")
	}

	// Transfer the list to the user map
	for _, item := range models.AdminUserList {
		user := models.NewUser(item.UserName, item.CurrentSimulationID, item.ApiKey)
		models.Users[item.UserName] = &user
	}
}
