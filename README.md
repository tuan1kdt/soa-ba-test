# Go SOA backend test

## Description

A simple RESTful web service written in Go programming language. The project follows the principles of Hexagonal Architecture.

It uses [Gin](https://gin-gonic.com/) as the HTTP framework and [PostgreSQL](https://www.postgresql.org/) as the database with [pgx](https://github.com/jackc/pgx/) as the driver and [Squirrel](https://github.com/Masterminds/squirrel/) as the query builder.

## Getting Started

### Fast Testing
I currently have a running instance of the application on my VPS. You can test the application by sending requests to the following URL:

```bash
http://159.223.62.240:8080/docs/index.html
```

### Running Locally
1. If you do not use devcontainer, ensure you have [Go](https://go.dev/dl/) 1.23 or higher and [Task](https://taskfile.dev/installation/) installed on your machine:

    ```bash
    go version && task --version
    ```

2. Create a copy of the `.env.example` file and rename it to `.env`:

    ```bash
    cp .env.example .env
    ```

    Update configuration values as needed.

3. Install all dependencies, run docker compose, create database schema, and run database migrations:

    ```bash
    task
    ```

4. Run the project in development mode:

    ```bash
    task dev
    ```

## Archived
- Optimize product loading for a smooth display on a scrollable board.
- Avoid initial loading of thousands of rows by retrieving products in batches.
- Dynamic filter and incremental search.
- Create an APi to calculate the distance in km between a location (ip) and a city where a
  product produced is located.
- Get percentage of products per category
- Get percentage of products per supplier
- Create an API that generates a formatted PDF file of product data in the back end.
## TODO/Improvements
- Add caching layer to reduce the number of requests to the database.
- Add cursor-based pagination to improve performance for List products API.
- Date Format Validation (Regex)
