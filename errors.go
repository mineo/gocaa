package gocaa

import (
	"fmt"
	"net/url"
)

type HTTPError struct {
	StatusCode int
	Url        *url.URL
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("%i on %s", e.StatusCode, e.Url.String())
}

type InvalidImageSizeError struct {
	Entitytype string
	Size       int
}

func (e InvalidImageSizeError) Error() string {
	return fmt.Sprintf("%s doesn't support image size %i", e.Entitytype, e.Size)
}
