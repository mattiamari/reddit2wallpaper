bin_name = reddit2wallpaper
version = 0.1.1
sources = cmd/reddit2wallpaper/main.go
ldflags = "-X main.appVersion=$(version)"

.PHONY: dist

all: linux osx win

clean:
	rm -rf dist

dist:
	mkdir -p dist

linux: dist
	GOOS=linux GOARCH=amd64 go build -ldflags=$(ldflags) -o dist/$(bin_name)-linux-amd64 $(sources)

osx: dist
	GOOS=darwin GOARCH=amd64 go build -ldflags=$(ldflags) -o dist/$(bin_name)-darwin-amd64 $(sources)

win: dist
	GOOS=windows GOARCH=amd64 go build -ldflags=$(ldflags) -o dist/$(bin_name)-windows-amd64.exe $(sources)
