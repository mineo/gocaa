// Package caa provides access to the Cover Art Archive (http://coverartarchive.org)
package caa

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"code.google.com/p/go-uuid/uuid"
)

const baseurl = "http://coverartarchive.org"

// CAAClient manages the communication with the Cover Art Archive.
type CAAClient struct {
	useragent string
	client    http.Client
}

// NewCAAClient returns a new CAAClient that uses the User-Agent useragent
func NewCAAClient(useragent string) (c *CAAClient) {
	c = &CAAClient{useragent: useragent, client: http.Client{}}
	return
}

func buildURL(path string) (url *url.URL) {
	url, err := url.Parse(baseurl)

	if err != nil {
		return
	}

	url.Path = path
	return
}

func (c *CAAClient) get(url *url.URL) (resp *http.Response, err error) {
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("User-Agent", c.useragent)

	resp, err = c.client.Do(req)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return
}

func (c *CAAClient) getAndJSON(url *url.URL) (info *CoverArtInfo, err error) {
	resp, err := c.get(url)

	defer resp.Body.Close()

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &info)

	return

}

func (c *CAAClient) getImage(entitytype string, mbid uuid.UUID, imageid string, size int) (image CoverArtImage, err error) {
	var extra string

	if size == ImageSizeSmall || size == 250 {
		extra = "-250"
	} else if size == ImageSizeLarge || size == 500 {
		extra = "-500"
	} else {
		extra = ""
	}

	url := buildURL(entitytype + "/" + mbid.String() + "/" + imageid + extra)
	resp, err := c.get(url)

	defer resp.Body.Close()

	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = HTTPError{StatusCode: resp.StatusCode, URL: url}
		return
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	image.Data = data

	ext := filepath.Ext(resp.Request.URL.String())
	mimetype := mime.TypeByExtension(ext)

	image.Mimetype = mimetype

	return

}

// GetReleaseInfo retrieves information about the images in the Cover Art Archive for the release with the MBID mbid
func (c *CAAClient) GetReleaseInfo(mbid uuid.UUID) (info *CoverArtInfo, err error) {
	url := buildURL("release/" + mbid.String())
	info, err = c.getAndJSON(url)
	return
}

// GetReleaseFront retrieves the front image of the release with the MBID mbid in the specified size
func (c *CAAClient) GetReleaseFront(mbid uuid.UUID, size int) (image CoverArtImage, err error) {
	image, err = c.getImage("release", mbid, "front", size)
	return
}

// GetReleaseBack retrieves the back image of the release with the MBID mbid in the specified size
func (c *CAAClient) GetReleaseBack(mbid uuid.UUID, size int) (image CoverArtImage, err error) {
	image, err = c.getImage("release", mbid, "back", size)
	return
}

// GetReleaseImage retrieves the image with the id imageid of the release with the MBID mbid in the specified size
func (c *CAAClient) GetReleaseImage(mbid uuid.UUID, imageid int, size int) (image CoverArtImage, err error) {
	id := strconv.Itoa(imageid)
	image, err = c.getImage("release", mbid, id, size)
	return
}

// GetReleaseGroupInfo retrieves information about the images in the Cover Art Archive for the release group with the MBID mbid
func (c *CAAClient) GetReleaseGroupInfo(mbid uuid.UUID) (info *CoverArtInfo, err error) {
	url := buildURL("release-group/" + mbid.String())
	info, err = c.getAndJSON(url)
	return
}

// GetReleaseGroupFront retrieves the front image of the release group with the MBID mbid in the specified size
func (c *CAAClient) GetReleaseGroupFront(mbid uuid.UUID, size int) (image CoverArtImage, err error) {
	if size != ImageSizeOriginal {
		err = InvalidImageSizeError{EntityType: "release-group", Size: size}
		return
	}
	image, err = c.getImage("release-group", mbid, "front", size)
	return
}
