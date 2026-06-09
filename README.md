# TravelSphere

TravelSphere is a Go web application built with the Beego framework. It provides destination browsing, country search, wishlist management, dashboard summary, and API support.

## What it includes

- Home page with popular destinations
- Country list page with search and region filter
- Country detail pages showing attractions and weather
- Session-based login and logout flows
- Wishlist CRUD API
- Dashboard summary API
- Clean project structure with `controllers`, `services`, `utils`, `filters`, and `routers`

## Project structure

- `controllers/` - server-side rendered pages and page controllers
- `controllers/api/` - JSON API controllers
- `services/` - business logic and API client integration
- `utils/` - HTTP helpers, response helpers, validation, and formatting
- `store/` - in-memory wishlist storage
- `views/` - templates for HTML pages
- `static/` - CSS and JavaScript assets
- `routers/` - application routing configuration
- `tests/` - unit test folders

## Before running

The project supports an optional `.env` file for environment variables.

### Do you need Beego installed?

No. You do not need to install Beego globally. This project uses Go modules (`go.mod`) and downloads Beego automatically when you run the app.

### Optional environment variables

- `APP_PORT` - HTTP port for the server (default `8080`)
- `DEMO_USERNAME` - demo username (password use is currently unused)
- `DEMO_PASSWORD` - demo password (currently unused)
- `REST_COUNTRIES_BASE_URL` - base URL for REST Countries API
- `OPENTRIPMAP_API_KEY` - API key for OpenTripMap
- `OPENTRIPMAP_BASE_URL` - custom base URL for OpenTripMap
- `WEATHER_API_BASE_URL` - base URL for the weather API
- `WEATHER_API_KEY` - API key for weather data

## Configuration

Before running the project, create your configuration file:

```bash
cp app.conf.example app.conf
```


## How to run

1. Clone the repository from GitHub:
   ```bash
   git clone https://github.com/mollah2022/TravelSphere.git
   ```
2. Change into the project directory:
   ```bash
   cd TravelSphere
   ```
3. Download Go module dependencies:
   ```bash
   go mod tidy
   ```
4. Create a `.env` file if needed or set environment variables.
5. Run the application with Go:
   ```bash
   go run main.go
   ```
6. Open the browser:
   ```
   http://localhost:8080
   ```

### Optional: run with Bee CLI

If you want to use the Beego `bee` command, install it first:

```bash
go install github.com/beego/bee/v2@latest
```

Then run:

```bash
bee run
```

## Docker

This repository includes a `Dockerfile` that uses `golang:1.25.5`.

## Run with Docker

Pull the image:

```bash
docker pull mollah2022/travelsphere

Build the image:

```bash
docker build -t travelsphere .
```

Run the container in the foreground:

```bash
docker run -p 8080:8080 travelsphere
```

Then open the browser:

```text
http://localhost:8080
```

Stop the container by pressing `Ctrl+C`.



### Run Docker in the background

Start it in detached mode:

```bash
docker run -d -p 8080:8080 travelsphere
```

Find the running container ID:

```bash
docker ps
```

Stop it:

```bash
docker stop <CONTAINER_ID>
```

### If Docker is not running

Start Docker on Linux:

```bash
sudo systemctl start docker
```

Stop Docker:

```bash
sudo systemctl stop docker
```

Open the browser after Docker is running:

```text
http://localhost:8080
```

## Testing

Run all tests:

```bash
go test ./...
```

Run only the API package tests:

```bash
go test ./controllers/api
```

## Main routes

- `GET /` - Home page
- `GET /countries` - Country explorer page
- `GET /countries/:slug` - Country detail page
- `GET /login`, `POST /login` - Login page and login action
- `GET /logout` - Logout
- `GET /wishlist` - Wishlist page
- `GET /dashboard` - Dashboard page

### API routes

- `GET /api/countries` - Search countries
- `GET /api/countries/:slug` - Country detail API
- `GET /api/attractions` - Attractions API
- `GET /api/suggestions` - Search suggestion API
- `GET /api/wishlist` - Get wishlist items
- `POST /api/wishlist` - Add wishlist item
- `PUT /api/wishlist/:id` - Update wishlist item
- `DELETE /api/wishlist/:id` - Delete wishlist item
- `GET /api/dashboard/summary` - Dashboard summary API

## Notes

- The app is built with Beego.
- Sessions are stored in memory.
- `services.InitServices()` initializes all service dependencies.
- If API keys are missing, the app will still run and safely fallback for attractions/weather.

---

This README is written so a trainer or reviewer can quickly understand how to run and test TravelSphere.
