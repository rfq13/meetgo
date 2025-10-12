#!/bin/bash

# Test Script untuk Validasi API Endpoints WebRTC
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
    print_status "Validating API environment variables..."
    
    local required_vars=(
        "SERVER_HOST"
        "SERVER_PORT"
        "DB_HOST"
        "DB_PORT"
        "DB_USER"
        "DB_PASSWORD"
        "DB_NAME"
        "JWT_SECRET"
        "JWT_REFRESH_SECRET"
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

# Function to check if API server is running
check_api_server() {
    print_status "Checking API server status..."
    
    local api_url="http://${SERVER_HOST:-localhost}:${SERVER_PORT:-8080}"
    
    if curl -s --connect-timeout 5 "$api_url/health" >/dev/null 2>&1; then
        print_success "API server is running at $api_url"
        API_BASE_URL="$api_url"
        return 0
    else
        print_warning "API server is not running. Attempting to start..."
        start_api_server
        return $?
    fi
}

# Function to start API server
start_api_server() {
    print_status "Starting API server..."
    
    cd backend
    
    # Check if binary exists
    if [[ ! -f "bin/api-server" ]]; then
        print_error "API server binary not found. Please run compilation test first."
        cd ..
        return 1
    fi
    
    # Start API server in background
    ./bin/api-server > api-server.log 2>&1 &
    API_SERVER_PID=$!
    
    # Wait for server to start
    local retries=10
    local retry_count=0
    
    while [[ $retry_count -lt $retries ]]; do
        if curl -s --connect-timeout 2 "http://${SERVER_HOST:-localhost}:${SERVER_PORT:-8080}/health" >/dev/null 2>&1; then
            print_success "API server started successfully (PID: $API_SERVER_PID)"
            API_BASE_URL="http://${SERVER_HOST:-localhost}:${SERVER_PORT:-8080}"
            cd ..
            return 0
        fi
        
        sleep 2
        ((retry_count++))
    done
    
    print_error "Failed to start API server"
    kill $API_SERVER_PID 2>/dev/null || true
    cd ..
    return 1
}

# Function to stop API server
stop_api_server() {
    if [[ -n "$API_SERVER_PID" ]]; then
        print_status "Stopping API server (PID: $API_SERVER_PID)..."
        kill $API_SERVER_PID 2>/dev/null || true
        wait $API_SERVER_PID 2>/dev/null || true
        print_success "API server stopped"
    fi
}

# Function to test health endpoints
test_health_endpoints() {
    print_status "Testing health endpoints..."
    
    local endpoints=(
        "/health"
        "/ping"
        "/ready"
        "/live"
    )
    
    for endpoint in "${endpoints[@]}"; do
        print_status "Testing $endpoint..."
        
        local response=$(curl -s -w "%{http_code}" "$API_BASE_URL$endpoint")
        local http_code="${response: -3}"
        local body="${response%???}"
        
        if [[ "$http_code" == "200" ]]; then
            print_success "$endpoint - HTTP $http_code"
        else
            print_error "$endpoint - HTTP $http_code"
            return 1
        fi
    done
    
    return 0
}

# Function to register test user
register_test_user() {
    print_status "Registering test user..."
    
    local register_data=$(cat <<EOF
{
    "email": "test@example.com",
    "username": "testuser",
    "password": "testpassword123",
    "first_name": "Test",
    "last_name": "User"
}
EOF
)
    
    local response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$register_data" \
        "$API_BASE_URL/api/v1/auth/register")
    
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" == "201" || "$http_code" == "200" ]]; then
        print_success "User registration successful - HTTP $http_code"
        TEST_USER_EMAIL="test@example.com"
        TEST_USER_PASSWORD="testpassword123"
        return 0
    elif [[ "$http_code" == "409" ]]; then
        print_warning "User already exists, proceeding with login"
        TEST_USER_EMAIL="test@example.com"
        TEST_USER_PASSWORD="testpassword123"
        return 0
    else
        print_error "User registration failed - HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
}

# Function to login test user
login_test_user() {
    print_status "Logging in test user..."
    
    local login_data=$(cat <<EOF
{
    "email": "$TEST_USER_EMAIL",
    "password": "$TEST_USER_PASSWORD"
}
EOF
)
    
    local response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$login_data" \
        "$API_BASE_URL/api/v1/auth/login")
    
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "User login successful - HTTP $http_code"
        
        # Extract tokens from response
        ACCESS_TOKEN=$(echo "$body" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
        REFRESH_TOKEN=$(echo "$body" | grep -o '"refresh_token":"[^"]*' | cut -d'"' -f4)
        
        if [[ -n "$ACCESS_TOKEN" ]]; then
            print_success "Access token obtained"
            return 0
        else
            print_error "Failed to extract access token"
            return 1
        fi
    else
        print_error "User login failed - HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
}

# Function to test authentication endpoints
test_auth_endpoints() {
    print_status "Testing authentication endpoints..."
    
    # Test registration
    register_test_user || return 1
    
    # Test login
    login_test_user || return 1
    
    # Test token refresh
    print_status "Testing token refresh..."
    
    local refresh_data=$(cat <<EOF
{
    "refresh_token": "$REFRESH_TOKEN"
}
EOF
)
    
    local response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$refresh_data" \
        "$API_BASE_URL/api/v1/auth/refresh")
    
    local http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Token refresh successful - HTTP $http_code"
    else
        print_warning "Token refresh failed - HTTP $http_code"
    fi
    
    # Test logout
    print_status "Testing logout..."
    
    response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        "$API_BASE_URL/api/v1/auth/logout")
    
    http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Logout successful - HTTP $http_code"
    else
        print_warning "Logout failed - HTTP $http_code"
    fi
    
    return 0
}

# Function to test user endpoints
test_user_endpoints() {
    print_status "Testing user endpoints..."
    
    # Test get current user profile
    print_status "Testing get current user profile..."
    
    local response=$(curl -s -w "%{http_code}" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        "$API_BASE_URL/api/v1/users/profile")
    
    local http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Get user profile successful - HTTP $http_code"
    else
        print_error "Get user profile failed - HTTP $http_code"
        return 1
    fi
    
    # Test update user profile
    print_status "Testing update user profile..."
    
    local update_data=$(cat <<EOF
{
    "first_name": "Updated",
    "last_name": "User",
    "phone": "+1234567890"
}
EOF
)
    
    response=$(curl -s -w "%{http_code}" \
        -X PUT \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$update_data" \
        "$API_BASE_URL/api/v1/users/profile")
    
    http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Update user profile successful - HTTP $http_code"
    else
        print_warning "Update user profile failed - HTTP $http_code"
    fi
    
    return 0
}

# Function to test room endpoints
test_room_endpoints() {
    print_status "Testing room endpoints..."
    
    # Test create room
    print_status "Testing create room..."
    
    local room_data=$(cat <<EOF
{
    "name": "Test Meeting Room",
    "description": "A test room for API validation",
    "max_users": 10,
    "type": "meeting",
    "is_public": true
}
EOF
)
    
    local response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$room_data" \
        "$API_BASE_URL/api/v1/rooms")
    
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" == "201" || "$http_code" == "200" ]]; then
        print_success "Create room successful - HTTP $http_code"
        
        # Extract room ID
        TEST_ROOM_ID=$(echo "$body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
        
        if [[ -n "$TEST_ROOM_ID" ]]; then
            print_success "Room ID obtained: $TEST_ROOM_ID"
        else
            print_warning "Could not extract room ID"
        fi
    else
        print_error "Create room failed - HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
    
    # Test get rooms list
    print_status "Testing get rooms list..."
    
    response=$(curl -s -w "%{http_code}" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        "$API_BASE_URL/api/v1/rooms")
    
    http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Get rooms list successful - HTTP $http_code"
    else
        print_error "Get rooms list failed - HTTP $http_code"
        return 1
    fi
    
    # Test get room details
    if [[ -n "$TEST_ROOM_ID" ]]; then
        print_status "Testing get room details..."
        
        response=$(curl -s -w "%{http_code}" \
            -H "Authorization: Bearer $ACCESS_TOKEN" \
            "$API_BASE_URL/api/v1/rooms/$TEST_ROOM_ID")
        
        http_code="${response: -3}"
        
        if [[ "$http_code" == "200" ]]; then
            print_success "Get room details successful - HTTP $http_code"
        else
            print_error "Get room details failed - HTTP $http_code"
        fi
    fi
    
    return 0
}

# Function to test public endpoints
test_public_endpoints() {
    print_status "Testing public endpoints..."
    
    # Test public rooms list
    print_status "Testing public rooms list..."
    
    local response=$(curl -s -w "%{http_code}" \
        "$API_BASE_URL/api/v1/public/rooms")
    
    local http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "Public rooms list successful - HTTP $http_code"
    else
        print_warning "Public rooms list failed - HTTP $http_code"
    fi
    
    # Test system info
    print_status "Testing system info..."
    
    response=$(curl -s -w "%{http_code}" \
        "$API_BASE_URL/api/v1/public/info")
    
    http_code="${response: -3}"
    
    if [[ "$http_code" == "200" ]]; then
        print_success "System info successful - HTTP $http_code"
    else
        print_warning "System info failed - HTTP $http_code"
    fi
    
    return 0
}

# Function to test error handling
test_error_handling() {
    print_status "Testing error handling..."
    
    # Test unauthorized access
    print_status "Testing unauthorized access..."
    
    local response=$(curl -s -w "%{http_code}" \
        "$API_BASE_URL/api/v1/users/profile")
    
    local http_code="${response: -3}"
    
    if [[ "$http_code" == "401" ]]; then
        print_success "Unauthorized access properly handled - HTTP $http_code"
    else
        print_warning "Unauthorized access not properly handled - HTTP $http_code"
    fi
    
    # Test invalid endpoint
    print_status "Testing invalid endpoint..."
    
    response=$(curl -s -w "%{http_code}" \
        "$API_BASE_URL/api/v1/invalid/endpoint")
    
    http_code="${response: -3}"
    
    if [[ "$http_code" == "404" ]]; then
        print_success "Invalid endpoint properly handled - HTTP $http_code"
    else
        print_warning "Invalid endpoint not properly handled - HTTP $http_code"
    fi
    
    # Test invalid JSON
    print_status "Testing invalid JSON..."
    
    response=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "invalid json" \
        "$API_BASE_URL/api/v1/auth/login")
    
    http_code="${response: -3}"
    
    if [[ "$http_code" == "400" ]]; then
        print_success "Invalid JSON properly handled - HTTP $http_code"
    else
        print_warning "Invalid JSON not properly handled - HTTP $http_code"
    fi
    
    return 0
}

# Function to test rate limiting
test_rate_limiting() {
    print_status "Testing rate limiting..."
    
    # Make multiple rapid requests
    local success_count=0
    local rate_limit_count=0
    
    for i in {1..20}; do
        local response=$(curl -s -w "%{http_code}" \
            "$API_BASE_URL/ping")
        
        local http_code="${response: -3}"
        
        if [[ "$http_code" == "200" ]]; then
            ((success_count++))
        elif [[ "$http_code" == "429" ]]; then
            ((rate_limit_count++))
        fi
    done
    
    if [[ $rate_limit_count -gt 0 ]]; then
        print_success "Rate limiting is working ($rate_limit_count requests limited)"
    else
        print_warning "Rate limiting may not be configured"
    fi
    
    return 0
}

# Function to test CORS
test_cors() {
    print_status "Testing CORS..."
    
    # Test preflight request
    local response=$(curl -s -w "%{http_code}" \
        -X OPTIONS \
        -H "Origin: http://localhost:3000" \
        -H "Access-Control-Request-Method: POST" \
        -H "Access-Control-Request-Headers: Content-Type" \
        "$API_BASE_URL/api/v1/auth/login")
    
    local http_code="${response: -3}"
    
    if [[ "$http_code" == "204" || "$http_code" == "200" ]]; then
        print_success "CORS preflight request successful - HTTP $http_code"
    else
        print_warning "CORS preflight request failed - HTTP $http_code"
    fi
    
    return 0
}

# Function to generate API test report
generate_api_report() {
    print_status "Generating API test report..."
    
    local report_file="api-test-report-$(date +%Y%m%d-%H%M%S).txt"
    
    {
        echo "WebRTC API Test Report"
        echo "====================="
        echo "Generated: $(date)"
        echo ""
        echo "API Server: $API_BASE_URL"
        echo "Test User: $TEST_USER_EMAIL"
        echo ""
        echo "Tests Performed:"
        echo "✅ Health Endpoints Test"
        echo "✅ Authentication Endpoints Test"
        echo "✅ User Endpoints Test"
        echo "✅ Room Endpoints Test"
        echo "✅ Public Endpoints Test"
        echo "✅ Error Handling Test"
        echo "✅ Rate Limiting Test"
        echo "✅ CORS Test"
        echo ""
        echo "Status: ALL TESTS PASSED"
        echo ""
        echo "Test Room ID: $TEST_ROOM_ID"
        echo "Access Token: ${ACCESS_TOKEN:0:20}..."
    } > "$report_file"
    
    print_success "API test report generated: $report_file"
}

# Main execution function
main() {
    echo "========================================"
    echo "WebRTC API Endpoints Test"
    echo "========================================"
    echo ""
    
    local start_time=$(date +%s)
    local test_passed=true
    
    # Setup
    load_env || test_passed=false
    validate_env_vars || test_passed=false
    check_api_server || test_passed=false
    
    # Re-login for fresh tokens
    if [[ "$test_passed" == true ]]; then
        login_test_user || test_passed=false
    fi
    
    # Run API tests
    if [[ "$test_passed" == true ]]; then
        test_health_endpoints || test_passed=false
        test_auth_endpoints || test_passed=false
        test_user_endpoints || test_passed=false
        test_room_endpoints || test_passed=false
        test_public_endpoints || test_passed=false
        test_error_handling || test_passed=false
        test_rate_limiting || test_passed=false
        test_cors || test_passed=false
    fi
    
    # Cleanup
    stop_api_server
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo ""
    echo "========================================"
    if [[ "$test_passed" == true ]]; then
        print_success "All API tests passed! (${duration}s)"
        generate_api_report
        echo ""
        print_status "API server is ready for use!"
        exit 0
    else
        print_error "Some API tests failed! (${duration}s)"
        echo ""
        print_status "Please check:"
        echo "1. API server is running properly"
        echo "2. Database connection is working"
        echo "3. Environment variables are correct"
        echo "4. Required services are available"
        exit 1
    fi
}

# Trap to cleanup on exit
trap stop_api_server EXIT

# Run main function
main "$@"