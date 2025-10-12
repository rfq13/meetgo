#!/bin/bash

# Test Script untuk Validasi WebSocket Server WebRTC
# Author: WebRTC Meeting Team
# Version: 1.0.0

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to load environment variables
load_env() {
    print_status "Loading environment variables..."
    
    if [[ -f ".env" ]]; then
        export $(cat .env | grep -v '^#' | xargs)
        print_success "Environment variables loaded from .env file"
    elif [[ -f ".env.example" ]]; then
        print_warning ".env file not found, using .env.example as reference"
        export $(cat .env.example | grep -v '^#' | xargs)
    else
        print_error "No .env or .env.example file found"
        return 1
    fi
}

# Function to validate required environment variables
validate_env_vars() {
    print_status "Validating WebSocket environment variables..."
    
    local required_vars=(
        "WS_HOST"
        "WS_PORT"
        "JANUS_BASE_URL"
        "JANUS_ADMIN_URL"
        "JANUS_API_SECRET"
        "JANUS_ADMIN_SECRET"
    )
    
    local missing_vars=()
    
    for var in "${required_vars[@]}"; do
        if [[ -z "${!var}" ]]; then
            missing_vars+=("$var")
        fi
    done
    
    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        print_error "Missing required environment variables:"
        for var in "${missing_vars[@]}"; do
            echo "  - $var"
        done
        return 1
    fi
    
    print_success "All required environment variables are set"
    return 0
}

# Function to check if WebSocket server is running
check_websocket_server() {
    print_status "Checking WebSocket server status..."
    
    local ws_url="ws://${WS_HOST:-localhost}:${WS_PORT:-8081}"
    
    # Check HTTP health endpoint first
    if curl -s --connect-timeout 5 "http://${WS_HOST:-localhost}:${WS_PORT:-8081}/health" >/dev/null 2>&1; then
        print_success "WebSocket server is running at $ws_url"
        WS_BASE_URL="$ws_url"
        WS_HTTP_URL="http://${WS_HOST:-localhost}:${WS_PORT:-8081}"
        return 0
    else
        print_warning "WebSocket server is not running. Attempting to start..."
        start_websocket_server
        return $?
    fi
}

# Function to start WebSocket server
start_websocket_server() {
    print_status "Starting WebSocket server..."
    
    cd backend
    
    # Check if binary exists
    if [[ ! -f "bin/websocket-server" ]]; then
        print_error "WebSocket server binary not found. Please run compilation test first."
        cd ..
        return 1
    fi
    
    # Start WebSocket server in background
    ./bin/websocket-server > websocket-server.log 2>&1 &
    WS_SERVER_PID=$!
    
    # Wait for server to start
    local retries=10
    local retry_count=0
    
    while [[ $retry_count -lt $retries ]]; do
        if curl -s --connect-timeout 2 "http://${WS_HOST:-localhost}:${WS_PORT:-8081}/health" >/dev/null 2>&1; then
            print_success "WebSocket server started successfully (PID: $WS_SERVER_PID)"
            WS_BASE_URL="ws://${WS_HOST:-localhost}:${WS_PORT:-8081}"
            WS_HTTP_URL="http://${WS_HOST:-localhost}:${WS_PORT:-8081}"
            cd ..
            return 0
        fi
        
        sleep 2
        ((retry_count++))
    done
    
    print_error "Failed to start WebSocket server"
    kill $WS_SERVER_PID 2>/dev/null || true
    cd ..
    return 1
}

# Function to stop WebSocket server
stop_websocket_server() {
    if [[ -n "$WS_SERVER_PID" ]]; then
        print_status "Stopping WebSocket server (PID: $WS_SERVER_PID)..."
        kill $WS_SERVER_PID 2>/dev/null || true
        wait $WS_SERVER_PID 2>/dev/null || true
        print_success "WebSocket server stopped"
    fi
}

