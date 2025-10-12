# Riset Teknologi WebRTC dan Integrasi dengan Golang

## 1. Pengenalan WebRTC

WebRTC (Web Real-Time Communication) adalah teknologi open-source yang memungkinkan aplikasi web dan mobile untuk melakukan komunikasi real-time (audio, video, dan data) antar browser tanpa memerlukan plugin tambahan. WebRTC dikembangkan oleh Google dan sekarang merupakan standar W3C dan IETF.

### 1.1 Komponen Utama WebRTC

#### 1.1.1 MediaStream (getUserMedia)
- API untuk mengakses kamera dan microphone pengguna
- Menghasilkan MediaStream object yang berisi audio dan video tracks
- Contoh penggunaan:
```javascript
navigator.mediaDevices.getUserMedia({ audio: true, video: true })
  .then(stream => {
    // Gunakan stream untuk WebRTC
  })
  .catch(error => {
    console.error('Error accessing media devices:', error);
  });
```

#### 1.1.2 RTCPeerConnection
- Komponen utama untuk melakukan koneksi WebRTC
- Menangani negosiasi, encoding/decoding, dan transmisi media
- Mengelola state machine untuk koneksi WebRTC
- Contoh penggunaan:
```javascript
const peerConnection = new RTCPeerConnection(configuration);

// Menambahkan local stream
localStream.getTracks().forEach(track => {
  peerConnection.addTrack(track, localStream);
});

// Menangani remote stream
peerConnection.ontrack = (event) => {
  remoteVideo.srcObject = event.streams[0];
};
```

#### 1.1.3 RTCDataChannel
- Memungkinkan pertukaran data antar peer
- Dapat digunakan untuk chat, file sharing, atau data custom
- Mendukung reliable dan unreliable data delivery
- Contoh penggunaan:
```javascript
const dataChannel = peerConnection.createDataChannel('chat');

dataChannel.onmessage = (event) => {
  console.log('Received message:', event.data);
};

dataChannel.send('Hello WebRTC!');
```

### 1.2 Proses Koneksi WebRTC

Proses koneksi WebRTC melibatkan beberapa langkah:

1. **Signaling**: Pertukaran informasi koneksi (SDP offer/answer, ICE candidates) melalui signaling server
2. **Connection Establishment**: Pembuatan koneksi peer-to-peer menggunakan ICE framework
3. **Media Streaming**: Transmisi audio/video antar peer
4. **Data Exchange**: Pertukaran data melalui DataChannel (opsional)

## 2. WebRTC Library untuk Golang

### 2.1 Pion

Pion adalah WebRTC library untuk Golang yang paling populer dan aktif dikembangkan.

#### 2.1.1 Keunggulan Pion
- Pure Go implementation (tidak ada dependency C/C++)
- Aktif dikembangkan dan maintained oleh komunitas
- Mendukung semua fitur WebRTC standar
- Mudah digunakan dan extensible
- Mendukung WebRTC DataChannels, MediaStream, dan lainnya

#### 2.1.2 Fitur Pion
- WebRTC API implementation
- STUN/TURN server support
- SFU (Selective Forwarding Unit) implementation
- Media processing capabilities
- Interoperability dengan browser dan platform lainnya

#### 2.1.3 Contoh Penggunaan Pion
```go
package main

import (
	"fmt"
	"time"

	"github.com/pion/webrtc/v3"
)

func main() {
	// Membuat RTCPeerConnection
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer peerConnection.Close()

	// Menangani ICE candidates
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			fmt.Println("ICE Candidate:", c.ToJSON())
		}
	})

	// Menangani track dari remote peer
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Println("Track received:", track.Codec().MimeType)
	})

	// Membuat data channel
	dataChannel, err := peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		panic(err)
	}
	defer dataChannel.Close()

	// Menangani pesan dari data channel
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel: %s\n", msg.Data)
	})

	// Membuat offer
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Set local description
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Offer created:", offer.SDP)

	// Simulasi penerimaan answer (dalam aplikasi nyata, ini akan melalui signaling server)
	// answer := webrtc.SessionDescription{
	//     Type: webrtc.SDPTypeAnswer,
	//     SDP:  "...",
	// }
	// err = peerConnection.SetRemoteDescription(answer)
	// if err != nil {
	//     panic(err)
	// }

	// Tunggu untuk keperluan demo
	time.Sleep(30 * time.Second)
}
```

### 2.2 go-webrtc

go-webrtc adalah binding Golang untuk libwebrtc, library WebRTC native yang dikembangkan oleh Google.

#### 2.2.1 Keunggulan go-webrtc
- Menggunakan libwebrtc yang sama dengan browser Chrome
- Performa tinggi karena menggunakan native implementation
- Mendukung semua fitur WebRTC terbaru

