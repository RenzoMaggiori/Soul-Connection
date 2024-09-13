# Backend
[![language](https://img.shields.io/badge/go-v1.22+-blue?logo=go)](https://go.dev/doc/)

This is a backend service for Soul Connection, designed to migrate information from their existing API to a new system. The service handles data retrieval, processing, and storage, ensuring a smooth migration and integration with Soul Connection’s new infrastructure.

## Prerequisites

Ensure you have the following installed on your machine:

- [Docker](https://docs.docker.com/get-started/get-docker/)

- [Go `v1.22`](https://go.dev/doc/install)

## Installation

1. **Clone the Repository**

    ``` bash
    git clone https://github.com/RenzoMaggiori/Soul-Connection.git
    cd Soul-Connection/backend
    ```
2. **From the backend folder run the launch script**

    ``` bash
    ./scripts/launch.sh
    ```
3. **After the API is running start the migration**

    ``` bash
    docker compose start migration
    ```

## Environment variables

- **Environment variables have to be set in a `.env` file:**

    ``` env
    POSTGRES_USER=<user>
    POSTGRES_PASSWORD=<password>
    POSTGRES_DB=<db-name>
    POSTGRES_HOST=<host>
    POSTGRES_PORT=<db-port>
    
    MONGO_USER=<user>
    MONGO_PASSWORD=<password>
    MONGO_HOST=<host>
    MONGO_PORT=<db-port>
    
    API_KEY=<key>
    API_EMAIL=<email>
    API_PASSWORD=<password>

    WEB_URL=<url>
    ```

## Running Tests

- **Go to the `api` folder and run the command:**

    ```
    ../scripts/tests.sh
    ```



## Authors

| [<img src="https://github.com/RenzoMaggiori.png?size=85" width=85><br><sub>Renzo Maggiori</sub>](https://github.com/RenzoMaggiori) | [<img src="https://github.com/oriollinan.png?size=85" width=85><br><sub>Oriol Liñan</sub>](https://github.com/oriollinan)
|:---:|:---:|

