package gocaa

import (
	"fmt"
	"net/url"
)

type HTTPError struct {
	StatusCode int
	URL        *url.URL
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("%d on %s", e.StatusCode, e.URL.String())
}

type InvalidImageSizeError struct {
	Entitytype string
	Size       int
}

func (e InvalidImageSizeError) Error() string {
	return fmt.Sprintf("%s doesn't support image size %d", e.Entitytype, e.Size)
}
