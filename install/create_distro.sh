#!/bin/bash

set -e  

APP_NAME="Commander"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$SCRIPT_DIR/../build"
BUNDLE_DIR="$BUILD_DIR/$APP_NAME.app"
PKG_DIR="$BUNDLE_DIR/Contents"
RESOURCES_DIR="$PKG_DIR/Resources"
MACOS_DIR="$RESOURCES_DIR/MacOS"
INSTALL_PATH="/Applications"
LAUNCH_AGENT_PLIST="$BUNDLE_DIR/com.askchriswatson.$APP_NAME.plist"
ICONSET_DIR="$BUILD_DIR/$APP_NAME.iconset"
ICNS_FILE="$RESOURCES_DIR/$APP_NAME.icns"

mkdir -p "$BUILD_DIR"
rm -rf "$BUNDLE_DIR"
mkdir -p "$PKG_DIR"
mkdir -p "$RESOURCES_DIR"
mkdir -p "$MACOS_DIR"

# build app
GOOS=darwin GOARCH=amd64 go build -o "$BUILD_DIR/$APP_NAME" ./cmd/main.go

# check if build was successful
if [ ! -f "$BUILD_DIR/$APP_NAME" ]; then
    echo "Error: Failed to build the application"
    exit 1
fi

# copy app to MACOS_DIR for packaging
cp "$BUILD_DIR/$APP_NAME" "$MACOS_DIR/"

# create Info.plist
cat <<EOF > "$PKG_DIR/Info.plist"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>$APP_NAME</string>
    <key>CFBundleIconFile</key>
    <string>$APP_NAME.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.askchriswatson.$APP_NAME</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundleVersion</key>
    <string>1.0</string>
</dict>
</plist>
EOF


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


# add script for installing the plist to the pkg dir
cat <<EOF > "$BUNDLE_DIR/install_launch_agent_plist.sh"
#!/bin/bash

APP_NAME="$APP_NAME"
PLIST_NAME="com.askchriswatson.\$APP_NAME.plist"
PLIST_SOURCE_PATH="/Volumes/\$APP_NAME Installer/\$PLIST_NAME"
PLIST_DEST_PATH="\$HOME/Library/LaunchAgents/\$PLIST_NAME"

# Copy the plist file to the LaunchAgents directory
cp "\$PLIST_SOURCE_PATH" "\$PLIST_DEST_PATH"

# Load the plist file using launchctl
launchctl load "\$PLIST_DEST_PATH"

echo "Launch agent installed and loaded successfully."
EOF

# create icon set
mkdir -p "$ICONSET_DIR"
cp $SCRIPT_DIR/icon-64.png "$ICONSET_DIR/icon_64x64.png"
cp $SCRIPT_DIR/icon-128.png "$ICONSET_DIR/icon_128x128.png"
cp $SCRIPT_DIR/icon-256.png "$ICONSET_DIR/icon_256x256.png"
cp $SCRIPT_DIR/icon-512.png "$ICONSET_DIR/icon_512x512.png"

iconutil -c icns -o "$ICNS_FILE" "$ICONSET_DIR"

# Clean up the iconset directory
rm -rf "$ICONSET_DIR"

# set executable 
chmod +x "$BUNDLE_DIR/install_launch_agent_plist.sh"

# create README 
cat <<EOF > "$BUNDLE_DIR/README.txt"
Installation Instructions
=======================

To install the launch agent that will automatically start $APP_NAME:

1. Open Terminal
2. Navigate to the mounted installer volume:
   cd "/Volumes/$APP_NAME Installer"

3. Run the installation script:
   ./install_launch_agent_plist.sh

The script will:
- Copy the launch agent plist to your user's LaunchAgents directory
- Load the launch agent so $APP_NAME starts automatically

To verify installation:
- Check that the plist exists in ~/Library/LaunchAgents/
- $APP_NAME should start automatically on your next login

To uninstall:
1. Unload the launch agent:
   launchctl unload ~/Library/LaunchAgents/com.askchriswatson.$APP_NAME.plist
   
2. Remove the plist file:
   rm ~/Library/LaunchAgents/com.askchriswatson.$APP_NAME.plist

EOF


# Remove existing DMG if it exists
rm -f "$BUILD_DIR/$APP_NAME.dmg"

TMP_DMG_DIR=$(mktemp -d)

cp -R "$BUNDLE_DIR" "$TMP_DMG_DIR/"

# make sure our addl files make it into the DMG
cp "$BUNDLE_DIR/install_launch_agent_plist.sh" "$TMP_DMG_DIR/"
cp "$BUNDLE_DIR/README.txt" "$TMP_DMG_DIR/"
cp "$LAUNCH_AGENT_PLIST" "$TMP_DMG_DIR/"


# create the DMG
create-dmg \
    --volname "$APP_NAME Installer" \
    --window-pos 200 120 \
    --window-size 800 600 \
    --icon-size 100 \
    --icon "$APP_NAME.app" 400 150  \
    --icon "install_launch_agent_plist.sh" 200 300 \
    --icon "README.txt"  200 150 \
    --icon "com.askchriswatson.$APP_NAME.plist" 400 300 \
    --app-drop-link 600 150 \
    --no-internet-enable \
    "$BUILD_DIR/$APP_NAME.dmg" \
    "$TMP_DMG_DIR"

# clean up tmp dir
rm -rf "$TMP_DMG_DIR"

# sign the app bundle
codesign --force --deep --sign - \
    "$BUNDLE_DIR"

# verify the signature
codesign --verify --deep --strict "$BUNDLE_DIR"
codesign --force -vv --deep "$BUNDLE_DIR" # Sign the app bundle
codesign --force --deep --sign - "$BUNDLE_DIR"

# a real app would have some notarization steps here

# instructions for users to manually install the plist
echo "Please copy the plist file to ~/Library/LaunchAgents/ and load it using launchctl."