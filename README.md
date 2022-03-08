# discord-emoji-extractor

[![time tracker](https://wakatime.com/badge/github/l1ving/discord-emoji-extractor.svg)](https://wakatime.com/badge/github/l1ving/discord-emoji-extractor)
[![CodeFactor](https://img.shields.io/codefactor/grade/github/5HT2/discord-emoji-extractor?logo=codefactor&logoColor=white)](https://www.codefactor.io/repository/github/5HT2/discord-emoji-extractor)

Download all the emojis you've ever sent inside messages on Discord. Supports skipping duplicates and resuming downloads.
![image](https://user-images.githubusercontent.com/17222512/133897865-45e69ee8-5214-43cf-87d8-01f96e43baac.png)

## Usage

These instructions are for building to Go project. 
The bash [equivalent](https://github.com/5HT2/discord-emoji-extractor/blob/master/extract.sh) does not need compiling.

```bash
git clone git@github.com:5HT2/discord-emoji-extractor.git
cd discord-emoji-extractor
make
./extract -h # Run the program with help arguments
```

## Running the Go version

1. Download your Discord data backup. You can get this by going to Discord Settings > Privacy & safety > Request all of my data
2. Extract the data somewhere. Doesn't matter.
3. Follow the above [usage](#Usage) instructions and run the program from anywhere.
4. Follow the interactive instructions. You can use the `-dir $DIR -dirconfirm` args with `DIR` set to a path to skip the prompts.
5. At any point you may cancel downloading and re-run the program, and it will resume downloading.

## Running the Bash version

1. Download your Discord data backup. You can get this by going to Discord Settings > Privacy & safety > Request all of my data
2. Extract the data somewhere. Doesn't matter.
3. `cd my_data/messages/`
4. Download [`extract.sh`](https://github.com/5HT2/discord-emoji-extractor/blob/master/extract.sh) to the `messages` directory.
```bash
wget https://github.com/5HT2/discord-emoji-extractor/raw/master/extract.sh
```
5. Make the script executable.
```bash
chmod +x extract.sh
```
6. Run the `extract.sh` script. The bash script does not support pausing and resuming downloads.

## How it works

I basically just wanted a way to grab old emojis from servers that I'd left, with a picture preview, so the gist of the
script is just reading each `messages.csv`, grepping for emojis and parsing the required info from each message 
(list of emojis, each emoji's name, ID, and type, etc).

The Go version works similarly, but is much faster and supports cancelling and resuming downloading.

## Contributing and Improvements

Feature-wise, maybe you could parse the events file to get emojis which you've also used as reactions?

A bot to dynamically upload to a bunch of selected servers (as per the default 50 emoji / server limit) and skip ones with the same name, would be neat. Bots have a stricted upload ratelimit than users, so letting it run in the background to wait out the timeout would be ideal.

Possible command syntax could be
```bash
./extract -dir $DIR -dirconfirm -upload -token $TOKEN -serverids 96230004047740928,785362280601616406,343525052332900352
```
