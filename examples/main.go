package main

import (
	"log"

	"github.com/mineo/gocaa"
	"sync"
)

var mbid = caa.StringToUUID("8b263c14-eb61-4a19-a784-6ebd084c46aa")

func getReleaseInfo(c *caa.CAAClient, wg *sync.WaitGroup) {
	defer wg.Done()

	info, err := c.GetReleaseInfo(mbid)

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	log.Printf("The URL for %s's front image is %s\n", mbid.String(), info.Images[0].Image)
}

func getReleaseImageByID(c *caa.CAAClient, wg *sync.WaitGroup) {
	defer wg.Done()
	image, err := c.GetReleaseImage(mbid, 7645498428, caa.ImageSizeSmall)

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	log.Printf("The MIME type of %s's front image is %s\n", mbid.String(), image.Mimetype)
}

func getReleaseGroupFrontImage(c *caa.CAAClient, wg *sync.WaitGroup) {
	defer wg.Done()
	rgid := caa.StringToUUID("f1498e3c-b179-46c3-9b82-1d011e725f30")
	image, err := c.GetReleaseGroupFront(rgid, caa.ImageSizeOriginal)

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	log.Printf("The MIME type for %s's front image is %s\n", rgid.String(), image.Mimetype)
}

func main() {
	c := caa.NewCAAClient("mineostestclient")
	wg := sync.WaitGroup{}

	go getReleaseInfo(c, &wg)
	wg.Add(1)
	go getReleaseImageByID(c, &wg)
	wg.Add(1)
	go getReleaseGroupFrontImage(c, &wg)
	wg.Add(1)

	wg.Wait()
}
