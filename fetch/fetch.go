package fetch

import (
	"capfront/models"
	"capfront/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// Iterates through ApiList to refresh all user objects for one user
//
// Returns: false if any table fails.
//
// Returns: true if all tables succeed.
func FetchUserObjects(ctx *gin.Context, username string) bool {
	// Comment for shorter diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		utils.Trace(utils.Cyan, fmt.Sprintf(" Fetch user objects was called from %s#%d\n", file, no))
	}

	user := models.Users[username]
	utils.Trace(utils.Cyan, "Testing Data Objects")
	if !user.Sim.Fetch() {
		utils.Trace(utils.Cyan, "Sim did not fetch\n")
	}
	if !user.Com.Fetch() {
		utils.Trace(utils.Cyan, "Com did not fetch\n")
	}
	if !user.Ind.Fetch() {
		utils.Trace(utils.Cyan, "Ind did not fetch\n")
	}
	if !user.Cla.Fetch() {
		utils.Trace(utils.Cyan, "Cla did not fetch\n")
	}
	if !user.Isl.Fetch() {
		utils.Trace(utils.Cyan, "Isl did not fetch\n")
	}
	if !user.Csl.Fetch() {
		utils.Trace(utils.Cyan, "Csl did not fetch\n")
	}
	if !user.Tra.Fetch() {
		utils.Trace(utils.Cyan, "Tra did not fetch\n")
	}

	// Comment for shorter diagnostics
	// s, _ := json.MarshalIndent(models.Users[username], "  ", "  ")
	// fmt.Printf("User record after creating the simulation is %s\n", string(s))

	utils.Trace(utils.Cyan, "Refresh complete\n")
	return true
}

// Populates an object.
// Currently used only by Initialise, but could be generalised
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
