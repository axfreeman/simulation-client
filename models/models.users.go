// usergo
// User models and related objects

package models

import "capfront/api"

// Full details of a user.
type User struct {
	UserName            string `json:"username"`              // Repeats the key in the map,for ease of use
	ApiKey              string `json:"api_key"`               // The api key allocated to this user
	CurrentSimulationID int    `json:"current_simulation_id"` // the id of the simulation that this user is currently using
	Message             string // store for messages to be displayed to the user when appropriate
	LastVisitedPage     string // Remember what the user was looking at (used when an action is requested)
	ViewedTimeStamp     int    // Indexes the History field. Selects what the user is viewing
	Dataset             map[string]api.DataObject
	Sim                 api.DataObject
	Com                 api.DataObject
	Ind                 api.DataObject
	Cla                 api.DataObject
	Isl                 api.DataObject
	Csl                 api.DataObject
	Tra                 api.DataObject
}

// Constructor for a standard initial User, containing empty entries.
// for fields to be populated from the server.
func NewUser(username string, current_simulation_id int, apiKey string) User {
	new_user := User{
		UserName:            username,
		ApiKey:              apiKey,
		CurrentSimulationID: current_simulation_id,
		LastVisitedPage:     "",
		ViewedTimeStamp:     0,
		Sim: api.DataObject{
			ApiUrl:   `simulations/current`,
			ApiKey:   apiKey,
			DataList: new([]Simulation),
		},
		Com: api.DataObject{
			ApiUrl:   `commodity`,
			ApiKey:   apiKey,
			DataList: new([]Commodity),
		},
		Ind: api.DataObject{
			ApiUrl:   `industry`,
			ApiKey:   apiKey,
			DataList: new([]Industry),
		},
		Cla: api.DataObject{
			ApiUrl:   `classes`,
			ApiKey:   apiKey,
			DataList: new([]Class),
		},
		Isl: api.DataObject{
			ApiUrl:   `stocks/industry`,
			ApiKey:   apiKey,
			DataList: new([]Industry_Stock),
		},
		Csl: api.DataObject{
			ApiUrl:   `stocks/class`,
			ApiKey:   apiKey,
			DataList: new([]Class_Stock),
		},
		Tra: api.DataObject{
			ApiUrl:   `trace`,
			ApiKey:   apiKey,
			DataList: new([]Trace),
		},
		Dataset: map[string]api.DataObject{
			"simulations": {
				ApiUrl:   `simulations/current`,
				ApiKey:   apiKey,
				DataList: new([]Simulation),
			},
			"commodities": {
				ApiUrl:   `commodity`,
				ApiKey:   apiKey,
				DataList: new([]Commodity),
			},
			"industries": {
				ApiUrl:   `industry`,
				ApiKey:   apiKey,
				DataList: new([]Industry),
			},
			"classes": {
				ApiUrl:   `classes`,
				ApiKey:   apiKey,
				DataList: new([]Class),
			},
			"industry stocks": {
				ApiUrl:   `stocks/industry`,
				ApiKey:   apiKey,
				DataList: new([]Industry_Stock),
			},
			"class stocks": {
				ApiUrl:   `stocks/class`,
				ApiKey:   apiKey,
				DataList: new([]Class_Stock),
			},
			"trace": {
				ApiUrl:   `trace`,
				ApiKey:   apiKey,
				DataList: new([]Trace),
			},
		},
	}
	return new_user
}

// Wrappers for the various lists.

// Special case, because the dashboard displays a list of the user's
// simulations. If the user has none we make up a fake list with nothing
// in it, to ensure the app does not crash when displaying the dashboard.
func (u User) Simulations() *[]Simulation {
	list := u.Sim.DataList.(*[]Simulation)
	if len(*list) == 0 {
		var fakeList []Simulation = *new([]Simulation)
		return &fakeList
	}
	return list
}

func (u User) Commodities() *[]Commodity {
	return u.Dataset["commodities"].DataList.(*[]Commodity)
}

func (u User) Industries() *[]Industry {
	return u.Dataset["industries"].DataList.(*[]Industry)
}

func (u User) Classes() *[]Class {
	return u.Dataset["classes"].DataList.(*[]Class)
}

// Wrapper for the IndustryStockList
func (u User) IndustryStocks() *[]Industry_Stock {
	return u.Dataset["industry stocks"].DataList.(*[]Industry_Stock)
}

// Wrapper for the ClassStockList
func (u User) ClassStocks() *[]Class_Stock {
	return u.Dataset["class stocks"].DataList.(*[]Class_Stock)
}

// Wrapper for the TraceList
func (u User) Traces() *[]Trace {
	return u.Dataset["trace"].DataList.(*[]Trace)
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
