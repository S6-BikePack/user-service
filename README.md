<div align="center">
  <a href="https://github.com/S6-BikePack">
    <img src="assets/logo.png" alt="logo" width="200" height="auto" />
  </a>
  <h1>BikePack - User-Service</h1>

  <p>
    Part of the S6 BikePack project.
  </p>


<!-- Badges -->
[![golangci-lint](https://github.com/S6-BikePack/user-service/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/S6-BikePack/user-service/actions/workflows/golangci-lint.yml)
[![Makefile CI](https://github.com/S6-BikePack/user-service/actions/workflows/run-tests.yml/badge.svg)](https://github.com/S6-BikePack/user-service/actions/workflows/run-tests.yml)
[![Build](https://github.com/S6-BikePack/user-service/actions/workflows/build.yml/badge.svg)](https://github.com/S6-BikePack/user-service/actions/workflows/build.yml)

<h4>
    <a href="https://github.com/S6-BikePack">Home</a>
  <span> ยท </span>
    <a href="https://github.com/S6-BikePack/user-service#-about-the-project">Documentation</a>
  </h4>
</div>

<br />

<!-- Table of Contents -->
# ๐ Table of Contents

- [About the Project](#-about-the-project)
    * [Architecture](#-architecture)
    * [Tech Stack](#%EF%B8%8F-tech-stack)
    * [Environment Variables](#-environment-variables)
    * [Messages](#-messages)
    * [Data](#-data)
- [Getting Started](%EF%B8%8F-getting-started)
    * [Prerequisites](%EF%B8%8F-prerequisites)
    * [Running Tests](#-running-tests)
    * [Run Locally](#-run-locally)
    * [Deployment](#-deployment)
- [Usage](#-usage)



<!-- About the Project -->
## โญ About the Project

The User-Service is the service for the BikePack project that handles all users in the system. 
A user is anyone with an account in for BikePack and can be both a customer and a rider.
Using the system users can register to the system with the first and last name and firebase authentication data.

<!-- Architecture -->
### ๐  Architecture
For this service I have chosen a Hexagonal architecture. This keeps the service loosely coupled and thus flexible when having to change parts of the system.

<!-- TechStack -->
### ๐ฐ๏ธ Tech Stack
#### Language
  <ul>
    <li><a href="https://go.dev/">GoLang</a></li>
</ul>

#### Dependencies
  <ul>
    <li><a href="https://github.com/gin-gonic/gin">Gin</a><span> - Web framework</span></li>
    <li><a href="https://github.com/gin-gonic/gin">Amqp091-go</a><span> - Go AMQP 0.9.1 client</span></li>
    <li><a href="https://github.com/swaggo/swag">Swag</a><span> - Swagger documentation</span></li>
    <li><a href="https://gorm.io/index.html">GORM</a><span> - ORM library</span></li>
  </ul>

#### Database
  <ul>
    <li><a href="https://www.postgresql.org/">PostgreSQL</a></li>
</ul>

<!-- Env Variables -->
### ๐ Environment Variables

This service can be configured using environment variables or a json file. The location of the file can be passed to the `config` environment variable.

The following configuration is available:

```json
{
    "server": {
      "service": "string",
      "port": "string",
      "description": "string"
    },
    "rabbitMQ": {
      "host": "string",
      "port": "int",
      "user": "string",
      "password": "string",
      "exchange": "string"
    },
    "database": {
      "host": "string",
      "port": "int",
      "user": "string",
      "password": "string",
      "database": "string",
      "debug": "bool"
    }
}
```

<!-- Messages -->
## ๐จ Messages

### Publishing
The service publishes the following messages to the RabbitMQ server:

---
**user.create**

Published when a new user is created in the system.
Sends the newly created user in the  body.

```json
{
  "id": "string",
  "name": "string",
  "last_name": "string",
  "email": "string"
}
```



---
**user.update**

Published when a delivery is updated in the system.
Sends the updated delivery in the  body.

```json
{
  "id": "string", 
  "name": "string",
  "last_name": "string",
  "email": "string"
}
```

<!-- Data -->

##  ๐๏ธ Data

This service stores the following data:

```json
{
  "id": "string", //primary key
  "name": "string",
  "last_name": "string",
  "email": "string"
}
```

<!-- Getting Started -->
## 	๐ ๏ธ Getting Started

<!-- Prerequisites -->
### โผ๏ธ Prerequisites

Building the project requires Go 1.18.

This project requires a PostgreSQL compatible database with a database named `user` and a RabbitMQ server.
The easiest way to setup the project is to use the Docker-Compose file from the infrastructure repository.

<!-- Running Tests -->
### ๐งช Running Tests

The tests in the project can easily be run using make and the `make run-tests` command. This will start the required docker containers and run all tests in the project.

<!-- Run Locally -->
### ๐ Run Locally

Clone the project

```bash
  git clone https://github.com/S6-BikePack/user-service
```

Go to the project directory

```bash
  cd user-service
```

Run the project (Rest)

```bash
  go run cmd/rest/main.go
```


<!-- Deployment -->
### ๐ Deployment

To build this project run (Rest)

```bash
  go build cmd/rest/main.go
```


<!-- Usage -->
## ๐ Usage

### REST
Once the service is running you can find its swagger documentation with all the endpoints at `/swagger`