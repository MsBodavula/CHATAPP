package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "net/http"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var clients = make(map[*websocket.Conn]bool)

func handleConnections(c *gin.Context) {
    ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer ws.Close()
    clients[ws] = true

    for {
        var msg map[string]string
        err := ws.ReadJSON(&msg)
        if err != nil {
            delete(clients, ws)
            break
        }

        // Broadcast the message to all connected clients
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                client.Close()
                delete(clients, client)
            }
        }
    }
}

func main() {
    r := gin.Default()
    r.GET("/ws", handleConnections)
    r.Run(":8080")
}
