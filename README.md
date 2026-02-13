# 🚗 NEON RACERS - Multiplayer Flying Car Game

A real-time multiplayer flying car game built with Three.js and Go WebSockets featuring a cyberpunk neon aesthetic.

## 🎮 Features

- **Real-time Multiplayer**: Race with friends using WebSocket technology
- **3D Graphics**: Powered by Three.js with stunning neon visuals
- **Physics Engine**: Realistic flying car physics with acceleration, braking, and vertical movement
- **Cyberpunk Aesthetic**: Neon lights, grid floors, and retro-futuristic design
- **Responsive Controls**: Smooth keyboard controls for an immersive experience

## 🎯 Controls

- **W / ↑** - Accelerate forward
- **S / ↓** - Brake/Reverse
- **A / ←** - Turn left
- **D / →** - Turn right
- **SPACE** - Ascend (fly up)
- **SHIFT** - Descend (fly down)
- **R** - Boost (2x speed)

## 🛠️ Technology Stack

### Frontend
- **Three.js** - 3D graphics rendering
- **WebSocket API** - Real-time multiplayer communication
- **Vanilla JavaScript** - No framework overhead for maximum performance

### Backend
- **Go** - High-performance WebSocket server
- **Gorilla WebSocket** - WebSocket library for Go
- **Concurrent architecture** - Handles multiple players simultaneously

## 📋 Prerequisites

- **Go 1.21+** installed ([Download Go](https://golang.org/dl/))
- Modern web browser with WebGL support
- Terminal/Command prompt

## 🚀 Quick Start

### 1. Install Dependencies

Navigate to the backend directory and install Go dependencies:

```bash
cd backend
go mod download
```

### 2. Run the Server

From the backend directory:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

### 3. Play the Game

Open your browser and navigate to:
```
http://localhost:8080
```

### 4. Multiplayer

Open the same URL in multiple browser tabs or share with friends on your local network!

## 🏗️ Project Structure

```
flying-car-game/
├── backend/
│   ├── main.go          # WebSocket server and game logic
│   ├── go.mod           # Go module dependencies
│   └── go.sum           # Dependency checksums (auto-generated)
├── frontend/
│   └── index.html       # Three.js game client
└── README.md            # This file
```

## 🎨 Game Features Explained

### Environment
- **Neon Grid Floor**: Cyberpunk-style glowing grid
- **Neon Pillars**: Randomly placed obstacles with colored lights
- **Particle System**: Ambient floating particles for atmosphere
- **Dynamic Fog**: Distance-based fog for depth perception

### Car Mechanics
- **Physics-based movement**: Realistic acceleration and friction
- **3D Flight**: Full vertical movement with ascend/descend controls
- **Boost system**: Temporary speed increase with R key
- **Boundary system**: Keeps players within the arena

### Multiplayer System
- **Real-time sync**: Player positions updated in real-time
- **Unique colors**: Each player gets a unique neon color
- **Player count**: Live display of connected players
- **Connection status**: Visual indicator of server connection

## 🔧 Customization

### Modify Game Physics

Edit these constants in `frontend/index.html`:

```javascript
const ACCELERATION = 0.3;      // How fast you accelerate
const MAX_SPEED = 2.5;         // Maximum speed
const BRAKE_FORCE = 0.15;      // Braking strength
const TURN_SPEED = 0.04;       // Turning rate
const VERTICAL_SPEED = 0.2;    // Up/down movement speed
const FRICTION = 0.98;         // Air resistance (0.98 = 2% friction)
const BOOST_MULTIPLIER = 2.0;  // Boost speed multiplier
```

### Change Server Port

Edit `main.go`:

```go
log.Fatal(http.ListenAndServe(":8080", nil))
// Change :8080 to your preferred port
```

### Modify Arena Size

In `frontend/index.html`, find `createEnvironment()`:

```javascript
const gridSize = 200;  // Arena size (200x200 units)
const gridDivisions = 40;  // Grid line count
```

## 🌐 Network Play

To play with friends on your local network:

1. Find your local IP address:
   - **Windows**: `ipconfig` (look for IPv4 Address)
   - **Mac/Linux**: `ifconfig` or `ip addr`

2. Share your IP with friends: `http://YOUR_IP:8080`

3. Make sure your firewall allows incoming connections on port 8080

## 🐛 Troubleshooting

### "Cannot connect to server"
- Ensure the Go server is running (`go run main.go`)
- Check that port 8080 isn't already in use
- Verify your browser supports WebSockets

### "Page not loading"
- Make sure you're in the backend directory when running the server
- The frontend files must be in `../frontend/` relative to the backend

### "Other players not visible"
- Check your browser console for WebSocket errors
- Ensure all players are connecting to the same server address
- Verify your firewall isn't blocking WebSocket connections

## 🎓 Learning Resources

- [Three.js Documentation](https://threejs.org/docs/)
- [Go WebSocket Tutorial](https://pkg.go.dev/github.com/gorilla/websocket)
- [WebSocket Protocol](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API)

## 🚀 Future Enhancements

Ideas for expanding the game:

- [ ] Power-ups and collectibles
- [ ] Lap-based racing system
- [ ] Player chat system
- [ ] Leaderboards and scoring
- [ ] Different car models and skins
- [ ] Collision detection between cars
- [ ] Sound effects and music
- [ ] Mobile touch controls
- [ ] Persistent player profiles
- [ ] Multiple arena maps

## 📝 License

This project is open source and available for educational purposes.

## 🤝 Contributing

Feel free to fork this project and add your own features! Some ideas:
- Add new car designs
- Create different arena themes
- Implement game modes (races, battles, etc.)
- Add visual effects (trails, explosions, etc.)

## 🎮 Have Fun!

Enjoy racing in the neon-lit skies! Remember: it's not about winning, it's about looking cool while flying. 🌟

---

**Made with ❤️ using Three.js and Go**
