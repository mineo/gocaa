package caa

// CoverArtInfo is the unmarshaled representation of a JSON file in the Cover Art Archive.
// See https://musicbrainz.org/doc/Cover_Art_Archive/API#Cover_Art_Archive_Metadata for an example.
type CoverArtInfo struct {
	Images  []CoverArtImageInfo
	Release string
}

// CoverArtImageInfo is the unmarshaled representation of a single images metadata in a CAA JSON file.
// See https://musicbrainz.org/doc/Cover_Art_Archive/API#Cover_Art_Archive_Metadata for an example.
type CoverArtImageInfo struct {
	Approved   bool
	Back       bool
	Comment    string
	Edit       int
	Front      bool
	ID         string
	Image      string
	Thumbnails thumbnailMap
	Types      []string
}

// CoverArtImage is a wrapper around an image from the CAA, containing its binary data and mimetype information.
type CoverArtImage struct {
	Data     []byte
	Mimetype string
}

type thumbnailMap map[string]string
