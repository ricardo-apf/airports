package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
)

// Results type to mimic the json from the query
type Results struct {
	Rows []struct {
		Fields struct {
			Lat      float64 `json:"lat"`
			Lon      float64 `json:"lon"`
			Name     string  `json:"name"`
			Distance float64
		} `json:"fields"`
	} `json:"rows"`
}

func main() {
	area := []float64{0, 0, 0, 0}
	//Some attempt at ASCII art, because why not !
	fmt.Println("********************************************\n***** Airports in a given area locator *****\n********************************************")
	// call the function to get the user inputs
	area[0] = getInput("Type the value for the Longitude start coordinate", "Longitude", 180)
	area[1] = getInput("Type the value for the Longitude end coordinate", "Longitude", 180)
	area[2] = getInput("Type the value for the Latitude start coordinate", "Latitude", 90)
	area[3] = getInput("Type the value for the Latitude end coordinate", "Latitude", 90)

	// call the function to ensure that the coordinates are in the right order
	area = sortCoordinates(area)
	//call the function to calculate the center of the search area
	center := calcCenter(area)
	//call the function to get the list of airports from the search area
	airportList := getJSON(area)
	//call the function to get the distance from every airport in the list to the center of the search area
	airportList = getDistance(airportList, center)
	//call the function to sort the results from the closest to the farthest of the center of the search area
	airportList = sortClosest(airportList)

	//validate if any results were found
	if len(airportList.Rows) == 0 {
		fmt.Println("No ariports found in the designated area")
	} else {
		// print the results
		for _, v := range airportList.Rows {
			fmt.Printf("Airport\t%v\nLongitude\t%v\nLatitude\t%v\n\n", v.Fields.Name, v.Fields.Lon, v.Fields.Lat)
		}
	}
}


//function to get the coordinates of the search area. Longitude = (lon1 + lon2)/2, latitude = (lat1 + lat2)/2
func calcCenter(calcArea []float64) []float64 {
	retCenter := []float64{0, 0}
	retCenter[0] = (calcArea[0] + calcArea[1]) / 2
	retCenter[1] = (calcArea[2] + calcArea[3]) / 2
	return retCenter
}

// function to build the url for the query and pass the body to the results type.
func getJSON(getArea []float64) Results {
	// declare a var airports of the Response type
	var retAirports Results
	// build the url with the user provided data
	unformatedURL := fmt.Sprintf("lon:[%v TO %v] AND lat:[%v TO %v]", getArea[0], getArea[1], getArea[2], getArea[3])
	// encode the url
	formatedURL := url.QueryEscape(unformatedURL)
	// send the query to the API
	response, _ := http.Get("https://mikerhodes.cloudant.com/airportdb/_design/view1/_search/geo?q=" + formatedURL)
	// get the body from the response
	body, _ := ioutil.ReadAll(response.Body)
	// unmarshal the json to the var airports of the Response type
	json.Unmarshal([]byte(body), &retAirports)
	return retAirports
}

// function to calculate the distance to the center of the search area
func getDistance(retAirportList Results, calCenter []float64) Results {
	for i, v := range retAirportList.Rows {
		//Euclidean distance: ((x, y), (a, b)) = √((x - a)² + (y - b)²)
		retAirportList.Rows[i].Fields.Distance = math.Sqrt(math.Pow((calCenter[0]-v.Fields.Lon), 2) + math.Pow((calCenter[1]-v.Fields.Lat), 2))
	}
	return retAirportList
}

// function sorting the results from the closest to the farthest of the center of the search area using the Slice function from the sort package
func sortClosest(retAirportList Results) Results {
	sort.Slice(retAirportList.Rows, func(i, j int) bool {
		return retAirportList.Rows[i].Fields.Distance < retAirportList.Rows[j].Fields.Distance
	})
	return retAirportList
}

//function to grab the user input for the search range and validate it
func getInput(question string, kind string, limits float64) float64 {
	var retCoordinate float64
	for {
		fmt.Printf("%v :\n", question)
		in := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanf(in, "%f\n", &retCoordinate)
		//validate if the previous line returned an error, or if the values are outside of the lon/lat range
		// if they are, ask for a new input
		if err != nil || retCoordinate > limits || retCoordinate < (limits*-1) {
			fmt.Printf("\nPlease insert a valid %v: Values between -%v and %v\n", kind, limits, limits) //input validation
			continue
		}
		break
	}
	return retCoordinate
}

// function to ensure that the coordinates are in the correct ortder without bothering the user with a new input
func sortCoordinates(retArea []float64) []float64 {
	if retArea[0] > retArea[1] {
		temp := retArea[1]
		retArea[1] = retArea[0]
		retArea[0] = temp
	}
	if retArea[2] > retArea[3] {
		temp := retArea[3]
		retArea[3] = retArea[2]
		retArea[2] = temp
	}
	return retArea
}
