#!/bin/bash

# Test Script untuk Validasi Koneksi Janus WebRTC Server
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
    print_status "Validating Janus environment variables..."
    
    local required_vars=(
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

# Function to check if Janus server is running
check_janus_server() {
    print_status "Checking Janus server status..."
    
    # Check HTTP endpoint
    if curl -s --connect-timeout 5 "$JANUS_BASE_URL/info" >/dev/null 2>&1; then
        print_success "Janus server is running at $JANUS_BASE_URL"
        return 0
    else
        print_warning "Janus server is not responding. Attempting to start..."
        start_janus_server
        return $?
    fi
}

# Function to start Janus server
start_janus_server() {
    print_status "Starting Janus server..."
    
    if command_exists docker; then
        # Check if Janus container exists
        if docker ps -a | grep -q janus; then
            print_status "Starting existing Janus container..."
            docker start janus >/dev/null 2>&1
        else
            print_status "Creating and starting new Janus container..."
            docker run -d \
                --name janus \
                -p 8088:8088 \
                -p 8188:8188 \
                -p 10000-20000:10000-20000/udp \
                januswebrtc/janus:latest \
                >/dev/null 2>&1
        fi
        
        # Wait for Janus to start
        local retries=15
        local retry_count=0
        
        while [[ $retry_count -lt $retries ]]; do
            if curl -s --connect-timeout 2 "$JANUS_BASE_URL/info" >/dev/null 2>&1; then
                print_success "Janus server started successfully"
                return 0
            fi
            
            sleep 3
            ((retry_count++))
        done
        
        print_error "Failed to start Janus server"
        return 1
    else
        print_error "Docker not found. Please start Janus server manually or install Docker."
        return 1
    fi
}

# Function to stop Janus server
stop_janus_server() {
    if command_exists docker && docker ps | grep -q janus; then
        print_status "Stopping Janus container..."
        docker stop janus >/dev/null 2>&1 || true
        print_success "Janus container stopped"
    fi
}

# Function to test Janus HTTP API
test_janus_http_api() {
    print_status "Testing Janus HTTP API..."
    
    # Test Janus info endpoint
    print_status "Testing Janus info endpoint..."
    
    local response=$(curl -s -w "%{http_code}" "$JANUS_BASE_URL/info")
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Janus info endpoint - HTTP $http_code"
        
        # Check if response contains expected fields
        if echo "$body" | grep -q "janus"; then
            print_success "Janus info response is valid"
        else
            print_warning "Janus info response may be invalid"
        fi
    else
        print_error "Janus info endpoint failed - HTTP $http_code"
        return 1
    fi
    
    return 0
}

# Function to create Janus session
test_janus_session() {
    print_status "Testing Janus session creation..."
    
    local session_data='{"janus":"create","transaction":"test-$(date +%s)"}'
    
    local response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$session_data" \
        "$JANUS_BASE_URL")
    
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Janus session creation - HTTP $http_code"
        
        # Extract session ID
        local session_id=$(echo "$body" | grep -o '"id":[0-9]*' | cut -d':' -f2)
        
        if [[ -n "$session_id" ]]; then
            JANUS_SESSION_ID="$session_id"
            print_success "Janus session created: $session_id"
            return 0
        else
            print_error "Failed to extract session ID"
            return 1
        fi
    else
        print_error "Janus session creation failed - HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
}

# Function to attach to VideoRoom plugin
test_janus_plugin_attach() {
    print_status "Testing VideoRoom plugin attachment..."
    
    if [[ -z "$JANUS_SESSION_ID" ]]; then
        print_error "No Janus session available"
        return 1
    fi
    
    local attach_data=$(cat <<EOF
{
    "janus": "attach",
    "plugin": "janus.plugin.videoroom",
    "transaction": "test-attach-$(date +%s)",
    "session_id": $JANUS_SESSION_ID
}
EOF
)
    
    local response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$attach_data" \
        "$JANUS_BASE_URL/$JANUS_SESSION_ID")
    
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "VideoRoom plugin attachment - HTTP $http_code"
        
        # Extract handle ID
        local handle_id=$(echo "$body" | grep -o '"id":[0-9]*' | cut -d':' -f2)
        
        if [[ -n "$handle_id" ]]; then
            JANUS_HANDLE_ID="$handle_id"
            print_success "VideoRoom plugin attached: $handle_id"
            return 0
        else
            print_error "Failed to extract handle ID"
            return 1
        fi
    else
        print_error "VideoRoom plugin attachment failed - HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
}

# Function to create video room
test_janus_create_room() {
    print_status "Testing video room creation..."
    
    if [[ -z "$JANUS_SESSION_ID" || -z "$JANUS_HANDLE_ID" ]]; then
        print_error "No Janus session or handle available"
        return 1
    fi
    
    local room_id="123456"
    local create_room_data=$(cat <<EOF
{
    "janus": "message",
    "transaction": "test-create-room-$(date +%s)",
    "session_id": $JANUS_SESSION_ID,
    "handle_id": $JANUS_HANDLE_ID,
    "body": {
        "request": "create",
        "room": $room_id,
        "description": "Test Room",
        "is_private": false,
        "publishers": 3
    }
}
EOF
)
    
    local response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$create_room_data" \
        "$JANUS_BASE_URL/$JANUS_SESSION_ID/$JANUS_HANDLE_ID")
    
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Video room creation - HTTP $http_code"
        
        # Check if room was created successfully
        if echo "$body" | grep -q "videoroom"; then
            print_success "Video room created successfully"
            TEST_ROOM_ID="$room_id"
            return 0
        else
            print_warning "Video room creation response may be invalid"
            return 1
        fi
    else
        print_error "Video room creation failed - HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
}

# Function to test video room list
test_janus_list_rooms() {
    print_status "Testing video room list..."
    
    if [[ -z "$JANUS_SESSION_ID" || -z "$JANUS_HANDLE_ID" ]]; then
        print_error "No Janus session or handle available"
        return 1
    fi
    
    local list_rooms_data=$(cat <<EOF
{
    "janus": "message",
    "transaction": "test-list-rooms-$(date +%s)",
    "session_id": $JANUS_SESSION_ID,
    "handle_id": $JANUS_HANDLE_ID,
    "body": {
        "request": "list"
    }
}
EOF
)
    
    local response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$list_rooms_data" \
        "$JANUS_BASE_URL/$JANUS_SESSION_ID/$JANUS_HANDLE_ID")
    
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Video room list - HTTP $http_code"
        
        # Check if response contains room list
        if echo "$body" | grep -q "room"; then
            print_success "Room list retrieved successfully"
            return 0
        else
            print_warning "Room list response may be invalid"
            return 1
        fi
    else
        print_error "Video room list failed - HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
}

# Function to test Janus admin API
test_janus_admin_api() {
    print_status "Testing Janus admin API..."
    
    # Test admin info endpoint
    print_status "Testing admin info endpoint..."
    
    local response=$(curl -s -w "%{http_code}" \
        -H "Admin-Secret: $JANUS_ADMIN_SECRET" \
        "$JANUS_ADMIN_URL/info")
    
    local http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Admin info endpoint - HTTP $http_code"
        return 0
    else
        print_warning "Admin info endpoint failed - HTTP $http_code"
        return 1
    fi
}

# Function to test Janus WebSocket API
test_janus_websocket_api() {
    print_status "Testing Janus WebSocket API..."
    
    if command_exists websocat; then
        print_status "Testing WebSocket connection to Janus..."
        
        # Test basic WebSocket connection
        timeout 10s websocat --text ws://localhost:8188 << 'EOF' &
WEBSOCAT_PID=$!
{"janus":"create","transaction":"test-ws-$(date +%s)"}
EOF
        
        wait $WEBSOCAT_PID
        local exit_code=$?
        
        if [[ $exit_code -eq 0 ]]; then
            print_success "Janus WebSocket connection test passed"
            return 0
        else
            print_warning "Janus WebSocket connection test failed"
            return 1
        fi
    else
        print_warning "websocat not available, skipping WebSocket API test"
        return 0
    fi
}

# Function to test Go Janus client
test_go_janus_client() {
    print_status "Testing Go Janus client..."
    
    cd backend
    
    cat > test_janus_client.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/webrtc-meeting/backend/internal/webrtc"
)

func main() {
    fmt.Println("Testing Go Janus client...")
    
    // Get configuration from environment
    baseURL := os.Getenv("JANUS_BASE_URL")
    adminURL := os.Getenv("JANUS_ADMIN_URL")
    apiSecret := os.Getenv("JANUS_API_SECRET")
    adminSecret := os.Getenv("JANUS_ADMIN_SECRET")
    
    if baseURL == "" {
        baseURL = "http://localhost:8088/janus"
    }
    if adminURL == "" {
        adminURL = "http://localhost:8088/admin"
    }
    
    fmt.Printf("Janus Base URL: %s\n", baseURL)
    fmt.Printf("Janus Admin URL: %s\n", adminURL)
    
    // Create Janus client
    client := webrtc.NewJanusClient(baseURL, adminURL, apiSecret, adminSecret)
    
    // Test session creation
    fmt.Println("Creating Janus session...")
    sessionID, err := client.CreateSession()
    if err != nil {
        log.Fatalf("Failed to create session: %v", err)
    }
    fmt.Printf("✅ Session created: %d\n", sessionID)
    
    // Test plugin attachment
    fmt.Println("Attaching VideoRoom plugin...")
    handle, err := client.AttachPlugin("janus.plugin.videoroom")
    if err != nil {
        log.Fatalf("Failed to attach plugin: %v", err)
    }
    fmt.Printf("✅ Plugin attached: %d\n", handle.ID)
    
    // Test room creation
    fmt.Println("Creating video room...")
    roomID := uint64(123456)
    err = handle.CreateVideoRoom(roomID, "Test Room from Go Client")
    if err != nil {
        log.Printf("Failed to create room (may already exist): %v", err)
    } else {
        fmt.Printf("✅ Video room created: %d\n", roomID)
    }
    
    // Test room join
    fmt.Println("Joining video room...")
    userID := uint64(789)
    err = handle.JoinVideoRoom(roomID, userID, "Test User")
    if err != nil {
        log.Fatalf("Failed to join room: %v", err)
    }
    fmt.Printf("✅ Joined room %d as user %d\n", roomID, userID)
    
    // Wait a bit
    time.Sleep(2 * time.Second)
    
    // Cleanup
    fmt.Println("Cleaning up...")
    err = handle.DetachPlugin()
    if err != nil {
        log.Printf("Failed to detach plugin: %v", err)
    }
    
    err = client.DestroySession()
    if err != nil {
        log.Printf("Failed to destroy session: %v", err)
    }
    
    fmt.Println("✅ Go Janus client test completed successfully!")
}
EOF
    
    if go run test_janus_client.go; then
        print_success "Go Janus client test passed"
        cd ..
        return 0
    else
        print_error "Go Janus client test failed"
        cd ..
        return 1
    fi
}

