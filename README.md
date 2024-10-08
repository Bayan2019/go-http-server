# go-http-server

It is guided project in [http-server](https://www.boot.dev/courses/learn-http-servers) to learn how to build server.

A web server is just a computer that serves data over a network, typically the Internet. Servers run software that listens for incoming requests from clients. When a request is received, the server responds with the requested data.

* Use JSON, headers, and status codes to communicate with clients via a RESTful API
* Use type safe SQL to store and retrieve data from a Postgres database
* Implement a secure authentication/authorization system with well-tested cryptography libraries
* Build and understand webhooks and API keys

## Setup

<!-- go get 
    
- `go get github.com/google/uuid`
- `go get github.com/lib/pq`
- `go get github.com/joho/godotenv` -->

1. Create a `.env` file with the setted values for variables
    - PORT (probably 8080)
    - FILEPATH ((in our case a dot: . which indicates the current directory) )
    - DB_URL
    - PLATFORM - set it equal to "dev"
    - JWT_SECRET
    - POLKA_KEY

2. Install Postgres (if it not installed already) \
`brew install postgresql@15` \
then start the Postgres server (if it not started already) in the background \
`brew services start postgresql@15` \
then enter pscl shell: \
`psql postgres` \
then create a new database \
`CREATE DATABASE chirpy;` \
then exit \
`exit` \
(for more details look into [Storage](https://www.boot.dev/lessons/c3e215a5-1d8f-437b-9f89-3606118800ec))

3. Install Goose \
`go install github.com/pressly/goose/v3/cmd/goose@latest`\
and then `cd` into the sql/schema directory and run: \
`goose postgres <connection_string> up`\
`<connection_string>` can be DB_URL (for more details look into [Goose Migrations](https://www.boot.dev/lessons/ea036a3f-6fa3-446a-ba20-c04cb913e12a)).

4. Install SQLC \
`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`\
then from the root of projecty run \
`sqlc generate` \
(for more details look into [SQLC](https://www.boot.dev/lessons/e5bddf3d-d96b-487e-97e6-7a5aa06b1ee1))

5. Run the command \
`go build -o server && ./server`


## Chirpy

Chirpy is a social network similar to Twitter.