bin_name = reddit2wallpaper
sources = cmd/reddit2wallpaper/main.go

.PHONY: dist

all: dist build_linux build_osx build_win

clean:
	rm -rf dist

dist:
	mkdir -p dist

build_linux:
	GOOS=linux GOARCH=amd64 go build -o dist/$(bin_name)-linux-amd64 $(sources)

build_osx:
	GOOS=darwin GOARCH=amd64 go build -o dist/$(bin_name)-darwin-amd64 $(sources)

build_win:
	GOOS=windows GOARCH=amd64 go build -o dist/$(bin_name)-windows-amd64.exe $(sources)