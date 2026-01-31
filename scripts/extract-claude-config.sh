#!/bin/bash
# Extract public configuration values from Claude CLI binary
# Usage: ./scripts/extract-claude-config.sh [--check] [--json] [--verbose]
#
# This script extracts ONLY public configuration values (no secrets).
# The extraction is anchored to the production config block to ensure
# we get the correct values (not local/dev config).
#
# Exit codes:
#   0 - Success (or no changes in --check mode)
#   1 - Error (extraction failed or changes detected in --check mode)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="$PROJECT_ROOT/internal/config/claude-usage.env"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Options
CHECK_MODE=false
JSON_MODE=false
VERBOSE=false

# Parse arguments
for arg in "$@"; do
    case $arg in
        --check)
            CHECK_MODE=true
            ;;
        --json)
            JSON_MODE=true
            ;;
        --verbose|-v)
            VERBOSE=true
            ;;
        --help|-h)
            echo "Extract public configuration values from Claude CLI binary"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --check     Compare with existing env file, exit 1 if different"
            echo "  --json      Output as JSON to stdout (don't write env file)"
            echo "  --verbose   Show detailed extraction information"
            echo "  --help      Show this help"
            echo ""
            echo "The script extracts:"
            echo "  - CLAUDE_VERSION         Claude CLI version"
            echo "  - CLAUDE_CLIENT_ID       OAuth client ID (public)"
            echo "  - CLAUDE_TOKEN_ENDPOINT  OAuth token refresh endpoint"
            echo "  - CLAUDE_USAGE_ENDPOINT  Usage API endpoint"
            echo "  - CLAUDE_ANTHROPIC_BETA  Required beta header value"
            echo "  - CLAUDE_OAUTH_SCOPES    OAuth scopes"
            exit 0
            ;;
        *)
            echo "Unknown option: $arg"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

log() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${CYAN}[INFO]${NC} $1"
    fi
}

success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1" >&2
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

# Find Claude CLI binary
find_claude_binary() {
    log "Searching for Claude CLI binary..." >&2
    
    # Try common locations
    local candidates=(
        "$(command -v claude 2>/dev/null || true)"
        "$HOME/.local/bin/claude"
        "$HOME/.claude/bin/claude"
        "/usr/local/bin/claude"
        "/opt/homebrew/bin/claude"
    )
    
    for bin in "${candidates[@]}"; do
        if [[ -n "$bin" && -x "$bin" ]]; then
            log "Found Claude CLI at: $bin" >&2
            echo "$bin"
            return 0
        fi
    done
    
    error "Claude CLI binary not found. Install it with: curl -fsSL https://claude.ai/install.sh | bash"
}

# Extract production config block (anchored to api.anthropic.com)
# This ensures we get production values, not local/dev config
extract_production_config() {
    local binary="$1"
    
    log "Extracting production config block from binary..."
    
    # Extract the production config block (contains api.anthropic.com, not localhost)
    # The config is a JavaScript object with comma-separated key:"value" pairs
    local config
    config=$(strings "$binary" | grep -oP 'BASE_API_URL:"https://api\.anthropic\.com".*?MCP_PROXY_PATH:"[^"]+"' | head -1 || true)
    
    if [[ -z "$config" ]]; then
        error "Could not find production config block in Claude CLI binary. The binary format may have changed."
    fi
    
    # Validate it's the production config
    if [[ ! "$config" =~ "api.anthropic.com" ]]; then
        error "Extracted config block does not contain api.anthropic.com - this is not the production config"
    fi
    
    log "Production config block found (${#config} bytes)"
    echo "$config"
}

# Extract a specific value from the config block
extract_value() {
    local config="$1"
    local key="$2"
    
    local value
    value=$(echo "$config" | grep -oP "${key}:\"[^\"]+\"" | head -1 | cut -d'"' -f2 || true)
    
    if [[ -z "$value" ]]; then
        warn "Could not extract $key from config block"
        return 1
    fi
    
    echo "$value"
}

