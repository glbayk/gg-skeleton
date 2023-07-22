# gg-skeleton - Generalized Boilerplate for Go REST APIs including Authentication and Authorization with JWT and Gorm

[Go fiber](https://github.com/gofiber/fiber) + [Gorm](https://github.com/go-gorm/gorm) + [GoDotEnv](https://github.com/joho/godotenv) + [RabbitMQ](https://www.rabbitmq.com/) + [MongoDB](https://www.mongodb.com/)

## Overview

This repository can be used as a starting point for building REST APIs in Go. Main focus is on Authentication and Authorization with JWT. It includes a main client, for verification and new password generation. A RabbitMQ client, that is responsible for queues related to failed confirmations or forgot emails. And a MongoDB client, that is responsible for storing logs. (This is a work in progress and the client is used with just examples in mind. It is not a complete flow.)

Main usecase for this repo is to be used as a API Gateway for microservices. It is not recommended to use this repo as a monolith.

> **Note:** The user model has been kept simple. It is recommended to add more fields to the user model or extend it via different models and relations.

## Endpoints

### Open to the world:

**Metrics**

- **GET** `/api/metrics` - get the metrics of the server

**Health Check**

- **GET** `/api/v1/ping` - check if the server is up and running

**Auth**

- **POST** `/api/v1/auth/register` - register a new user
- **POST** `/api/v1/auth/login` - login a user and get access and refresh tokens
- **GET** `/api/v1/auth/activate/:token` - activate a user's account (requires activation token)
- **POST** `/api/v1/auth/forgot-password` - send a password reset email to the user (requires email)

### Protected

**User**

- **GET** `/api/v1/user/me` - get the user's profile information
- **GET** `/api/v1/user/refresh-token` - get a new access token using the refresh token (requires refresh token in the authorization header)
- **POST** `/api/v1/user/change-password` - change the user's password (requires old and new password)

## Getting started

1. Go install the package:

```bash
go install github.com/glbayk/gg-skeleton
```

2. Copy the `.env.example` file to `.env` and replace the values with your own:
3. Build the project:

```bash
go build
```

4. Run the project:

```bash
./gg-skeleton
```

## QnA (Questions and Answers)
