# Go Config Management Service

## Description

A simple REST Config Management web service written in Go programming language.

It uses [Gin](https://gin-gonic.com/) as the HTTP framework. Currently configurations are stored in memory, which will be lost when the application is stop. 

## Try with docker container
1. Create docker image
    ```bash
    docker build -t gocms .
    ```
2. Create and run docker container
    ```bash
    docker run --name gocms_app  -p 8080:8080 gocms:latest
    ```
3. Open API documentation in the browser
    `http://localhost:8080/docs/index.html`

## Getting Started

1. If you do not use devcontainer, ensure you have [Go](https://go.dev/dl/) 1.23 or higher and [Task](https://taskfile.dev/installation/) installed on your machine:

    ```bash
    go version && task --version
    ```

2. Create a copy of the `.env.example` file, put and rename it to `/cmd/http/.env`:

    ```bash
    cp .env.example /cmd/http/.env
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

## API Documentation

API documentation (swagger v2.0) can be found in `docs/` directory. To view the documentation, open the browser and go to `http://localhost:8080/docs/index.html`. The documentation is generated using [swaggo](https://github.com/swaggo/swag/) with [gin-swagger](https://github.com/swaggo/gin-swagger/) middleware.

API documentation (openapi v3.0) is `openapy.yaml`. This document is created by converting from swagger v2.0 document using [SwaggerEditor](https://editor.swagger.io/).