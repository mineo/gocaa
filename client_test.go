package caa

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"code.google.com/p/go-uuid/uuid"
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
				"thumbnails":{
					"large":"http://coverartarchive.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd/829521842-500.jpg",
					"small":"http://coverartarchive.org/release/76df3287-6cda-33eb-8e9a-044b5e15ffdd/829521842-250.jpg"
				},
				"id":"829521842"
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
			   "images": [
				   {
					   "approved": true,
					   "back": false,
					   "comment": "",
					   "edit": 20202510,
					   "front": true,
					   "id": "2860563776",
					   "image": "http://coverartarchive.org/release/f7638b9b-a9aa-4c03-8734-9e692699f8b1/2860563776.jpg",
					   "thumbnails": {
						   "large": "http://coverartarchive.org/release/f7638b9b-a9aa-4c03-8734-9e692699f8b1/2860563776-500.jpg",
						   "small": "http://coverartarchive.org/release/f7638b9b-a9aa-4c03-8734-9e692699f8b1/2860563776-250.jpg"
					   },
					   "types": [
						   "Front"
					   ]
				   }
			   ],
			   "release": "http://musicbrainz.org/release/f7638b9b-a9aa-4c03-8734-9e692699f8b1"
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
	d.Assert(info.Release, Equals, "http://musicbrainz.org/release/f7638b9b-a9aa-4c03-8734-9e692699f8b1")

	i := info.Images[0]
	d.Assert(i.Comment, Equals, "")
	d.Assert(i.Edit, Equals, 20202510)
	d.Assert(i.ID, Equals, "2860563776")
	d.Assert(i.Image, Equals, "http://coverartarchive.org/release/f7638b9b-a9aa-4c03-8734-9e692699f8b1/2860563776.jpg")
	d.Assert(i.Types[0], Equals, "Front")
	d.Assert(len(i.Types), Equals, 1)
	d.AssertFalse(i.Back)
	d.AssertTrue(i.Approved)
	d.AssertTrue(i.Front)
}
