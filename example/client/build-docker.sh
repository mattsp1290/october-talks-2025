#!/bin/bash

# Default values
DEFAULT_TAG="ag-ui-protocol/ag-ui-client:latest"
SDK_PATH="/Users/punk1290/git/ag-ui/go-sdk"

# Parse command line arguments
TAG="${1:-$DEFAULT_TAG}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building Docker image with tag: ${TAG}${NC}"

# Check if SDK exists
if [ ! -d "$SDK_PATH" ]; then
    echo -e "${RED}Error: SDK directory not found at $SDK_PATH${NC}"
    exit 1
fi

# Create temp directory for SDK
TEMP_SDK_DIR="./ag-ui-go-sdk"
echo -e "${YELLOW}Creating temporary SDK copy at ${TEMP_SDK_DIR}${NC}"

# Clean up function
cleanup() {
    if [ -d "$TEMP_SDK_DIR" ]; then
        echo -e "${YELLOW}Cleaning up temporary SDK directory${NC}"
        rm -rf "$TEMP_SDK_DIR"
    fi
}

# Set trap to ensure cleanup happens on exit
trap cleanup EXIT

# Copy SDK to build context
cp -r "$SDK_PATH" "$TEMP_SDK_DIR"
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to copy SDK${NC}"
    exit 1
fi

echo -e "${GREEN}SDK copied successfully${NC}"

# Build Docker image
echo -e "${YELLOW}Starting Docker build...${NC}"
docker build -t "$TAG" .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Docker build successful!${NC}"
    echo -e "${GREEN}Image tagged as: ${TAG}${NC}"
else
    echo -e "${RED}Docker build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}Build complete!${NC}"