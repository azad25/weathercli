package main

import (
	"fmt"
	"weathercli/utils"
	"weathercli/weather"
)

func main() {
	utils.LoadEnv()

	var city string
	fmt.Print("Enter city name: ")
	fmt.Scanln(&city)

	coordinates, err := weather.GetCoordinates(city)
	if err != nil {
		fmt.Printf("Error fetching coordinates: %v\n", err)
		return
	}

	currentWeather, err := weather.GetCurrentWeather(coordinates.Lat, coordinates.Lon)
	if err != nil {
		fmt.Printf("Error fetching current weather data: %v\n", err)
		return
	}

	fmt.Printf("Current weather in %s:\n", city)
	fmt.Printf("Temperature: %.2f°C\n", currentWeather.CurrentTemperature)
	fmt.Printf("Description: %s\n", currentWeather.Description)

	forecast, err := weather.GetForecast(coordinates.Lat, coordinates.Lon)
	if err != nil {
		fmt.Printf("Error fetching weather forecast: %v\n", err)
		return
	}

	fmt.Printf("\nHourly Weather Forecast for %s:\n", city)
	for _, hour := range forecast.Hourly {
		fmt.Printf("%s: Temperature: %.2f°C, Description: %s\n",
			hour.Time, hour.Temperature, hour.Description)
	}

	fmt.Printf("\n7-Day Weather Forecast for %s:\n", city)
	for _, day := range forecast.Daily {
		fmt.Printf("%s: Temperature: %.2f°C, Description: %s\n",
			day.Date, day.Temperature, day.Description)
	}
}
