package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// WriteJSONIndented formats JSON and writes it to the http.ResponseWriter.
func WriteJSONIndented(w http.ResponseWriter, data interface{}) error {
	// Marshal data to JSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Format the JSON data
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, dataJSON, "", "  ")
	if err != nil {
		return err
	}

	log.Println("Response", string(prettyJSON.Bytes()))
	// Write the formatted JSON to the response
	_, err = w.Write(prettyJSON.Bytes())
	return err
}