#### 2.2.2 Kekurangan go-webrtc
- Memerlukan CGO dan dependency C/C++
- Lebih kompleks dalam instalasi dan deployment
- Kurang aktif dikembangkan dibandingkan Pion

#### 2.2.3 Contoh Penggunaan go-webrtc
```go
package main

import (
	"fmt"
	"time"

	"github.com/keroserene/go-webrtc"
)

func main() {
	// Membuat RTCPeerConnection
	config := webrtc.NewConfiguration()
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer peerConnection.Close()

	// Menambahkan ICE server
	iceServers := []*webrtc.ICEServer{
		&webrtc.ICEServer{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	}
	peerConnection.SetICEServers(iceServers)

	// Membuat data channel
	dataChannel, err := peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		panic(err)
	}
	defer dataChannel.Close()

	// Menangani pesan dari data channel
	dataChannel.OnMessage = func(payload []byte) {
		fmt.Printf("Message from DataChannel: %s\n", payload)
	}

	// Membuat offer
	offer, err := peerConnection.CreateOffer()
	if err != nil {
		panic(err)
	}

	// Set local description
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Offer created:", offer.SDP)

	// Tunggu untuk keperluan demo
	time.Sleep(30 * time.Second)
}
```

### 2.3 Perbandingan Pion vs go-webrtc

| Aspek | Pion | go-webrtc |
|-------|------|-----------|
| Implementasi | Pure Go | Binding C/C++ (libwebrtc) |
| Dependency | Tidak ada CGO | Memerlukan CGO |
| Performa | Baik | Sangat Baik |
| Komunitas | Aktif | Kurang aktif |
| Deployment | Mudah | Kompleks |
| Fitur | Lengkap | Lengkap |
| Stabilitas | Stabil | Stabil |

**Rekomendasi**: Untuk aplikasi WebRTC meeting ini, saya merekomendasikan menggunakan **Pion** karena:
- Pure Go implementation yang memudahkan deployment
- Komunitas yang aktif dan dokumentasi yang baik
- Fitur yang lengkap untuk kebutuhan WebRTC meeting
- Tidak ada dependency C/C++ yang mempersulit deployment

## 3. Arsitektur WebRTC dengan Golang

### 3.1 Komponen Utama

#### 3.1.1 Signaling Server
- Bertanggung jawab untuk pertukaran SDP offer/answer dan ICE candidates
- Menggunakan WebSocket untuk komunikasi real-time antara client dan server
- Mengelola room dan participant management
- Implementasi menggunakan Gorilla WebSocket atau WebSocket library lainnya

#### 3.1.2 STUN/TURN Server
- STUN (Session Traversal Utilities for NAT): Membantu client menemukan public IP address
- TURN (Traversal Using Relays around NAT): Relay server jika koneksi peer-to-peer tidak memungkinkan
- Dapat menggunakan public STUN/TURN server atau meng-host sendiri

#### 3.1.3 WebRTC Server
- Mengelola koneksi WebRTC menggunakan Pion
- Menangani media processing jika diperlukan
- Mengelola state koneksi dan error handling

### 3.2 Arsitektur Sistem

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client 1  │     │   Client 2  │     │   Client N  │
│             │     │             │     │             │
│  ┌─────────┐│     │  ┌─────────┐│     │  ┌─────────┐│
│  │Browser  ││     │  │Browser  ││     │  │Browser  ││
│  │         ││     │  │         ││     │  │         ││
│  │WebRTC   ││     │  │WebRTC   ││     │  │WebRTC   ││
│  └─────────┘│     │  └─────────┘│     │  └─────────┘│
│      │      │     │      │      │     │      │      │
│      │      │     │      │      │     │      │      │
│  ┌───▼──────┼┐    │  ┌───▼──────┼┐    │  ┌───▼──────┼┐
│  │WebSocket ││    │  │WebSocket ││    │  │WebSocket ││
│  └──────────┘│    │  └──────────┘│    │  └──────────┘│
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │
       └────────────────────┼────────────────────┘
                            │
                 ┌──────────▼──────────┐
                 │   Signaling Server  │
                 │      (Golang)       │
                 │                      │
                 │  ┌─────────────────┐ │
                 │  │   Room Manager  │ │
                 │  └─────────────────┘ │
                 │                      │
                 │  ┌─────────────────┐ │
                 │  │ WebRTC Handler  │ │
                 │  │     (Pion)      │ │
                 │  └─────────────────┘ │
                 └──────────────────────┘
                            │
       ┌────────────────────┼────────────────────┐
       │                    │                    │
