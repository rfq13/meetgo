#!/bin/bash

# Test Script untuk Validasi Kompilasi Backend WebRTC
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

# Function to validate Go installation
validate_go() {
    print_status "Validating Go installation..."
    
    if ! command_exists go; then
        print_error "Go is not installed or not in PATH"
        return 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go version: $GO_VERSION"
    
    # Check minimum Go version (1.21+)
    if ! go version | grep -E "go1\.(2[1-9]|[3-9][0-9]|[0-9]{3,})" >/dev/null; then
        print_warning "Go version should be 1.21 or higher for optimal compatibility"
    fi
    
    return 0
}

# Function to validate project structure
validate_project_structure() {
    print_status "Validating project structure..."
    
    local required_files=(
        "backend/go.mod"
        "backend/go.sum"
        "backend/cmd/api/main.go"
        "backend/cmd/websocket/main.go"
        "backend/internal/config/config.go"
        "backend/internal/database/connection.go"
        "backend/internal/api/router.go"
        "backend/internal/auth/service.go"
        "backend/internal/auth/handler.go"
        "backend/internal/user/service.go"
        "backend/internal/user/handler.go"
        "backend/internal/room/service.go"
        "backend/internal/room/handler.go"
        "backend/internal/webrtc/janus_client.go"
        "backend/internal/websocket/handler.go"
        "backend/models/user.go"
        "backend/models/room.go"
    )
    
    local missing_files=()
    
    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            missing_files+=("$file")
        fi
    done
    
    if [[ ${#missing_files[@]} -gt 0 ]]; then
        print_error "Missing required files:"
        for file in "${missing_files[@]}"; do
            echo "  - $file"
        done
        return 1
    fi
    
    print_success "All required files are present"
    return 0
}

# Function to validate go.mod
validate_go_mod() {
    print_status "Validating go.mod file..."
    
    cd backend
    
    # Check if go.mod is valid
    if ! go mod verify >/dev/null 2>&1; then
        print_error "go.mod verification failed"
        cd ..
        return 1
    fi
    
    # Check for required dependencies
    local required_deps=(
        "github.com/gin-gonic/gin"
        "github.com/golang-jwt/jwt/v5"
        "github.com/google/uuid"
        "github.com/gorilla/websocket"
        "github.com/joho/godotenv"
        "github.com/sirupsen/logrus"
        "gorm.io/driver/postgres"
        "gorm.io/gorm"
    )
    
    local missing_deps=()
    
    for dep in "${required_deps[@]}"; do
        if ! grep -q "$dep" go.mod; then
            missing_deps+=("$dep")
        fi
    done
    
    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        print_error "Missing required dependencies:"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        cd ..
        return 1
    fi
    
    print_success "go.mod file is valid"
    cd ..
    return 0
}

# Function to download dependencies
download_dependencies() {
    print_status "Downloading Go dependencies..."
    
    cd backend
    
    if ! go mod download; then
        print_error "Failed to download dependencies"
        cd ..
        return 1
    fi
    
    if ! go mod tidy; then
        print_error "Failed to tidy dependencies"
        cd ..
        return 1
    fi
    
    print_success "Dependencies downloaded successfully"
    cd ..
    return 0
}

# Function to compile API server
compile_api_server() {
    print_status "Compiling API server..."
    
    cd backend
    
    # Create build directory
    mkdir -p bin
    
    # Compile API server
    if ! go build -o bin/api-server ./cmd/api; then
        print_error "Failed to compile API server"
        cd ..
        return 1
    fi
    
    # Check if binary was created
    if [[ ! -f "bin/api-server" ]]; then
        print_error "API server binary was not created"
        cd ..
        return 1
    fi
    
    # Check binary size
    BINARY_SIZE=$(stat -c%s "bin/api-server" 2>/dev/null || stat -f%z "bin/api-server" 2>/dev/null)
    print_success "API server compiled successfully (size: $BINARY_SIZE bytes)"
    
    cd ..
    return 0
}

# Function to compile WebSocket server
compile_websocket_server() {
    print_status "Compiling WebSocket server..."
    
    cd backend
    
    # Compile WebSocket server
    if ! go build -o bin/websocket-server ./cmd/websocket; then
        print_error "Failed to compile WebSocket server"
        cd ..
        return 1
    fi
    
    # Check if binary was created
    if [[ ! -f "bin/websocket-server" ]]; then
        print_error "WebSocket server binary was not created"
        cd ..
        return 1
    fi
    
    # Check binary size
    BINARY_SIZE=$(stat -c%s "bin/websocket-server" 2>/dev/null || stat -f%z "bin/websocket-server" 2>/dev/null)
    print_success "WebSocket server compiled successfully (size: $BINARY_SIZE bytes)"
    
    cd ..
    return 0
}

# Function to run basic syntax validation
validate_syntax() {
    print_status "Validating Go syntax..."
    
    cd backend
    
    # Check for syntax errors in all Go files
    if ! go vet ./...; then
        print_error "Go vet found issues"
        cd ..
        return 1
    fi
    
    # Run go fmt check
    UNFORMATTED=$(gofmt -l .)
    if [[ -n "$UNFORMATTED" ]]; then
        print_warning "Following files are not properly formatted:"
        echo "$UNFORMATTED"
        print_status "Running go fmt..."
        gofmt -w .
    fi
    
    print_success "Syntax validation passed"
    cd ..
    return 0
}

# Function to run tests
run_tests() {
    print_status "Running unit tests..."
    
    cd backend
    
    # Run tests with coverage
    if ! go test -v -cover ./...; then
        print_error "Some tests failed"
        cd ..
        return 1
    fi
    
    print_success "All tests passed"
    cd ..
    return 0
}

# Function to check for security issues
check_security() {
    print_status "Running basic security checks..."
    
    cd backend
    
    # Check for hardcoded secrets (basic check)
    if grep -r -i "password\|secret\|key" --include="*.go" . | grep -v "\.env\|config\|test" | head -5; then
        print_warning "Potential hardcoded secrets found. Please review."
    fi
    
    # Check for SQL injection vulnerabilities (basic check)
    if grep -r "fmt\.Sprintf.*%s.*SELECT\|fmt\.Sprintf.*%s.*INSERT\|fmt\.Sprintf.*%s.*UPDATE\|fmt\.Sprintf.*%s.*DELETE" --include="*.go" .; then
        print_warning "Potential SQL injection vulnerabilities found. Please use parameterized queries."
    fi
    
    print_success "Basic security checks completed"
    cd ..
    return 0
}

# Function to generate build report
generate_report() {
    print_status "Generating build report..."
    
    local report_file="build-report-$(date +%Y%m%d-%H%M%S).txt"
    
    {
        echo "WebRTC Backend Build Report"
        echo "==========================="
        echo "Generated: $(date)"
        echo ""
        echo "Go Version: $(go version)"
        echo ""
        echo "Build Status: SUCCESS"
        echo ""
        echo "Binaries Created:"
        if [[ -f "backend/bin/api-server" ]]; then
            echo "- API Server: backend/bin/api-server ($(stat -c%s backend/bin/api-server 2>/dev/null || stat -f%z backend/bin/api-server 2>/dev/null) bytes)"
        fi
        if [[ -f "backend/bin/websocket-server" ]]; then
            echo "- WebSocket Server: backend/bin/websocket-server ($(stat -c%s backend/bin/websocket-server 2>/dev/null || stat -f%z backend/bin/websocket-server 2>/dev/null) bytes)"
        fi
        echo ""
        echo "Dependencies:"
        go list -m all | head -20
        echo "..."
    } > "$report_file"
    
    print_success "Build report generated: $report_file"
}

# Main execution function
main() {
    echo "========================================"
    echo "WebRTC Backend Compilation Test"
    echo "========================================"
    echo ""
    
    local start_time=$(date +%s)
    local test_passed=true
    
    # Run all validation steps
    validate_go || test_passed=false
    validate_project_structure || test_passed=false
    validate_go_mod || test_passed=false
    download_dependencies || test_passed=false
    validate_syntax || test_passed=false
    compile_api_server || test_passed=false
    compile_websocket_server || test_passed=false
    run_tests || test_passed=false
    check_security || test_passed=false
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo ""
    echo "========================================"
    if [[ "$test_passed" == true ]]; then
        print_success "All compilation tests passed! (${duration}s)"
        generate_report
        echo ""
        print_status "Next steps:"
        echo "1. Run: ./backend/bin/api-server (to start API server)"
        echo "2. Run: ./backend/bin/websocket-server (to start WebSocket server)"
        echo "3. Configure environment variables in .env file"
        exit 0
    else
        print_error "Some compilation tests failed! (${duration}s)"
        echo ""
        print_status "Please fix the errors above and try again."
        exit 1
    fi
}

# Run main function
main "$@"