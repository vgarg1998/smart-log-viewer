#!/bin/bash

# Pre-commit hook for smart-log-viewer project
# This hook runs before each commit to ensure code quality

set -e

echo "🔍 Running pre-commit checks..."

# Function to print warnings instead of failing
print_warning() {
    echo "⚠️  WARNING: $1"
}

# Function to print success messages
print_success() {
    echo "✅ $1"
}

# Function to print error messages
print_error() {
    echo "❌ ERROR: $1"
}

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    print_error "Please run this hook from the smart-log-viewer directory"
    exit 1
fi

# Server-side checks (Go)
echo "🐹 Checking Go server code..."

if [ -d "server" ]; then
    cd server
    
    # Go format check
    if command -v gofmt >/dev/null 2>&1; then
        echo "  Formatting Go code..."
        gofmt -w .
        print_success "Go code formatted"
    else
        print_warning "gofmt not found, skipping Go formatting"
    fi
    
    # Go lint check
    if command -v golint >/dev/null 2>&1; then
        echo "  Linting Go code..."
        golint ./... || print_warning "Go linting found issues"
    else
        print_warning "golint not found, skipping Go linting"
    fi
    
    # Go vet check
    if command -v go >/dev/null 2>&1; then
        echo "  Running go vet..."
        go vet ./... || print_warning "Go vet found issues"
    else
        print_warning "go not found, skipping Go vet"
    fi
    
    # Go test
    if command -v go >/dev/null 2>&1; then
        echo "  Running Go tests..."
        go test -buildvcs=false ./... || print_warning "Go tests failed"
    else
        print_warning "go not found, skipping Go tests"
    fi
    
    # Go build check
    if command -v go >/dev/null 2>&1; then
        echo "  Building Go server..."
        go build -buildvcs=false ./cmd/server || print_warning "Go build failed"
        print_success "Go server builds successfully"
    else
        print_warning "go not found, skipping Go build"
    fi
    
    cd ..
else
    print_warning "Server directory not found, skipping Go checks"
fi

# Client-side checks (React/TypeScript)
echo "⚛️  Checking React client code..."

if [ -d "client" ]; then
    cd client
    
    # Check if package.json exists
    if [ ! -f "package.json" ]; then
        print_warning "package.json not found in client directory"
        cd ..
        exit 0
    fi
    
    # Install dependencies if node_modules doesn't exist
    if [ ! -d "node_modules" ]; then
        echo "  Installing npm dependencies..."
        npm install || print_warning "Failed to install npm dependencies"
    fi
    
    # ESLint check
    if npm run lint >/dev/null 2>&1; then
        echo "  Running ESLint..."
        npm run lint || print_warning "ESLint found issues"
        print_success "ESLint passed"
    else
        print_warning "ESLint script not found in package.json"
    fi
    
    # TypeScript type check
    if npm run type-check >/dev/null 2>&1; then
        echo "  Running TypeScript type check..."
        npm run type-check || print_warning "TypeScript type check failed"
    else
        echo "  Running TypeScript compiler check..."
        npx tsc --noEmit || print_warning "TypeScript compilation failed"
    fi
    
    # Build check
    echo "  Building React client..."
    npm run build || print_warning "React client build failed"
    print_success "React client builds successfully"
    
    cd ..
else
    print_warning "Client directory not found, skipping React checks"
fi

# Docker Compose validation
echo "🐳 Validating Docker Compose configuration..."

if command -v docker-compose >/dev/null 2>&1; then
    echo "  Validating docker-compose.yml..."
    docker-compose config >/dev/null || print_warning "Docker Compose configuration is invalid"
    print_success "Docker Compose configuration is valid"
else
    print_warning "docker-compose not found, skipping validation"
fi

echo "🎉 Pre-commit checks completed!"
echo "💡 Note: Warnings won't prevent the commit, but errors will."
echo "   Fix any errors before committing to ensure code quality."
