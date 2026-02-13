# 🚗 NEON RACERS - Multiplayer Flying Car Game

A real-time multiplayer flying car game built with Three.js and Go WebSockets featuring a cyberpunk neon aesthetic and **secure user authentication**.

## 🎮 Features

- **User Authentication**: Secure registration and login system with bcrypt password hashing
- **Session Management**: Token-based authentication for secure gameplay
- **Real-time Multiplayer**: Race with authenticated friends using WebSocket technology
- **3D Graphics**: Powered by Three.js with stunning neon visuals
- **Physics Engine**: Realistic flying car physics with acceleration, braking, and vertical movement
- **Cyberpunk Aesthetic**: Neon lights, grid floors, and retro-futuristic design
- **Responsive Controls**: Smooth keyboard controls for an immersive experience
- **Player Profiles**: Track your username and racing stats

## 🔐 Authentication System

### Registration
- Create a unique pilot ID (username)
- Secure password (minimum 6 characters)
- Password confirmation validation
- Real-time form validation

### Login
- Authenticate with username and password
- Secure session tokens (24-hour expiration)
- Automatic redirect if already logged in
- Token verification on page load

### Security Features
- **bcrypt Password Hashing**: Passwords are hashed with bcrypt before storage
- **Token-based Sessions**: Secure session management with random tokens
- **Protected WebSocket**: Game WebSocket requires valid authentication token
- **Auto-logout**: Invalid/expired tokens automatically redirect to login
- **Client-side Validation**: Input validation before server requests

## 🎯 Controls

- **W / ↑** - Accelerate forward
- **S / ↓** - Brake/Reverse
- **A / ←** - Turn left
- **D / →** - Turn right
- **SPACE** - Ascend (fly up)
- **SHIFT** - Descend (fly down)
- **R** - Boost (2x speed)
- **Logout Button** - Exit the arena and return to login

## 🛠️ Technology Stack

### Frontend
- **Three.js** - 3D graphics rendering
- **WebSocket API** - Real-time multiplayer communication
- **LocalStorage** - Secure token storage
- **Vanilla JavaScript** - No framework overhead for maximum performance

### Backend
- **Go** - High-performance WebSocket server
- **Gorilla WebSocket** - WebSocket library for Go
- **bcrypt** - Password hashing and verification
- **Concurrent architecture** - Handles multiple authenticated players simultaneously
- **RESTful API** - Authentication endpoints

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

This will install:
- `github.com/gorilla/websocket` - WebSocket support
- `golang.org/x/crypto` - bcrypt password hashing

### 2. Run the Server

From the backend directory:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

You should see:
```
Server starting on :8080
Authentication system enabled
```

### 3. Access the Game

Open your browser and navigate to:
```
http://localhost:8080
```

You'll be automatically redirected to the login page.

### 4. Create an Account

1. Click "CREATE ACCOUNT →" on the login page
2. Choose a unique username (3-20 characters)
3. Create a secure password (minimum 6 characters)
4. Confirm your password
5. Click "CREATE PILOT"

### 5. Login and Play

1. Enter your username and password
2. Click "ENTER ARENA"
3. Start racing!

### 6. Multiplayer

Share the URL with friends! Each player needs to:
1. Register their own account
2. Login with their credentials
3. Join the same arena

## 🏗️ Project Structure

```
flying-car-game/
├── backend/
│   ├── main.go          # WebSocket server, game logic, and auth API
│   ├── go.mod           # Go module dependencies
│   └── go.sum           # Dependency checksums (auto-generated)
├── frontend/
│   ├── index.html       # Auto-redirect page
│   ├── login.html       # Login page
│   ├── register.html    # Registration page
│   └── game.html        # Three.js game client
└── README.md            # This file
```

## 🔌 API Endpoints

### Authentication Endpoints

#### POST `/api/register`
Register a new user account.

**Request:**
```json
{
  "username": "pilot123",
  "password": "securepass"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Account created successfully",
  "user": {
    "username": "pilot123",
    "createdAt": "2024-02-13T10:30:00Z",
    "totalRaces": 0,
    "bestSpeed": 0
  }
}
```

#### POST `/api/login`
Authenticate and receive a session token.

**Request:**
```json
{
  "username": "pilot123",
  "password": "securepass"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "base64-encoded-token",
  "user": {
    "username": "pilot123",
    "totalRaces": 0,
    "bestSpeed": 0
  }
}
```

#### POST `/api/logout`
Invalidate the current session.

**Headers:**
```
Authorization: your-session-token
```

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

#### GET `/api/verify`
Verify if a session token is valid.

**Headers:**
```
Authorization: your-session-token
```

**Response:**
```json
{
  "success": true,
  "user": {
    "username": "pilot123",
    "totalRaces": 0,
    "bestSpeed": 0
  }
}
```

### WebSocket Endpoint

#### WS `/ws?token={session-token}`
Connect to the game server with authenticated WebSocket.

