# 🚀 QUICK START GUIDE - WITH AUTHENTICATION

## Get Started in 5 Steps:

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

You should see:
```
Server starting on :8080
Authentication system enabled
```

### Step 3: Open Your Browser
Navigate to: `http://localhost:8080`

You'll be automatically redirected to the login page.

### Step 4: Create Your Account

**First Time Users:**
1. Click "CREATE ACCOUNT →"
2. Choose a username (3-20 characters)
3. Create a password (minimum 6 characters)
4. Confirm your password
5. Click "CREATE PILOT"
6. You'll be redirected to the login page

### Step 5: Login and Race!

1. Enter your username and password
2. Click "ENTER ARENA"
3. Start flying!

## 🔐 Authentication Features

- **Secure Login**: Your password is hashed with bcrypt
- **Session Tokens**: 24-hour automatic login
- **Logout Button**: Top-right corner in the game
- **Auto-redirect**: Already logged in? Goes straight to game

## 🎮 Controls Reminder

- **W/↑** - Accelerate
- **S/↓** - Brake
- **A/←** - Turn Left
- **D/→** - Turn Right
- **SPACE** - Ascend
- **SHIFT** - Descend
- **R** - Boost

## 🌐 Multiplayer Setup

### Local Network (Same WiFi)

1. **Start server** on one computer
2. **Find host IP address:**
   - Windows: `ipconfig` (look for IPv4)
   - Mac/Linux: `ifconfig` or `ip addr`
3. **Friends visit:** `http://HOST_IP:8080`
4. **Each friend must:**
   - Create their own account
   - Login with their credentials
   - Race together!

### Example:
```
Host IP: 192.168.1.100
Friends visit: http://192.168.1.100:8080
Everyone registers → Everyone logs in → Race!
```

## 📁 Project Files

```
flying-car-game/
├── backend/
│   ├── main.go       # Server + Authentication
│   └── go.mod        # Dependencies
├── frontend/
│   ├── index.html    # Auto-redirect
│   ├── login.html    # Login page
│   ├── register.html # Registration page  
│   └── game.html     # The actual game
├── start.sh          # Mac/Linux launcher
├── start.bat         # Windows launcher
└── README.md         # Full documentation
```

## 🎨 What You'll See

### Login Screen
- Cyberpunk cyan/magenta theme
- Username and password fields
- Link to registration
- Animated background effects

### Registration Screen
- Create pilot ID (username)
- Password with confirmation
- Real-time validation
- Success message before redirect

### Game Arena
- Neon cyberpunk flying car
- Your username displayed (top-right)
- Logout button
- Other players with their names
- Speed and altitude HUD
- Real-time player count

## ⚡ Pro Tips

### Security
- Choose a unique username
- Use a secure password (can use same as other games, it's hashed!)
- Don't share your password
- Logout when done playing

### Gameplay
- Use boost (R) for extra speed on straightaways
- Combine turning with vertical movement for aerial maneuvers
- Stay within the arena boundaries (pink pillars mark edges)
- Higher altitude = better overview of other players

### Multiplayer
- All players need individual accounts
- Sessions last 24 hours
- If disconnected, just log back in
- Your progress/stats are tied to your username

## 🐛 Common Issues

**Q: "Page keeps redirecting to login"**
A: Your session expired or token is invalid. Just login again!

**Q: "Username already taken"**
A: Someone else (or you!) already registered that name. Try another!

**Q: "Can't see other players"**
A: Make sure they're logged in and connected to the same server.

**Q: "Server won't start"**
A: Check if port 8080 is already in use. Close other applications using it.

**Q: "Password doesn't work"**
A: Passwords are case-sensitive. Make sure Caps Lock is off!

## 🎯 Quick Reference

| Page | URL | Purpose |
|------|-----|---------|
| Home | `http://localhost:8080` | Auto-redirects based on login |
| Login | `http://localhost:8080/login.html` | Authenticate |
| Register | `http://localhost:8080/register.html` | Create account |
| Game | `http://localhost:8080/game.html` | Play (requires login) |

## 🚀 Ready to Race?

1. ✅ Start server
2. ✅ Create account
3. ✅ Login
4. ✅ Invite friends (they each register too!)
5. ✅ Race in the neon skies!

Enjoy the game! 🏁✨
