# discord-emoji-extractor

Download all the emojis you've ever sent inside messages on Discord. Skips duplicates and supports resuming downloads.

## Usage

1. Download your Discord data backup. You can get this by going to Discord Settings > Privacy & safety > Request all of my data
2. Extract the data somewhere. Doesn't matter.
3. `cd my_data/messages/`
4. Download [`run.sh`](https://github.com/l1ving/discord-emoji-extractor/blob/master/run.sh)
and [`extract_emojis.pl`](https://github.com/l1ving/discord-emoji-extractor/blob/master/extract_emojis.pl) to the messages directory.
```bash
wget https://github.com/l1ving/discord-emoji-extractor/blob/master/run.sh
wget https://github.com/l1ving/discord-emoji-extractor/blob/master/extract_emojis.pl
```
5. Make both scripts executable.
```bash
chmod +x run.sh extract_emojis.pl
```
6. Run the `run.sh` script. At any point you may cancel downloading and re-run the script from the `messages/` directory, and it will resume downloading.

#### Requirements

The following packages are required:
- [`ripgrep`](https://github.com/BurntSushi/ripgrep/)
- `bash`
- `sed`
- `perl`
- `wget`

## How it works

Both scripts are well documented with comments. 
I basically just wanted a way to grab old emojis from servers that I'd left, with a picture preview.

The `run.sh` script mainly parses the messages containing emojis, captured with `rg`. 
The `extract_emojis.pl` script is used for extracting the emojis from the parsed messages, and turning them into `wget` links.

## Contributing and Improvements

Contributions are welcome, this was a half hour hacky project I whipped up quickly, so the code isn't really the best.
Feature-wise, maybe you could parse the events file to get emojis which you've also used as reactions? 
