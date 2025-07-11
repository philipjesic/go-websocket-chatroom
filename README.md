# WebSocket Chat App (Go + HTML/JavaScript)

A simple real-time chat application using a Go WebSocket server and a frontend written in plain HTML and JavaScript.

---

## 📦 Features

- Connect multiple clients via WebSockets  
- Broadcast messages to all connected users  
- Usernames passed as query parameters  
- Simple HTML/JS interface  
- JSON-based message protocol  

---

## 📁 Project Structure

```
/go-chatroom/
│
├── main.go          # Go WebSocket server
├── frontend
    |-- index.html   # Frontend HTML
    |-- main.js      # Frontend JavaScript
└── README.md        # This file
```

---

## 🚀 Getting Started

### ✅ 1. Run the WebSocket Server

Make sure you have [Go installed](https://golang.org/dl/).

```bash
go run main.go
```

This starts the WebSocket server on `localhost:8080`.

---

### ✅ 2. Serve the HTML Frontend

To avoid browser CSP and `file://` restrictions, use a local HTTP server:

#### Option 1: Node.js (if you have it)

```bash
npx http-server -p 3000
```

Then open your browser and go to:

```
http://localhost:3000/index.html
```

---

## 💬 How It Works

### WebSocket URL

```
ws://localhost:8080/ws?username=YourName
```

### Message Format (JSON)

```json
{
  "username": "YourName",
  "content": "Hello world"
}
```

Each client sends messages using this format. The server broadcasts them to all connected clients.
