build:
	CGO_ENABLED=0 go build -ldflags "-s -w -extldflags=-static" -o cpbx main.go