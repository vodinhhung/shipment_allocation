# Shipment Allocation

A Go application that optimally allocates shipments to vehicles across different zones using Integer Linear Programming techniques.

## Overview

This project provides a solution for the vehicle-shipment allocation problem, which aims to minimize the total cost of shipment delivery while satisfying various constraints:

- Each zone requires a specific number of shipments to be delivered
- Vehicles have minimum and maximum shipment capacities
- Not all vehicles can serve all zones
- Each vehicle-zone pair has a specific cost per shipment

## Architecture

```
shipment_allocation/
├── bin/             # Compiled binary
├── cmd/             # Application entrypoints
│   ├── main.go      # Main application entry
│   └── router/      # HTTP routing
├── internal/        # Private application code
│   ├── api/         # Core API implementation
│   ├── common/      # Utility functions
│   ├── dependency/  # Business dependencies
│   │   └── sql/     # Database scripts
│   └── model/       # Data models
```

## Features

- Optimal shipment allocation using backtracking algorithm
- Cost minimization across vehicle fleet
- Constraint satisfaction for vehicle capacities
- Zone-specific shipment requirements

## Getting Started

### Prerequisites

- Go 1.18 or higher
- MySQL database

### Installation

```bash
go mod tidy
```

### Building

```bash
go build -o bin/main cmd/main.go
```

### Running

```bash
./bin/main
```

### Mysql create table / insert data

Running the files `internal/dependency/sql/create.sql` for creating the tables; and `internal/dependency/sql/add.sql` for inserting the data.

## API Usage

The application provides an API for shipment allocation. Input requires:

- List of zones with shipment demands
- List of vehicles with capacity constraints
- Cost matrix for vehicle-zone combinations

### Example request structure
## API Usage

### POST /allocate

Allocates shipments to vehicles optimally based on the provided constraints.

#### Request

```json
{
  "zones": [
    {"zone_id": "Z1", "shipments": 2000},
    {"zone_id": "Z2", "shipments": 1500}
  ],
  "vehicles": ["V1", "V2"]
}
```

#### Response

```json
{
  "status": "success",
  "assignments": {
    "V1": {
      "Z1": 2000,
      "Z2": 300
    },
    "V2": {
      "Z2": 1200
    }
  },
  "total_cost": 21700
}
```

This endpoint takes your zone demands, vehicle constraints, then returns the optimal allocation of shipments to minimize overall delivery costs while satisfying all constraints.

## Algorithm

The allocation uses a backtracking algorithm that:

1. Processes zones sequentially
2. For each zone, attempts to assign its shipment demand across available vehicles
3. Minimizes total cost while ensuring all constraints are satisfied
4. Prunes search paths that exceed the best cost found so far
