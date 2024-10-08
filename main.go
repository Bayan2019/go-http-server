package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Bayan2019/go-http-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Create a struct that will hold any stateful,
// in-memory data we'll need to keep track of.
type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	jwtSecret      string
	// Load POLKA_KEY into your server and store it in your apiConfig
	polkaKey string
}

func main() {
	// call godotenv.Load() at the beginning of your main() function
	// to load the .env file into your environment variables
	godotenv.Load(".env")
	// const port = "8080"
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}
	filepath := os.Getenv("FILEPATH")
	if filepath == "" {
		log.Fatal("FILEPATH is not found in the environment")
	}
	// use os.Getenv to get the DB_URL from the environment:
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	// Use a standard http.FileServer as the handler
	// Use http.Dir to convert a filepath
	// (in our case a dot: . which indicates the current directory)
	// to a directory for the http.FileServer.
	fileServer := http.FileServer(http.Dir(filepath))

	//sql.Open() a connection to your database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	// Use your SQLC generated database package
	// to create a new *database.Queries
	db := database.New(conn)

	// This is the secret used to sign and verify JWTs.
	// By keeping it safe, no other servers will be able
	// to create valid JWTs for your server.
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Add a new secret value to your .env file called POLKA_KEY.
	// This is the api key that polka will send so that
	// we know it's them (and not someone else trying to get free Chirpy red).
	// Load it into your server and store it in your apiConfig.
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		// and store db in your apiConfig struct so
		// that handlers can access it:
		DB:       db,
		Platform: platform,
		// store JWT Secret in your apiConfig struct.
		jwtSecret: jwtSecret,
		// Load POLKA_KEY into your server and store it in your apiConfig.
		polkaKey: polkaKey,
	}

	// Create a new http.ServeMux
	mux := http.NewServeMux()
	// Use the http.NewServeMux's .Handle() method
	// to add a handler for the root path (/).
	// mux.Handle("/", fileServer)
	// Update the fileserver path
	// Update the fileserver to use the /app/ path instead of /
	// to strip the /app prefix
	// from the request path before passing it to
	// the fileserver handler.
	// mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	// Wrap the http.FileServer handler with the middleware method
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	// Add the readiness endpoint
	// using the mux.HandleFunc to register your handler.
	// to only accept GET requests
	// we'll be serving the API from the /api path
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// Register that handler with the serve mux on the /metrics path.
	// to only accept GET requests
	// we'll be serving the API from the /api path
	// Swap out the GET /api/metrics endpoint,
	// which just returns plain text, for a GET /admin/metrics
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	// create and register a handler on the /reset path
	// that, when hit, will reset your fileserverHits back to 0
	// Update the /reset endpoint to only accept POST requests
	// we'll be serving the API from the /api path
	// pdate the POST /api/reset to POST /admin/reset.
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	// Add a new endpoint to the Chirpy API that accepts a POST request at /api/validate_chirp
	// Delete the /api/validate_chirp endpoint that we created before
	// but port all that logic into POST /api/chirps.
	// Users should not be allowed to create invalid chirps!
	// mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	// Add a new endpoint to your server POST /api/users
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	// Add a POST /api/chirps handler.
	// It accepts a JSON payload with a body field:
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	// Add a GET /api/chirps endpoint that returns all chirps in the database.
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	// Add a GET /api/chirps/{chirpID} endpoint
	// that returns a single chirp by its ID.
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	// Add a POST /api/login endpoint.
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	// Create a POST /api/refresh endpoint.
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	// Create a new POST /api/revoke endpoint.
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	// Add a PUT /api/users endpoint
	mux.HandleFunc("PUT /api/users", apiCfg.handlerEditUser)
	// Add a new DELETE /api/chirps/{chirpID} route to your server
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	// Add a POST /api/polka/webhooks endpoint.
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhookRedChirpy)
	// http.HandleFunc("/form", formHandler)
	// http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting Server at port %s\n", port)

	// Create a new http.Server struct
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	// Use the server's ListenAndServe method to start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
