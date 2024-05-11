// display.objects.go
// handlers to display the objects of the simulation on the user's browser

package display

import (
	"capfront/api"
	"capfront/fetch"
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

var Router *gin.Engine = gin.New()

// Helper function for most display handlers. Main purpose is to maintain
// synchronisation between the data on the server and the local copy.
//
//	Returns the username and any error.
//	Sets 'LastVisitedPage' so we can return here after an action.
func synchWithServer(ctx *gin.Context) (string, error) {
	// Comment for less detailed diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		utils.Trace(utils.Yellow, fmt.Sprintf(" UserStatus was called from %s#%d\n", file, no))
	}

	// find out what the browser knows
	username := utils.GUESTUSER
	user := models.Users[username]
	if user.ApiKey == "" {
		utils.Trace(utils.Red, "ERROR: User has no api key\n")
		return username, nil
	}

	// find out what the server knows
	body, err := api.ServerRequest(user.ApiKey, `admin/user/`+username)
	if err != nil {
		log.Printf("The server was upset, and replied :%v\n", err)
		return username, err
	}
	synched_user := new(models.User)          // Isolated user record, not added to the list of users.
	err = json.Unmarshal(body, &synched_user) // it's only there as a receptacle for the server data.
	if err != nil {
		log.Printf("We couldn't make sense of what the server says about user %s", username)
		return username, err
	}
	// userDetails, _ := json.MarshalIndent(synched_user, " ", " ")
	// utils.Trace(utils.Yellow, fmt.Sprintf("The server sent this user record: %s\n", string(userDetails)))

	utils.Trace(utils.Yellow, fmt.Sprintf(
		"The server sent a user record which says the current simulation is %d; the client says it is %d\n",
		synched_user.CurrentSimulationID,
		models.Users[username].CurrentSimulationID,
	))

	// Update client data from the server if need be.
	if models.Users[username].CurrentSimulationID != synched_user.CurrentSimulationID {
		utils.Trace(utils.Yellow, fmt.Sprintf("We are out of synch. Server thinks our simulation is %d and client says it is %d\n",
			synched_user.CurrentSimulationID,
			models.Users[username].CurrentSimulationID))
		// Resynchronise
		if !fetch.FetchUserObjects(ctx, username) {
			utils.Trace(utils.Red, fmt.Sprintf("ERROR: Could not retrieve data for user %s\n", username))
			return username, nil
		}
	}

	models.Users[username].LastVisitedPage = ctx.Request.URL.Path
	utils.Trace(utils.Yellow, fmt.Sprintf("User %s is good to go\n", username))
	return username, nil
}

// Helper function to list out users and templates
func ListData() {
	fmt.Printf("\nTemplateList has %d elements which are:\n", len(models.TemplateList))
	for i := 0; i < len(models.TemplateList); i++ {
		fmt.Println(models.TemplateList[i])
	}

	fmt.Printf("\nAdminUserList has %d elements which are:\n", len(models.AdminUserList))
	for i := 0; i < len(models.AdminUserList); i++ {
		fmt.Println(models.AdminUserList[i])
	}

	fmt.Println("\nUsers", len(models.Users))
	m, _ := json.MarshalIndent(models.Users, " ", " ")
	fmt.Println(string(m))

}

// helper function to obtain the state of the current simulation
// to be replaced by inline call
func Get_current_state(username string) string {
	return models.Users[username].Get_current_state()
}

// helper function to set the state of the current simulation.
// To be replaced by inline call
func set_current_state(username string, new_state string) {
	u, ok := models.Users[username]
	if !ok {
		return
	}
	u.Set_current_state(new_state)
}

// display all commodities in the current simulation
func ShowCommodities(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Commodities")
		return
	}
	user := models.Users[username]
	state := user.Get_current_state()

	// diagnostic to see what gin.H does

	fmt.Println(
		gin.H{
			"Title":       "Commodities",
			"commodities": user.Commodities(),
			"username":    username,
			"state":       state,
		},
	)

	ctx.HTML(http.StatusOK, "commodities.html", gin.H{
		"Title":       "Commodities",
		"commodities": user.Commodities(),
		"username":    username,
		"state":       state,
	})
}

// display all industries in the current simulation
func ShowIndustries(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Industries ")
		return
	}

	state := models.Users[username].Get_current_state()
	ctx.HTML(http.StatusOK, "industries.html", gin.H{
		"Title":      "Industries",
		"industries": models.Users[username].Industries(),
		"username":   username,
		"state":      state,
	})
}

// display all classes in the current simulation
func ShowClasses(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Classes ")
		return
	}
	state := models.Users[username].Get_current_state()
	ctx.HTML(http.StatusOK, "classes.html", gin.H{
		"Title":    "Classes",
		"classes":  models.Users[username].Classes(),
		"username": username,
		"state":    state,
	})
}

// Display one specific commodity
func ShowCommodity(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display one Commodity ")
		return
	}

	state := models.Users[username].Get_current_state()
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
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display one industry ")
		return
	}

	state := models.Users[username].Get_current_state()
	id, _ := strconv.Atoi(ctx.Param("id"))
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
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display one class ")
		return
	}

	state := models.Users[username].Get_current_state()
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
		utils.Trace(utils.Yellow, fmt.Sprintf(" ShowIndexPage was called from %s#%d\n", file, no))
	}
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display the Index Page ")
		return
	}
	state := models.Users[username].Get_current_state()
	clist := *models.Users[username].Commodities()
	ilist := *models.Users[username].Industries()
	cllist := *models.Users[username].Classes()

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Title":       "Economy",
		"industries":  ilist,
		"commodities": clist,
		"classes":     cllist,
		"username":    username,
		"state":       state,
	})
}

// Fetch the trace from the local database
func ShowTrace(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Trace records ")
		return
	}

	state := models.Users[username].Get_current_state()
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

// Display all templates, and all simulations belonging to this user,
// in the user dashboard.
func UserDashboard(ctx *gin.Context) {
	if _, file, no, ok := runtime.Caller(1); ok {
		utils.Trace(utils.Yellow, fmt.Sprintf(" User Dashboard was called from %s line #%d\n", file, no))
	}

	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display the User dashboard ")
		return
	}

	state := models.Users[username].Get_current_state()
	slist := *models.Users[username].Simulations()

	ctx.HTML(http.StatusOK, "user-dashboard.html", gin.H{
		"Title":       "Dashboard",
		"simulations": slist,
		"templates":   models.TemplateList,
		"username":    username,
		"state":       state,
	})
}

// Diagnostic endpoint to display the data in the system.
func DataHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, models.Users)
}

// TODO not working yet
func SwitchSimulation(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Data Listing ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to switch to simulation %d", username, id))
	ctx.HTML(http.StatusOK, "notready.html", gin.H{
		"Title": "Not Ready",
	})
}

// TODO not working yet
func DeleteSimulation(ctx *gin.Context) {
}

// TODO not working yet
func RestartSimulation(ctx *gin.Context) {
}

// display all industry stocks in the current simulation
func ShowIndustryStocks(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display industry Stocks ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to show industry stocks %d", username, id))

	state := models.Users[username].Get_current_state()
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
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Class Stocks ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to show class stocks %d", username, id))
	state := models.Users[username].Get_current_state()
	cslist := *models.Users[username].ClassStocks()

	ctx.HTML(http.StatusOK, "class_stocks.html", gin.H{
		"Title":    "Class Stocks",
		"stocks":   cslist,
		"username": username,
		"state":    state,
	})
}
