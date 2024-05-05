// display.objects.go
// handlers to display the objects of the simulation on the user's browser

package display

import (
	"capfront/api"
	"capfront/colour"
	"capfront/logging"
	"capfront/models"
	"capfront/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Helper function for most display handlers. Main purpose is to maintain
// synchronisation between the data on the server and the local copy.
//
//	Returns the username and any error.
//	Sets 'LastVisitedPage' so we can return here after an action.
func userStatus(ctx *gin.Context) (string, error) {
	// Uncomment for more detailed diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logging.Trace(colour.Cyan, fmt.Sprintf(" UserStatus was called from %s#%d\n", file, no))
	}

	// find out what the browser knows
	username := utils.GUESTUSER

	if models.Users[username].ApiKey == "" {
		logging.Trace(colour.Red, "ERROR: User has no api key\n")
		return username, nil
	}

	// find out what the server knows
	body, err := api.ServerRequest(
		username,
		fmt.Sprintf("Ask the server about user %s", username),
		`admin/user/`+username)
	if err != nil {
		log.Printf("The server was upset, and replied :%v\n", err)
		return username, err
	}
	log.Printf("The server has a record for user %s", username)

	synched_user := models.NewUser(username)  // Isolated user record, not added to the list of users.
	err = json.Unmarshal(body, &synched_user) // it's only there as a receptacle for the server data.
	if err != nil {
		log.Printf("We couldn't make sense of what the server says about user %s", username)
		return username, err
	}

	// If something changed on the server, update all the client data for this user
	log.Printf("The server has told us about user %s", username)
	if models.Users[username].CurrentSimulationID != synched_user.CurrentSimulationID {
		log.Printf("We are out of synch. Server thinks our simulation is %d and client says it is %d",
			synched_user.CurrentSimulationID,
			models.Users[username].CurrentSimulationID)
		if !api.FetchUserObjects(ctx, username) {
			log.Printf("Could not retrieve data for this user")
			return username, nil
		}
	}

	models.Users[username].LastVisitedPage = ctx.Request.URL.Path
	models.Users[username].CurrentSimulationID = synched_user.CurrentSimulationID

	log.Printf("User %s is good to go", username)
	return username, nil
}

// helper function to obtain the state of the current simulation
// if no user is logged in, return null state
func get_current_state(username string) string {
	this_user := models.Users[username]
	if this_user == nil {
		return "NO KNOWN USER"
	}
	this_user_history := this_user.History
	if len(this_user_history) == 0 {
		// User doesn't yet have any simulations
		return "NO SIMULATION YET"
	}

	id := this_user.CurrentSimulationID
	sims := *this_user.Simulations()
	if sims == nil {
		return "UNKNOWN"
	}

	for i := 0; i < len(sims); i++ {
		s := sims[i]
		if s.Id == id {
			return s.State
		}
	}
	return "UNKNOWN"
}

// helper function to set the state of the current simulation
// if we fail it's a programme error so we don't test for that
func set_current_state(username string, new_state string) {
	this_user := models.Users[username]
	id := this_user.CurrentSimulationID
	sims := *this_user.Simulations()
	log.Output(1, fmt.Sprintf("resetting state to %s for user %s", new_state, this_user.UserName))
	for i := 0; i < len(sims); i++ {
		s := &sims[i]
		if (*s).Id == id {
			(*s).State = new_state
			return
		}
		log.Output(1, fmt.Sprintf("simulation with id %d not found", id))
	}
}

// display all commodities in the current simulation
func ShowCommodities(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve Commodities ")
		return
	}
	state := get_current_state(username)

	ctx.HTML(http.StatusOK, "commodities.html", gin.H{
		"Title":       "Commodities",
		"commodities": models.Users[username].Commodities(),
		"username":    username,
		"state":       state,
	})
}

// display all industries in the current simulation
func ShowIndustries(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve Industries ")
		return
	}

	state := get_current_state(username)
	ctx.HTML(http.StatusOK, "industries.html", gin.H{
		"Title":      "Industries",
		"industries": models.Users[username].Industries(),
		"username":   username,
		"state":      state,
	})
}

// display all classes in the current simulation
func ShowClasses(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve Classes ")
		return
	}
	state := get_current_state(username)
	ctx.HTML(http.StatusOK, "classes.html", gin.H{
		"Title":    "Classes",
		"classes":  models.Users[username].Classes(),
		"username": username,
		"state":    state,
	})
}

// Display one specific commodity
func ShowCommodity(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve one Commodity ")
		return
	}

	state := get_current_state(username)
	id, _ := strconv.Atoi(ctx.Param("id"))

	clist := *models.Users[username].Commodities()
	for i := 0; i < len(clist); i++ {
		if id == clist[i].Id {
			ctx.HTML(http.StatusOK, "commodity.html", gin.H{
				"Title":     "Commodity",
				"commodity": clist[i],
				"username":  username,
				"state":     state,
			})
		}
	}
}

