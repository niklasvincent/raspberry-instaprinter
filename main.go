package main

import (
  "fmt"
  "github.com/jung-kurt/gofpdf"
  "github.com/gedex/go-instagram/instagram"
)

func main() {
  var _ = instagram.NewClient(nil)

  fmt.Println("Hello World\n")

  pdf := gofpdf.New("P", "mm", "A4", "")
  pdf.AddPage()
  pdf.SetFont("Arial", "B", 16)
  pdf.Cell(40, 10, "Hello, world")
  pdf.OutputFileAndClose("hello.pdf")
}
