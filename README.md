# Virtual Orb Project

This project simulates the functionalities of an Orb, providing simplified capabilities similar to the real-life counterpart. It's developed using Golang, embracing the libraries and design patterns to ensure scalability, reliability, and maintainability.

# High Level Flow

+--------------------+       /status         +-----------------------------+
|    Virtual-ORB     |---------------------->|  Mock Uniqueness Service    |
|     (Periodic      |                       |       (External)            |
|       Job)         |       /sign-up        |                             |
|                    |---------------------->|                             |
+--------------------+                       +-----------------------------+

Virtual-ORB: Positioned at the left side. Represents the service that has a job running periodically every 5 seconds. This job sends HTTP requests to the Mock Uniqueness Service.
submitting battery, cpu usage, cpu temp, disk space as well as simulates signups, submitting random iris codes to the api. 

## How the Project is Organized

virtual-orb/
├── cmd/
│ └── virtual-orb/ # The primary application's directory
├── pkg/
│ ├── domain/ # Domain logic and types
│ ├── platform/ # Platform specific code (e.g., system info retrieval)
│ └── service/ # Core services of the application, includes business logic
├── mock/ # Mock implementations for testing and development
├── test_helper/ # Helper functions and utilities for tests

## Prerequisites

- Docker: Ensure that you have Docker installed on your machine. If not, download and install from [Docker official website](https://www.docker.com/products/docker-desktop).

## How to Run the Server

1. Navigate to the root directory of the project.
2. Run the command `make docker`. This will pull the necessary Docker images and run the server.

If everything goes well, you should be greeted logs similar to below in the terminal, which shows that we are successfully sending status and sign-ups correctly.

```
...
virtual-orb_1               | {"level":"info","ts":1692990351.1448863,"caller":"virtual-orb/main.go:88","msg":"Reporting status succeeded"}
virtual-orb_1               | {"level":"info","ts":1692990351.172803,"caller":"virtual-orb/main.go:103","msg":"Signing up succeeded"}

```

## How to Run Tests

1. Navigate to the root directory of the project.
2. Run the command `make test`.

If the tests pass successfully, you should be greeted with a test summary at the end. similar to the below:

```
....
    --- PASS: TestStatusService/should_handle_post_request_error (0.00s)
PASS
        virtual-orb/pkg/service coverage: 86.8% of statements
ok      virtual-orb/pkg/service 1.611s  coverage: 86.8% of statements

```

## Decisions on External Libraries

- **[goimagehash](https://github.com/corona10/goimagehash)**: I chose `goimagehash` to simulate how the Orb functions in real-life. Although the Orb uses a sophisticated neural network, `goimagehash` provides the advantage of utilizing the Hamming distance to account for variations in the same iris image due to factors like age, time, weather, etc.

- **[Snowflake](https://github.com/bwmarrin/snowflake)**: This was used to generate IDs that are unlikely to collide. While uniqueness wasn't a strict requirement, given the projected billion-scale data handling, it felt right to ensure ID distinctiveness.

### Others:

- **Circuit Breaker Logic**: Introduced to provide resilience in cases of network failures or other unpredictable issues. It ensures that the system gracefully handles external failures.

- **Graceful Shutdown Logic**: This ensures that our server can handle shutdown signals in a way that it finishes processing existing requests and resources are properly released, reducing the chance of data corruption and ensuring the durability of the system.

- **Mock Uniqueness Service (Mockoon CLI)**:
  An external mock service which the Virtual-ORB service communicates with.
  Has an /sign-up endpoint.
  Has a /status endpoint.
  Has a /health-check endpoint for health checks.

## Versions used.

go 1.20
docker 4.5.0 and docker-compose 3.4
