# STAD KV Server

This is a simple key-value store server implemented in Go.

## Overview

The STAD KV Server provides a RESTful API for storing, retrieving, and deleting key-value pairs. It uses an in-memory data structure that is persisted to disk for data durability.

## Main Components

### Engine (`engine.go`)

The `Engine` struct is the core of the key-value store:

- It maintains a map of key-value pairs in memory.
- Data is persisted to disk in JSON format.
- Provides methods for CRUD operations: Set, Get, Delete, and DeleteAll.
- Uses a read-write mutex for thread-safe operations.

### Server (`main.go`)

The main server implementation:

- Sets up HTTP routes for /set, /get, /delete, and /deleteAll endpoints.
- Implements an authentication middleware.
- Handles incoming HTTP requests and interacts with the Engine.

## API Endpoints

1. `POST /set`: Set one or more key-value pairs
2. `GET /get?key=<key>`: Retrieve a value by key
3. `DELETE /delete?key=<key>`: Delete a key-value pair
4. `DELETE /deleteAll`: Delete all key-value pairs

## Authentication

All endpoints are protected by an authentication middleware. Clients must provide a valid authentication token in the `Authorization` header.

## Usage

1. Set the `AUTH_TOKEN` environment variable.
2. Run the server: `go run .`
3. The server will start on port 8080.

## Data Persistence

Data is automatically saved to `engine_data.json` in the same directory as the server executable.
