# Weather CLI Application

This is a command-line interface (CLI) application written in Go (Golang) that fetches current weather data, hourly forecasts, and 7-day forecasts for a given city using the Open-Meteo API.

## Features

- Fetches current weather data including temperature and weather description.
- Provides hourly forecast data for the next 24 hours.
- Displays a 7-day weather forecast.
- Supports any city worldwide.

## Project Structure

The project is structured as follows:

```
weather-cli/
│
├── main.go          # Entry point of the application
├── weather/
│   ├── weather.go   # Weather API and data handling functions
│   └── utils.go     # Utility functions
│
├── .env.sample      # Sample environment variable file
├── README.md        # Project documentation
└── go.mod
```

## How to Install

1. Clone this repository:

    ```bash
    git clone https://github.com/yourusername/weather-cli.git
    ```

2. Navigate to the project directory:

    ```bash
    cd weather-cli
    ```

3. Create a `.env` file based on the provided `.env.sample` file and add your Open-Meteo API key.

4. Build the application:

    ```bash
    go build -o weather-cli
    ```

5. Run the application:

    ```bash
    ./weather-cli
    ```

## API and Packages Used

- **Open-Meteo API**: Used to fetch weather data.
- **Opencage-Geolocation**: Used to fetch coordinates based on city.
- **net/http**: Standard Go library for making HTTP requests.
- **encoding/json**: Standard Go library for JSON encoding and decoding.
- **fmt**: Standard Go library for formatted I/O.
- **time**: Standard Go library for working with dates and times.

## Sample Input and Output

### Input

```
Enter city name: London
```

### Output

```
Current weather in London:
Temperature: 19.00°C
Description: Partly cloudy

Hourly Weather Forecast for London:
Time: 5 PM, Temperature: 19.00°C, Description: Partly cloudy
Time: 6 PM, Temperature: 18.50°C, Description: Partly cloudy
Time: 7 PM, Temperature: 17.00°C, Description: Partly cloudy
...

7-Day Weather Forecast for London:
Date: MONDAY, Temperature: 23.00°C, Description: Partly cloudy
Date: TUESDAY, Temperature: 25.00°C, Description: Clear sky
Date: WEDNESDAY, Temperature: 26.00°C, Description: Clear sky
...
```
