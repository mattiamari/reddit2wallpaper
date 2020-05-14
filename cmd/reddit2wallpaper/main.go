package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mattiamari/reddit2wallpaper/pkg/downloader"
)

var appVersion string

func main() {
	downloader.AppVersion = appVersion
	fmt.Printf("Reddit2Wallpaper v%s https://github.com/mattiamari/reddit2wallpaper\n\n", appVersion)

	var subreddit string
	var downloadDir string
	var minHeight int
	var minWidth int
	var topPosts bool

	flag.StringVar(&subreddit, "subreddit", "EarthPorn", "Name of the subreddit to use")
	flag.StringVar(&downloadDir, "download_dir", "", "Directory in which to save wallpapers")
	flag.IntVar(&minWidth, "minwidth", -1, "Minimum width for the photo to download")
	flag.IntVar(&minHeight, "minheight", -1, "Minimum height for the photo to download")
	flag.BoolVar(&topPosts, "top_posts", false, "Fetch 'top' posts instead of 'new' posts")

	flag.Parse()

	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		log.Fatalf("Download directory '%s' does not exist\n", downloadDir)
	}

	fmt.Printf("Looking for %dx%d photos on r/%s \n", minWidth, minHeight, subreddit)

	sort := downloader.SortDefault
	if topPosts {
		sort = downloader.SortNew
	}

	posts, err := downloader.GetPosts(subreddit, sort, 100)
	if err != nil {
		log.Fatal(err)
	}

	posts = downloader.FilterImages(posts)
	posts = downloader.FilterResolution(posts, minWidth, minHeight)

	downloader.DownloadAll(posts, downloadDir)
}
