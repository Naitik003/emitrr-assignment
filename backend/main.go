package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	leaderboard      = make(map[string]int)
	leaderboardMutex sync.Mutex
	redisClient      *redis.Client
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Assuming Redis is running locally on default port
		DB:   0,                // Use default DB
	})
}

func signUpHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	exists, err := redisClient.Exists(ctx, user.Username).Result()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists == 1 {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Store username and password in Redis
	err = redisClient.Set(ctx, user.Username, user.Password, 0).Err()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return success message
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s created successfully", user.Username)
}

func signInHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Retrieve password from Redis
	password, err := redisClient.Get(ctx, user.Username).Result()
	if err != nil {
		http.Error(w, "Username not found", http.StatusNotFound)
		return
	}

	// Check if password matches
	if user.Password != password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	// Add additional claims or data to the token if needed

	// Sign the token with a secret key
	secretKey := []byte("kitten_secret")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Append token to response cookies
	cookie := http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s authenticated successfully", user.Username)
}

func createAccount(w http.ResponseWriter, r *http.Request) {

	fmt.Println("one")
	ctx := r.Context()
	fmt.Println(ctx)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("two")

	// Parse request body
	type profile struct {
		Username string `json:"username"`
		Points   int    `json:"points"`
	}
	fmt.Println("three")

	var receivedProfile profile
	if err := json.NewDecoder(r.Body).Decode(&receivedProfile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(receivedProfile)

	profileJSON, err := json.Marshal(receivedProfile)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store the profile JSON in Redis
	token := uuid.New().String()

	err1 := redisClient.Set(ctx, token, profileJSON, 0).Err()
	if err1 != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Created Successfully")
	}
	expiry := time.Now().Add(24 * time.Hour)

	// Append token to response cookies
	cookie := http.Cookie{
		Name:    "token",
		Value:   token,
		Path:    "/",
		Expires: expiry,
	}
	fmt.Println(token)
	fmt.Println(cookie)

	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s authenticated successfully", receivedProfile.Username)
}

func updateLeaderboard(ctx context.Context, username string, points int) {
	leaderboardMutex.Lock()
	defer leaderboardMutex.Unlock()

	leaderboard[username] += points

	// Update leaderboard in Redis
	redisClient.Set(ctx, username, leaderboard[username], 0)
}

func getValues(ctx context.Context, username string, points int) {
	leaderboardMutex.Lock()
	defer leaderboardMutex.Unlock()

	leaderboard[username] += points

	// Update leaderboard in Redis
	redisClient.Set(ctx, username, leaderboard[username], 0)
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	leaderboardMutex.Lock()
	defer leaderboardMutex.Unlock()

	json.NewEncoder(w).Encode(leaderboard)
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the API!")
}

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func wsHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel() // Make sure to call cancel to release resources associated with the context

// 	for {
// 		messageType, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		type profile struct {
// 			ID       string `json:"id"`
// 			Username string `json:"username"`
// 			Points   int    `json:"points"`
// 		}

// 		var receivedProfile profile
// 		err = json.Unmarshal(message, &receivedProfile)
// 		if err != nil {
// 			log.Println("Error unmarshaling message:", err)
// 			continue
// 		}
// 		receivedProfile.ID = uuid.New().String()

// 		fmt.Println("Received message - Username:", receivedProfile.Username, "Points:", receivedProfile.Points, "id", receivedProfile.ID)

// 		token := jwt.New(jwt.SigningMethodHS256)
// 		claims := token.Claims.(jwt.MapClaims)
// 		claims["id"] = receivedProfile.ID

// 		// Sign the token with a secret key
// 		secretKey := []byte("kitten_secret")
// 		tokenString, err := token.SignedString(secretKey)
// 		fmt.Println(tokenString)
// 		if err != nil {
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}

// 		err = conn.WriteMessage(websocket.TextMessage, []byte(tokenString))
// 		if err != nil {
// 			log.Println("Error sending token:", err)
// 			return
// 		}

// 		profileJSON, err := json.Marshal(receivedProfile)
// 		if err != nil {
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}

// 		// Store the profile JSON in Redis
// 		err = redisClient.Set(ctx, receivedProfile.ID, profileJSON, 0).Err()
// 		if err != nil {
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		} else {
// 			fmt.Println("Created Successfully")
// 		}

// 		err = conn.WriteMessage(messageType, message)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }

// func wsHandler1(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	// defer cancel()

// 	for {
// 		messageType, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		type profile struct {
// 			Points int `json:"points"`
// 		}

// 		var receivedProfile profile
// 		err = json.Unmarshal(message, &receivedProfile)
// 		if err != nil {
// 			log.Println("Error unmarshaling message:", err)
// 			continue
// 		}

// 		fmt.Println("Received message", "id", receivedProfile.Points)

// 		err = conn.WriteMessage(messageType, message)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }

// func getRandomCards(deck []string, numCards int) []string {
// 	rand.Seed(time.Now().UnixNano())
// 	var randomCards []string

// 	for i := 0; i < numCards; i++ {
// 		index := rand.Intn(len(deck))
// 		randomCards = append(randomCards, deck[index])
// 	}

// 	return randomCards
// }

// func wsFetchHandler(w http.ResponseWriter, r *http.Request) {
// 	deck := []string{"Cat", "Exploding Kitten", "Defuse", "Shuffle"}
// 	numCards := 5

// 	randomCards := getRandomCards(deck, numCards)

// 	// Convert the random selection of cards to JSON
// 	jsonData, err := json.Marshal(randomCards)
// 	if err != nil {
// 		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the response headers and write the JSON response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(jsonData)
// }

func main() {
	fmt.Println("Starting Exploding Kitten backend...")

	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/create", createAccount)
	// router.HandleFunc("/ws", wsHandler)
	// router.HandleFunc("/game", wsHandler1)
	// router.HandleFunc("/fetch", wsFetchHandler)
	router.HandleFunc("/api/leaderboard", getLeaderboard)
	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		// Use the request's context directly
		signUpHandler(r.Context(), w, r)
	})
	router.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		// Use the request's context directly
		signInHandler(r.Context(), w, r)
	})

	// Start HTTP server
	http.ListenAndServe(":8000",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.MaxAge(86400),
		)(router))
}
