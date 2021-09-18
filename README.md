# discord-emoji-extractor

Download all the emojis you've ever sent inside messages on Discord. Skips duplicates and supports resuming downloads.

## Usage

1. Download your Discord data backup. You can get this by going to Discord Settings > Privacy & safety > Request all of my data
2. Extract the data somewhere. Doesn't matter.
3. `cd my_data/messages/`
4. Download [`extract.sh`](https://github.com/l1ving/discord-emoji-extractor/blob/master/extract.sh) to the `messages` directory.
```bash
wget https://github.com/l1ving/discord-emoji-extractor/raw/master/extract.sh
```
5. Make the script executable.
```bash
chmod +x extract.sh
```
6. Run the `extract.sh` script. At any point you may cancel downloading and re-run the script from the `messages/` directory, and it will resume downloading.

## How it works

I basically just wanted a way to grab old emojis from servers that I'd left, with a picture preview, so the gist of the
script is just reading each `messages.csv`, grepping for emojis and parsing the required info from each message 
(list of emojis, each emoji's name, ID, and type, etc).

## Contributing and Improvements

Feature-wise, maybe you could parse the events file to get emojis which you've also used as reactions? 
