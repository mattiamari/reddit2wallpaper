package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/mattiamari/reddit2wallpaper/pkg/downloader"
)

const appVersion = "0.2.3"

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Couldn't find the current user")
	}

	downloader.AppVersion = appVersion
	fmt.Printf("Reddit2Wallpaper v%s https://github.com/mattiamari/reddit2wallpaper\n\n", appVersion)

	var subreddit string
	downloadDir := filepath.Join(usr.HomeDir, "media", "Pictures", "Wallpapers", "reddit2wallpaper")
	var minHeight int
	var minWidth int
	var ratioW int
	var ratioH int
	var topPosts bool

	flag.StringVar(&subreddit, "subreddit", "EarthPorn", "Name of the subreddit to use")
	flag.StringVar(&downloadDir, "download_dir", downloadDir, "Directory in which to save wallpapers")
	flag.IntVar(&minWidth, "minwidth", 0, "Minimum width for the photo to download")
	flag.IntVar(&minHeight, "minheight", 0, "Minimum height for the photo to download")
	flag.IntVar(&ratioW, "ratio_width", 0, "Aspect ratio width")
	flag.IntVar(&ratioH, "ratio_height", 0, "Aspect ratio height")
	flag.BoolVar(&topPosts, "top_posts", false, "Fetch 'top' posts instead of 'new' posts")

	flag.Parse()

	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		fmt.Printf("Download directory '%s' does not exist\nSo I'll just make it lolol\n", downloadDir)
		err = os.MkdirAll(downloadDir, 0777)
		if err != nil {
			log.Fatal("I had some trouble creating a directory to put the photos")
		}
	}

	if (ratioW != 0 || ratioH != 0) && (ratioW < 0 || ratioH < 0 || ratioW*ratioH == 0) {
		log.Fatal("Invalid aspect ratio")
	}

	fmt.Printf("Looking for %dx%d %d:%d photos on r/%s \n", minWidth, minHeight, ratioW, ratioH, subreddit)

	sort := downloader.SortDefault
	if topPosts {
		sort = downloader.SortNew
	}

	posts, err := downloader.GetPosts(subreddit, sort, 100)
	if err != nil {
		log.Fatal(err)
	}

	posts = posts.Filter(downloader.FileExtensionFilter([]string{"jpg", "jpeg", "png"}))
	posts = posts.Filter(downloader.ResolutionFilter(minWidth, minHeight))

	if ratioW != 0 && ratioH != 0 {
		posts = posts.Filter(downloader.AspectRatioFilter(ratioW, ratioH))
	}

	downloader.DownloadAll(posts, downloadDir)
}