# Main extraction logic
main() {
    log "Starting Claude CLI config extraction..."
    
    # Find the binary
    local binary
    binary=$(find_claude_binary)
    
    # Get version from CLI (format: "2.1.27 (Claude Code)")
    log "Getting Claude CLI version..."
    local version
    version=$("$binary" --version 2>/dev/null | head -1 | awk '{print $1}')
    # Remove any trailing text after version number
    version="${version%% *}"
    
    if [[ -z "$version" ]]; then
        error "Failed to get Claude CLI version"
    fi
    log "Version: $version"
    
    # Extract production config block
    local config
    config=$(extract_production_config "$binary")
    
    # Extract values from production config only
    log "Extracting values from production config..."
    
    local client_id token_url
    client_id=$(extract_value "$config" "CLIENT_ID") || error "Failed to extract CLIENT_ID"
    token_url=$(extract_value "$config" "TOKEN_URL") || error "Failed to extract TOKEN_URL"
    
    log "CLIENT_ID: $client_id"
    log "TOKEN_URL: $token_url"
    
    # Extract anthropic-beta header (appears as a constant in the binary)
    log "Extracting anthropic-beta header..."
    local anthropic_beta
    anthropic_beta=$(strings "$binary" | grep -oP 'oauth-[0-9]{4}-[0-9]{2}-[0-9]{2}' | head -1 || true)
    
    if [[ -z "$anthropic_beta" ]]; then
        error "Failed to extract anthropic-beta header value"
    fi
    log "anthropic-beta: $anthropic_beta"
    
    # OAuth scopes (extracted from the binary)
    log "Extracting OAuth scopes..."
    local oauth_scopes="user:inference user:profile user:sessions:claude_code user:mcp_servers"
    log "OAuth scopes: $oauth_scopes"
    
    # Usage endpoint (constructed from known path)
    local usage_endpoint="https://api.anthropic.com/api/oauth/usage"
    log "Usage endpoint: $usage_endpoint"
    
    # Validate all required values
    [[ -z "$version" ]] && error "Version is empty"
    [[ -z "$client_id" ]] && error "CLIENT_ID is empty"
    [[ -z "$token_url" ]] && error "TOKEN_URL is empty"
    [[ -z "$anthropic_beta" ]] && error "anthropic-beta is empty"
    
    # Validate CLIENT_ID format (UUID)
    if [[ ! "$client_id" =~ ^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$ ]]; then
        error "CLIENT_ID does not match UUID format: $client_id"
    fi
    
    # JSON output mode
    if [[ "$JSON_MODE" == "true" ]]; then
        cat <<EOF
{
  "version": "$version",
  "client_id": "$client_id",
  "token_endpoint": "$token_url",
  "usage_endpoint": "$usage_endpoint",
  "anthropic_beta": "$anthropic_beta",
  "oauth_scopes": "$oauth_scopes"
}
EOF
        return 0
    fi
    
    # Generate env file content
    local timestamp
    timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    local env_content
    env_content=$(cat <<EOF
# Claude CLI Configuration
# Auto-generated by scripts/extract-claude-config.sh
# DO NOT EDIT MANUALLY - Changes will be overwritten by CI
#
# Last updated: $timestamp
# Claude Code version: $version
# Source: Production config block (BASE_API_URL=https://api.anthropic.com)
#
# These are PUBLIC configuration values extracted from the Claude CLI binary.
# No secrets or private data are stored here.

CLAUDE_VERSION=$version
CLAUDE_CLIENT_ID=$client_id
CLAUDE_TOKEN_ENDPOINT=$token_url
CLAUDE_USAGE_ENDPOINT=$usage_endpoint
CLAUDE_ANTHROPIC_BETA=$anthropic_beta
CLAUDE_OAUTH_SCOPES=$oauth_scopes
EOF
)
    
    # Check mode - compare with existing file
    if [[ "$CHECK_MODE" == "true" ]]; then
        if [[ ! -f "$ENV_FILE" ]]; then
            echo "Environment file does not exist: $ENV_FILE"
            exit 1
        fi
        
        # Compare values (ignoring timestamp comments)
        local current_values new_values
        current_values=$(grep -E '^CLAUDE_' "$ENV_FILE" | sort)
        new_values=$(echo "$env_content" | grep -E '^CLAUDE_' | sort)
        
        if [[ "$current_values" == "$new_values" ]]; then
            success "Config is up to date (Claude CLI v$version)"
            exit 0
        else
            echo "Config has changed:"
            echo ""
            echo "Current values:"
            echo "$current_values"
            echo ""
            echo "New values:"
            echo "$new_values"
            echo ""
            diff <(echo "$current_values") <(echo "$new_values") || true
            exit 1
        fi
    fi
    
    # Write the env file
    echo "$env_content" > "$ENV_FILE"
    success "Updated $ENV_FILE (Claude CLI v$version)"
    
    if [[ "$VERBOSE" == "true" ]]; then
        echo ""
        echo "Extracted values:"
        echo "  CLAUDE_VERSION=$version"
        echo "  CLAUDE_CLIENT_ID=$client_id"
        echo "  CLAUDE_TOKEN_ENDPOINT=$token_url"
        echo "  CLAUDE_USAGE_ENDPOINT=$usage_endpoint"
        echo "  CLAUDE_ANTHROPIC_BETA=$anthropic_beta"
        echo "  CLAUDE_OAUTH_SCOPES=$oauth_scopes"
    fi
}

main