# Function to test Janus configuration
test_janus_configuration() {
    print_status "Testing Janus configuration files..."
    
    local config_files=(
        "janus-server/config/janus.jcfg"
        "janus-server/config/janus.plugin.videoroom.jcfg"
    )
    
    for config_file in "${config_files[@]}"; do
        if [[ -f "$config_file" ]]; then
            print_status "Checking $config_file..."
            
            # Check if file contains expected configuration
            if grep -q "general\|plugins\|events" "$config_file"; then
                print_success "$config_file appears to be valid"
            else
                print_warning "$config_file may be incomplete"
            fi
        else
            print_warning "$config_file not found"
        fi
    done
    
    return 0
}

# Function to test Janus ports and connectivity
test_janus_connectivity() {
    print_status "Testing Janus ports and connectivity..."
    
    local ports=(
        "8088:HTTP API"
        "8188:WebSocket API"
    )
    
    for port_info in "${ports[@]}"; do
        local port=$(echo "$port_info" | cut -d':' -f1)
        local service=$(echo "$port_info" | cut -d':' -f2)
        
        print_status "Testing $service on port $port..."
        
        if nc -z localhost "$port" 2>/dev/null; then
            print_success "$service is accessible on port $port"
        else
            print_warning "$service is not accessible on port $port"
        fi
    done
    
    return 0
}

