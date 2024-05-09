// api.server.go
// container for interacting with remote server

package api

import (
	"bytes"
	"capfront/colour"
	"capfront/logging"
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

// Prepare and send a request for a protected service to the server
// using the user's api key.
//
//	ctx is the context of a handler.
//	username is the name of the user requesting the service.
//	description is a user-friendly name for the requested action.
//	relativePath is appended to the server URL to tell the server what to do.
func ServerRequest(username string, description string, relativePath string) ([]byte, error) {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logging.Trace(colour.Cyan, fmt.Sprintf(" ServerRequest was called from %s#%d\n", file, no))
		logging.Trace(colour.Cyan, fmt.Sprintf(colour.Cyan+" Username was %s, description %s, relativePath %s, APIURL %s\n"+colour.Reset, username, description, relativePath, utils.APISOURCE))
	}

	user, ok := models.Users[username]
	if !ok {
		return nil, fmt.Errorf(" User %s is not in the local database", username)
	}

	url := utils.APISOURCE + relativePath
	body, _ := json.Marshal(models.RequestData{User: username}) // (overkill diagnostic? - not actually needed)
	resp, err := http.NewRequest("GET", url, bytes.NewBuffer(body))

	if err != nil {
		log.Output(1, fmt.Sprintf("Error %v for user %s from URL %s for resource %s \n", err, username, url, description))
		return nil, err
	}

	// logging.Trace(colour.Cyan, fmt.Sprintf(" Api key is [%s]\n", user.ApiKey))
	resp.Header.Add("Content-Type", "application/json")
	resp.Header.Set("User-Agent", "Capitalism reader")
	resp.Header.Add("x-api-key", user.ApiKey)

	client := &http.Client{Timeout: time.Second * 5} // Timeout after 5 seconds
	res, _ := client.Do(resp)
	if res == nil {
		log.Output(1, "Server is down or misbehaving")
		return nil, nil
	}

	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		log.Output(1, fmt.Sprintf("Server rejected request '%s' with status %s\n", description, res.Status))
		logging.Trace(colour.Red, fmt.Sprintf("It said %s\n", string(b)))
		return nil, fmt.Errorf(string(b))
	}

	// Comment for fewer diagnostics
	fmt.Println(colour.Cyan + " Leaving ProtectedServerRequest. Everything seems to have worked." + colour.Reset)

	return b, nil
}

// Contains the information needed to fetch data for one model from the remote server.
// Name is a description, just for diagnostic purposes.
// ApiURL is the endpoint to get the data from the server.
type ApiItem struct {
	Name   string // the data to be obtained
	ApiUrl string // the url to be used in accessing the backend
}

// a list of items needed to fetch data from the remote server
var ApiList = [7]ApiItem{
	{`simulation`, `simulations/`},
	{`commodity`, `commodity/`},
	{`industry`, `industry/`},
	{`class`, `classes/`},
	{`industry_stock`, `stocks/industry`},
	{`class_stock`, `stocks/class`},
	{`trace`, `trace/`},
}

// Iterates through ApiList to refresh all user objects for one user
//
// Returns: false if any table fails.
//
// Returns: true if all tables succeed.
func FetchUserObjects(ctx *gin.Context, username string) bool {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logging.Trace(colour.Cyan, fmt.Sprintf(" Fetch user objects was called from %s#%d\n", file, no))
	}

	// (miss out trace for now - it's too big)
	for i := 0; i < len(ApiList)-1; i++ {
		a := ApiList[i]
		logging.Trace(colour.Cyan, fmt.Sprintf(" FetchUserObjects is fetching API item %d with name %s from URL %s\n", i, a.Name, a.ApiUrl))
		if !FetchAPI(&a, username) {
			logging.Trace(colour.Cyan, "There are no objects to retrieve from the remote server. Do not continue \n")
			return false
		}
	}
	// Comment for shorter diagnostics
	// s, _ := json.MarshalIndent(models.Users[username], "  ", "  ")
	// fmt.Printf("User record after creating the simulation is %s\n", string(s))

	logging.Trace(colour.Cyan, "Refresh complete")
	return true
}

// Fetch the data specified by item for user.
//
//	 item: specifies what is to be retrieved, using which URL
//	 username: the name of the user - serves as an index into the Users map.
//
//		if we got something, return true.
//		if not, for whatever reason, return false.
func FetchAPI(item *ApiItem, username string) (result bool) {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logging.Trace(colour.Cyan, fmt.Sprintf("fetch API was called from %s#%d\n", file, no))
		log.Output(1, fmt.Sprintf("User %s asked to fetch the table named %s from the URL %s\n", username, item.Name, item.ApiUrl))
	}

	var jsonErr error
	body, err := ServerRequest(username, "Fetch Table", item.ApiUrl)

	if err != nil {
		log.Output(1, "ERROR: The server did not send a response; this is a programming error")
		return false
	}

	if len(string(body)) == 0 {
		log.Output(1, "INFORMATION: The server response was empty")
		return false
	}

	log.Output(1, fmt.Sprintf("INFORMATION: The server sent a table of length %d\n", len(string(body))))

	// check for '[]' response (a list with no elements in it)
	if body[0] == 91 && body[1] == 93 {
		log.Output(1, "INFORMATION: The server sent an empty table; this means the user has no simulations yet.")
		return false
	}

	// Populate the user record.
	logging.Trace(colour.Cyan, fmt.Sprintf("Unmarshalling data for user %s into %v\n", username, item.Name))

	switch item.Name {

	case `simulation`:
		jsonErr = json.Unmarshal(body, &models.Users[username].SimulationList)
	case `commodity`:
		jsonErr = json.Unmarshal(body, &models.Users[username].CommodityList)
	case `industry`:
		jsonErr = json.Unmarshal(body, &models.Users[username].IndustryList)
	case `class`:
		jsonErr = json.Unmarshal(body, &models.Users[username].ClassList)
	case `industry_stock`:
		jsonErr = json.Unmarshal(body, &models.Users[username].IndustryStockList)
	case `class_stock`:
		jsonErr = json.Unmarshal(body, &models.Users[username].ClassStockList)
	case `trace`:
		jsonErr = json.Unmarshal(body, &models.Users[username].TraceList)
	default:
		logging.Trace(colour.Red, fmt.Sprintf("Unknown dataset%s ", item.Name))
	}

	if jsonErr != nil {
		logging.Trace(colour.Red, fmt.Sprintf("Failed to unmarshal template json because: %s", jsonErr))
		return false
	}

	logging.Trace(colour.Red, fmt.Sprintf("Data refreshed for user %s\n", username))
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
	log.Output(1, fmt.Sprintf("Request for %s data accepted", target))
	return true
}

// Runs once at startup.
// Retrieve users and templates from the server database.
func Initialise() {
	// Retrieve users on the server
	if !FetchGlobalObject(utils.APISOURCE+`admin/users`, &models.AdminUserList) {
		log.Fatal("Could not retrieve user information from the server. Stopping")
	}
	// transfer the list to the user map
	for _, item := range models.AdminUserList {
		user := models.User{UserName: item.UserName, CurrentSimulationID: item.CurrentSimulationID, ApiKey: item.ApiKey}
		models.Users[item.UserName] = &user
	}

	// Retrieve the templates on the server
	if !FetchGlobalObject(utils.APISOURCE+`templates/templates`, &models.TemplateList) {
		log.Fatal("Could not retrieve templates information from the server. Stopping")
	}
}