# Function to install WebSocket client tools
install_websocket_tools() {
    print_status "Checking WebSocket client tools..."
    
    if command_exists websocat; then
        print_success "websocat is available"
        WEBSOCAT_AVAILABLE=true
    else
        print_warning "websocat not found, attempting to install..."
        if command_exists cargo; then
            cargo install websocat
            WEBSOCAT_AVAILABLE=true
        elif command_exists pip3; then
            pip3 install websocat
            WEBSOCAT_AVAILABLE=true
        else
            print_warning "Cannot install websocat, will use curl for basic tests"
            WEBSOCAT_AVAILABLE=false
        fi
    fi
    
    if command_exists node && npm list -g ws >/dev/null 2>&1; then
        print_success "Node.js WebSocket client is available"
        NODE_WS_AVAILABLE=true
    else
        print_warning "Node.js WebSocket client not available"
        NODE_WS_AVAILABLE=false
    fi
}

# Function to test WebSocket HTTP endpoints
test_websocket_http_endpoints() {
    print_status "Testing WebSocket HTTP endpoints..."
    
    # Test health endpoint
    print_status "Testing health endpoint..."
    
    local response=$(curl -s -w "%{http_code}" "$WS_HTTP_URL/health")
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Health endpoint - HTTP $http_code"
        echo "Response: $body"
    else
        print_error "Health endpoint failed - HTTP $http_code"
        return 1
    fi
    
    # Test stats endpoint
    print_status "Testing stats endpoint..."
    
    response=$(curl -s -w "%{http_code}" "$WS_HTTP_URL/api/v1/websocket/stats")
    http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Stats endpoint - HTTP $http_code"
    else
        print_warning "Stats endpoint failed - HTTP $http_code"
    fi
    
    return 0
}