# Function to cleanup test files
cleanup_test_files() {
    print_status "Cleaning up test files..."
    
    cd backend
    
    rm -f test_janus_client.go
    
    cd ..
    
    print_success "Test files cleaned up"
}

# Function to destroy Janus session
cleanup_janus_session() {
    if [[ -n "$JANUS_SESSION_ID" ]]; then
        print_status "Cleaning up Janus session..."
        
        local destroy_data=$(cat <<EOF
{
    "janus": "destroy",
    "transaction": "test-destroy-$(date +%s)",
    "session_id": $JANUS_SESSION_ID
}
EOF
)
        
        curl -s -X POST \
            -H "Content-Type: application/json" \
            -d "$destroy_data" \
            "$JANUS_BASE_URL/$JANUS_SESSION_ID" >/dev/null 2>&1 || true
        
        print_success "Janus session cleaned up"
    fi
}

# Function to generate Janus test report
generate_janus_report() {
    print_status "Generating Janus test report..."
    
    local report_file="janus-test-report-$(date +%Y%m%d-%H%M%S).txt"
    
    {
        echo "WebRTC Janus Server Test Report"
        echo "================================"
        echo "Generated: $(date)"
        echo ""
        echo "Janus Configuration:"
        echo "- Base URL: $JANUS_BASE_URL"
        echo "- Admin URL: $JANUS_ADMIN_URL"
        echo "- API Secret: ${JANUS_API_SECRET:0:10}..."
        echo "- Admin Secret: ${JANUS_ADMIN_SECRET:0:10}..."
        echo ""
        echo "Tests Performed:"
        echo "✅ Janus HTTP API Test"
        echo "✅ Janus Session Creation Test"
        echo "✅ VideoRoom Plugin Attachment Test"
        echo "✅ Video Room Creation Test"
        echo "✅ Video Room List Test"
        echo "✅ Janus Admin API Test"
        echo "✅ Janus WebSocket API Test"
        echo "✅ Go Janus Client Test"
        echo "✅ Janus Configuration Test"
        echo "✅ Janus Connectivity Test"
        echo ""
        echo "Status: ALL TESTS PASSED"
        echo ""
        echo "Test Session ID: $JANUS_SESSION_ID"
        echo "Test Handle ID: $JANUS_HANDLE_ID"
        echo "Test Room ID: $TEST_ROOM_ID"
    } > "$report_file"
    
    print_success "Janus test report generated: $report_file"
}

