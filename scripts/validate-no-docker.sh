#!/bin/bash

# WebRTC Backend Validation Script (No Docker Required)
# This script runs validation tests without Docker dependencies

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

# Function to run test and capture result
run_test() {
    local test_name="$1"
    local test_script="$2"
    
    print_status "Running $test_name test..."
    
    if [ -f "$test_script" ] && [ -x "$test_script" ]; then
        if "$test_script"; then
            print_success "$test_name test passed"
            return 0
        else
            print_error "$test_name test failed"
            return 1
        fi
    else
        print_error "$test_name script not found or not executable"
        return 1
    fi
}

# Function to generate summary report
generate_report() {
    local passed=$1
    local total=$2
    local failed=$((total - passed))
    
    echo "========================================"
    echo "VALIDATION SUMMARY"
    echo "========================================"
    echo "Total tests: $total"
    echo -e "Passed: ${GREEN}$passed${NC}"
    echo -e "Failed: ${RED}$failed${NC}"
    
    if [ $passed -eq $total ]; then
        echo -e "\n${GREEN}🎉 All tests passed!${NC}"
        print_success "Backend is ready for deployment"
    else
        echo -e "\n${YELLOW}⚠️  Some tests failed${NC}"
        print_warning "Please review failed tests before deployment"
    fi
    
    echo "========================================"
}

# Function to cleanup
cleanup() {
    print_status "Cleaning up..."
    # Kill any running processes
    pkill -f "api-server" 2>/dev/null || true
    pkill -f "websocket-server" 2>/dev/null || true
    print_success "Cleanup completed"
}

# Main execution
main() {
    echo "========================================"
    echo "WEBRTC BACKEND VALIDATION (NO DOCKER)"
    echo "========================================"
    echo
    
    # Check basic prerequisites
    print_status "Checking prerequisites..."
    
    if ! command_exists go; then
        print_error "Go is not installed"
        exit 1
    fi
    
    if ! command_exists curl; then
        print_warning "curl is not installed, some tests may be limited"
    fi
    
    print_success "Basic prerequisites met"
    echo
    
    # Initialize counters
    total_tests=0
    passed_tests=0
    
    # Test 1: Compilation
    ((total_tests++))
    if run_test "Compilation" "scripts/test-compilation.sh"; then
        ((passed_tests++))
    fi
    echo
    
    # Test 2: Database (configuration only)
    ((total_tests++))
    if run_test "Database" "scripts/test-database.sh"; then
        ((passed_tests++))
    fi
    echo
    
    # Test 3: WebSocket (compilation only)
    ((total_tests++))
    if run_test "WebSocket" "scripts/test-websocket.sh"; then
        ((passed_tests++))
    fi
    echo
    
    # Generate final report
    generate_report $passed_tests $total_tests
    
    # Cleanup
    cleanup
    
    # Exit with appropriate code
    if [ $passed_tests -eq $total_tests ]; then
        exit 0
    else
        exit 1
    fi
}

# Handle script interruption
trap cleanup EXIT

# Run main function
main "$@"