// api.data.go
// DataObject is the intermediary between the client and the server.

package api

import (
	"capfront/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Defines a data object to be synchronised with the server
// ApiUrl is the endpoint on the server
// DataList is the local client storage for the data
type DataObject struct {
	ApiUrl   string
	ApiKey   string
	DataList interface{}
}

// Retrieves the data for a single simulation object from the server.
//
//	Unmarshals the server response into the object
//
//	Return true if it worked
//
//	Return false if there was an error of any kind
//
//	TODO return an error code instead of a boolean
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

// Currently used only by Initialise.
func FetchGlobalObject(url string, target any) bool {
	resp, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Trace(utils.BrightWhite, fmt.Sprint("Error constructing server request", err))
		return false
	}

	resp.Header.Add("x-api-key", utils.ADMINKEY)
	client := &http.Client{Timeout: time.Second * 2} // Timeout after 2 seconds
	res, _ := client.Do(resp)
	if res == nil {
		utils.Trace(utils.BrightWhite, "Server did not respond")
		return false
	}

	if res.StatusCode != 200 {
		utils.Trace(utils.BrightWhite, fmt.Sprintf("Server rejected admin request with status code %d\n", res.StatusCode))
		return false
	}

	body_as_string, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	jsonErr := json.Unmarshal(body_as_string, target)
	if jsonErr != nil {
		utils.Trace(utils.BrightWhite, fmt.Sprint("Could not unmarshal the server response:\n", string(body_as_string)))
		return false
	}
	utils.Trace(utils.BrightWhite, fmt.Sprintf("Request for data from endpoint %s accepted\n", url))
	return true
}
