#!/bin/bash

# Build RPM script for kidsh
set -e

# Configuration
PACKAGE_NAME="kidsh"
VERSION="1.0.0"
RELEASE="1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building RPM package for ${PACKAGE_NAME}${NC}"

# Check if we're in the right directory
if [ ! -f "kidsh.spec" ]; then
    echo -e "${RED}Error: kidsh.spec not found. Please run this script from the project root.${NC}"
    exit 1
fi

# Create build directory
BUILD_DIR="rpmbuild"
mkdir -p ${BUILD_DIR}/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

# Create source tarball
echo -e "${YELLOW}Creating source tarball...${NC}"
tar --exclude='.git' --exclude='rpmbuild' --exclude='*.rpm' \
    --exclude='kidsh' --exclude='go.sum' \
    -czf ${BUILD_DIR}/SOURCES/${PACKAGE_NAME}-${VERSION}.tar.gz .

# Copy spec file
cp kidsh.spec ${BUILD_DIR}/SPECS/

# Build RPM
echo -e "${YELLOW}Building RPM...${NC}"
rpmbuild --define "_topdir $(pwd)/${BUILD_DIR}" -bb ${BUILD_DIR}/SPECS/kidsh.spec

# Find and display the built RPM
RPM_FILE=$(find ${BUILD_DIR}/RPMS -name "*.rpm" | head -1)
if [ -n "$RPM_FILE" ]; then
    echo -e "${GREEN}RPM built successfully:${NC}"
    echo -e "${GREEN}  $(realpath $RPM_FILE)${NC}"
    
    # Show package info
    echo -e "${YELLOW}Package information:${NC}"
    rpm -qip "$RPM_FILE"
else
    echo -e "${RED}Error: RPM file not found${NC}"
    exit 1
fi

echo -e "${GREEN}RPM build completed successfully!${NC}" 