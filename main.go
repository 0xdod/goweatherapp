package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/urfave/cli.v1"
)

const (
	geoCodeKey = "2fbfc74d49f540bd8fb968932a46e6ee"
	darkSkyKey = "4d055fbc5a091059c7b10273c2929555"
)

type geometry struct {
	Lat, Lng float64
}

type result struct {
	Formatted string `json:"formatted"`
	geometry  `json:"geometry"`
}

type status struct {
	Code    int
	Message string
}

type geoData struct {
	Results []result `json:"results"`
	status  `json:"status"`
	Total   int `json:"total_results"`
}

type currently struct {
	Temperature, ApparentTemperature float64
}
type weatherReport struct {
	currently `json:"currently"`
}

func main() {
	app := cli.NewApp()
	app.Name = "Weather App"
	app.Usage = "Fetches weather details of a location"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address, a",
			Value: "Lagos Nigeria",
			Usage: "Address of Location to fetch weather data",
		},
	}
	app.Action = func(c *cli.Context) error {
		addr := url.QueryEscape(c.GlobalString("address"))
		geocodeURL := "https://api.opencagedata.com/geocode/v1/json?q=" + addr + "&key=" + geoCodeKey
		geoCodeRequest(geocodeURL)
		return nil
	}
	app.Run(os.Args)
}

func geoCodeRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var data geoData
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		log.Fatal("Unable to decode response data")
	}
	if data.Total == 0 || data.status.Code == 400 {
		log.Fatal("Bad Request, Unable to process")
	} else {
		var (
			lat = fmt.Sprint(data.Results[0].Lat)
			lng = fmt.Sprint(data.Results[0].Lng)
		)
		fmt.Println("Fetching weather result for " + data.Results[0].Formatted)
		getWeather(lat, lng)
	}
}

func getWeather(lat string, lng string) {
	resp, err := http.Get("https://api.darksky.net/forecast/" + darkSkyKey + "/" + lat + "," + lng)
	if err != nil {
		log.Fatal("Unable to get weather details: ", err)
	}
	defer resp.Body.Close()
	var wr weatherReport
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&wr); err != nil {
		log.Fatal("Unable to decode response data")
	}
	var (
		msg     string
		temp    = wr.Temperature
		appTemp = wr.ApparentTemperature
	)
	if temp == appTemp {
		msg = "Temperature is " + fmt.Sprintf("%.2f", ftc(temp))
	} else {
		msg = "Temperature currently is " + fmt.Sprintf("%.2f", ftc(temp)) + " but it feels like " + fmt.Sprintf("%.2f", ftc(appTemp)) + " here!"
	}
	fmt.Println(msg)
}

func ftc(t float64) float64 {
	return ((t - 32) / 1.8)
}
