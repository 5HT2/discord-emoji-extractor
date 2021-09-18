#!/bin/bash

download_emoji () {
    FILE="emojis/$2.$3"
    if [[ ! -s "$FILE" ]]; then
        echo "Downloading $FILE ($1)"
        curl -# -o "$FILE" "https://cdn.discordapp.com/emojis/$1.$3"
    fi
}

parse_emoji () {
    for emoji in $(echo "$1" | grep -oEi '<(a|):[A-z0-9_]+:[0-9]+>'); do
        EMOJI="${emoji:1: -1}"
        EMOJI_ID="${EMOJI##*:}"
        EMOJI_NAME="$(cut -d ":" -f2 <<< "$EMOJI")"
        case "$EMOJI" in
            a:*) download_emoji "$EMOJI_ID" "$EMOJI_NAME" "gif";;
            :*) download_emoji "$EMOJI_ID" "$EMOJI_NAME" "png";;
        esac
    done
}

parse_channel_messages () {
    while IFS="" read -r message || [ -n "$message" ]; do
        parse_emoji "$message"
    done < "$1"/messages.csv
}

# For each directory, parse it's messages
for dir in */; do
    if [[ -f "$dir/messages.csv" ]]; then
        parse_channel_messages "$dir"
    fi
done
