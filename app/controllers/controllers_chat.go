package controllers

import (
	"log"
	"net/http"
	"sync"
	"tanya_dokter_app/app/repository"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/app/utils"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
)

// HandleWebSocket godoc
// @Summary Handle WebSocket connection
// @Description Handle WebSocket connection
// @Tags HandleWebSocket
// @Produce json
// @Param sender_id path integer true "sender_id"
// @Param receiver_id path integer true "receiver_id"
// @Param Body body reqres.GlobalChatRequest true "Handle body"
// @Success 200
// @Router /v1/chat/ws/{sender_id}/{receiver_id} [get]
// @Security ApiKeyAuth
// @Security JwtToken
func HandleWebSocket(c echo.Context) error {
	userPengirim := c.Param("sender_id")
	userPenerima := c.Param("receiver_id")

	if userPengirim == "" || userPenerima == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "sender_id and receiver_id are required",
		})
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return err
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
	}()

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var message reqres.GlobalChatRequest
		if err := c.Bind(&message); err != nil {
			return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
		}

		message.SenderID = userPengirim
		message.Message = string(messageBytes)
		message.ReceiverID = userPenerima
		message.Status = 0
		message.Timestamp = time.Now().Unix()

		_, err = repository.NewChatRepository(&message)
		if err != nil {
			log.Println("Error saving chat:", err)
			continue
		}

		broadcastMessage(message)
	}

	return nil
}

func SendMessageHandler(c echo.Context) error {
	// id := c.Get("user_id").(int)
	var message reqres.GlobalChatRequest

	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validasi sender_id dan receiver_id
	if message.SenderID == "" || message.ReceiverID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "SenderID and ReceiverID are required",
		})
	}

	message.Timestamp = time.Now().Unix()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Message sent successfully",
		"data":    message,
	})
}

func GetMessagesByUsersHandler(c echo.Context) error {

	id := c.Get("user_id")

	param := utils.PopulatePaging(c, "status")

	// userPengirim := id
	// userPenerima := c.Param("user_penerima")

	userPengirim := id.(string)
	userPenerima := "8"

	if userPengirim == "" || userPenerima == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user_pengirim and user_penerima are required",
		})
	}

	messages := repository.GetAllMessages(userPengirim, param)

	// Kirim respons dengan pesan yang ditemukan
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages": messages,
	})
}

func broadcastMessage(message reqres.GlobalChatRequest) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		if err := client.WriteJSON(message); err != nil {
			log.Println("Error writing JSON to client:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
