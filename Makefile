dev:
	./bin/go-bindata web/...
	go build
	./raspberry-instaprinter