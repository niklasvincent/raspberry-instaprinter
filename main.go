package main

import (
  "time"
  "fmt"
//   "github.com/jung-kurt/gofpdf"
)

func main() {
  setupLogging()

  configuration := ReadConfiguration("configuration.json")

  // var instagramClient = instagram.NewClient(nil)
  // instagramClient.ClientID = configuration.InstagramClientId

  go startEmbeddedWebServer(":10443")

  EventLoop(&configuration)

  // pdf := gofpdf.New("P", "mm", "A4", "")
  // pdf.AddPage()
  // pdf.SetFont("Arial", "B", 16)
  // pdf.Cell(40, 10, "Hello, world")
  // pdf.OutputFileAndClose("hello.pdf")
}

func EventLoop(configuration *Configuration) {
  var currentMaxID = "0"
  currentLastCreatedTime := int64(0)
  t := time.NewTicker(10 * time.Second)
  for _ = range t.C {
      for _, hashTag := range configuration.HashTags {
        log.Info(fmt.Sprintf("Looking for new photos with tag #%s", hashTag))

        instagramPhotos, nextMaxID, lastCreatedTime := PhotosWithHashTag(hashTag, currentMaxID, currentLastCreatedTime, configuration.InstagramClientId)
        currentMaxID = nextMaxID
        currentLastCreatedTime = lastCreatedTime
        for i, p := range instagramPhotos {
          fmt.Printf("%d. %v - %s, %s\n", i, p.CreatedTime, p.ID, p.URL)
        }
        // opt := &instagram.Parameters{
        //   MinID: CurrentMaxID,
        //   Count: 100,
        // }

        // media, _, err := instagramClient.Tags.RecentMedia(HashTag, opt)

        // if err != nil {
        //   log.Error(fmt.Sprintf("Tags.RecentMedia returned error: %v", err))
        // }

        // if (len(media) > 0) {
        //   CurrentMaxID = media[0].ID
        // }

        // log.Debug(fmt.Sprintf("Current Max ID = %v", CurrentMaxID))

        // for _, m := range media {
        //   //var ImageURL = m.Images.StandardResolution.URL
        //   fmt.Printf("ID: %v\n\n", m.ID)
        // }
      }
  }
}