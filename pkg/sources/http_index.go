package sources

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"sort"
	"strings"
)

type directoryIndex map[string]int

// GetIndex Returns a directory index for the given response and URL.
func (r Routes) GetIndex(c *gin.Context) render.Render {
	index := make(directoryIndex)

	requestURL := strings.TrimSuffix(c.Request.URL.Path, "/") + "/"

	for url := range r {
		// Skip all non-matching registered routes
		if !strings.HasPrefix(url, requestURL) {
			continue
		}

		// Trim prefix and split path
		path := strings.TrimPrefix(url, requestURL)
		pathParts := strings.Split(path, "/")

		// next item
		nextItem := pathParts[0]

		itemCount, ok := index[nextItem]
		if !ok {
			index[nextItem] = 0
		}

		if itemCount < len(pathParts) {
			index[nextItem] = len(pathParts)
		}
	}

	results := make([]string, 0)
	for url, levelCount := range index {
		item := url
		if levelCount > 1 {
			item += "/"
		}

		results = append(results, item)
	}

	if len(results) > 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i] > results[j]
		})

		return render.String{Format: strings.Join(results, "\n")}
	}

	return nil
}
