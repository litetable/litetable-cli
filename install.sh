#!/bin/bash
set -e

# Configuration
INSTALL_DIR="${HOME}/.litetable"
LITETABLE_BIN_PATH="${INSTALL_DIR}/bin"
BIN_DIR="/usr/local/bin"
GITHUB_REPO="litetable/litetable-cli"
CLI_NAME="litetable"
VERSION=${1:-"latest"}  # Use provided version or "latest" if not specified
REQUIRED_PATH_LINE="export PATH=\"${LITETABLE_BIN_PATH}:\$PATH\""


# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}Installing LiteTable CLI${NC} ${VERSION}..."

# Create installation directory
mkdir -p "${INSTALL_DIR}/bin"
echo -e "${GREEN}✓${NC} Created installation directory: ${INSTALL_DIR}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
  x86_64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    echo -e "${YELLOW}Warning: Architecture $ARCH might not be supported${NC}"
    ;;
esac

echo -e "${GREEN}✓${NC} Detected OS: ${OS}, Architecture: ${ARCH}"

# Get release info based on version
echo "Fetching release information..."
if [ "$VERSION" = "latest" ]; then
  RELEASE_URL=$(curl -s https://api.github.com/repos/${GITHUB_REPO}/releases/latest | grep "browser_download_url.*${OS}_${ARCH}" | cut -d '"' -f 4)
  VERSION_TAG=$(curl -s https://api.github.com/repos/${GITHUB_REPO}/releases/latest | grep "tag_name" | cut -d '"' -f 4)
else
  RELEASE_URL=$(curl -s https://api.github.com/repos/${GITHUB_REPO}/releases/tags/${VERSION} | grep "browser_download_url.*${OS}_${ARCH}" | cut -d '"' -f 4)
  VERSION_TAG=$VERSION
fi

if [ -z "$RELEASE_URL" ]; then
  echo -e "${YELLOW}No release found for ${OS}_${ARCH} with version ${VERSION}. Falling back to building from source.${NC}"

  # Check if go is installed
  if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go to continue."
    exit 1
  fi

  # Create a temporary directory
  TMP_DIR=$(mktemp -d)
  cd "$TMP_DIR"

  # Clone the repository
  if [ "$VERSION" = "latest" ]; then
    git clone "https://github.com/${GITHUB_REPO}.git" .
  else
    git clone -b "${VERSION}" "https://github.com/${GITHUB_REPO}.git" . || {
      echo "Error: Version ${VERSION} not found."
      exit 1
    }
  fi

  # Build the binary
  echo "Building from source..."
  go build -o "${INSTALL_DIR}/bin/${CLI_NAME}"

  # Clean up
  cd - > /dev/null
  rm -rf "$TMP_DIR"
else
  # Download the binary
  echo "Downloading ${VERSION_TAG} from ${RELEASE_URL}..."
  curl -L "${RELEASE_URL}" -o "${INSTALL_DIR}/bin/${CLI_NAME}"
fi

# Make binary executable
chmod +x "${INSTALL_DIR}/bin/${CLI_NAME}"
echo -e "${GREEN}✓${NC} Downloaded and made executable"

# Determine appropriate shell config file
SHELL_FILE=""
if [ -f "${HOME}/.zshrc" ]; then
  SHELL_FILE="${HOME}/.zshrc"
elif [ -f "${HOME}/.bashrc" ]; then
  SHELL_FILE="${HOME}/.bashrc"
else
  SHELL_FILE="${HOME}/.profile"
fi

# Only append the line if it's not already present
if ! grep -Fxq "$REQUIRED_PATH_LINE" "$SHELL_FILE"; then
  echo "" >> "$SHELL_FILE"
  echo "# Added by LiteTable CLI installer" >> "$SHELL_FILE"
  echo "$REQUIRED_PATH_LINE" >> "$SHELL_FILE"
  echo "✓ Added LiteTable bin to PATH in $SHELL_FILE"
  echo "→ Run 'source $SHELL_FILE' or restart your terminal to apply the change."
else
  echo "✓ PATH export already present in $SHELL_FILE"
fi

echo -e "${GREEN}LiteTable CLI${NC} ${VERSION_TAG} ${GREEN}installed successfully!${NC}"
echo "Run 'litetable --help' to get started"
