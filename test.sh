#!/bin/bash
# test-techstack.sh - Lightweight test script for popular cloud-native repos

ORIGINAL_DIR=$(pwd)
TEST_DIR="test-repos-temp"
PASSED=0
FAILED=0
ERRORS=()

# Cleanup on exit
cleanup() {
    cd "$ORIGINAL_DIR"
    rm -rf "$TEST_DIR"
    echo ""
    echo "========================================="
    echo "TEST SUMMARY"
    echo "========================================="
    echo "✅ Passed: $PASSED"
    echo "❌ Failed: $FAILED"
    if [ $FAILED -gt 0 ]; then
        echo ""
        echo "Failed tests:"
        for error in "${ERRORS[@]}"; do
            echo "  - $error"
        done
    fi
    echo "========================================="
}
trap cleanup EXIT

# Check if techstack is available
if ! command -v techstack &> /dev/null; then
    echo "❌ ERROR: techstack command not found"
    echo ""
    echo "Please build first:"
    echo "  make build && export PATH=\$(pwd)/build:\$PATH"
    exit 1
fi

# Create test directory
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Test function with shallow clones
test_repo() {
    local name=$1
    local url=$2
    
    echo ""
    echo "========================================="
    echo "Testing: $name"
    echo "========================================="
    
    # Shallow clone (--depth 1) for minimal download
    echo "⬇️  Cloning $name..."
    if ! git clone --depth 1 --single-branch --quiet "$url" "$name" 2>/dev/null; then
        echo "❌ Failed to clone $name"
        FAILED=$((FAILED + 1))
        ERRORS+=("$name - Clone failed")
        return 1
    fi
    
    # Run detection
    if techstack get --path "$name" --quiet 2>/dev/null; then
        echo "✅ $name passed"
        PASSED=$((PASSED + 1))
    else
        echo "❌ $name failed"
        FAILED=$((FAILED + 1))
        ERRORS+=("$name - Detection failed")
    fi
    
    # Clean up immediately
    rm -rf "$name"
}

echo "🔍 Testing StackRadar"
echo "Using shallow clones (--depth 1) for lightweight testing"
echo ""

# Cloud-native Go projects
test_repo "kubernetes" "https://github.com/kubernetes/kubernetes.git" || true
test_repo "prometheus" "https://github.com/prometheus/prometheus.git" || true
test_repo "traefik" "https://github.com/traefik/traefik.git" || true
test_repo "etcd" "https://github.com/etcd-io/etcd.git" || true

# Python
test_repo "flask" "https://github.com/pallets/flask.git" || true
test_repo "fastapi" "https://github.com/tiangolo/fastapi.git" || true

# Go
test_repo "gin" "https://github.com/gin-gonic/gin.git" || true
test_repo "hugo" "https://github.com/gohugoio/hugo.git" || true

# Java
test_repo "spring-boot" "https://github.com/spring-projects/spring-boot.git" || true

# Node.js
test_repo "express" "https://github.com/expressjs/express.git" || true
test_repo "vue" "https://github.com/vuejs/core.git" || true

# .NET
test_repo "aspnetcore" "https://github.com/dotnet/aspnetcore.git" || true

# Rust
test_repo "actix-web" "https://github.com/actix/actix-web.git" || true
