package caa

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pborman/uuid"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var (
	caaclient *CAAClient
	server    *httptest.Server
	_         = Suite(&MySuite{})
)

const TESTUSERAGENT = "useragent"

func setup(f func(http.ResponseWriter, *http.Request)) {
	caaclient = NewCAAClient(TESTUSERAGENT)
	server = httptest.NewServer(http.HandlerFunc(f))
	caaclient.BaseURL = server.URL
}

type D struct {
	*C
}

func (d *D) AssertTrue(b bool) {
	d.Assert(b, Equals, true)
}

func (d *D) AssertFalse(b bool) {
	d.Assert(b, Equals, false)
}

func (s *MySuite) TestUserAgent(c *C) {
	f := func(w http.ResponseWriter, req *http.Request) {
		ua := req.Header.Get("User-Agent")
		c.Assert(ua, Equals, TESTUSERAGENT)
		c.Assert(req.URL.Path, Equals, "/test")
	}

	setup(f)
	defer server.Close()

	u, _ := url.Parse(fmt.Sprintf("%s/test", server.URL))
	caaclient.get(u)
}

func (s *MySuite) TestStatusCodes(c *C) {
	possibleStatus := []int{400, 404, 405, 406, 502, 503}
	for _, statusCode := range possibleStatus {
		f := func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(statusCode)
			return
		}

		setup(f)

		u, _ := url.Parse(fmt.Sprintf("%s", server.URL))
		_, err := caaclient.get(u)

		switch t := err.(type) {
		case HTTPError:
			c.Assert(t.StatusCode, Equals, statusCode)
		default:
			c.Fail()
		}

		server.Close()
	}
}

func (s *MySuite) TestGetReleaseInfo(c *C) {
	d := D{c}
	mbid := "76df3287-6cda-33eb-8e9a-044b5e15ffdd"
	f := func(w http.ResponseWriter, req *http.Request) {
		path := fmt.Sprintf("/release/%s", mbid)
		c.Assert(req.URL.Path, Equals, path)

		// Taken from https://musicbrainz.org/doc/Cover_Art_Archive/API#.2Frelease.2F.7Bmbid.7D.2F
		jsonresp := `
                {
                  "images":[
                     {
                        "types":[
                           "Front"
                        ],
                        "front":true,
                        "back":false,
                        "edit":17462565,
                        "image":"http://coverartarchive.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd/829521842.jpg",
                        "comment":"",
                        "approved":true,
                        "id":"829521842",
                        "thumbnails":{
                          "250":"http://coverartarchive.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd/829521842-250.jpg",
                          "500":"http://coverartarchive.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd/829521842-500.jpg",
                          "1200":"http://coverartarchive.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd/829521842-1200.jpg",
                          "small":"http://coverartarchive.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd/829521842-250.jpg",
                          "large":"http://coverartarchive.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd/829521842-500.jpg"
                        }
                     }
                  ],
                  "release":"http://musicbrainz.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd"
                }
		`
		w.Write([]byte(jsonresp))
	}

	setup(f)
	defer server.Close()

	info, err := caaclient.GetReleaseInfo(uuid.Parse(mbid))

	if err != nil {
		c.Fail()
	}

	d.Assert(len(info.Images), Equals, 1)
	d.Assert(info.Release, Equals, "http://musicbrainz.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd")

	i := info.Images[0]
	d.Assert(i.Comment, Equals, "")
	d.Assert(i.Edit, Equals, 17462565)
	d.Assert(i.ID, Equals, "829521842")
	d.Assert(i.Image, Equals, "http://coverartarchive.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd/829521842.jpg")
	d.Assert(i.Types[0], Equals, "Front")
	d.Assert(len(i.Types), Equals, 1)
	d.Assert(len(i.Thumbnails), Equals, 5)
	d.AssertFalse(i.Back)
	d.AssertTrue(i.Approved)
	d.AssertTrue(i.Front)
}

func (s *MySuite) TestGetReleaseGroupInfo(c *C) {
	d := D{c}
	mbid := "c31a5e2b-0bf8-32e0-8aeb-ef4ba9973932"
	f := func(w http.ResponseWriter, req *http.Request) {
		path := fmt.Sprintf("/release-group/%s", mbid)
		c.Assert(req.URL.Path, Equals, path)

		// Taken from https://musicbrainz.org/doc/Cover_Art_Archive/API#.2Frelease-group.2F.7Bmbid.7D.2F
		jsonresp := `
                {
                    "release":"https://musicbrainz.org/release/f268b8bc-2768-426b-901b-c7966e76de29",
                    "images":[
                        {
                            "edit":37284546,
                            "id":"12750224075",
                            "image":"http://coverartarchive.org/release/f268b8bc-2768-426b-901b-c7966e76de29/12750224075.png",
                            "thumbnails":{
                                "250":"http://coverartarchive.org/release/f268b8bc-2768-426b-901b-c7966e76de29/12750224075-250.jpg",
                                "500":"http://coverartarchive.org/release/f268b8bc-2768-426b-901b-c7966e76de29/12750224075-500.jpg",
                                "1200":"http://coverartarchive.org/release/f268b8bc-2768-426b-901b-c7966e76de29/12750224075-1200.jpg",
                                "small":"http://coverartarchive.org/release/f268b8bc-2768-426b-901b-c7966e76de29/12750224075-250.jpg",
                                "large":"http://coverartarchive.org/release/f268b8bc-2768-426b-901b-c7966e76de29/12750224075-500.jpg"
                            },
                            "comment":"",
                            "approved":true,
                            "front":false,
                            "types":[
                                "Back"
                            ],
                            "back":true
                        }
                    ]
                }
		`
		w.Write([]byte(jsonresp))
	}

	setup(f)
	defer server.Close()

	info, err := caaclient.GetReleaseGroupInfo(uuid.Parse(mbid))

	if err != nil {
		c.Fail()
	}

	d.Assert(len(info.Images), Equals, 1)
	d.Assert(info.Release, Equals, "https://musicbrainz.org/release/f268b8bc-2768-426b-901b-c7966e76de29")

	i := info.Images[0]
	d.Assert(i.Comment, Equals, "")
	d.Assert(i.Edit, Equals, 37284546)
	d.Assert(i.ID, Equals, "12750224075")
	d.Assert(i.Image, Equals, "http://coverartarchive.org/release/f268b8bc-2768-426b-901b-c7966e76de29/12750224075.png")
	d.Assert(i.Types[0], Equals, "Back")
	d.Assert(len(i.Types), Equals, 1)
	d.Assert(len(i.Thumbnails), Equals, 5)
	d.AssertTrue(i.Back)
	d.AssertTrue(i.Approved)
	d.AssertFalse(i.Front)
}