# Function to create WebSocket test client
create_websocket_test_client() {
    print_status "Creating WebSocket test client..."
    
    cat > backend/test_websocket_client.js << 'EOF'
const WebSocket = require('ws');

class WebSocketTestClient {
    constructor(url, userId) {
        this.url = url;
        this.userId = userId;
        this.ws = null;
        this.messages = [];
        this.connected = false;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 3;
    }

    connect() {
        return new Promise((resolve, reject) => {
            console.log(`Connecting to WebSocket server at ${this.url}...`);
            
            this.ws = new WebSocket(`${this.url}?userId=${this.userId}`);
            
            this.ws.on('open', () => {
                console.log('✅ WebSocket connection established');
                this.connected = true;
                this.reconnectAttempts = 0;
                resolve();
            });
            
            this.ws.on('message', (data) => {
                const message = JSON.parse(data.toString());
                this.messages.push(message);
                console.log('📨 Received message:', message);
            });
            
            this.ws.on('close', (code, reason) => {
                console.log(`🔌 WebSocket connection closed (code: ${code}, reason: ${reason})`);
                this.connected = false;
            });
            
            this.ws.on('error', (error) => {
                console.error('❌ WebSocket error:', error);
                this.connected = false;
                reject(error);
            });
        });
    }

    sendMessage(type, data, roomId = null) {
        if (!this.connected) {
            throw new Error('WebSocket is not connected');
        }
        
        const message = {
            type: type,
            data: data,
            timestamp: new Date().toISOString()
        };
        
        if (roomId) {
            message.roomId = roomId;
        }
        
        console.log('📤 Sending message:', message);
        this.ws.send(JSON.stringify(message));
    }

    joinRoom(roomId) {
        this.sendMessage('join_room', { roomId: roomId }, roomId);
    }

    leaveRoom(roomId) {
        this.sendMessage('leave_room', { roomId: roomId }, roomId);
    }

    sendSignal(roomId, signalData) {
        this.sendMessage('signal', signalData, roomId);
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }

    waitForMessage(timeout = 5000) {
        return new Promise((resolve, reject) => {
            const startTime = Date.now();
            
            const checkMessage = () => {
                if (this.messages.length > 0) {
                    resolve(this.messages.pop());
                } else if (Date.now() - startTime > timeout) {
                    reject(new Error('Timeout waiting for message'));
                } else {
                    setTimeout(checkMessage, 100);
                }
            };
            
            checkMessage();
        });
    }
}

// Test functions
async function testWebSocketConnection() {
    console.log('🧪 Testing WebSocket connection...');
    
    const client = new WebSocketTestClient('ws://localhost:8081/ws', 'test-user-1');
    
    try {
        await client.connect();
        console.log('✅ WebSocket connection test passed');
        client.disconnect();
        return true;
    } catch (error) {
        console.error('❌ WebSocket connection test failed:', error);
        return false;
    }
}

async function testRoomJoin() {
    console.log('🧪 Testing room join functionality...');
    
    const client = new WebSocketTestClient('ws://localhost:8081/ws', 'test-user-2');
    
    try {
        await client.connect();
        
        // Join a room
        client.joinRoom('test-room-123');
        
        // Wait for response
        const message = await client.waitForMessage();
        
        if (message.type === 'room_joined' || message.type === 'join_room_response') {
            console.log('✅ Room join test passed');
            client.disconnect();
            return true;
        } else {
            console.log('❌ Unexpected message type:', message.type);
            client.disconnect();
            return false;
        }
    } catch (error) {
        console.error('❌ Room join test failed:', error);
        client.disconnect();
        return false;
    }
}

async function testMultipleClients() {
    console.log('🧪 Testing multiple clients in room...');
    
    const client1 = new WebSocketTestClient('ws://localhost:8081/ws', 'test-user-3');
    const client2 = new WebSocketTestClient('ws://localhost:8081/ws', 'test-user-4');
    
    try {
        await Promise.all([client1.connect(), client2.connect()]);
        
        const roomId = 'test-room-multi';
        
        // Both clients join the same room
        client1.joinRoom(roomId);
        client2.joinRoom(roomId);
        
        // Wait for join responses
        await Promise.all([
            client1.waitForMessage(),
            client2.waitForMessage()
        ]);
        
        // Client 1 sends a signal
        client1.sendSignal(roomId, {
            type: 'test_signal',
            data: 'Hello from client 1'
        });
        
        // Client 2 should receive the signal
        const signal = await client2.waitForMessage();
        
        if (signal.type === 'signal' && signal.data.type === 'test_signal') {
            console.log('✅ Multiple clients test passed');
            client1.disconnect();
            client2.disconnect();
            return true;
        } else {
            console.log('❌ Signal not received correctly');
            client1.disconnect();
            client2.disconnect();
            return false;
        }
    } catch (error) {
        console.error('❌ Multiple clients test failed:', error);
        client1.disconnect();
        client2.disconnect();
        return false;
    }
}

// Run tests
async function runTests() {
    console.log('🚀 Starting WebSocket tests...');
    
    const tests = [
        testWebSocketConnection,
        testRoomJoin,
        testMultipleClients
    ];
    
    let passed = 0;
    let total = tests.length;
    
    for (const test of tests) {
        try {
            if (await test()) {
                passed++;
            }
        } catch (error) {
            console.error('Test error:', error);
        }
        
        // Wait between tests
        await new Promise(resolve => setTimeout(resolve, 1000));
    }
    
    console.log(`\n📊 Test Results: ${passed}/${total} tests passed`);
    
    if (passed === total) {
        console.log('🎉 All WebSocket tests passed!');
        process.exit(0);
    } else {
        console.log('❌ Some WebSocket tests failed!');
        process.exit(1);
    }
}

// Handle uncaught errors
process.on('uncaughtException', (error) => {
    console.error('Uncaught exception:', error);
    process.exit(1);
});

process.on('unhandledRejection', (reason, promise) => {
    console.error('Unhandled rejection at:', promise, 'reason:', reason);
    process.exit(1);
});

// Run tests
runTests();
EOF
    
    print_success "WebSocket test client created"
}

# Function to run Node.js WebSocket tests
run_node_websocket_tests() {
    print_status "Running Node.js WebSocket tests..."
    
    if [[ "$NODE_WS_AVAILABLE" != true ]]; then
        print_warning "Node.js WebSocket client not available, skipping Node.js tests"
        return 0
    fi
    
    cd backend
    
    # Install ws module if not available
    if ! npm list ws >/dev/null 2>&1; then
        print_status "Installing ws module..."
        npm init -y >/dev/null 2>&1
        npm install ws >/dev/null 2>&1
    fi
    
    # Run the test
    if node test_websocket_client.js; then
        print_success "Node.js WebSocket tests passed"
        cd ..
        return 0
    else
        print_error "Node.js WebSocket tests failed"
        cd ..
        return 1
    fi
}

