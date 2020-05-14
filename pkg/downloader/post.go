package downloader

import (
	"regexp"
	"strconv"
	"strings"
)

var resolutionRegex *regexp.Regexp = regexp.MustCompile(`.*\[(\d+)[xX](\d+)\]`)
var imageExtensionRegex *regexp.Regexp = regexp.MustCompile(`\.(jpeg|jpg|png)$`)

// Post is a Reddit post
type Post struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type PostList []Post

type Filter func(Post) bool

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

func (posts PostList) Filter(filter Filter) PostList {
	res := []Post{}

	for _, p := range posts {
		if filter(p) {
			res = append(res, p)
		}
	}

	return res
}

func ResolutionFilter(minWidth int, minHeight int) Filter {
	return func(p Post) bool {
		w, h := p.GetPictureResolution()
		return w >= minWidth && h >= minHeight
	}
}

func AspectRatioFilter(a int, b int) Filter {
	return func(p Post) bool {
		w, h := p.GetPictureResolution()
		return w/h >= a/b
	}
}

func FileExtensionFilter(extensions []string) Filter {
	return func(p Post) bool {
		for _, ext := range extensions {
			if strings.HasSuffix(p.URL, "."+ext) {
				return true
			}
		}

		return false
	}
}
