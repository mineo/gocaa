package gocaa

type CoverArtInfo struct {
	Images  []CoverArtImageInfo
	Release string
}

type CoverArtImageInfo struct {
	Types      []string
	Front      bool
	Back       bool
	Comment    string
	Thumbnails ThumbnailMap
	Approved   bool
	Edit       int
}

type CoverArtImage struct {
	Mimetype string
	Data     []byte
}

type ThumbnailMap map[string]string
