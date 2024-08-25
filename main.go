package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"golang.org/x/net/html"
)

type HocrConverter struct {
	Hocr       *html.Node
	Namespace  string
	BoxPattern *regexp.Regexp
}

type Args struct {
	HocrFile      string
	InputImage    string
	OutputPdfFile string
}

func NewHocrConverter(hocrFileName string) (*HocrConverter, error) {
	hocr := &HocrConverter{
		BoxPattern: regexp.MustCompile(`bbox((\s+\d+){4})`),
	}
	if hocrFileName != "" {
		err := hocr.ParseHocr(hocrFileName)
		if err != nil {
			return nil, err
		}
	}
	return hocr, nil
}

func (h *HocrConverter) ParseHocr(hocrFileName string) error {
	file, err := os.Open(hocrFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return err
	}

	h.Hocr = doc

	var findNamespace func(*html.Node)
	findNamespace = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "html" {
			for _, attr := range n.Attr {
				if strings.HasPrefix(attr.Key, "xmlns") {
					h.Namespace = attr.Val
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findNamespace(c)
		}
	}
	findNamespace(doc)

	return nil
}

func (h *HocrConverter) ElementCoordinates(element *html.Node) (int, int, int, int, error) {
	var out = [4]int{0, 0, 0, 0}
	var err error
	// this level of indentation is making me wish to delete this repo
	if element.Attr != nil {
		for _, attr := range element.Attr {
			if attr.Key == "title" {
				matches := h.BoxPattern.FindStringSubmatch(attr.Val)
				if len(matches) > 0 {
					coords := strings.Fields(matches[1])
					for i, coord := range coords {
						out[i], err = strconv.Atoi(coord)
						if err != nil {
							return 0, 0, 0, 0, err
						}
					}
				}
			}
		}
	}
	return out[0], out[1], out[2], out[3], nil
}

func (h *HocrConverter) ToString() string {
	var text strings.Builder
	var getString func(*html.Node)
	getString = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getString(c)
		}
	}
	if h.Hocr != nil {
		body := findBody(h.Hocr)
		if body != nil {
			getString(body)
		}
	}
	return text.String()
}

func findBody(n *html.Node) *html.Node {
	var body *html.Node
	var findBodyNode func(*html.Node)
	findBodyNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBodyNode(c)
		}
	}
	findBodyNode(n)
	return body
}

func main() {
	var args Args

	flag.StringVar(&args.HocrFile, "hocr", "", "The HOCR file to process")
	flag.StringVar(&args.InputImage, "image", "", "The image file to process")
	flag.StringVar(&args.OutputPdfFile, "pdf", "", "The output PDF file")

	flag.Parse()

	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	hocr, err := NewHocrConverter(args.HocrFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when creating HOCR object: %v\n", err)
		os.Exit(1)
	}
	hocr.ToPDF(args.InputImage, args.OutputPdfFile)
}

func (h *HocrConverter) ToPDF(imageFileName, outFileName string) error {
	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.AddPage()
	width, height, err := h.getImageDimensions(imageFileName)
	if err != nil {
		return fmt.Errorf("Error getting image dimensions from file '%s': %v", imageFileName, err)
	}

	pdf.ImageOptions(imageFileName, 0, 0, width, height, false, gofpdf.ImageOptions{}, 0, "")
	text := h.ToString()

	pdf.MoveTo(0, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(width, 10, text, "", 0, "CM", false, 0, "")
	err = pdf.OutputFileAndClose(outFileName)
	if err != nil {
		return fmt.Errorf("Error when writing to output PDF file '%s': %v", outFileName, err)
	}
	return nil
}

func (h *HocrConverter) getImageDimensions(imageFileName string) (float64, float64, error) {
	file, err := os.Open(imageFileName)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return float64(img.Width), float64(img.Height), nil
}
