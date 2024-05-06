// usergo
// User models and related objects

package models

// Full details of a user.
type User struct {
	UserName            string `json:"username"`              // Repeats the key in the map,for ease of use
	ApiKey              string `json:"api_key"`               // The api key allocated to this user
	CurrentSimulationID int    `json:"current_simulation_id"` // the id of the simulation that this user is currently using
	Message             string // store for messages to be displayed to the user when appropriate
	LastVisitedPage     string // Remember what the user was looking at (used when an action is requested)
	ViewedTimeStamp     int    // Indexes the History field. Selects what the user is viewing
	SimulationList      []Simulation
	CommodityList       []Commodity
	IndustryList        []Industry
	ClassList           []Class
	IndustryStockList   []Industry_Stock
	ClassStockList      []Class_Stock
	TraceList           []Trace
	History             []HistoryItem // the history of the current simulation
}

// Constructor for a standard initial User, containing empty entries.
// for fields to be populated from the server.
func NewUser(username string) User {
	new_datum := User{
		UserName:            username,
		ApiKey:              "",
		CurrentSimulationID: 0,
		LastVisitedPage:     "",
		ViewedTimeStamp:     0,
		History:             make([]HistoryItem, 0),
	}
	return new_datum
}

// Wrappers for the various lists.
//
// Main purpose is to provide for future development.
//
// SimulationList is a special case, because the dashboard displays a list
// of the user's simulations. As a workaround, if the user has none we
// make up a fake list with nothing in it, to ensure the app does not
// crash when displaying the dashboard.
func (u User) Simulations() *[]Simulation {
	if len(u.SimulationList) == 0 {
		var fakeList []Simulation = make([]Simulation, 0)
		return &fakeList
	}
	return &u.SimulationList
}

// Wrapper for the CommodityList
func (u User) Commodities() *[]Commodity {
	return &u.CommodityList
}

// Wrapper for the IndustryList
func (u User) Industries() *[]Industry {
	return &u.IndustryList
}

// Wrapper for the ClassList
func (u User) Classes() *[]Class {
	return &u.ClassList
}

// Wrapper for the IndustryStockList
func (u User) IndustryStocks() *[]Industry_Stock {
	return &u.IndustryStockList
}

// Wrapper for the ClassStockList
func (u User) ClassStocks() *[]Class_Stock {
	return &u.ClassStockList
}

// Wrapper for the TraceList
func (u User) Traces() *[]Trace {
	return &u.TraceList
}

// Format of responses from the server for post requests
// Specifically (so far), login or register.
type ServerMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// contains the details of every user's simulations and their status, accessed by username
var Users = make(map[string]*User)

// List of basic user data, for use by the administrator
var AdminUserList []User
