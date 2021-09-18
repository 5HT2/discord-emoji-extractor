discord-emoji-extractor: clean
	go get -u "github.com/schollz/progressbar/v3"
	go build -o extract

clean:
	rm -f extract
