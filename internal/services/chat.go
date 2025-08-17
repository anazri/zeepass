package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type ChatService struct {
	rooms       map[string]*ChatRoom
	roomMutex   sync.RWMutex
	upgrader    websocket.Upgrader
	redisClient *redis.Client
	rateLimiter map[string]*RateLimiter
	limiterMutex sync.RWMutex
}

type RateLimiter struct {
	tokens    int
	capacity  int
	lastRefill time.Time
	mutex     sync.Mutex
}

type ChatRoom struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Clients   map[*Client]bool   `json:"-"`
	Messages  []EncryptedMessage `json:"messages"`
	CreatedAt time.Time          `json:"created_at"`
	mutex     sync.RWMutex
}

type Client struct {
	Room     *ChatRoom
	Conn     *websocket.Conn
	UserID   string
	UserName string
	Send     chan []byte
}

type EncryptedMessage struct {
	Type      string    `json:"type"`
	Room      string    `json:"room"`
	User      string    `json:"user"`
	Encrypted string    `json:"encrypted"`
	IV        string    `json:"iv"`
	Timestamp time.Time `json:"timestamp"`
	MessageID string    `json:"message_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Size      int       `json:"size"`
}

type MessageConfig struct {
	MaxMessageSize    int           // Maximum message size in bytes
	MessageExpiration time.Duration // Message expiration time
	RateLimit         int           // Messages per minute per user
	MaxRoomMessages   int           // Maximum messages stored per room
}

type WSMessage struct {
	Type      string `json:"type"`
	Room      string `json:"room"`
	User      string `json:"user"`
	Encrypted string `json:"encrypted"`
	IV        string `json:"iv"`
	Timestamp string `json:"timestamp"`
}

var chatService *ChatService
var messageConfig = MessageConfig{
	MaxMessageSize:    4096,        // 4KB max message size
	MessageExpiration: 24 * time.Hour, // Messages expire after 24 hours
	RateLimit:         30,          // 30 messages per minute per user
	MaxRoomMessages:   1000,        // Store max 1000 messages per room
}

func init() {
	chatService = &ChatService{
		rooms: make(map[string]*ChatRoom),
		rateLimiter: make(map[string]*RateLimiter),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for now (should be more restrictive in production)
				return true
			},
		},
	}
	
	// Initialize Redis client - will be set after InitRedis() is called
	chatService.redisClient = nil
	
	// Start cleanup routines
	go chatService.cleanupRoutine()
	go chatService.expiredMessageCleanup()
}

func GetChatService() *ChatService {
	return chatService
}

// SetRedisClient sets the Redis client for chat service
func (cs *ChatService) SetRedisClient(client *redis.Client) {
	cs.redisClient = client
}

// HandleWebSocket upgrades HTTP connection to WebSocket
func (cs *ChatService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := cs.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	
	client := &Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	
	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump(cs)
}

// CreateRoom creates a new chat room
func (cs *ChatService) CreateRoom(roomID, roomName string) *ChatRoom {
	cs.roomMutex.Lock()
	defer cs.roomMutex.Unlock()
	
	room := &ChatRoom{
		ID:        roomID,
		Name:      roomName,
		Clients:   make(map[*Client]bool),
		Messages:  make([]EncryptedMessage, 0),
		CreatedAt: time.Now(),
	}
	
	cs.rooms[roomID] = room
	log.Printf("Created room: %s (%s)", roomName, roomID)
	return room
}

// GetRoom retrieves a room by ID
func (cs *ChatService) GetRoom(roomID string) *ChatRoom {
	cs.roomMutex.RLock()
	defer cs.roomMutex.RUnlock()
	return cs.rooms[roomID]
}

// JoinRoom adds a client to a room
func (cs *ChatService) JoinRoom(client *Client, roomID, userID, userName string) error {
	cs.roomMutex.RLock()
	room := cs.rooms[roomID]
	cs.roomMutex.RUnlock()
	
	if room == nil {
		// Room doesn't exist, create it
		room = cs.CreateRoom(roomID, "Chat Room")
	}
	
	room.mutex.Lock()
	defer room.mutex.Unlock()
	
	client.Room = room
	client.UserID = userID
	client.UserName = userName
	room.Clients[client] = true
	
	log.Printf("User %s (%s) joined room %s", userName, userID, roomID)
	
	// Send recent messages to new client
	go cs.sendRecentMessages(client, room)
	
	// Notify other clients
	cs.broadcastUserJoined(room, userName)
	
	return nil
}

// LeaveRoom removes a client from a room
func (cs *ChatService) LeaveRoom(client *Client) {
	if client.Room == nil {
		return
	}
	
	room := client.Room
	room.mutex.Lock()
	defer room.mutex.Unlock()
	
	if _, ok := room.Clients[client]; ok {
		delete(room.Clients, client)
		close(client.Send)
		
		log.Printf("User %s left room %s", client.UserName, room.ID)
		
		// Notify other clients
		cs.broadcastUserLeft(room, client.UserName)
		
		// Clean up empty room
		if len(room.Clients) == 0 {
			cs.roomMutex.Lock()
			delete(cs.rooms, room.ID)
			cs.roomMutex.Unlock()
			log.Printf("Deleted empty room: %s", room.ID)
		}
	}
}

// BroadcastMessage sends encrypted message to all clients in a room
func (cs *ChatService) BroadcastMessage(room *ChatRoom, message EncryptedMessage, userID string) error {
	// Check rate limit
	if !cs.checkRateLimit(userID) {
		return fmt.Errorf("rate limit exceeded")
	}
	
	// Check message size
	if len(message.Encrypted) > messageConfig.MaxMessageSize {
		return fmt.Errorf("message too large: %d bytes (max: %d)", len(message.Encrypted), messageConfig.MaxMessageSize)
	}
	
	room.mutex.Lock()
	defer room.mutex.Unlock()
	
	// Prepare message
	message.MessageID = generateMessageID()
	message.Timestamp = time.Now()
	message.ExpiresAt = time.Now().Add(messageConfig.MessageExpiration)
	message.Size = len(message.Encrypted)
	
	// Store message in Redis (if available)
	if cs.redisClient != nil {
		if err := cs.storeMessageInRedis(message); err != nil {
			log.Printf("Failed to store message in Redis: %v", err)
			// Continue without Redis - messages will be stored in memory
		}
	}
	
	// Also keep in memory for active clients (fallback)
	room.Messages = append(room.Messages, message)
	
	// Keep only recent messages in memory
	if len(room.Messages) > 100 {
		room.Messages = room.Messages[len(room.Messages)-100:]
	}
	
	// Broadcast to all clients
	messageData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return fmt.Errorf("failed to marshal message")
	}
	
	for client := range room.Clients {
		select {
		case client.Send <- messageData:
		default:
			// Client channel is full, remove client
			delete(room.Clients, client)
			close(client.Send)
		}
	}
	
	return nil
}

func (cs *ChatService) sendRecentMessages(client *Client, room *ChatRoom) {
	var messages []EncryptedMessage
	
	// Try to get messages from Redis first (if Redis is available)
	if cs.redisClient != nil {
		redisMessages, err := cs.getMessagesFromRedis(room.ID, 50)
		if err != nil {
			log.Printf("Failed to get messages from Redis: %v", err)
		} else {
			messages = redisMessages
		}
	}
	
	// Fallback to in-memory messages if Redis failed or unavailable
	if len(messages) == 0 {
		room.mutex.RLock()
		messages = room.Messages
		room.mutex.RUnlock()
		
		// Send last 50 messages
		start := 0
		if len(messages) > 50 {
			start = len(messages) - 50
		}
		messages = messages[start:]
	}
	
	for _, message := range messages {
		// Check if message has expired
		if time.Now().After(message.ExpiresAt) {
			continue
		}
		
		messageData, err := json.Marshal(message)
		if err != nil {
			continue
		}
		
		select {
		case client.Send <- messageData:
		default:
			return
		}
	}
}

func (cs *ChatService) broadcastUserJoined(room *ChatRoom, userName string) {
	notification := EncryptedMessage{
		Type:      "user_joined",
		Room:      room.ID,
		User:      "system",
		Encrypted: "",
		IV:        "",
		Timestamp: time.Now(),
	}
	
	messageData, _ := json.Marshal(notification)
	
	for client := range room.Clients {
		select {
		case client.Send <- messageData:
		default:
		}
	}
}

func (cs *ChatService) broadcastUserLeft(room *ChatRoom, userName string) {
	notification := EncryptedMessage{
		Type:      "user_left",
		Room:      room.ID,
		User:      "system",
		Encrypted: "",
		IV:        "",
		Timestamp: time.Now(),
	}
	
	messageData, _ := json.Marshal(notification)
	
	for client := range room.Clients {
		select {
		case client.Send <- messageData:
		default:
		}
	}
}

// Client methods
func (c *Client) readPump(cs *ChatService) {
	defer func() {
		cs.LeaveRoom(c)
		c.Conn.Close()
	}()
	
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	
	for {
		_, messageData, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		
		var wsMsg WSMessage
		if err := json.Unmarshal(messageData, &wsMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}
		
		// Handle different message types
		switch wsMsg.Type {
		case "join":
			cs.JoinRoom(c, wsMsg.Room, generateUserID(), wsMsg.User)
		case "message":
			if c.Room != nil {
				timestamp, _ := time.Parse(time.RFC3339, wsMsg.Timestamp)
				encMsg := EncryptedMessage{
					Type:      "message",
					Room:      wsMsg.Room,
					User:      wsMsg.User,
					Encrypted: wsMsg.Encrypted,
					IV:        wsMsg.IV,
					Timestamp: timestamp,
				}
				if err := cs.BroadcastMessage(c.Room, encMsg, c.UserID); err != nil {
					log.Printf("Failed to broadcast message: %v", err)
					// Send error back to client
					errorMsg := map[string]string{
						"type": "error",
						"message": err.Error(),
					}
					if errorData, err := json.Marshal(errorMsg); err == nil {
						select {
						case c.Send <- errorData:
						default:
						}
					}
				}
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			
			// Add queued messages
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}
			
			if err := w.Close(); err != nil {
				return
			}
			
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Utility functions
func (cs *ChatService) cleanupRoutine() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		cs.roomMutex.Lock()
		for roomID, room := range cs.rooms {
			room.mutex.RLock()
			isEmpty := len(room.Clients) == 0
			isOld := time.Since(room.CreatedAt) > 24*time.Hour
			room.mutex.RUnlock()
			
			if isEmpty && isOld {
				delete(cs.rooms, roomID)
				log.Printf("Cleaned up old room: %s", roomID)
			}
		}
		cs.roomMutex.Unlock()
	}
}

func generateMessageID() string {
	return time.Now().Format("20060102150405") + "-" + generateRandomString(6)
}

func generateUserID() string {
	return "user-" + generateRandomString(8)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

// Redis storage functions
func (cs *ChatService) storeMessageInRedis(message EncryptedMessage) error {
	ctx := context.Background()
	
	// Store individual message with expiration
	messageKey := fmt.Sprintf("msg:%s:%s", message.Room, message.MessageID)
	messageData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	
	// Set message with expiration
	if err := cs.redisClient.Set(ctx, messageKey, messageData, messageConfig.MessageExpiration).Err(); err != nil {
		return err
	}
	
	// Add to room message list (sorted set with timestamp as score)
	roomKey := fmt.Sprintf("room:%s:messages", message.Room)
	score := float64(message.Timestamp.Unix())
	
	if err := cs.redisClient.ZAdd(ctx, roomKey, &redis.Z{
		Score:  score,
		Member: message.MessageID,
	}).Err(); err != nil {
		return err
	}
	
	// Set expiration for room message list
	cs.redisClient.Expire(ctx, roomKey, messageConfig.MessageExpiration)
	
	// Trim to keep only recent messages
	cs.redisClient.ZRemRangeByRank(ctx, roomKey, 0, int64(-messageConfig.MaxRoomMessages-1))
	
	return nil
}

func (cs *ChatService) getMessagesFromRedis(roomID string, limit int) ([]EncryptedMessage, error) {
	ctx := context.Background()
	
	// Get recent message IDs from sorted set
	roomKey := fmt.Sprintf("room:%s:messages", roomID)
	messageIDs, err := cs.redisClient.ZRevRange(ctx, roomKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	
	var messages []EncryptedMessage
	for _, messageID := range messageIDs {
		messageKey := fmt.Sprintf("msg:%s:%s", roomID, messageID)
		messageData, err := cs.redisClient.Get(ctx, messageKey).Result()
		if err != nil {
			if err == redis.Nil {
				continue // Message expired
			}
			log.Printf("Error getting message %s: %v", messageID, err)
			continue
		}
		
		var message EncryptedMessage
		if err := json.Unmarshal([]byte(messageData), &message); err != nil {
			log.Printf("Error unmarshaling message %s: %v", messageID, err)
			continue
		}
		
		messages = append(messages, message)
	}
	
	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	
	return messages, nil
}

// Rate limiting functions
func (cs *ChatService) checkRateLimit(userID string) bool {
	cs.limiterMutex.Lock()
	defer cs.limiterMutex.Unlock()
	
	limiter, exists := cs.rateLimiter[userID]
	if !exists {
		limiter = &RateLimiter{
			tokens:     messageConfig.RateLimit,
			capacity:   messageConfig.RateLimit,
			lastRefill: time.Now(),
		}
		cs.rateLimiter[userID] = limiter
	}
	
	return limiter.allow()
}

func (rl *RateLimiter) allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	
	// Refill tokens based on elapsed time (token bucket algorithm)
	tokensToAdd := int(elapsed.Minutes())
	if tokensToAdd > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
		rl.lastRefill = now
	}
	
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Cleanup functions
func (cs *ChatService) expiredMessageCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		// Only run Redis cleanup if Redis is available
		if cs.redisClient != nil {
			ctx := context.Background()
			
			// Clean up expired messages
			pattern := "msg:*:*"
			iter := cs.redisClient.Scan(ctx, 0, pattern, 0).Iterator()
			
			for iter.Next(ctx) {
				key := iter.Val()
				// Check if key exists (expired keys are automatically removed)
				exists, err := cs.redisClient.Exists(ctx, key).Result()
				if err != nil || exists == 0 {
					// Remove from room message lists
					parts := strings.Split(key, ":")
					if len(parts) >= 3 {
						roomID := parts[1]
						messageID := parts[2]
						roomKey := fmt.Sprintf("room:%s:messages", roomID)
						cs.redisClient.ZRem(ctx, roomKey, messageID)
					}
				}
			}
		}
		
		// Clean up rate limiters for inactive users
		cs.limiterMutex.Lock()
		for userID, limiter := range cs.rateLimiter {
			limiter.mutex.Lock()
			if time.Since(limiter.lastRefill) > 24*time.Hour {
				delete(cs.rateLimiter, userID)
			}
			limiter.mutex.Unlock()
		}
		cs.limiterMutex.Unlock()
		
		log.Printf("Completed expired message cleanup")
	}
}