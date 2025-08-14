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

## Code Coverage, Code Quality, and Test Coverage
### Code Coverage

1. Unit testing have been applied to core service module (94.7% coverage) and memory-based storage module (100 % coverage). These modules contain most of the application logics. Unit tests use mocking library [Mockery](https://github.com/vektra/mockery)

2. TODO: Unit test for http handler

### Code Quality

1.  **TODO**: Integrate static code analysis (linting & sonar)

2. Basic error handling have been applied for standard flow. **TODO**: Implementation to edge cases.

3. **TODO** To add proper logging using lib like Zap, and other obsertabilities (Metric & Tracing) using lib like Open Telemetry or calls  monitoring tool's API

4. **TODO** To use proper configuration (using lib like Viper)

### Test Coverage

1. Standard functionalities have been tested, including:

- Call every endpoints when the configuration store is empty

- Call every endpoints when the configuration store is filled, for both existing and not-existing config

2.  **TODO**: More exploratory testing

  

## Functional and Technical Design

### Functional Design

All basic functionalities have been applied with notes

1. Create or replace (update) configuration only accepts hardcoded config types. Each type has json schema. Currently only config type = 'person' is accepted.

2. A replace (update) may have different config type. It allows configs of a particular type migrated to new type one by one.

3. A rollback (revert) increase config version. It also have a reference to the original copied version.

4. Each version has creation timestamp

5.  **IDEA**: Add authorization process, then each version should store the creator of the version.

6.  **IDEA**: Add configuration folder/bucket/vault, a container that groups configurations. Each container may have access control (permission)

  

### Hexagonal Architecture

In the Hexagonal architecture core module does not have dependency to external module/party. Interaction to external module is through adapter which implement contract in the core module.

In this project, core module contains logic to validate configuration schema. Where the storage and http handler is the external modules.

Benefits:

1. Extendability. For example, core service not changes while implementing a new storage (e.g. db, filesystem).

2. As like other coding to interface concept, an implementation is selected using dependency inversion/injection

3. As like other coding to interface concept, unit testing is simple since the driven adapter is mockable

  

### API Design

We use resource oriented rest architectural style.

1. The URL is resource based

>/cms/configs/{name}/versions/{version}/

2. We use standard http verbs to a resource

3. For specific purpose like rollback a version, we use custom methods instead of standard PATCH. Each custom methods may have their own permission. The request payload is also simpler.

> /cms/configs/{name}/versions/{version}/rollback

4. We support skip/offset - limit based pagination. **IDEA** Add support to cursor-based pagination. It prevents performance degradation when using skip-limit with relational databases.

> // to retrieve page 3 records

> /cms/configs?skip=20&limit=10

5.  **TODO** Configs is not sorted yet. It creates un-deterministic behavior

6.  **TODO** Enhance error response, currently we don't have error code

  

### Notes

1.  **TODO** Create or replace config is not thread safe yet. Concurrent access may causes some historical version having a same version number.
2. Schema validations are stored as map entries in the core service.