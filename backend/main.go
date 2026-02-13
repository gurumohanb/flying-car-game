package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type Player struct {
	ID       string                 `json:"id"`
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

var gameState = &GameState{
	Players: make(map[string]*Player),
}

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
			Position: player.Position,
			Rotation: player.Rotation,
			Velocity: player.Velocity,
			Color:    player.Color,
		}
	}
	return players
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	playerID := generateID()
	player := &Player{
		ID: playerID,
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
	log.Printf("Player %s connected. Total players: %d", playerID, len(gameState.Players))

	// Send initial state to new player
	initMsg := Message{
		Type: "init",
		Data: map[string]interface{}{
			"id":      playerID,
			"players": gameState.GetPlayers(),
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
	log.Printf("Player %s disconnected. Total players: %d", playerID, len(gameState.Players))
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

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + string(rune(time.Now().Nanosecond()%1000))
}

func getRandomColor() string {
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A", 
		"#98D8C8", "#FFD93D", "#6BCF7F", "#C792EA",
		"#FF8C94", "#A8E6CF", "#FF6F91", "#5DADE2",
	}
	return colors[time.Now().Nanosecond()%len(colors)]
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	
	// Serve static files
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
