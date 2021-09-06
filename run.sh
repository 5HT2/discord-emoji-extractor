#!/bin/bash

MISSING_PKGS=false
check_installed () {
    PKG_PATH="$(which "$1" 2>/dev/null)"
    if [[ -z "$PKG_PATH" ]]; then
        echo "$1 not found, please install it!"
        MISSING_PKGS=true
    fi
}

check_installed rg
check_installed sed
check_installed perl
check_installed wget 

if [[ "$MISSING_PKGS" == "true" ]]; then
    echo "Missing one or more packages, exiting"
    exit 1
fi

# Capture all messages containing emojis
echo "Finding emojis in messages"
rg "<(a|):([A-z0-9_]+):([0-9]+)>" > emojis-0.txt

# Remove file names from beginning of messages
echo "Reformatting found messages"
sed -E "s/^c[0-9]+\/messages.csv://g" emojis-0.txt > emojis-1.txt

# Remove message info (for non-dm channels) from beginning of messages
sed -E "s/^[0-9]+,[0-9]+-[0-9]+-[0-9]+ [0-9]+:[0-9]+:[0-9]+\.[0-9]+\+[0-9]+:[0-9]+,//g" emojis-1.txt > emojis-2.txt

# Split any messages containing multiple emojis in one line, into multiple lines
# The reason we do this is because when the perl script is processing
# the lines, it will only take the first matched group, for some reason
echo "Splitting found messages"
perl -pe 's/<(a|):([A-z0-9_]+):([0-9]+)>/<$1:$2:$3>\n/g' emojis-2.txt | grep -v -E "^,$" > emojis-3.txt

# Extract and reformat all unique emojis into wget commands
echo "Extracting emoji links from messages"
./extract_emojis.pl < emojis-3.txt | sort | uniq -u > emojis-4.txt

# Run wget
echo "Downloading all emojis"
mkdir emojis || {
    echo "Failed to make emojis folder"
    # Don't exit here, because we might already have an emojis folder 
}
mv emojis-4.txt emojis/emojis-4.txt || {
    echo "Failed to move file to emojis folder"
    exit 1
}
# Cleanup caches
rm emojis-*.txt
cd emojis/ || {
    echo "Failed to cd to emojis folder"
    exit 1
}
bash < emojis-4.txt
rm emojis-4.txt