**Authentication:** Required via query parameter `token`

## 🎨 Authentication Flow

```
┌─────────────┐
│  index.html │ ─── Check token ───┐
└─────────────┘                     │
                                    ▼
                            Token exists & valid?
                                    │
                   ┌────────────────┴────────────────┐
                   │                                  │
                  YES                                NO
                   │                                  │
                   ▼                                  ▼
            ┌─────────────┐                   ┌─────────────┐
            │  game.html  │                   │ login.html  │
            └─────────────┘                   └─────────────┘
                   │                                  │
                   │                          Click "Create Account"
                   │                                  │
                   │                                  ▼
                   │                          ┌─────────────┐
                   │                          │register.html│
                   │                          └─────────────┘
                   │                                  │
                   │                          Submit registration
                   │                                  │
                   │                                  ▼
                   │                          POST /api/register
                   │                                  │
                   │                          Account created
                   │                                  │
                   │                                  ▼
                   │◄─────────────────────────  Back to login
                   │
            WebSocket with token
                   │
                   ▼
            Multiplayer Arena!
```

## 🔧 Customization

### Session Duration

Edit `backend/main.go` to change token expiration:

```go
session := &Session{
    Token:     token,
    Username:  req.Username,
    ExpiresAt: time.Now().Add(24 * time.Hour), // Change duration here
}
```

### Password Requirements

Edit `backend/main.go` to modify validation:

```go
if len(req.Password) < 6 { // Change minimum length
    sendJSONError(w, "Password must be at least 6 characters", http.StatusBadRequest)
    return
}
```

### Game Physics

Edit `frontend/game.html`:

```javascript
const ACCELERATION = 0.3;      // How fast you accelerate
const MAX_SPEED = 2.5;         // Maximum speed
const BRAKE_FORCE = 0.15;      // Braking strength
const TURN_SPEED = 0.04;       // Turning rate
const VERTICAL_SPEED = 0.2;    // Up/down movement speed
const FRICTION = 0.98;         // Air resistance
const BOOST_MULTIPLIER = 2.0;  // Boost speed multiplier
```

## 🌐 Network Play

To play with friends on your local network:

1. **Start the server** as normal
2. **Find your local IP** address:
   - **Windows**: `ipconfig` (look for IPv4 Address)
   - **Mac/Linux**: `ifconfig` or `ip addr`
3. **Share the URL**: `http://YOUR_IP:8080`
4. **Friends register/login** on their devices
5. **Play together** in the same arena!

**Note**: All players must create their own accounts and login individually.

## 🐛 Troubleshooting

### "Cannot connect to server"
- Ensure the Go server is running (`go run main.go`)
- Check that port 8080 isn't already in use
- Verify "Authentication system enabled" appears in logs

### "Invalid or expired token"
- Your session has expired (24 hours)
- Click logout and login again
- Clear browser localStorage if issues persist

### "Username already taken"
- Choose a different username
- Usernames are case-sensitive

### "Page keeps redirecting to login"
- Token may be invalid/expired
- Clear browser cache and localStorage
- Try registering a new account

### "Other players not visible"
- Ensure all players are logged in
- Check WebSocket connection status
- Verify firewall isn't blocking connections

## 🔒 Security Notes

### Current Implementation
- **Development mode**: Allows all origins for WebSocket
- **In-memory storage**: User data resets on server restart
- **HTTP protocol**: Uses ws:// instead of wss://

### Production Recommendations
1. **Use HTTPS/WSS**: Secure WebSocket connections
2. **Database storage**: Persist user accounts
3. **Rate limiting**: Prevent brute force attacks
4. **Origin checking**: Restrict WebSocket origins
5. **CSRF protection**: Add tokens to API requests
6. **Input sanitization**: Validate all user inputs
7. **Password strength**: Enforce stronger requirements
8. **Account verification**: Add email verification
9. **Password reset**: Implement forgot password flow
10. **Session management**: Redis for distributed sessions

## 🚀 Future Enhancements

Authentication & Accounts:
- [ ] Email verification
- [ ] Password reset functionality
- [ ] OAuth integration (Google, GitHub)
- [ ] Two-factor authentication
- [ ] Persistent user profiles (database)
- [ ] Player statistics and leaderboards
- [ ] Achievement system
- [ ] Friend system

Gameplay:
- [ ] Power-ups and collectibles
- [ ] Lap-based racing system
- [ ] Player chat system
- [ ] Different car models and skins
- [ ] Collision detection between cars
- [ ] Sound effects and music
- [ ] Mobile touch controls
- [ ] Multiple arena maps
- [ ] Race modes (time trial, battle, etc.)

## 📝 License

This project is open source and available for educational purposes.

## 🎮 Have Fun!

Enjoy racing in the neon-lit skies with your friends! Remember to create your pilot account and race responsibly! 🌟

---

**Made with ❤️ using Three.js, Go, and bcrypt**
