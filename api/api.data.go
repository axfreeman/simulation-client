// api.data.go
// DataObject is the intermediary between the client and the server.

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
	DataList interface{}
}

func (d *DataObject) Fetch() bool {
	response, err := ServerRequest(d.ApiKey, d.ApiUrl)

	if err != nil {
		utils.Trace(utils.Red, fmt.Sprintf("ServerRequest produced the error %v\n", err))
		return false
	}

	if len(string(response)) == 0 {
		log.Output(1, "INFORMATION: The server response was empty")
		return false
	}

	// Uncomment for more diagnostics
	// fmt.Printf("Type of the data list is %T\n", d.DataList)

	// Populate the data object
	jsonErr := json.Unmarshal(response, &d.DataList)
	if jsonErr != nil {
		utils.Trace(utils.Red, fmt.Sprintf("Server response could not be unmarshalled: Unmarshal produced the error %v\n", jsonErr))
		utils.Trace(utils.Red, fmt.Sprintf("The server response was %s\n", response))
		return false
	}

	// Uncomment for more diagnostics
	// utils.Trace(utils.Cyan, "Server response was unmarshalled\n")
	return true
}
