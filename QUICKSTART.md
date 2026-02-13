# 🚀 QUICK START GUIDE

## Get Started in 3 Steps:

### Step 1: Install Go
If you don't have Go installed:
- Download from: https://golang.org/dl/
- Follow the installation instructions for your OS
- Verify installation: `go version`

### Step 2: Start the Server

**On Mac/Linux:**
```bash
cd flying-car-game
./start.sh
```

**On Windows:**
```batch
cd flying-car-game
start.bat
```

**Or manually:**
```bash
cd flying-car-game/backend
go mod download
go run main.go
```

### Step 3: Open Your Browser
Navigate to: `http://localhost:8080`

## 🎮 Controls Reminder

- **W/↑** - Accelerate
- **S/↓** - Brake
- **A/←** - Turn Left
- **D/→** - Turn Right
- **SPACE** - Ascend
- **SHIFT** - Descend
- **R** - Boost

## 🌐 Multiplayer

1. Start the server on one computer
2. Find your local IP address (ipconfig/ifconfig)
3. Share `http://YOUR_IP:8080` with friends
4. Everyone connects and races together!

## 📁 Project Files

```
flying-car-game/
├── backend/
│   ├── main.go       # Go WebSocket server
│   └── go.mod        # Dependencies
├── frontend/
│   └── index.html    # Three.js game
├── start.sh          # Mac/Linux launcher
├── start.bat         # Windows launcher
└── README.md         # Full documentation
```

## 🎨 What You'll See

- Neon cyberpunk arena with glowing grid floor
- Your flying car with cyan neon lights
- Other players in different neon colors
- Speed and altitude HUD
- Real-time player count
- Smooth 60 FPS gameplay

## ⚡ Tips

- Use boost (R) for extra speed on straightaways
- Combine turning with vertical movement for aerial maneuvers
- Stay within the arena boundaries (pink pillars mark the edges)
- The higher you fly, the better overview of other players

Enjoy the race! 🏁