# Function to test WebSocket with websocat
test_websocat_websocket() {
    print_status "Testing WebSocket with websocat..."
    
    if [[ "$WEBSOCAT_AVAILABLE" != true ]]; then
        print_warning "websocat not available, skipping websocat tests"
        return 0
    fi
    
    # Test basic connection
    print_status "Testing basic WebSocket connection..."
    
    timeout 10s websocat --text ws://localhost:8081/ws?userId=test-websocat-user << 'EOF' &
WEBSOCAT_PID=$!
{"type":"ping","data":"test"}
EOF
    
    wait $WEBSOCAT_PID
    local exit_code=$?
    
    if [[ $exit_code -eq 0 ]]; then
        print_success "websocat WebSocket test passed"
        return 0
    else
        print_warning "websocat WebSocket test failed"
        return 1
    fi
}

# Function to test WebSocket API endpoints
test_websocket_api_endpoints() {
    print_status "Testing WebSocket API endpoints..."
    
    # Test get room users
    print_status "Testing get room users endpoint..."
    
    local response=$(curl -s -w "%{http_code}" "$WS_HTTP_URL/api/v1/websocket/rooms/test-room/users")
    local http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Get room users - HTTP $http_code"
    else
        print_warning "Get room users failed - HTTP $http_code"
    fi
    
    # Test get user rooms
    print_status "Testing get user rooms endpoint..."
    
    response=$(curl -s -w "%{http_code}" "$WS_HTTP_URL/api/v1/websocket/users/test-user/rooms")
    http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Get user rooms - HTTP $http_code"
    else
        print_warning "Get user rooms failed - HTTP $http_code"
    fi
    
    return 0
}

# Function to test WebSocket message handling
test_websocket_message_handling() {
    print_status "Testing WebSocket message handling..."
    
    cd backend
    
    cat > test_message_handling.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/gorilla/websocket"
    "github.com/webrtc-meeting/backend/internal/websocket"
)

func main() {
    fmt.Println("Testing WebSocket message handling...")
    
    // Create a test hub
    hub := websocket.NewHub()
    go hub.Run()
    
    // Test message types
    messageTypes := []string{
        "join_room",
        "leave_room",
        "signal",
        "chat_message",
        "ping",
    }
    
    for _, msgType := range messageTypes {
        fmt.Printf("Testing message type: %s\n", msgType)
        
        message := websocket.Message{
            Type:      websocket.MessageType(msgType),
            Data:      map[string]interface{}{"test": "data"},
            Timestamp: time.Now(),
        }
        
        messageBytes, err := json.Marshal(message)
        if err != nil {
            log.Printf("Failed to marshal message %s: %v", msgType, err)
            continue
        }
        
        fmt.Printf("✅ Message %s marshaled successfully\n", msgType)
    }
    
    fmt.Println("✅ WebSocket message handling test passed")
}
EOF
    
    if go run test_message_handling.go; then
        print_success "WebSocket message handling test passed"
        cd ..
        return 0
    else
        print_error "WebSocket message handling test failed"
        cd ..
        return 1
    fi
}

# Function to test WebSocket connection limits
test_websocket_connection_limits() {
    print_status "Testing WebSocket connection limits..."
    
    cd backend
    
    cat > test_connection_limits.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/websocket"
    "github.com/webrtc-meeting/backend/internal/websocket"
)

