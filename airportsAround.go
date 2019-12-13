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
	area := []float64{0, 0, 0, 0, 0, 0, 0}
	//Some attempt at ASCII art, because why not !
	fmt.Println("*********************************\n***** Airports in your area *****\n*********************************")
	// call the function to get the user inputs
	area[4] = getInput("Type the value for your Longitude position", "Longitude", 180)
	area[5] = getInput("Type the value for your Latitude position", "Latitude", 90)
	area[6] = getInput("Type the value for the search distance (in Degrees)", "Distance", 360)

	// Converts the distance to positive value in case the user typed in a negative
	area[6] = math.Abs(area[6])

	//call the function to calculate coordinates of the search area
	area = calcCenter(area)

	//call the function to get the list of airports from the search area
	airportList := getJSON(area)

	//call the function to get the distance from every airport in the list to user location
	airportList = getDistance(airportList, area)

	//call the function to sort the results from the closest to the farthest of the search area
	airportList = sortClosest(airportList)

	//validate if any results were found
	if len(airportList.Rows) == 0 {
		fmt.Println("No ariports found in the designated area")
	} else {
		// print the results
		for _, v := range airportList.Rows {
			// validate if the distance is within the search area - some results might exceed since the search area returned is a square
			if v.Fields.Distance <= area[6] {
				fmt.Printf("Airport\t%v\nLongitude\t%v\nLatitude\t%v\nDistance\t%v\n\n", v.Fields.Name, v.Fields.Lon, v.Fields.Lat, v.Fields.Distance)
			}
		}
	}
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

//function to get the coordinates for the search range.
func calcCenter(calcArea []float64) []float64 {

	calcArea[1] = (calcArea[4] + calcArea[6])
	calcArea[0] = (calcArea[4] - calcArea[6])
	calcArea[3] = (calcArea[5] + calcArea[6])
	calcArea[2] = (calcArea[5] - calcArea[6])

	return calcArea
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
		retAirportList.Rows[i].Fields.Distance = math.Sqrt(math.Pow((calCenter[4]-v.Fields.Lon), 2) + math.Pow((calCenter[5]-v.Fields.Lat), 2))
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
