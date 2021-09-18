package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	dirFlag    = flag.String("dir", "", "Directory to search")
	dirConfirm = flag.Bool("dirconfirm", false, "Automatically confirm dir is correct")
	msgFile    = "messages.csv"
	baseUrl    = "https://cdn.discordapp.com/emojis/"
	fileMode   = os.FileMode(0644)
)

type Emoji struct {
	name     string
	id       int
	strId    string
	filename string
	ext      string
}

func main() {
	flag.Parse()

	// Let the user select a directory
	dir := selectDir(true)
	dir = formatDir(dir)
	fmt.Printf("Searching directory \"%s\"\n", dir)

	// Find all messages.csv files in said dir
	validFiles := getFileList(dir)
	validFilesAmt := len(validFiles)
	fmt.Printf("Found %v %s files\n", validFilesAmt, msgFile)

	if validFilesAmt == 0 {
		fmt.Printf("Couldn't find any %s files, maybe try another directory? "+
			"Make sure you are selecting the messages directory which contains c<number> folders, "+
			"or a c<number> folder itself. Exiting...", msgFile)
		return
	}

	uniqueEmojis := extractUniqueEmojis(validFiles)

	// Download all unique emojis found
	emojisDir := dir + "emojis/"
	err := os.Mkdir(emojisDir, fileMode)
	if err != nil && !strings.HasSuffix(err.Error(), "file exists") {
		// The file exists error is fine to ignore
		fmt.Printf("Error creating emojis directory: %v\n", err)
		return
	}

	downloadAllEmojis(uniqueEmojis, emojisDir)
}

// downloadAllEmojis will download each emoji to messages/emojis/name-id.ext
func downloadAllEmojis(emojis []Emoji, dir string) {
	emojisToDownload := make([]Emoji, 0)
	current := 0
	skipped := 0

	// Make a list of non-downloaded emojis (not found on disk)
	// These are split into two for loops to avoid messing with the download progress bar
	for _, emoji := range emojis {
		path := dir + emoji.filename

		if !checkFileExists(path) {
			emojisToDownload = append(emojisToDownload, emoji)
		} else {
			skipped += 1
		}
	}

	if skipped != 0 {
		fmt.Printf("Skipping %v emojis because they are already downloaded\n", skipped)
	}

	total := len(emojisToDownload)

	// Download all missing emojis
	for _, emoji := range emojisToDownload {
		current += 1
		path := dir + emoji.filename

		err := downloadFile(path, baseUrl+emoji.strId+emoji.ext, current, total)
		if err != nil {
			fmt.Printf("Error downloading emoji: %v\n", err)
		}
	}

	fmt.Println("") // for whatever reason the bar doesn't do this for you
}

// downloadFile will download a given file from url, and display the progress
func downloadFile(filepath string, url string, current int, total int) error {
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	f, _ := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, fileMode)
	defer f.Close()

	count := "(" + strconv.Itoa(current) + "/" + strconv.Itoa(total) + ")"
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("downloading emojis "+count),
		progressbar.OptionSpinnerType(4),
		progressbar.OptionShowBytes(true),
	)

	_, err := io.Copy(io.MultiWriter(f, bar), res.Body)
	defer res.Body.Close()

	return err
}

// extractUniqueEmojis will search all given files, and extract unique emojis by ID
func extractUniqueEmojis(files []string) []Emoji {
	messages := make([]string, 0)
	for _, file := range files {
		csvF, err := os.Open(file)
		if err != nil {
			fmt.Printf("Couldn't open file \"%s\": %v\n", file, err)
			continue
		}

		r := csv.NewReader(csvF)

		for {
			column, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("Error reading csv from \"%s\": %v\n", file, err)
				break
			}
			messages = append(messages, column[2]) // Third column is the message content
		}
	}

	fmt.Printf("Found %v messages to parse\n", len(messages))

	// Find emojis in all messages now

	parsedEmojis := make([]Emoji, 0)

	for _, message := range messages {
		re := regexp.MustCompile("<(a|):([A-z0-9_]+):([0-9]+)>")
		emojis := re.FindAllStringSubmatch(message, -1)

		// Parse all emojis in a message. emoji[0] is the full match, 1 is group 1 and so on
		for _, emoji := range emojis {
			id, err := strconv.Atoi(emoji[3])
			if err != nil {
				fmt.Printf("Error extracting emoji \"%s\": %v\n", emoji, err)
				continue
			}

			ext := ".png"
			if len(emoji[1]) != 0 {
				ext = ".gif"
			}

			filename := emoji[2] + "-" + emoji[3] + ext

			parsedEmoji := Emoji{emoji[2], id, emoji[3], filename, ext}
			parsedEmojis = append(parsedEmojis, parsedEmoji)
		}
	}

	fmt.Printf("Found %v emojis inside messages\n", len(parsedEmojis))

	uniqueEmojis := make([]Emoji, 0)

	// Append unique emojis by their ID. This is slow, but I am unsure how to make it faster
	for _, pEmoji := range parsedEmojis {
		found := false
		for _, uEmoji := range uniqueEmojis {
			if uEmoji.id == pEmoji.id {
				found = true
				break
			}
		}

		if !found {
			uniqueEmojis = append(uniqueEmojis, pEmoji)
		}
	}

	fmt.Printf("Found %v unique emojis\n", len(uniqueEmojis))

	return uniqueEmojis
}

// getFileList will look for all the messages.csv files that exist in dir
func getFileList(dir string) []string {
	files := make([]string, 0)

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, f := range fileInfos {
		msgFilePath := dir + f.Name() + "/" + msgFile

		if !f.IsDir() {
			// If f is not a directory, but is a messages.csv
			if f.Name() == msgFile {
				msgFilePath = dir + "/" + msgFile
				if checkFileExists(msgFilePath) {
					files = append(files, msgFilePath)
				}
			}
		} else {
			// If f is a directory (hopefully the channel directory)
			if checkFileExists(msgFilePath) {
				files = append(files, msgFilePath)
			}
		}
	}

	return files
}

// checkFileExists will check if a certain path exists, and ensures it is not a directory
func checkFileExists(path string) bool {
	if f, err := os.Stat(path); err == nil {
		return !f.IsDir() // path exists, return true if it is not a directory
	} else if os.IsNotExist(err) {
		return false // file does not exist
	} else { // schrodinger's file
		log.Printf("Error: %v", err)
		panic(err)
	}
}

// formatDir will append a / to the end of a dir path if it is missing
func formatDir(dir string) string {
	last, _ := utf8.DecodeLastRuneInString(dir)
	if last != '/' {
		dir += "/"
	}
	return dir
}

// selectDir will ask the user for a directory and confirm the directory they chose
func selectDir(firstRun bool) string {
	// If dirConfirm is selected and there's a dir set, choose it automatically
	if *dirConfirm && len(*dirFlag) != 0 {
		return *dirFlag
	}

	var dir string
	if firstRun {
		dir = *dirFlag
	}

	if len(*dirFlag) == 0 {
		fmt.Printf("Select a directory to scan (use . for current): ")
		fmt.Scan(&dir)
	}

	fmt.Printf("Selected directory: \"%s\"\n", dir)
	fmt.Printf("Is this correct? (Y/N): ")

	var correct string
	fmt.Scan(&correct)

	first := strings.ToLower(correct[0:1])
	if first != "y" {
		fmt.Printf("Selected No, trying again.\n")
		return selectDir(false)
	}

	return dir
}
