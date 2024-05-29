# Trip Service

Welcome to the Trip Service project! This project offers a REST API for managing vehicle trips and integrates with the Mapbox API for reverse geocoding.

## Requirements

To run the Trip Service project, ensure you have the following dependencies installed:

- Go 1.16 or later
- Gin framework
- Mapbox API token

## Setup

Follow these steps to set up and run the Trip Service project:

## Clone the repository:

```bash
git clone https://github.com/skartikey/trip-service.git
cd trip-service
go mod tidy
```

## Running the Server

To start the server, run the following command:

```bash
go run cmd/trip-service/main.go
```

The server will start running on port 8080 by default.

## Running Unit Tests

To run the unit tests for the Trip Service project, execute the following command:

```bash
go test ./...
```
This command will run all the unit tests in the project and display the results.

## Endpoints

### Add Trips
- **Description**: Endpoint that takes a list of trips as JSON, each one associated with a vehicle identifier, and stores these trips in memory.
- **HTTP Method**: POST
- **Endpoint**: `/trips`
- **Payload Structure**: Each trip should follow a specific JSON structure. See [`trip.json`](testdata/trip.json) for an example of the payload structure.

### Get Trip Postcodes
- **Description**: Retrieve the postcode for the first and last coordinates of a trip identified by its unique identifier.
- **HTTP Method**: GET
- **Endpoint**: `/trips/:id/postcodes`

### Get Trip Speeds
- **Description**: Retrieve the speed for each GPS point of a trip in kilometers per hour (KPH), identified by its unique identifier.
- **HTTP Method**: GET
- **Endpoint**: `/trips/:id/speeds`

### Get Vehicle Trips
- **Description**: Retrieve a list of trips associated with a vehicle identifier, along with the start and end postcode for each trip, and the average speed across the trip.
- **HTTP Method**: GET
- **Endpoint**: `/vehicles/:id/trips`

## Design Decisions, Compromises, and Assumptions

- **Choice of Gin Framework**: I chose the Gin framework for building the REST API due to its lightweight nature, high performance, and ease of use. Gin provides robust routing, middleware support, and efficient request handling, making it suitable for building scalable APIs.

## Pending Items

Although the Trip Service project is functional, there are some pending items that could enhance its functionality and performance:

- **Authentication**: Implement authentication mechanisms such as JWT (JSON Web Tokens) or OAuth for securing the API endpoints and ensuring that only authorized users can access the resources.
- **Configuration Management**: Add support for configuration files or environment variables to configure application settings such as server port, database connection details, and API keys.
- **Benchmark Testing**: Conduct benchmark tests to evaluate the performance of the API under different load conditions and optimize performance bottlenecks if any.
- **Enhanced Logging**: Improve logging mechanisms to provide more informative and structured logs, including request/response details, error handling, and application metrics.

## Future Improvements

Despite the functional implementation of the Trip Service project, there are areas for future improvement and expansion:

- **Enhanced Error Handling**: Implement more robust error handling mechanisms to gracefully handle exceptions, validate input data, and provide informative error messages to clients.
- **Data Validation**: Implement data validation and sanitization techniques to ensure the integrity and security of the application data, preventing common security vulnerabilities such as injection attacks and data corruption.
- **Database Integration**: Integrate a database management system such as MongoDB or PostgreSQL to persist trip data and provide more advanced querying and analytics capabilities.
- **API Versioning**: Implement API versioning to manage changes and updates to the API endpoints without affecting existing clients. This ensures backward compatibility and facilitates the evolution of the API over time.