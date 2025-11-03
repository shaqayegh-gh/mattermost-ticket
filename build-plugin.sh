#!/bin/bash

set -e

echo "🔨 Building Mattermost Ticket Plugin..."

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Cleaning previous builds...${NC}"
rm -rf dist bundle *.tar.gz

mkdir -p dist

echo -e "${BLUE}Building for Linux AMD64...${NC}"
GOOS=linux GOARCH=amd64 go build -o dist/plugin-linux-amd64 ./server


echo -e "${BLUE}Creating plugin bundle...${NC}"
mkdir -p bundle/server/dist
cp plugin.json bundle/
cp dist/* bundle/server/dist/

# Create tar.gz
echo -e "${BLUE}Creating tar.gz archive...${NC}"
cd bundle
tar -czf ../mattermost-ticket-plugin.tar.gz .
cd ..

rm -rf bundle

echo -e "${GREEN}✅ Plugin built successfully!${NC}"
echo -e "${GREEN}📦 Package: mattermost-ticket-plugin.tar.gz${NC}"
echo ""
echo "To install:"
echo "1. Go to Mattermost System Console → Plugins → Plugin Management"
echo "2. Click 'Upload Plugin'"
echo "3. Select mattermost-ticket-plugin.tar.gz"
echo "4. Enable the plugin"

