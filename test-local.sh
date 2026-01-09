#!/bin/bash
# Local acceptance testing script for TrueNAS Terraform Provider
# This script sets up environment variables and runs acceptance tests

set -e

# Check if TrueNAS credentials are provided
if [ -z "$TRUENAS_HOST" ]; then
    echo "Error: TRUENAS_HOST environment variable is not set"
    echo "Usage: TRUENAS_HOST=192.168.1.100 TRUENAS_TOKEN=your-token ./test-local.sh"
    exit 1
fi

if [ -z "$TRUENAS_TOKEN" ]; then
    echo "Error: TRUENAS_TOKEN environment variable is not set"
    echo "Usage: TRUENAS_HOST=192.168.1.100 TRUENAS_TOKEN=your-token ./test-local.sh"
    exit 1
fi

# Enable acceptance tests
export TF_ACC=1

echo "Running acceptance tests against TrueNAS at $TRUENAS_HOST"
echo "=================================================="

# Run acceptance tests one at a time with delays to avoid rate limiting
for test in TestAccVmResource_basic TestAccVmResource_update TestAccVmResource_startOnCreate; do
    echo ""
    echo "Running $test..."
    go test ./internal/provider -v -run "^${test}$" -timeout 5m
    if [ $? -eq 0 ]; then
        echo "✓ $test passed"
    else
        echo "✗ $test failed"
    fi
    echo "Waiting 5 seconds to avoid rate limiting..."
    sleep 5
done

echo ""
echo "=================================================="
echo "Acceptance tests completed"
