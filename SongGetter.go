package main

/*
	Comp510
	Project 2 part 1 written and completed by Huy Le
		Take in an input and use input for a search query. Pull data from API and decode.
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {

	api := "https://searchly.asuarez.dev/api/v1/song/search"

	// ask for user input and build a query search url
	fmt.Println("Enter an artist or song name: ")
	var userInput string
	fmt.Scanln(&userInput)
	queryURl := api + "?query=" + userInput

	// retrieve data from API
	response, err := http.Get(queryURl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()     // defer executes at end of function
	if response.StatusCode != 200 { //		.Body implements io.Reader
		log.Fatal("Didn't get 200") // status code 200 is ok success status response code
	}

	// read data from api and assign data to variable to contain []byte values
	rawData, err := ioutil.ReadAll(response.Body) // read response variable which contains API data
	if err != nil {
		log.Fatal(err)
	}

	// decode []byte data into struct
	var data SongData
	json.Unmarshal(rawData, &data) // decode body of byte[] into interface

	reg, _ := regexp.Compile("-") // variable function to remove dash primarily between the artist and song

	// print out results
	for i := 0; i < len(data.Response.Result); i++ {
		fmt.Println(reg.ReplaceAllString(data.Response.Result[i].Name, ""))
	}
}
