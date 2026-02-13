package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// User represents a registered user
type User struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never send to client
	CreatedAt    time.Time `json:"createdAt"`
	TotalRaces   int       `json:"totalRaces"`
	BestSpeed    int       `json:"bestSpeed"`
}

// Session represents an active user session
type Session struct {
	Token     string
	Username  string
	ExpiresAt time.Time
}

type Player struct {
	ID       string                 `json:"id"`
	Username string                 `json:"username"`
	Position map[string]float64     `json:"position"`
	Rotation map[string]float64     `json:"rotation"`
	Velocity map[string]float64     `json:"velocity"`
	Color    string                 `json:"color"`
	Conn     *websocket.Conn        `json:"-"`
}

type GameState struct {
	Players map[string]*Player `json:"players"`
	mu      sync.RWMutex
}

type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// Auth stores
var (
	users    = make(map[string]*User)
	sessions = make(map[string]*Session)
	usersMu  sync.RWMutex
	sessMu   sync.RWMutex
)

var gameState = &GameState{
	Players: make(map[string]*Player),
}

// Authentication handlers

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validation
	if len(req.Username) < 3 || len(req.Username) > 20 {
		sendJSONError(w, "Username must be 3-20 characters", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		sendJSONError(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	usersMu.Lock()
	defer usersMu.Unlock()

	// Check if user already exists
	if _, exists := users[req.Username]; exists {
		sendJSONError(w, "Username already taken", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendJSONError(w, "Error creating account", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		TotalRaces:   0,
		BestSpeed:    0,
	}

	users[req.Username] = user
	log.Printf("New user registered: %s", req.Username)

	sendJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "Account created successfully",
		"user": map[string]interface{}{
			"username":   user.Username,
			"createdAt":  user.CreatedAt,
			"totalRaces": user.TotalRaces,
			"bestSpeed":  user.BestSpeed,
		},
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	usersMu.RLock()
	user, exists := users[req.Username]
	usersMu.RUnlock()

	if !exists {
		sendJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		sendJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create session
	token := generateToken()
	session := &Session{
		Token:     token,
		Username:  req.Username,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	sessMu.Lock()
	sessions[token] = session
	sessMu.Unlock()

	log.Printf("User logged in: %s", req.Username)

	sendJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"token":   token,
		"user": map[string]interface{}{
			"username":   user.Username,
			"totalRaces": user.TotalRaces,
			"bestSpeed":  user.BestSpeed,
		},
	})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		sendJSONError(w, "No token provided", http.StatusBadRequest)
		return
	}

	sessMu.Lock()
	delete(sessions, token)
	sessMu.Unlock()

	sendJSONResponse(w, map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

func handleVerifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		sendJSONError(w, "No token provided", http.StatusUnauthorized)
		return
	}

	sessMu.RLock()
	session, exists := sessions[token]
	sessMu.RUnlock()

	if !exists || time.Now().After(session.ExpiresAt) {
		sendJSONError(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	usersMu.RLock()
	user, exists := users[session.Username]
	usersMu.RUnlock()

	if !exists {
		sendJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	sendJSONResponse(w, map[string]interface{}{
		"success": true,
		"user": map[string]interface{}{
			"username":   user.Username,
			"totalRaces": user.TotalRaces,
			"bestSpeed":  user.BestSpeed,
		},
	})
}

// Game handlers

func (gs *GameState) AddPlayer(player *Player) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Players[player.ID] = player
}

func (gs *GameState) RemovePlayer(id string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	delete(gs.Players, id)
}

func (gs *GameState) UpdatePlayer(id string, position, rotation, velocity map[string]float64) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if player, exists := gs.Players[id]; exists {
		player.Position = position
		player.Rotation = rotation
		player.Velocity = velocity
	}
}

func (gs *GameState) GetPlayers() map[string]*Player {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	players := make(map[string]*Player)
	for id, player := range gs.Players {
		players[id] = &Player{
			ID:       player.ID,
			Username: player.Username,
			Position: player.Position,
			Rotation: player.Rotation,
			Velocity: player.Velocity,
			Color:    player.Color,
		}
	}
	return players
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	sessMu.RLock()
	session, exists := sessions[token]
	sessMu.RUnlock()

	if !exists || time.Now().After(session.ExpiresAt) {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	playerID := generateID()
	player := &Player{
		ID:       playerID,
		Username: session.Username,
		Position: map[string]float64{
			"x": float64((len(gameState.Players) % 5) * 50),
			"y": 10.0,
			"z": float64((len(gameState.Players) / 5) * 50),
		},
		Rotation: map[string]float64{
			"x": 0.0,
			"y": 0.0,
			"z": 0.0,
		},
		Velocity: map[string]float64{
			"x": 0.0,
			"y": 0.0,
			"z": 0.0,
		},
		Color: getRandomColor(),
		Conn:  conn,
	}

	gameState.AddPlayer(player)
	log.Printf("Player %s (%s) connected. Total players: %d", session.Username, playerID, len(gameState.Players))

	// Send initial state to new player
	initMsg := Message{
		Type: "init",
		Data: map[string]interface{}{
			"id":       playerID,
			"username": session.Username,
			"players":  gameState.GetPlayers(),
		},
	}
	conn.WriteJSON(initMsg)

	// Notify all players about new player
	broadcastNewPlayer(player)

	// Handle messages from this player
	go handlePlayerMessages(player)

	// Keep connection alive
	for {
		time.Sleep(1 * time.Second)
		if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			break
		}
	}

	gameState.RemovePlayer(playerID)
	broadcastPlayerLeft(playerID)
	log.Printf("Player %s (%s) disconnected. Total players: %d", session.Username, playerID, len(gameState.Players))
}

func handlePlayerMessages(player *Player) {
	for {
		var msg Message
		err := player.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		switch msg.Type {
		case "update":
			if pos, ok := msg.Data["position"].(map[string]interface{}); ok {
				position := make(map[string]float64)
				for k, v := range pos {
					if val, ok := v.(float64); ok {
						position[k] = val
					}
				}

				rotation := make(map[string]float64)
				if rot, ok := msg.Data["rotation"].(map[string]interface{}); ok {
					for k, v := range rot {
						if val, ok := v.(float64); ok {
							rotation[k] = val
						}
					}
				}

				velocity := make(map[string]float64)
				if vel, ok := msg.Data["velocity"].(map[string]interface{}); ok {
					for k, v := range vel {
						if val, ok := v.(float64); ok {
							velocity[k] = val
						}
					}
				}

				gameState.UpdatePlayer(player.ID, position, rotation, velocity)
				broadcastPlayerUpdate(player.ID, position, rotation, velocity)
			}
		}
	}
}

func broadcastNewPlayer(player *Player) {
	msg := Message{
		Type: "playerJoined",
		Data: map[string]interface{}{
			"id":       player.ID,
			"username": player.Username,
			"position": player.Position,
			"rotation": player.Rotation,
			"velocity": player.Velocity,
			"color":    player.Color,
		},
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()
	for id, p := range gameState.Players {
		if id != player.ID {
			p.Conn.WriteJSON(msg)
		}
	}
}

func broadcastPlayerUpdate(id string, position, rotation, velocity map[string]float64) {
	msg := Message{
		Type: "playerUpdate",
		Data: map[string]interface{}{
			"id":       id,
			"position": position,
			"rotation": rotation,
			"velocity": velocity,
		},
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()
	for playerID, player := range gameState.Players {
		if playerID != id {
			player.Conn.WriteJSON(msg)
		}
	}
}

func broadcastPlayerLeft(id string) {
	msg := Message{
		Type: "playerLeft",
		Data: map[string]interface{}{
			"id": id,
		},
	}

	gameState.mu.RLock()
	defer gameState.mu.RUnlock()
	for _, player := range gameState.Players {
		player.Conn.WriteJSON(msg)
	}
}

// Utility functions

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + string(rune(time.Now().Nanosecond()%1000))
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func getRandomColor() string {
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A",
		"#98D8C8", "#FFD93D", "#6BCF7F", "#C792EA",
		"#FF8C94", "#A8E6CF", "#FF6F91", "#5DADE2",
	}
	return colors[time.Now().Nanosecond()%len(colors)]
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func main() {
	// Auth endpoints
	http.HandleFunc("/api/register", handleRegister)
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/logout", handleLogout)
	http.HandleFunc("/api/verify", handleVerifyToken)

	// Game endpoint
	http.HandleFunc("/ws", handleWebSocket)

	// Serve static files
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	log.Println("Server starting on :8080")
	log.Println("Authentication system enabled")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
