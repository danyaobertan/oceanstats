# OceanStats

OceanStats is a backend API for a set of underwater sensors. It provides simulated data for temperature, transparency, and detected fish species. The API is implemented in Go language and uses PostgreSQL as the database. It also includes a Swagger specification for the API.

## Table of Contents
- [Requirements](#requirements)
- [Setup](#setup)
- [Data Generation](#data-generation)
- [PostgresSQL Database](#postgressql-database)
- [API Endpoints](#api-endpoints)
- [Caching](#caching)
- [Swagger Documentation](#swagger-documentation)
- [Tools Used](#tools-used)

## Requirements

The requirements for the OceanStats project are as follows:
- Generate fake data for the sensors, including temperature, transparency, and fish species counts.
- Implement an API to access the sensor data.
- Use Docker Compose to provide a convenient way to set up the project locally.
- Store the data in PostgresSQL.
- Implement a caching for some API endpoint using Redis.
- Include a Swagger specification for the API.
- Optionally, add end-to-end tests.

## Setup

To set up the OceanStats project locally, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/oceanstats.git
    ```
   
2. Navigate to the project directory:
   ```bash
   cd oceanstats
   ```
   
3. Start the project using Docker Compose:
   ```bash
    docker-compose build
    ```
4. Run the project with air:
   ```bash
    air
    ```
The API will be accessible at http://localhost:8000.

## Data Generation
1. One-time "Kickoff" Phase:
   + 24 sensor groups for each greek letter names (alpha, beta, gamma, delta, epsilon, zeta, eta, theta, iota, kappa, lambda, mu, nu, xi, omicron, pi, rho, sigma, tau, upsilon, phi, chi, psi, omega).
   + Sensor groups cover all Ocean zone Epipelagic, Mesopelagic, Bathypelagic, Abyssopelagic and Trenches.
   + Fishes are taken from this page: https://oceana.org/ocean-fishes
   + Sensors coordinates are generated based on the zone they are in.
2. Regularly Repeated Phase:
   + Temperature and transparency are generated based on the zone, depth and time of the day.
   + Fish are generated based on the zone of their habitat.

## PostgresSQL Database
Migrations are made using gorm.
Database Schema:
![Database Schema](assets/Postgres_DB_Schema.PNG)

Data are written in batches when the buffer is full\every 10 minutes:

## API Endpoints
The OceanStats API provides the following endpoints for accessing the sensor data and gathering relevant statistics:

```
GET /group/<groupName>/transparency/average: Retrieves the current average transparency inside the group.
```

```
GET /group/<groupName>/temperature/average: Retrieves the current average temperature inside the group.
```

```
GET /group/<groupName>/species: Retrieves the full list of species (with counts) currently detected inside the group.
```

```
GET /group/<groupName>/species/top/<N>: Retrieves the top N species (with counts) currently detected inside the group.
```

```
GET /region/temperature/min: Retrieves the current minimum temperature inside the specified region.
```

```
GET /region/temperature/max: Retrieves the current maximum temperature inside the specified region.
```

```
GET /sensor/<codeName>/temperature/average: Retrieves the average temperature detected by a particular sensor between the specified date/time pairs (UNIX timestamps).
```
Please refer to the Swagger specification for more details on the request/response formats of each endpoint.

## Caching
The API caches the results of the following endpoints in Redis with a TTL of 10 seconds:

```
GET /group/<groupName>/transparency/average
```

```
GET /group/<groupName>/temperature/average
```

## Swagger Documentation
The Swagger documentation for the API is available at http://localhost:8000/api/docs/index.html.
It covers statistics endpoints.

## Tools Used
- Go (for the API implementation)
- Docker (for containerization)
- Docker Compose (for local development)
- PostgresSQL (for data storage)
- Redis (for caching)
- Swagger (for API documentation)
- Air (for live reloading)
- Gin (for routing)
- Gorm (for ORM)
- Viper (for configuration)