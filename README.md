# go-http-server

It is guided project in [http-server](https://www.boot.dev/courses/learn-http-servers) to learn how to build server.

A web server is just a computer that serves data over a network, typically the Internet. Servers run software that listens for incoming requests from clients. When a request is received, the server responds with the requested data.

* Use JSON, headers, and status codes to communicate with clients via a RESTful API
* Use type safe SQL to store and retrieve data from a Postgres database
* Implement a secure authentication/authorization system with well-tested cryptography libraries
* Build and understand webhooks and API keys

## Setup

Create a `.env` file with the setted values for variables
- PORT (probably 8080)
- FILEPATH ((in our case a dot: . which indicates the current directory) )
- DB_URL
- PLATFORM (=dev)
- JWT_SECRET
- POLKA_KEY

1. `cd` into the sql/schema directory and run: \
`goose postgres <connection_string> up`\
`<connection_string>` can be DB_URL


2. Run the command \
`go build -o server && ./server`

## Chirpy

Chirpy is a social network similar to Twitter.