// Display one specific industry
func ShowIndustry(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve one industry ")
		return
	}

	state := get_current_state(username)
	id, _ := strconv.Atoi(ctx.Param("id")) //TODO check user didn't do something stupid
	// TODO here and elsewhere create a method to get the simulation
	ilist := *models.Users[username].Industries()
	for i := 0; i < len(ilist); i++ {
		if id == ilist[i].Id {
			ctx.HTML(http.StatusOK, "industry.html", gin.H{
				"Title":    "Industry",
				"industry": ilist[i],
				"username": username,
				"state":    state,
			})
		}
	}
}

// Display one specific class
func ShowClass(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve one class ")
		return
	}

	state := get_current_state(username)
	id, _ := strconv.Atoi(ctx.Param("id")) //TODO check user didn't do something stupid
	list := *models.Users[username].Classes()

	for i := 0; i < len(list); i++ {
		if id == list[i].Id {
			ctx.HTML(http.StatusOK, "class.html", gin.H{
				"Title":    "Class",
				"class":    list[i],
				"username": username,
				"state":    state,
			})
		}
	}
}

// Displays snapshot of the economy

func ShowIndexPage(ctx *gin.Context) {
	// Uncomment for more detailed diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logging.Trace(colour.Cyan, fmt.Sprintf(" ShowIndexPage was called from %s#%d\n", file, no))
	}
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Problem retrieving index page ")
		return
	}
	state := get_current_state(username)

	api.UserMessage = `This is the home page`

	clist := *models.Users[username].Commodities()
	ilist := *models.Users[username].Industries()
	cllist := *models.Users[username].Classes()

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Title":       "Economy",
		"industries":  ilist,
		"commodities": clist,
		"message":     models.Users[username].Message,
		"classes":     cllist,
		"username":    username,
		"state":       state,
	})
}

// Fetch the trace from the local database
func ShowTrace(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve Trace records ")
		return
	}

	state := get_current_state(username)
	tlist := *models.Users[username].Traces()

	ctx.HTML(
		http.StatusOK,
		"trace.html",
		gin.H{
			"Title":    "Simulation Trace",
			"trace":    tlist,
			"username": username,
			"state":    state,
		},
	)
}

// Retrieve all templates, and all simulations belonging to this user, from the local database
// Display them in the user dashboard
func UserDashboard(ctx *gin.Context) {

	if _, file, no, ok := runtime.Caller(1); ok {
		logging.Trace(colour.Cyan, fmt.Sprintf(" User Dashboard was called from %s line #%d\n", file, no))
	}

	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve User Dashboard ")
		return
	}

	state := get_current_state(username)
	slist := *models.Users[username].Simulations()

	ctx.HTML(http.StatusOK, "user-dashboard.html", gin.H{
		"Title":       "Dashboard",
		"simulations": slist,
		"templates":   models.TemplateList,
		"username":    username,
		"state":       state,
	})
}

// a diagnostic endpoint to display the data in the system
func DataHandler(ctx *gin.Context) {
	// username, loginStatus, _ := userStatus(ctx)
	// b, err := json.Marshal(models.Users)
	// if err != nil {
	// 	fmt.Println("Could not marshal the Users object")
	// 	return
	// }
	ctx.JSON(http.StatusOK, models.Users)
}

func SwitchSimulation(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve Data Listing ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to switch to simulation %d", username, id))
	ctx.HTML(http.StatusOK, "notready.html", gin.H{
		"Title": "Not Ready",
	})
}

func DeleteSimulation(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve Data to delete a simulation ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to delete simulation %d", username, id))
	api.ServerRequest(username, "Delete simulation", "simulations/delete/"+ctx.Param("id"))
	api.FetchUserObjects(ctx, username)
	UserDashboard(ctx)
}

func RestartSimulation(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve Data to restart a simulation ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to restart simulation %d", username, id))
	ctx.HTML(http.StatusOK, "notready.html", gin.H{
		"Title": "Not Ready",
	})
}

// display all industry stocks in the current simulation
func ShowIndustryStocks(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve industry Stocks ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to delete simulation %d", username, id))

	state := get_current_state(username)
	islist := *models.Users[username].IndustryStocks()

	ctx.HTML(http.StatusOK, "industry_stocks.html", gin.H{
		"Title":    "Industry Stocks",
		"stocks":   islist,
		"username": username,
		"state":    state,
	})
}

// display all the class stocks in the current simulation
func ShowClassStocks(ctx *gin.Context) {
	username, err := userStatus(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve Class Stocks ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to delete simulation %d", username, id))
	state := get_current_state(username)
	cslist := *models.Users[username].ClassStocks()

	ctx.HTML(http.StatusOK, "class_stocks.html", gin.H{
		"Title":    "Class Stocks",
		"stocks":   cslist,
		"username": username,
		"state":    state,
	})
}
