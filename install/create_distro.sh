#!/bin/bash

APP_NAME="CommanderPrototypev0"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$SCRIPT_DIR/../build"
PKG_DIR="$SCRIPT_DIR/../build/dist"
INSTALL_PATH="/Applications"
LAUNCH_AGENT_PLIST="$PKG_DIR/com.askchriswatson.$APP_NAME.plist"

mkdir -p "$BUILD_DIR"
mkdir -p "$PKG_DIR"

# build app
GOOS=darwin GOARCH=amd64 go build -o "$BUILD_DIR/$APP_NAME" ./cmd/main.go

# check if build was successful
if [ ! -f "$BUILD_DIR/$APP_NAME" ]; then
    echo "Error: Failed to build the application"
    exit 1
fi

# copy app to PKG_DIR for packaging
cp "$BUILD_DIR/$APP_NAME" "$PKG_DIR/"

# create launch agent plist
cat <<EOF > "$LAUNCH_AGENT_PLIST"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.askchriswatson.$APP_NAME</string>
    <key>ProgramArguments</key>
    <array>
        <string>$INSTALL_PATH/$APP_NAME</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
EOF

# check if plist file was created
if [ ! -f "$LAUNCH_AGENT_PLIST" ]; then
    echo "Error: Failed to create plist file at $LAUNCH_AGENT_PLIST"
    exit 1
fi

# copy plist to PKG_DIR for packaging
cp "$LAUNCH_AGENT_PLIST" "$PKG_DIR/"

# create the DMG
create-dmg \
    --volname "$APP_NAME Installer" \
    --window-pos 200 120 \
    --window-size 800 400 \
    --icon-size 100 \
    --icon "$APP_NAME" 200 190 \
    --hide-extension "$APP_NAME" \
    --app-drop-link 600 185 \
    "$PKG_DIR/$APP_NAME.dmg" \
    "$PKG_DIR"

# instructions for users to manually install the plist
echo "Please copy the plist file to ~/Library/LaunchAgents/ and load it using launchctl."