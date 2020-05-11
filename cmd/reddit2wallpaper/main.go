package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
)

const AppVersion string = "0.1.0"
const RedditUrl string = "https://www.reddit.com/r/"
const UserAgent string = "Reddit2Wallpaper/" + AppVersion

var ResolutionRegex *regexp.Regexp = regexp.MustCompile(`.*\[(\d+)[xX](\d+)\]`)
var ImageExtensionRegex *regexp.Regexp = regexp.MustCompile(`\.(jpeg|jpg|png)$`)

//var errlog *log.Logger = log.New(os.Stderr, "[error]", log.LstdFlags)

type Response struct {
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Children []PostWrapper `json:"children"`
}

type PostWrapper struct {
	Data Post `json:"data"`
}

type Post struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func (p Post) GetPictureResolution() *Resolution {
	res := ResolutionRegex.FindStringSubmatch(p.Title)

	if len(res) < 3 {
		return nil
	}

	width, errH := strconv.Atoi(res[1])
	height, errW := strconv.Atoi(res[2])

	if errH != nil || errW != nil {
		return nil
	}

	return &Resolution{width, height}
}

type Resolution struct {
	Width  int
	Height int
}

func main() {
	//log.SetOutput(os.Stdout)
	fmt.Printf("Reddit2Wallpaper v%s\n", AppVersion)

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
		fmt.Printf("Download directory '%s' does not exist\n", downloadDir)
		return
	}

	fmt.Printf("Looking for %dx%d photos on r/%s \n", minWidth, minHeight, subreddit)

	url := RedditUrl + subreddit + "/new.json"
	if topPosts {
		url = RedditUrl + subreddit + "/top.json"
	}

	fmt.Println("Fetching " + url)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Couldn't get '%s': %s\n", url, err.Error())
		return
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var body Response
	err = dec.Decode(&body)
	if err != nil {
		fmt.Printf("Couldn't parse json: %s\n", err.Error())
		return
	}

	for _, post := range body.Data.Children {
		resolution := post.Data.GetPictureResolution()

		if resolution == nil {
			continue
		}

		if minWidth > 0 && resolution.Width < minWidth {
			continue
		}

		if minHeight > 0 && resolution.Height < minHeight {
			continue
		}

		filename := path.Base(post.Data.Url)
		destFilename := path.Join(downloadDir, filename)

		if !ImageExtensionRegex.MatchString(filename) {
			continue
		}

		if _, err := os.Stat(destFilename); err == nil {
			fmt.Printf("%s already exists\n", post.Data.Url)
			continue
		}

		destFile, err := os.Create(destFilename)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		req, _ := http.NewRequest(http.MethodGet, post.Data.Url, nil)
		req.Header.Set("User-Agent", UserAgent)

		fmt.Printf("Downloading %s (%dx%d)... ", post.Data.Url, resolution.Width, resolution.Height)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer res.Body.Close()

		io.Copy(destFile, res.Body)
		fmt.Println("done")
	}
}
