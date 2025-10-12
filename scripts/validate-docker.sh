#!/bin/bash

# =============================================================================
# Docker Configuration Validation Script
# =============================================================================
# This script validates all Docker-related configuration files
# =============================================================================

echo "🔍 Validating Docker configuration files..."
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print success
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Function to print error
print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to print info
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if required files exist
check_file_exists() {
    local file="$1"
    local description="$2"
    
    if [ -f "$file" ]; then
        print_success "$description exists"
        return 0
    else
        print_error "$description not found"
        return 1
    fi
}

# Validate YAML syntax
validate_yaml() {
    local file="$1"
    local description="$2"
    
    if command -v python3 &> /dev/null; then
        if python3 -c "import yaml; yaml.safe_load(open('$file'))" 2>/dev/null; then
            print_success "$description syntax is valid"
            return 0
        else
            print_error "$description syntax error"
            return 1
        fi
    else
        print_warning "Python3 not found, skipping YAML validation"
        return 0
    fi
}

# Count lines in file
count_lines() {
    local file="$1"
    if [ -f "$file" ]; then
        wc -l < "$file" | tr -d ' '
    fi
}

# Validate Dockerfile syntax
validate_dockerfile() {
    local file="$1"
    local description="$2"
    
    # Basic Dockerfile validation
    if grep -q "FROM" "$file" && grep -q "CMD\|ENTRYPOINT" "$file"; then
        print_success "$description structure is valid"
        return 0
    else
        print_error "$description structure is invalid (missing FROM or CMD/ENTRYPOINT)"
        return 1
    fi
}

# Start validation
echo "📋 Checking required files..."
echo ""

# Check main configuration files
check_file_exists "docker-compose.yml" "Docker Compose configuration"
check_file_exists ".env.example" "Environment variables template"
check_file_exists "Makefile" "Makefile"

echo ""

# Check Dockerfiles
echo "📦 Checking Dockerfiles..."
dockerfiles=(
    "backend/Dockerfile:Backend API Dockerfile"
    "backend/Dockerfile.websocket:WebSocket Server Dockerfile"
    "janus-server/Dockerfile:Janus WebRTC Server Dockerfile"
)

for dockerfile_info in "${dockerfiles[@]}"; do
    IFS=':' read -r file description <<< "$dockerfile_info"
    if check_file_exists "$file" "$description"; then
        validate_dockerfile "$file" "$description"
        lines=$(count_lines "$file")
        print_info "   - Lines: $lines"
    fi
done

echo ""

# Check Janus configuration files
echo "⚙️  Checking Janus configuration files..."
janus_configs=(
    "janus-server/config/janus.jcfg:Janus main configuration"
    "janus-server/config/janus.plugin.videoroom.jcfg:Janus VideoRoom plugin configuration"
)

for config_info in "${janus_configs[@]}"; do
    IFS=':' read -r file description <<< "$config_info"
    if check_file_exists "$file" "$description"; then
        lines=$(count_lines "$file")
        print_info "   - Lines: $lines"
    fi
done

echo ""

# Validate docker-compose.yml
echo "🔧 Validating Docker Compose configuration..."
if [ -f "docker-compose.yml" ]; then
    validate_yaml "docker-compose.yml" "Docker Compose"
    
    # Extract and display service information
    if command -v python3 &> /dev/null; then
        echo ""
        python3 -c "
import yaml
import sys

try:
    with open('docker-compose.yml', 'r') as f:
        config = yaml.safe_load(f)
    
    if 'services' in config:
        services = list(config['services'].keys())
        print(f'📋 Services ({len(services)}):')
        for service in services:
            print(f'   - {service}')
    
    if 'networks' in config:
        networks = list(config['networks'].keys())
        print(f'🌐 Networks ({len(networks)}):')
        for network in networks:
            print(f'   - {network}')
    
    if 'volumes' in config:
        volumes = list(config['volumes'].keys())
        print(f'💾 Volumes ({len(volumes)}):')
        for volume in volumes:
            print(f'   - {volume}')
            
except Exception as e:
    print(f'Error parsing docker-compose.yml: {e}')
    sys.exit(1)
"
    fi
fi

echo ""

# Check environment variables
echo "🌍 Checking environment variables..."
if [ -f ".env.example" ]; then
    env_count=$(grep -v '^#' .env.example | grep -v '^$' | wc -l | tr -d ' ')
    print_info "Environment variables defined: $env_count"
    
    # Check for critical variables
    critical_vars=(
        "POSTGRES_DB"
        "POSTGRES_USER"
        "POSTGRES_PASSWORD"
        "REDIS_PASSWORD"
        "JWT_SECRET"
        "JANUS_API_SECRET"
        "JANUS_ADMIN_SECRET"
    )
    
    echo ""
    print_info "Critical environment variables:"
    for var in "${critical_vars[@]}"; do
        if grep -q "^$var=" .env.example; then
            print_success "   $var is defined"
        else
            print_warning "   $var is not defined"
        fi
    done
fi

echo ""

# Check Makefile targets
echo "🎯 Checking Makefile targets..."
if [ -f "Makefile" ]; then
    target_count=$(grep -E '^[a-zA-Z_-]+:' Makefile | grep -v '^#' | wc -l | tr -d ' ')
    print_info "Make targets available: $target_count"
    
    # Check for essential targets
    essential_targets=(
        "help"
        "dev"
        "docker-up"
        "docker-down"
        "build"
        "test"
        "clean"
    )
    
    echo ""
    print_info "Essential Make targets:"
    for target in "${essential_targets[@]}"; do
        if grep -q "^$target:" Makefile; then
            print_success "   $target"
        else
            print_warning "   $target not found"
        fi
    done
fi

echo ""

# Check directory structure
echo "📁 Checking directory structure..."
directories=(
    "backend:Backend source code"
    "frontend:Frontend source code"
    "janus-server:Janus WebRTC server"
    "janus-server/config:Janus configuration"
    "scripts:Utility scripts"
)

for dir_info in "${directories[@]}"; do
    IFS=':' read -r dir description <<< "$dir_info"
    if [ -d "$dir" ]; then
        print_success "$description directory exists"
    else
        print_warning "$description directory not found"
    fi
done

echo ""

# Summary
echo "📊 Validation Summary"
echo "===================="

# Count total files checked
total_files=0
valid_files=0

# Count configuration files
config_files=("docker-compose.yml" ".env.example" "Makefile")
for file in "${config_files[@]}"; do
    ((total_files++))
    if [ -f "$file" ]; then
        ((valid_files++))
    fi
done

# Count Dockerfiles
for dockerfile_info in "${dockerfiles[@]}"; do
    IFS=':' read -r file _ <<< "$dockerfile_info"
    ((total_files++))
    if [ -f "$file" ]; then
        ((valid_files++))
    fi
done

# Count Janus configs
for config_info in "${janus_configs[@]}"; do
    IFS=':' read -r file _ <<< "$config_info"
    ((total_files++))
    if [ -f "$file" ]; then
        ((valid_files++))
    fi
done

if [ $valid_files -eq $total_files ]; then
    print_success "All $total_files configuration files are present!"
    echo ""
    print_info "🚀 Your Docker configuration is ready to use!"
    echo ""
    print_info "Next steps:"
    echo "   1. Copy .env.example to .env and update the values"
    echo "   2. Run 'make dev' to start the development environment"
    echo "   3. Access the services at their respective ports"
else
    print_error "$((total_files - valid_files)) of $total_files files are missing!"
    echo ""
    print_warning "Please fix the missing files before proceeding."
fi

echo ""
print_info "For more information, see README-DOCKER.md"
echo ""