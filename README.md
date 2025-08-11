# Go Config Management Service

## Description

A simple REST Config Management web service written in Go programming language.

It uses [Gin](https://gin-gonic.com/) as the HTTP framework. Currently configurations are stored in memory, which will be lost when the application is stop. 

## Getting Started

1. If you do not use devcontainer, ensure you have [Go](https://go.dev/dl/) 1.23 or higher and `Make` are installed on your machine:

    ```bash
    go version && make --version
    ```
2. Go to project root directory

3. Create a copy of the `.env.example` file and rename it to `.env`:

    ```bash
    cp .env.example .env
    ```

    Update configuration values as needed.

4. Run the project in development mode:

    ```bash
    make run
    ```
5. Open API documentation in the browser
    `http://localhost:8080/docs/index.html`

## Run as a docker container

1. Go to project root directory
2. Create a copy of the `.env.example` file and rename it to `.env`:

    ```bash
    cp .env.example .env
    ```

    Update configuration values as needed.
3. Create docker image
    ```bash
    docker build -t gocms .
    ```
4. Create and run docker container
    ```bash
    docker run --name gocms_app  -p 8080:8080 gocms:latest
    ```
5. Open API documentation in the browser
    `http://localhost:8080/docs/index.html`

## API Documentation

API documentation (swagger v2.0) can be found in `docs/` directory. To view the documentation, open the browser and go to `http://localhost:8080/docs/index.html`. The documentation is generated using [swaggo](https://github.com/swaggo/swag/) with [gin-swagger](https://github.com/swaggo/gin-swagger/) middleware.

API documentation (openapi v3.0) is `openapy.yaml`. This document is created by converting from swagger v2.0 document using [SwaggerEditor](https://editor.swagger.io/).