┌──────▼───────┐    ┌──────▼───────┐    ┌──────▼───────┐
│  STUN/TURN   │    │   Database   │    │   Redis      │
│    Server    │    │  (PostgreSQL)│    │  (Cache)     │
└──────────────┘    └──────────────┘    └──────────────┘
```

### 3.3 Flow Aplikasi

1. **Client Connect**: Client terhubung ke signaling server via WebSocket
2. **Create/Join Room**: Client membuat atau bergabung ke room
3. **Room Management**: Server mengelola participant dalam room
4. **WebRTC Negotiation**: 
   - Client 1 membuat offer dan mengirim ke server
   - Server meneruskan offer ke Client 2
   - Client 2 membuat answer dan mengirim ke server
   - Server meneruskan answer ke Client 1
   - Kedua client menukar ICE candidates melalui server
5. **Connection Established**: Koneksi WebRTC terbentuk antara client
6. **Media Streaming**: Audio dan video mulai mengalir antar client
7. **Data Exchange**: Data channel digunakan untuk chat dan fitur lainnya

## 4. Implementasi Signaling Server dengan Golang

### 4.1 WebSocket Handler

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Untuk development, di production perlu di-setting dengan benar
	},
}

type Client struct {
	conn     *websocket.Conn
	roomID   string
	userID   string
	send     chan []byte
}

type Room struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mutex      sync.RWMutex
}

type Hub struct {
	rooms    map[string]*Room
	register chan *Client
	mutex    sync.RWMutex
}

type Message struct {
	Type      string      `json:"type"`
	RoomID    string      `json:"roomId"`
	UserID    string      `json:"userId"`
	Payload   interface{} `json:"payload"`
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// Cek apakah room sudah ada
			h.mutex.RLock()
			room, exists := h.rooms[client.roomID]
			h.mutex.RUnlock()

			if !exists {
				// Buat room baru jika belum ada
				room = &Room{
					clients:    make(map[*Client]bool),
					register:   make(chan *Client),
					unregister: make(chan *Client),
					broadcast:  make(chan []byte),
				}
				h.mutex.Lock()
				h.rooms[client.roomID] = room
				h.mutex.Unlock()
				go room.run()
			}

			// Register client ke room
			room.register <- client

		}
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.mutex.Lock()
			r.clients[client] = true
			r.mutex.Unlock()

		case client := <-r.unregister:
			r.mutex.Lock()
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
			r.mutex.Unlock()

		case message := <-r.broadcast:
			r.mutex.RLock()
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
			r.mutex.RUnlock()
		}
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		// Unregister client dari room
		hub.mutex.RLock()
		room, exists := hub.rooms[c.roomID]
		hub.mutex.RUnlock()

		if exists {
			room.unregister <- c
		}

		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Handle message berdasarkan type
		switch msg.Type {
		case "offer", "answer", "ice-candidate":
			// Forward WebRTC signaling messages ke room
			hub.mutex.RLock()
			room, exists := hub.rooms[c.roomID]
			hub.mutex.RUnlock()

			if exists {
				room.broadcast <- message
			}

		case "chat":
			// Handle chat message
			hub.mutex.RLock()
			room, exists := hub.rooms[c.roomID]
			hub.mutex.RUnlock()

			if exists {
				room.broadcast <- message
			}
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	roomID := r.URL.Query().Get("roomId")
	userID := r.URL.Query().Get("userId")

	if roomID == "" || userID == "" {
		log.Println("Missing roomId or userId")
		conn.Close()
		return
	}

	client := &Client{
		conn:   conn,
		roomID: roomID,
		userID: userID,
		send:   make(chan []byte, 256),
	}

	hub.register <- client

	// Start pumps
	go client.writePump()
	go client.readPump(hub)
}

func main() {
	hub := &Hub{
		rooms:    make(map[string]*Room),
		register: make(chan *Client),
	}

	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 4.2 WebRTC Handler dengan Pion

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/pion/webrtc/v3"
)

type WebRTCHandler struct {
	peerConnections map[string]*webrtc.PeerConnection
	mutex           sync.RWMutex
}

type WebRTCMessage struct {
	Type      string      `json:"type"`
	SDP       *string     `json:"sdp,omitempty"`
	Candidate interface{} `json:"candidate,omitempty"`
	UserID    string      `json:"userId"`
	RoomID    string      `json:"roomId"`
}

func NewWebRTCHandler() *WebRTCHandler {
	return &WebRTCHandler{
		peerConnections: make(map[string]*webrtc.PeerConnection),
	}
}

func (h *WebRTCHandler) CreatePeerConnection(userID, roomID string) (*webrtc.PeerConnection, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
	}

	// Set up handlers
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidateJSON := c.ToJSON()
		message := WebRTCMessage{
			Type:      "ice-candidate",
			Candidate: candidateJSON,
			UserID:    userID,
			RoomID:    roomID,
		}

		// Send candidate to other peers via WebSocket
		h.broadcastMessage(roomID, message)
	})

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("Connection state for user %s in room %s: %s", userID, roomID, state.String())
	})

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Track received from user %s in room %s: %s", userID, roomID, track.Codec().MimeType)
		// Handle incoming track (forward to other peers)
		h.handleTrack(track, receiver, userID, roomID)
	})

	// Store peer connection
	h.mutex.Lock()
	h.peerConnections[userID] = peerConnection
	h.mutex.Unlock()

	return peerConnection, nil
}

func (h *WebRTCHandler) HandleOffer(offer webrtc.SessionDescription, userID, roomID string) (*webrtc.SessionDescription, error) {
	h.mutex.RLock()
	peerConnection, exists := h.peerConnections[userID]
	h.mutex.RUnlock()

	if !exists {
		var err error
		peerConnection, err = h.CreatePeerConnection(userID, roomID)
		if err != nil {
			return nil, err
		}
	}

	// Set remote description
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		return nil, fmt.Errorf("failed to set remote description: %v", err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create answer: %v", err)
	}

	// Set local description
	if err := peerConnection.SetLocalDescription(answer); err != nil {
		return nil, fmt.Errorf("failed to set local description: %v", err)
	}

	return &answer, nil
}

func (h *WebRTCHandler) HandleAnswer(answer webrtc.SessionDescription, userID, roomID string) error {
	h.mutex.RLock()
	peerConnection, exists := h.peerConnections[userID]
	h.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("peer connection not found for user %s", userID)
	}

	// Set remote description
	if err := peerConnection.SetRemoteDescription(answer); err != nil {
		return fmt.Errorf("failed to set remote description: %v", err)
	}

	return nil
}

func (h *WebRTCHandler) HandleICECandidate(candidate webrtc.ICECandidateInit, userID, roomID string) error {
	h.mutex.RLock()
	peerConnection, exists := h.peerConnections[userID]
	h.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("peer connection not found for user %s", userID)
	}

	// Add ICE candidate
	if err := peerConnection.AddICECandidate(candidate); err != nil {
		return fmt.Errorf("failed to add ICE candidate: %v", err)
	}

	return nil
}

func (h *WebRTCHandler) handleTrack(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver, userID, roomID string) {
	// Forward track to other peers in the room
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for otherUserID, otherPeerConnection := range h.peerConnections {
		if otherUserID == userID {
			continue // Skip sender
		}

		// Create sender for this track
		sender, err := otherPeerConnection.AddTrack(track)
		if err != nil {
			log.Printf("Failed to add track to peer %s: %v", otherUserID, err)
			continue
		}

		// Handle RTCP feedback
		go func() {
			rtcpBuf := make([]byte, 1500)
			for {
				n, _, rtcpErr := sender.Read(rtcpBuf)
				if rtcpErr != nil {
					return
				}
				// Handle RTCP packets if needed
				_ = n
			}
		}()
	}
}

func (h *WebRTCHandler) broadcastMessage(roomID string, message WebRTCMessage) {
	// This would be implemented to send messages via WebSocket to other peers
	// Implementation depends on your WebSocket hub implementation
	messageJSON, _ := json.Marshal(message)
	log.Printf("Broadcasting message to room %s: %s", roomID, string(messageJSON))
}

func (h *WebRTCHandler) ClosePeerConnection(userID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if peerConnection, exists := h.peerConnections[userID]; exists {
		if err := peerConnection.Close(); err != nil {
			log.Printf("Error closing peer connection for user %s: %v", userID, err)
		}
		delete(h.peerConnections, userID)
	}
}

func main() {
	webrtcHandler := NewWebRTCHandler()

	http.HandleFunc("/webrtc/offer", func(w http.ResponseWriter, r *http.Request) {
		// Handle WebRTC offer
	})

	http.HandleFunc("/webrtc/answer", func(w http.ResponseWriter, r *http.Request) {
		// Handle WebRTC answer
	})

	http.HandleFunc("/webrtc/ice-candidate", func(w http.ResponseWriter, r *http.Request) {
		// Handle ICE candidate
	})

	log.Println("WebRTC server started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
```

## 5. Kesimpulan

WebRTC adalah teknologi yang powerful untuk komunikasi real-time di web. Dengan Golang dan library Pion, kita dapat membangun aplikasi WebRTC yang scalable dan performant. Beberapa poin penting:

1. **Pion adalah pilihan terbaik** untuk implementasi WebRTC di Golang karena pure Go implementation dan komunitas yang aktif.

2. **Arsitektur yang baik** memisahkan antara signaling server, WebRTC handler, dan aplikasi utama.

3. **Signaling server** menggunakan WebSocket untuk komunikasi real-time antara client dan server.

4. **WebRTC negotiation** melibatkan pertukaran SDP offer/answer dan ICE candidates melalui signaling server.

5. **Error handling dan monitoring** sangat penting untuk aplikasi WebRTC yang stabil.

Dengan pemahaman ini, kita dapat merancang arsitektur backend yang solid untuk aplikasi WebRTC meeting menggunakan Golang.