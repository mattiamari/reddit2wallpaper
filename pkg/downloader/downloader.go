package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

// AppVersion is the version number to be used in user agent string. Defaults to 1.0.0
var AppVersion string = "1.0.0"

var userAgent string = "Reddit2Wallpaper/" + AppVersion

const redditURL string = "https://www.reddit.com"

const (
	// SortDefault default post sorting
	SortDefault = ""

	// SortNew newer posts first
	SortNew = "new"
)

type FileExistsError struct {
	Filename string
}

func (e FileExistsError) Error() string {
	return fmt.Sprintf("File '%s' already exists", e.Filename)
}

type response struct {
	Data responseData `json:"data"`
}

type responseData struct {
	Children []postWrapper `json:"children"`
}

type postWrapper struct {
	Data Post `json:"data"`
}

func GetPosts(subreddit string, sort string, limit int) (PostList, error) {
	url := fmt.Sprintf("%s/r/%s/%s.json?limit=%d", redditURL, subreddit, sort, limit)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get '%s': %s\n", url, err.Error())
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var body response
	err = dec.Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse json: %s\n", err.Error())
	}

	posts := []Post{}

	for _, post := range body.Data.Children {
		posts = append(posts, post.Data)
	}

	for i := range posts {
		posts[i].CacheResolution()
	}

	return posts, nil
}

func Download(post Post, outputDirectory string) error {
	filename := path.Base(post.URL)
	destFilename := path.Join(outputDirectory, filename)

	if _, err := os.Stat(destFilename); err == nil {
		return FileExistsError{destFilename}
	}

	destFile, err := os.Create(destFilename)
	if err != nil {
		return err
	}
	defer destFile.Close()

	req, _ := http.NewRequest(http.MethodGet, post.URL, nil)
	req.Header.Set("User-Agent", userAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	io.Copy(destFile, res.Body)
	return nil
}

func DownloadAll(posts []Post, outputDirectory string) {
	for _, post := range posts {
		fmt.Printf("Downloading '%s'... ", post.Title)
		err := Download(post, outputDirectory)

		if _, ok := err.(FileExistsError); ok {
			fmt.Println("already exists, skipped")
		} else if err != nil {
			log.Printf("Couldn't download %s: %s", post.URL, err.Error())
		} else {
			fmt.Println("done")
		}
	}
}
