// fetch.go
// separate package for getting data from the api.
// seems to be necessary to keep this separate to avoid circular imports.
// Arrived at pragmatically - haven't really thought it through.

package fetch

import (
	"capfront/models"
	"capfront/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Iterates through ApiList to refresh all user objects for one user
//
// Returns: false if any table fails.
//
// Returns: true if all tables succeed.
func FetchUserObjects(ctx *gin.Context, username string) bool {
	user := models.Users[username]
	if !user.Sim.Fetch() {
		utils.Trace(utils.Red, "Sim did not fetch\n")
	}
	for key, value := range user.Datasets {
		fmt.Println("Processing", key)
		if !value.Fetch() {
			utils.Trace(utils.Red, "Something went wrong\n")
		}
	}

	utils.Trace(utils.Cyan, "Refresh complete\n")
	return true
}

// Currently used only by Initialise.
// TODO replace with FetchUserObjects.
func FetchGlobalObject(url string, target any) bool {
	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Output(1, fmt.Sprint("Error constructing server request", err))
		return false
	}

	resp.Header.Add("x-api-key", utils.ADMINKEY)
	client := &http.Client{Timeout: time.Second * 2} // Timeout after 2 seconds
	res, _ := client.Do(resp)
	if res == nil {
		log.Output(1, "Server did not respond")
		return false
	}

	if res.StatusCode != 200 {
		log.Output(1, "Server rejected admin request")
		return false
	}

	body_as_string, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	jsonErr := json.Unmarshal(body_as_string, target)
	if jsonErr != nil {
		log.Output(1, fmt.Sprint("Could not unmarshal the server response:\n", string(body_as_string)))
		return false
	}
	log.Output(1, "Request for server data accepted")
	return true
}

// Runs once at startup.
// Retrieve users and templates from the server database.
func Initialise() {
	// Retrieve the templates on the server
	if !FetchGlobalObject(utils.APISOURCE+`templates/templates`, &models.TemplateList) {
		log.Fatal("Could not retrieve templates information from the server. Stopping")
	}

	// Retrieve users on the server
	if !FetchGlobalObject(utils.APISOURCE+`admin/users`, &models.AdminUserList) {
		log.Fatal("Could not retrieve user information from the server. Stopping")
	}

	// transfer the list to the user map
	for _, item := range models.AdminUserList {
		user := models.NewUser(item.UserName, item.CurrentSimulationID, item.ApiKey)
		models.Users[item.UserName] = &user
	}

}
