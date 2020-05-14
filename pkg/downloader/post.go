package downloader

import (
	"regexp"
	"strconv"
)

var resolutionRegex *regexp.Regexp = regexp.MustCompile(`.*\[(\d+)[xX](\d+)\]`)
var imageExtensionRegex *regexp.Regexp = regexp.MustCompile(`\.(jpeg|jpg|png)$`)

// Post is a Reddit post
type Post struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// GetPictureResolution parses the resolution from post title
// e.g. "An awesome post title [<height>x<width>]"
// like some subreddits do
func (p Post) GetPictureResolution() (width int, height int) {
	res := resolutionRegex.FindStringSubmatch(p.Title)

	if len(res) < 3 {
		return -1, -1
	}

	width, errH := strconv.Atoi(res[1])
	height, errW := strconv.Atoi(res[2])

	if errH != nil || errW != nil {
		return -1, -1
	}

	return width, height
}

// IsImage returns true if post url is an image
func (p Post) IsImage() bool {
	return imageExtensionRegex.MatchString(p.URL)
}

func FilterResolution(posts []Post, minWidth int, minHeight int) []Post {
	res := []Post{}

	for _, p := range posts {
		if w, h := p.GetPictureResolution(); w >= minWidth && h >= minHeight {
			res = append(res, p)
		}
	}

	return res
}

func FilterAspectRatio(posts []Post, a int, b int) []Post {
	res := []Post{}

	for _, p := range posts {
		if w, h := p.GetPictureResolution(); w/h >= a/b {
			res = append(res, p)
		}
	}

	return res
}

func FilterImages(posts []Post) []Post {
	res := []Post{}

	for _, p := range posts {
		if p.IsImage() {
			res = append(res, p)
		}
	}

	return res
}
