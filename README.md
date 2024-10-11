# log-producer-aggregator

## Overview
- This project exposes a server with two endpoints: POST `logs/batch` and GET `logs/retrieve`.
- Upon recieving a batch of logs, the server pushes the log batch to a pool of workers which one of them will pick them and process into the database. A response is recieved directly.
- Upon recieving a request for logs, the server pushes the request to a pool of workers. One will make a database request to fetch them based upon the query params that are passed, startTime, endTime and logLevel.

## Structure - aggregator
### api
- **`handlers.go`***
- Holds the logic for each endpoint.

-**`server.go`**
- Server, database and workerpool setup.

### cmd
- **`main.go`**
- Main entry point to the application, handles starting up the server and closing it based on signal.

### internal
- **`circuitbreaker.go`**
- Circuit breaker logic

- **`workerpool.go`**
- Workerpool logic, for workers and the pool.
- Logic for processing the job types that are passed through

### storage
- **`database.go`**
- Logic for setting up the database connection.
- Closing the database connection also.

- **`models.go`**
- Holds the database structure for mongodb

- **`operations.go`**
- A selection of functions for different database operations e.g. log fetching and storing to the database

### utils
- **`log.go`**
- Utils for HTML logic, e.g. decoding body, passing a response back

- **`shared.go`**
- Shared strucs to use throughout the application

## producer
### cmd
- **`main.go`**
- Entry point into the application, runs the log producing script

### logs
- **`sender.go`**
- Produces random logs and sends a request to the server to recieve

## Setup

1. Build the docker containers
```bash
docker-compose build
```

2. Start the container
```bash
docker-compose up -d
```

example retrival endpoint:
```bash
curl -X GET "http://localhost:8005/logs/retrieve?startTime=2024-10-07T20:00:00Z&endTime=2024-10-08T08:00:00Z&logLevel=WARNING"
```

example batch endpoint:
```bash
curl -X POST "http://localhost:8005/logs/batch" \
-H "Content-Type: application/json" \
-d '[
    {
        "timestamp": "2024-10-08T00:00:00Z",
        "level": "ERROR",
        "message": "Database connection failed."
    },
    {
        "timestamp": "2024-10-08T01:00:00Z",
        "level": "INFO",
        "message": "User login successful."
    },
    {
        "timestamp": "2024-10-08T02:00:00Z",
        "level": "WARNING",
        "message": "High memory usage detected."
    }
]'
```