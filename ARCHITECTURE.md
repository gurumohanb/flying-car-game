# 🏗️ SYSTEM ARCHITECTURE

## High-Level Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     NEON RACERS GAME                        │
└─────────────────────────────────────────────────────────────┘

┌─────────────────┐         WebSocket          ┌──────────────┐
│   FRONTEND      │◄──────────────────────────►│   BACKEND    │
│   (Three.js)    │    Real-time Updates       │     (Go)     │
└─────────────────┘                             └──────────────┘
```

## Detailed Architecture

### Frontend (Three.js)
```
┌────────────────────────────────────────────────┐
│            Browser (index.html)                │
├────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────┐ │
│  │         Three.js Scene                   │ │
│  │  • Camera (Follow player)                │ │
│  │  • Renderer (WebGL)                      │ │
│  │  • Player Car (Local)                    │ │
│  │  • Other Cars (Remote)                   │ │
│  │  • Environment (Grid, Pillars, Fog)      │ │
│  └──────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────┐ │
│  │         Physics Engine                   │ │
│  │  • Velocity calculations                 │ │
│  │  • Collision detection                   │ │
│  │  • Input handling                        │ │
│  └──────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────┐ │
│  │       WebSocket Client                   │ │
│  │  • Send position updates                 │ │
│  │  • Receive other players                 │ │
│  │  • Handle join/leave events              │ │
│  └──────────────────────────────────────────┘ │
└────────────────────────────────────────────────┘
```

### Backend (Go)
```
┌────────────────────────────────────────────────┐
│            Go Server (main.go)                 │
├────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────┐ │
│  │      HTTP Server (:8080)                 │ │
│  │  • Serve static files                    │ │
│  │  • WebSocket upgrade endpoint            │ │
│  └──────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────┐ │
│  │      WebSocket Handler                   │ │
│  │  • Accept new connections                │ │
│  │  • Manage player sessions                │ │
│  │  • Handle disconnections                 │ │
│  └──────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────┐ │
│  │        Game State Manager                │ │
│  │  • Player registry (thread-safe)         │ │
│  │  • Position tracking                     │ │
│  │  • Broadcast system                      │ │
│  └──────────────────────────────────────────┘ │
└────────────────────────────────────────────────┘
```

## Data Flow

### Player Joins
```
1. Browser connects to ws://localhost:8080/ws
2. Server creates unique player ID
3. Server sends initial state (all current players)
4. Server broadcasts new player to others
5. Frontend spawns player car in scene
```

### Movement Update
```
1. Player presses key (W/A/S/D/SPACE/SHIFT/R)
2. Frontend calculates new position/rotation
3. Frontend sends update via WebSocket
   {
     type: "update",
     data: {
       position: {x, y, z},
       rotation: {x, y, z},
       velocity: {x, y, z}
     }
   }
4. Server receives update
5. Server updates player state
6. Server broadcasts to all other players
7. Other clients update remote car position
```

### Player Leaves
```
1. WebSocket connection closes
2. Server detects disconnection
3. Server removes player from game state
4. Server broadcasts playerLeft event
5. Other clients remove car from scene
6. Player count updates
```

## Message Types

### Client → Server
```javascript
{
  type: "update",
  data: {
    position: {x: float, y: float, z: float},
    rotation: {x: float, y: float, z: float},
    velocity: {x: float, y: float, z: float}
  }
}
```

### Server → Client

**Init (on connection)**
```javascript
{
  type: "init",
  data: {
    id: "player-id",
    players: {
      "player1-id": {...},
      "player2-id": {...}
    }
  }
}
```

**Player Joined**
```javascript
{
  type: "playerJoined",
  data: {
    id: "new-player-id",
    position: {x, y, z},
    rotation: {x, y, z},
    color: "#00ffff"
  }
}
```

**Player Update**
```javascript
{
  type: "playerUpdate",
  data: {
    id: "player-id",
    position: {x, y, z},
    rotation: {x, y, z},
    velocity: {x, y, z}
  }
}
```

**Player Left**
```javascript
{
  type: "playerLeft",
  data: {
    id: "player-id"
  }
}
```

## Concurrency Model (Go Backend)

```
Main Goroutine
    │
    ├─► HTTP Server (port 8080)
    │
    ├─► WebSocket Handler (per connection)
    │   │
    │   ├─► Read Goroutine
    │   │   └─► Listens for client messages
    │   │
    │   └─► Write Goroutine
    │       └─► Sends server messages
    │
    └─► Game State (mutex-protected)
        └─► Shared player data
```

## Performance Considerations

### Frontend
- **60 FPS target**: RequestAnimationFrame loop
- **Smooth camera**: Lerp interpolation (0.1 factor)
- **Efficient rendering**: Three.js scene graph optimization
- **Bounded updates**: Only send when position changes

### Backend
- **Concurrent connections**: Goroutines per client
- **Thread-safe state**: RWMutex for game state
- **Efficient broadcast**: Only send to relevant players
- **Memory management**: Cleanup on disconnect

## Network Protocol

- **Transport**: WebSocket (ws://)
- **Encoding**: JSON
- **Update rate**: ~60 updates/second (client-driven)
- **Latency tolerance**: Linear interpolation on client

## Security Notes (For Production)

Current implementation is for **local/educational use**. For production:

1. **Origin checking**: Currently allows all origins
2. **Rate limiting**: Add message rate limits
3. **Input validation**: Validate all position data
4. **Authentication**: Add player authentication
5. **HTTPS/WSS**: Use secure WebSocket
6. **DDoS protection**: Add connection limits

## Scalability

Current design supports:
- **Players**: ~100 simultaneous (local network)
- **Update rate**: 60 Hz per player
- **Arena size**: 200x200 units
- **Message size**: ~200 bytes per update

To scale further:
- Add spatial partitioning (only send nearby players)
- Implement server-side physics validation
- Use binary protocol (instead of JSON)
- Add load balancing for multiple servers
