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
	Title  string `json:"title"`
	URL    string `json:"url"`
	Width  int
	Height int
}

type PostList []Post

type Filter func(Post) bool

// CacheResolution parses the resolution from post title
// e.g. "An awesome post title [<height>x<width>]"
// like some subreddits do, and caches it in the post itself
func (p *Post) CacheResolution() {
	res := resolutionRegex.FindStringSubmatch(p.Title)

	if len(res) < 3 {
		p.Width = 0
		p.Height = 0
	}

	width, errH := strconv.Atoi(res[1])
	height, errW := strconv.Atoi(res[2])

	if errH != nil || errW != nil {
		p.Width = 0
		p.Height = 0
	}

	p.Width = width
	p.Height = height
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
		return p.Width >= minWidth && p.Height >= minHeight
	}
}

func AspectRatioFilter(a int, b int) Filter {
	return func(p Post) bool {
		return p.Width/p.Height >= a/b
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
