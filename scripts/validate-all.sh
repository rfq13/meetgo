#!/bin/bash

# Master Validation Script untuk WebRTC Backend Server
# Author: WebRTC Meeting Team
# Version: 1.0.0

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
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

print_header() {
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}========================================${NC}"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to make scripts executable
make_scripts_executable() {
    print_status "Making test scripts executable..."
    
    local scripts=(
        "test-compilation.sh"
        "test-database.sh"
        "test-api.sh"
        "test-websocket.sh"
        "test-janus.sh"
    )
    
    for script in "${scripts[@]}"; do
        if [[ -f "scripts/$script" ]]; then
            chmod +x "scripts/$script"
            print_success "Made $script executable"
        else
            print_error "Script $script not found"
            return 1
        fi
    done
    
    return 0
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    local missing_tools=()
    
    # Check required tools
    if ! command_exists go; then
        missing_tools+=("go")
    fi
    
    if ! command_exists curl; then
        missing_tools+=("curl")
    fi
    
    if ! command_exists docker; then
        missing_tools+=("docker")
    fi
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        print_error "Missing required tools:"
        for tool in "${missing_tools[@]}"; do
            echo "  - $tool"
        done
        print_status "Please install the missing tools and try again."
        return 1
    fi
    
    # Check optional tools
    local optional_tools=()
    
    if ! command_exists psql; then
        optional_tools+=("psql - PostgreSQL client")
    fi
    
    if ! command_exists node; then
        optional_tools+=("node - Node.js runtime")
    fi
    
    if ! command_exists websocat; then
        optional_tools+=("websocat - WebSocket client")
    fi
    
    if [[ ${#optional_tools[@]} -gt 0 ]]; then
        print_warning "Optional tools not found (some tests may be limited):"
        for tool in "${optional_tools[@]}"; do
            echo "  - $tool"
        done
    fi
    
    print_success "Prerequisites check completed"
    return 0
}

# Function to run individual test
run_test() {
    local test_name="$1"
    local test_script="$2"
    local start_time=$(date +%s)
    
    print_header "Running $test_name Test"
    
    if ./scripts/"$test_script"; then
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        print_success "$test_name test passed! (${duration}s)"
        
        # Record success
        echo "$(date '+%Y-%m-%d %H:%M:%S') - $test_name: PASSED (${duration}s)" >> validation-results.log
        
        return 0
    else
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        print_error "$test_name test failed! (${duration}s)"
        
        # Record failure
        echo "$(date '+%Y-%m-%d %H:%M:%S') - $test_name: FAILED (${duration}s)" >> validation-results.log
        
        return 1
    fi
}

# Function to generate summary report
generate_summary_report() {
    local total_tests=$1
    local passed_tests=$2
    local failed_tests=$3
    local total_duration=$4
    
    local report_file="validation-summary-$(date +%Y%m%d-%H%M%S).txt"
    
    {
        echo "WebRTC Backend Validation Summary Report"
        echo "========================================"
        echo "Generated: $(date)"
        echo ""
        echo "Test Results:"
        echo "------------"
        echo "Total Tests: $total_tests"
        echo "Passed: $passed_tests"
        echo "Failed: $failed_tests"
        echo "Success Rate: $(( passed_tests * 100 / total_tests ))%"
        echo "Total Duration: ${total_duration}s"
        echo ""
        echo "Detailed Results:"
        echo "----------------"
        
        while IFS= read -r line; do
            if [[ "$line" == *"PASSED"* ]]; then
                echo "✅ $line"
            elif [[ "$line" == *"FAILED"* ]]; then
                echo "❌ $line"
            fi
        done < validation-results.log
        
        echo ""
        echo "Generated Reports:"
        echo "-----------------"
        ls -la *report-*.txt 2>/dev/null || echo "No detailed reports found"
        
        echo ""
        echo "Next Steps:"
        echo "-----------"
        if [[ $failed_tests -eq 0 ]]; then
            echo "🎉 All tests passed! The system is ready for deployment."
            echo ""
            echo "To start the servers:"
            echo "1. cd backend"
            echo "2. ./bin/api-server (Terminal 1)"
            echo "3. ./bin/websocket-server (Terminal 2)"
            echo ""
            echo "Or use Docker Compose:"
            echo "docker-compose up -d"
        else
            echo "⚠️  Some tests failed. Please review the detailed reports and fix the issues."
            echo ""
            echo "Common fixes:"
            echo "1. Check environment variables in .env file"
            echo "2. Ensure all required services are running (PostgreSQL, Janus)"
            echo "3. Verify network connectivity and port availability"
            echo "4. Review individual test reports for specific error details"
        fi
        
    } > "$report_file"
    
    print_success "Summary report generated: $report_file"
    
    # Display summary to console
    echo ""
    print_header "VALIDATION SUMMARY"
    echo -e "${BLUE}Total Tests:${NC} $total_tests"
    echo -e "${GREEN}Passed:${NC} $passed_tests"
    echo -e "${RED}Failed:${NC} $failed_tests"
    echo -e "${BLUE}Success Rate:${NC} $(( passed_tests * 100 / total_tests ))%"
    echo -e "${BLUE}Total Duration:${NC} ${total_duration}s"
    echo ""
    
    if [[ $failed_tests -eq 0 ]]; then
        echo -e "${GREEN}🎉 ALL TESTS PASSED! System is ready for deployment.${NC}"
    else
        echo -e "${RED}⚠️  SOME TESTS FAILED! Please review the reports.${NC}"
    fi
}

# Function to cleanup on exit
cleanup() {
    print_status "Cleaning up..."
    
    # Stop any running services
    if command_exists docker; then
        docker stop janus 2>/dev/null || true
        docker stop postgres-webrtc 2>/dev/null || true
    fi
    
    # Kill any background processes
    jobs -p | xargs -r kill 2>/dev/null || true
    
    print_success "Cleanup completed"
}

# Function to display help
show_help() {
    echo "WebRTC Backend Validation Script"
    echo "================================"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -c, --clean    Clean up before running tests"
    echo "  -q, --quick    Run only essential tests (compilation + database)"
    echo "  -s, --skip     Skip specific tests (comma-separated)"
    echo "  -v, --verbose  Enable verbose output"
    echo ""
    echo "Available tests:"
    echo "  compilation    - Go code compilation and dependency validation"
    echo "  database       - Database connection and migration tests"
    echo "  api            - REST API endpoint tests"
    echo "  websocket      - WebSocket server tests"
    echo "  janus          - Janus WebRTC server tests"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run all tests"
    echo "  $0 --quick            # Run only essential tests"
    echo "  $0 --skip janus       # Skip Janus tests"
    echo "  $0 --skip api,janus   # Skip API and Janus tests"
}

# Main execution function
main() {
    local clean_mode=false
    local quick_mode=false
    local skip_tests=""
    local verbose_mode=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--clean)
                clean_mode=true
                shift
                ;;
            -q|--quick)
                quick_mode=true
                shift
                ;;
            -s|--skip)
                skip_tests="$2"
                shift 2
                ;;
            -v|--verbose)
                verbose_mode=true
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Set up cleanup trap
    trap cleanup EXIT
    
    # Initialize results log
    echo "WebRTC Backend Validation Results - $(date)" > validation-results.log
    echo "========================================" >> validation-results.log
    
    print_header "WEBRTC BACKEND VALIDATION SUITE"
    echo ""
    
    local start_time=$(date +%s)
    
    # Cleanup if requested
    if [[ "$clean_mode" == true ]]; then
        print_status "Cleaning up previous runs..."
        rm -f *report-*.txt validation-results.log
        cleanup
    fi
    
    # Check prerequisites
    check_prerequisites || exit 1
    
    # Make scripts executable
    make_scripts_executable || exit 1
    
    # Define test suite
    local tests=(
        "compilation:test-compilation.sh:true"
        "database:test-database.sh:true"
        "api:test-api.sh:false"
        "websocket:test-websocket.sh:false"
        "janus:test-janus.sh:false"
    )
    
    # Filter tests based on options
    local filtered_tests=()
    for test in "${tests[@]}"; do
        local test_name=$(echo "$test" | cut -d':' -f1)
        local test_script=$(echo "$test" | cut -d':' -f2)
        local is_essential=$(echo "$test" | cut -d':' -f3)
        
        # Skip tests if requested
        if [[ -n "$skip_tests" && "$skip_tests" == *"$test_name"* ]]; then
            print_warning "Skipping $test_name test"
            continue
        fi
        
        # In quick mode, only run essential tests
        if [[ "$quick_mode" == true && "$is_essential" != "true" ]]; then
            print_warning "Skipping $test_name test (quick mode)"
            continue
        fi
        
        filtered_tests+=("$test_name:$test_script")
    done
    
    # Run tests
    local total_tests=${#filtered_tests[@]}
    local passed_tests=0
    local failed_tests=0
    
    for test in "${filtered_tests[@]}"; do
        local test_name=$(echo "$test" | cut -d':' -f1)
        local test_script=$(echo "$test" | cut -d':' -f2)
        
        if run_test "$test_name" "$test_script"; then
            ((passed_tests++))
        else
            ((failed_tests++))
            
            # In non-verbose mode, stop on first failure
            if [[ "$verbose_mode" != true && "$quick_mode" != true ]]; then
                print_error "Stopping validation due to test failure. Use --verbose to continue."
                break
            fi
        fi
        
        echo ""
    done
    
    local end_time=$(date +%s)
    local total_duration=$((end_time - start_time))
    
    # Generate summary report
    generate_summary_report $total_tests $passed_tests $failed_tests $total_duration
    
    # Exit with appropriate code
    if [[ $failed_tests -eq 0 ]]; then
        exit 0
    else
        exit 1
    fi
}

# Run main function with all arguments
main "$@"