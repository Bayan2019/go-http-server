# go-http-server

It is guided project in [http-server](https://www.boot.dev/courses/learn-http-servers) to learn how to build server.

A web server is just a computer that serves data over a network, typically the Internet. Servers run software that listens for incoming requests from clients. When a request is received, the server responds with the requested data.

* Use JSON, headers, and status codes to communicate with clients via a RESTful API
* Use type safe SQL to store and retrieve data from a Postgres database
* Implement a secure authentication/authorization system with well-tested cryptography libraries
* Build and understand webhooks and API keys

## Setup

Create a `.env` file with the setted values for variables
- PORT
- FILEPATH
- DB_URL
- PLATFORM (=dev)
- JWT_SECRET
- POLKA_KEY


Run the command
`go build -o server && ./server`

## Chirpy

Chirpy is a social network similar to Twitter.