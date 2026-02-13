@echo off
echo.
echo ============================================
echo   NEON RACERS - Multiplayer Flying Car Game
echo ============================================
echo.
echo Starting Go WebSocket Server on :8080
echo.

cd backend

if not exist go.sum (
    echo Installing dependencies...
    go mod download
)

echo.
echo Server ready! Open your browser to:
echo   http://localhost:8080
echo.
echo Press Ctrl+C to stop the server
echo.

go run main.go

pause
