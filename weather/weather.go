package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"weathercli/utils"
)

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type CurrentWeatherResponse struct {
	CurrentTemperature float64 `json:"current_temperature"`
	Description        string  `json:"description"`
}

type HourlyForecast struct {
	Time        string  `json:"time"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

type DailyForecast struct {
	Date        string  `json:"date"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

type ForecastResponse struct {
	Hourly []HourlyForecast `json:"hourly"`
	Daily  []DailyForecast  `json:"daily"`
}

func GetCoordinates(city string) (Coordinates, error) {
	apiKey := utils.GetEnv("OPENCAGE_API_KEY")
	baseUrl := "https://api.opencagedata.com/geocode/v1/json"

	url := fmt.Sprintf("%s?q=%s&key=%s", baseUrl, url.QueryEscape(city), apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return Coordinates{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("failed to get coordinates: %s", resp.Status)
	}

	var result struct {
		Results []struct {
			Geometry struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"geometry"`
		} `json:"results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return Coordinates{}, err
	}

	if len(result.Results) == 0 {
		return Coordinates{}, fmt.Errorf("no results found for city: %s", city)
	}

	coordinates := Coordinates{
		Lat: result.Results[0].Geometry.Lat,
		Lon: result.Results[0].Geometry.Lng,
	}

	return coordinates, nil
}

func GetCurrentWeather(lat, lon float64) (CurrentWeatherResponse, error) {
	var weatherResponse struct {
		CurrentWeatherResponse
		Hourly struct {
			Temperature_2m []float64 `json:"temperature_2m"`
			WeatherCode    []int     `json:"weathercode"`
		} `json:"hourly"`
	}

	baseUrl := utils.GetEnv("BASE_URL")

	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current_weather=true&hourly=temperature_2m,weathercode", baseUrl, lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return weatherResponse.CurrentWeatherResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weatherResponse.CurrentWeatherResponse, fmt.Errorf("failed to get weather data: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return weatherResponse.CurrentWeatherResponse, err
	}

	weatherResponse.CurrentTemperature = weatherResponse.Hourly.Temperature_2m[0]
	weatherResponse.Description = mapWeatherCodeToDescription(weatherResponse.Hourly.WeatherCode[0])

	return weatherResponse.CurrentWeatherResponse, nil
}

func GetForecast(lat, lon float64) (ForecastResponse, error) {
	var forecastResponse struct {
		Hourly struct {
			Time           []string  `json:"time"`
			Temperature_2m []float64 `json:"temperature_2m"`
			WeatherCode    []int     `json:"weathercode"`
		} `json:"hourly"`
		Daily struct {
			Time               []string  `json:"time"`
			Temperature_2m_max []float64 `json:"temperature_2m_max"`
			WeatherCode        []int     `json:"weathercode"`
		} `json:"daily"`
	}

	baseUrl := utils.GetEnv("BASE_URL")

	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&hourly=temperature_2m,weathercode&daily=temperature_2m_max,weathercode&timezone=auto", baseUrl, lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return ForecastResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ForecastResponse{}, fmt.Errorf("failed to get weather forecast: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&forecastResponse)
	if err != nil {
		return ForecastResponse{}, err
	}

	var hourlyForecasts []HourlyForecast
	for i := 0; i < len(forecastResponse.Hourly.Time); i++ {
		hourlyForecast := HourlyForecast{
			Time:        forecastResponse.Hourly.Time[i],
			Temperature: forecastResponse.Hourly.Temperature_2m[i],
			Description: mapWeatherCodeToDescription(forecastResponse.Hourly.WeatherCode[i]),
		}
		hourlyForecasts = append(hourlyForecasts, hourlyForecast)
	}

	var dailyForecasts []DailyForecast
	for i := 0; i < len(forecastResponse.Daily.Time); i++ {
		parsedDate, _ := time.Parse("2006-01-02", forecastResponse.Daily.Time[i])
		humanReadableDate := parsedDate.Weekday().String()
		dailyForecast := DailyForecast{
			Date:        humanReadableDate,
			Temperature: forecastResponse.Daily.Temperature_2m_max[i],
			Description: mapWeatherCodeToDescription(forecastResponse.Daily.WeatherCode[i]),
		}
		dailyForecasts = append(dailyForecasts, dailyForecast)
	}

	return ForecastResponse{Hourly: hourlyForecasts, Daily: dailyForecasts}, nil
}

func mapWeatherCodeToDescription(code int) string {
	weatherDescriptions := map[int]string{
		0:  "Clear sky",
		1:  "Mainly clear",
		2:  "Partly cloudy",
		3:  "Overcast",
		45: "Fog",
		48: "Depositing rime fog",
		51: "Light drizzle",
		53: "Moderate drizzle",
		55: "Dense drizzle",
		56: "Light freezing drizzle",
		57: "Dense freezing drizzle",
		61: "Slight rain",
		63: "Moderate rain",
		65: "Heavy rain",
		66: "Light freezing rain",
		67: "Heavy freezing rain",
		71: "Slight snow fall",
		73: "Moderate snow fall",
		75: "Heavy snow fall",
		77: "Snow grains",
		80: "Slight rain showers",
		81: "Moderate rain showers",
		82: "Violent rain showers",
		85: "Slight snow showers",
		86: "Heavy snow showers",
		95: "Thunderstorm",
		96: "Thunderstorm with slight hail",
		99: "Thunderstorm with heavy hail",
	}

	return weatherDescriptions[code]
}