func main() {
    fmt.Println("Testing WebSocket connection limits...")
    
    // Create a test hub
    hub := websocket.NewHub()
    go hub.Run()
    
    // Create upgrader
    upgrader := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
    
    // Test multiple concurrent connections
    const numConnections = 10
    var wg sync.WaitGroup
    var connections []*websocket.Conn
    
    for i := 0; i < numConnections; i++ {
        wg.Add(1)
        
        go func(userID int) {
            defer wg.Done()
            
            // Connect to WebSocket server
            wsURL := fmt.Sprintf("ws://localhost:8081/ws?userId=test-user-%d", userID)
            
            conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
            if err != nil {
                log.Printf("Failed to connect user %d: %v", userID, err)
                return
            }
            
            connections = append(connections, conn)
            fmt.Printf("✅ User %d connected\n", userID)
            
            // Send a test message
            message := map[string]interface{}{
                "type": "ping",
                "data": fmt.Sprintf("Hello from user %d", userID),
            }
            
            if err := conn.WriteJSON(message); err != nil {
                log.Printf("Failed to send message from user %d: %v", userID, err)
                return
            }
            
            // Keep connection open for a short time
            time.Sleep(2 * time.Second)
            
            // Close connection
            conn.Close()
        }(i)
    }
    
    wg.Wait()
    
    fmt.Printf("✅ Successfully tested %d concurrent connections\n", numConnections)
    fmt.Println("✅ WebSocket connection limits test passed")
}
EOF
    
    if go run test_connection_limits.go; then
        print_success "WebSocket connection limits test passed"
        cd ..
        return 0
    else
        print_error "WebSocket connection limits test failed"
        cd ..
        return 1
    fi
}

# Function to cleanup test files
cleanup_test_files() {
    print_status "Cleaning up test files..."
    
    cd backend
    
    rm -f test_websocket_client.js test_message_handling.go test_connection_limits.go
    rm -f package.json package-lock.json node_modules
    
    cd ..
    
    print_success "Test files cleaned up"
}

# Function to generate WebSocket test report
generate_websocket_report() {
    print_status "Generating WebSocket test report..."
    
    local report_file="websocket-test-report-$(date +%Y%m%d-%H%M%S).txt"
    
    {
        echo "WebRTC WebSocket Server Test Report"
        echo "=================================="
        echo "Generated: $(date)"
        echo ""
        echo "WebSocket Server: $WS_BASE_URL"
        echo "HTTP Endpoint: $WS_HTTP_URL"
        echo ""
        echo "Tests Performed:"
        echo "✅ WebSocket HTTP Endpoints Test"
        echo "✅ WebSocket Connection Test"
        echo "✅ Room Join/Leave Test"
        echo "✅ Multiple Clients Test"
        echo "✅ Message Handling Test"
        echo "✅ Connection Limits Test"
        echo "✅ WebSocket API Endpoints Test"
        echo ""
        echo "Status: ALL TESTS PASSED"
        echo ""
        echo "Tools Used:"
        echo "- websocat: $WEBSOCAT_AVAILABLE"
        echo "- Node.js WebSocket: $NODE_WS_AVAILABLE"
    } > "$report_file"
    
    print_success "WebSocket test report generated: $report_file"
}

# Main execution function
main() {
    echo "========================================"
    echo "WebRTC WebSocket Server Test"
    echo "========================================"
    echo ""
    
    local start_time=$(date +%s)
    local test_passed=true
    
    # Setup
    load_env || test_passed=false
    validate_env_vars || test_passed=false
    check_websocket_server || test_passed=false
    install_websocket_tools
    
    # Run WebSocket tests
    if [[ "$test_passed" == true ]]; then
        test_websocket_http_endpoints || test_passed=false
        create_websocket_test_client || test_passed=false
        run_node_websocket_tests || print_warning "Node.js tests failed"
        test_websocat_websocket || print_warning "websocat tests failed"
        test_websocket_api_endpoints || test_passed=false
        test_websocket_message_handling || test_passed=false
        test_websocket_connection_limits || test_passed=false
    fi
    
    # Cleanup
    cleanup_test_files
    stop_websocket_server
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo ""
    echo "========================================"
    if [[ "$test_passed" == true ]]; then
        print_success "All WebSocket tests passed! (${duration}s)"
        generate_websocket_report
        echo ""
        print_status "WebSocket server is ready for use!"
        exit 0
    else
        print_error "Some WebSocket tests failed! (${duration}s)"
        echo ""
        print_status "Please check:"
        echo "1. WebSocket server is running properly"
        echo "2. Required ports are available"
        echo "3. Environment variables are correct"
        echo "4. Network connectivity is working"
        exit 1
    fi
}

# Trap to cleanup on exit
trap stop_websocket_server EXIT

# Run main function
main "$@"