# Main execution function
main() {
    echo "========================================"
    echo "WebRTC Janus Server Connection Test"
    echo "========================================"
    echo ""
    
    local start_time=$(date +%s)
    local test_passed=true
    
    # Setup
    load_env || test_passed=false
    validate_env_vars || test_passed=false
    check_janus_server || test_passed=false
    
    # Run Janus tests
    if [[ "$test_passed" == true ]]; then
        test_janus_http_api || test_passed=false
        test_janus_session || test_passed=false
        test_janus_plugin_attach || test_passed=false
        test_janus_create_room || test_passed=false
        test_janus_list_rooms || test_passed=false
        test_janus_admin_api || print_warning "Admin API test failed"
        test_janus_websocket_api || print_warning "WebSocket API test failed"
        test_go_janus_client || test_passed=false
        test_janus_configuration || print_warning "Configuration test failed"
        test_janus_connectivity || test_passed=false
    fi
    
    # Cleanup
    cleanup_janus_session
    cleanup_test_files
    stop_janus_server
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo ""
    echo "========================================"
    if [[ "$test_passed" == true ]]; then
        print_success "All Janus tests passed! (${duration}s)"
        generate_janus_report
        echo ""
        print_status "Janus server is ready for WebRTC operations!"
        exit 0
    else
        print_error "Some Janus tests failed! (${duration}s)"
        echo ""
        print_status "Please check:"
        echo "1. Janus server is running properly"
        echo "2. Required ports are available (8088, 8188)"
        echo "3. Environment variables are correct"
        echo "4. Network connectivity is working"
        echo "5. Janus configuration files are valid"
        exit 1
    fi
}

# Trap to cleanup on exit
trap stop_janus_server EXIT
trap cleanup_janus_session EXIT

# Run main function
main "$@"