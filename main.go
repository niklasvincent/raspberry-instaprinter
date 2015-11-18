package main

import (
  "time"
//   "github.com/jung-kurt/gofpdf"
//   "github.com/gedex/go-instagram/instagram"
)

func main() {
  setupLogging()
  go startEmbeddedWebServer(":10443")

  t := time.NewTicker(10 * time.Second)
  for _ = range t.C {
      log.Info("Better look for photos!")
  }

  // var _ = instagram.NewClient(nil)

  // fmt.Println("Hello World\n")

  // pdf := gofpdf.New("P", "mm", "A4", "")
  // pdf.AddPage()
  // pdf.SetFont("Arial", "B", 16)
  // pdf.Cell(40, 10, "Hello, world")
  // pdf.OutputFileAndClose("hello.pdf")
}
