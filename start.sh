#!/bin/bash

echo "🚗 NEON RACERS - Starting multiplayer server..."
echo ""
echo "=============================================="
echo "  Starting Go WebSocket Server on :8080"
echo "=============================================="
echo ""

cd backend

# Download dependencies if needed
if [ ! -f "go.sum" ]; then
    echo "📦 Installing dependencies..."
    go mod download
fi

echo "🚀 Launching server..."
echo ""
echo "✅ Server ready! Open your browser to:"
echo "   http://localhost:8080"
echo ""
echo "🌐 For multiplayer on local network, use:"
echo "   http://$(hostname -I | awk '{print $1}'):8080"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

go run main.go
