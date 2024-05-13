// models.users.go
// User models and related objects

package models

import "capfront/api"

// Full details of a user.
type User struct {
	UserName            string         `json:"username"`              // Repeats the key in the map,for ease of use
	ApiKey              string         `json:"api_key"`               // The api key allocated to this user
	CurrentSimulationID int            `json:"current_simulation_id"` // the id of the simulation that this user is currently using
	LastVisitedPage     string         // Remember what the user was looking at (used when an action is requested)
	Datasets            []*Dataset     // Repository for the data objects generated during the simulation
	TimeStamp           int            // Indexes Datasets. Selects the stage that the simulation has reached
	ViewedTimeStamp     int            // Indexes Datasets. Selects what the user is viewing
	ComparatorTimeStamp int            // Indexes Datasets. Selects what Viewed items are compared with.
	Sim                 api.DataObject // Details of the current simulation
}

var Users = make(map[string]*User) // Every user's simulation data
var AdminUserList []User           // Basic user data, for use by the administrator
type Dataset map[string]api.DataObject

// Constructor for a dataset object.
// Contains and defines all the standard objects of a simulation.
func NewDataset(apiKey string) Dataset {
	return map[string]api.DataObject{
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
	}
}

// Constructor for a standard initial User.
func NewUser(username string, current_simulation_id int, apiKey string) User {
	new_user := User{
		UserName:            username,
		ApiKey:              apiKey,
		CurrentSimulationID: current_simulation_id,
		LastVisitedPage:     "",
		TimeStamp:           0,
		ViewedTimeStamp:     0,
		ComparatorTimeStamp: 0,
		Datasets:            []*Dataset{},
		Sim: api.DataObject{
			ApiUrl:   `simulations/current`,
			ApiKey:   apiKey,
			DataList: new([]Simulation),
		},
	}
	new_dataset := NewDataset(new_user.ApiKey)
	new_user.Datasets = append(new_user.Datasets, &new_dataset)
	return new_user
}

// Wrappers for the object lists.
// The Simulations wrapper is a special case, because the dashboard
// displays a list of user simulations which may be empty.
// If the user has no simulationsm, we make up a fake list with nothing
// in it, to ensure the app can display the dashboard.
func (u User) Simulations() *[]Simulation {
	list := u.Sim.DataList.(*[]Simulation)
	if len(*list) == 0 {
		var fakeList []Simulation = *new([]Simulation)
		return &fakeList
	}
	return list
}

func (u User) Commodities() *[]Commodity {
	return (*u.Datasets[u.ViewedTimeStamp])["commodities"].DataList.(*[]Commodity)
}

func (u User) CommodityViews() *[]CommodityView {
	v := (*u.Datasets[u.ViewedTimeStamp])["commodities"].DataList.(*[]Commodity)
	c := (*u.Datasets[u.ComparatorTimeStamp])["commodities"].DataList.(*[]Commodity)
	return NewCommodityViews(v, c)
}

func (u User) Industries() *[]Industry {
	return (*u.Datasets[u.ViewedTimeStamp])["industries"].DataList.(*[]Industry)
}

func (u User) Classes() *[]Class {
	return (*u.Datasets[u.ViewedTimeStamp])["classes"].DataList.(*[]Class)
}

// Wrapper for the IndustryStockList
func (u User) IndustryStocks() *[]Industry_Stock {
	return (*u.Datasets[u.ViewedTimeStamp])["industry stocks"].DataList.(*[]Industry_Stock)
}

// Wrapper for the ClassStockList
func (u User) ClassStocks() *[]Class_Stock {
	return (*u.Datasets[u.ViewedTimeStamp])["class stocks"].DataList.(*[]Class_Stock)
}

// Wrapper for the TraceList
func (u User) Traces() *[]Trace {
	return (*u.Datasets[u.ViewedTimeStamp])["trace"].DataList.(*[]Trace)
}
