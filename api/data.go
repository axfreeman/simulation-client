package api

import (
	"capfront/utils"
	"encoding/json"
	"fmt"
	"log"
)

// Defines a data object to be synchronised with the server
// ApiUrl is the endpoint on the server
// DataList is the local client storage for the data
type DataObject struct {
	ApiUrl   string
	ApiKey   string
	DataList any
}

func (d *DataObject) Fetch() bool {
	body, err := ServerRequest(d.ApiKey, d.ApiUrl)

	if err != nil {
		utils.Trace(utils.Cyan, fmt.Sprintf("ServerRequest produced the error %v\n", err))
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

	// Populate the data object.
	utils.Trace(utils.Cyan, fmt.Sprintf("Unmarshalling data %s\n", string(body)))
	jsonErr := json.Unmarshal(body, &d.DataList)
	if jsonErr != nil {
		utils.Trace(utils.Cyan, fmt.Sprintf("Server response could not be unmarshalled: it produced the error %v\n", jsonErr))
		return false
	}
	utils.Trace(utils.Cyan, "Server response was unmarshalled")
	fmt.Println(d.DataList)
	return true
}
