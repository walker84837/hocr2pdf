package main

import (
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

// HocrConverter converts HOCR data into different formats.
type HocrConverter struct {
	Hocr       *html.Node
	Namespace  string
	boxPattern *regexp.Regexp
}

// New creates a new HocrConverter instance. If hocrFileName is provided,
// it parses the HOCR file at that path.
func newHocrConverter(hocrFileName string) (*HocrConverter, error) {
	hocr := &HocrConverter{
		boxPattern: regexp.MustCompile(`bbox((\s+\d+){4})`),
	}
	if hocrFileName != "" {
		if err := hocr.parseHocr(hocrFileName); err != nil {
			return nil, fmt.Errorf("failed to parse HOCR file: %w", err)
		}
	}
	return hocr, nil
}

// parseHocr parses the HOCR file and sets the document node and namespace.
func (h *HocrConverter) parseHocr(hocrFileName string) error {
	file, err := os.Open(hocrFileName)
	if err != nil {
		return fmt.Errorf("failed to open HOCR file: %w", err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		return fmt.Errorf("failed to parse HOCR file: %w", err)
	}

	h.Hocr = doc
	h.findNamespace(doc)
	return nil
}

// findNamespace finds and sets the XML namespace from the HOCR document.
func (h *HocrConverter) findNamespace(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "html" {
		for _, attr := range n.Attr {
			if strings.HasPrefix(attr.Key, "xmlns") {
				h.Namespace = attr.Val
				return
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		h.findNamespace(c)
	}
}

// Coordinates represents the bounding box coordinates of an element.
type Coordinates struct {
	X1, Y1 int // Top-left corner (x, y)
	X2, Y2 int // Bottom-right corner (x, y)
}

// ElementCoordinates extracts the bounding box coordinates from an HTML element.
func (h *HocrConverter) ElementCoordinates(element *html.Node) (Coordinates, error) {
	var coords Coordinates

	if element.Attr == nil {
		return coords, nil
	}

	for _, attr := range element.Attr {
		if attr.Key == "title" {
			if err := h.parseBoundingBox(attr.Val, &coords); err != nil {
				return coords, fmt.Errorf("failed to parse bounding box: %w", err)
			}
			break
		}
	}

	return coords, nil
}

// parseBoundingBox parses the bounding box string and assigns coordinates.
func (h *HocrConverter) parseBoundingBox(value string, coords *Coordinates) error {
	matches := h.boxPattern.FindStringSubmatch(value)
	if len(matches) == 0 {
		return nil
	}

	values := strings.Fields(matches[1])
	if len(values) != 4 {
		return fmt.Errorf("invalid number of coordinates found")
	}

	for i, v := range values {
		val, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("failed to convert coordinate to integer: %w", err)
		}
		switch i {
		case 0:
			coords.X1 = val
		case 1:
			coords.Y1 = val
		case 2:
			coords.X2 = val
		case 3:
			coords.Y2 = val
		}
	}
	return nil
}

// ToString converts the HOCR document to a plain text string.
func (h *HocrConverter) ToString() string {
	var text strings.Builder
	h.extractText(h.Hocr, &text)
	return text.String()
}

// extractText recursively extracts text from HTML nodes.
func (h *HocrConverter) extractText(n *html.Node, text *strings.Builder) {
	if n.Type == html.TextNode {
		text.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		h.extractText(c, text)
	}
}

// ToPDF generates a PDF from the HOCR document and an image file.
func (h *HocrConverter) ToPDF(imageFileName, outFileName string, overwrite bool) error {
	if !overwrite {
		if _, err := os.Stat(outFileName); !os.IsNotExist(err) {
			return fmt.Errorf("output file '%s' already exists. Use overwrite flag to overwrite", outFileName)
		}
	}

	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.AddPage()

	width, height, err := h.getImageDimensions(imageFileName)
	if err != nil {
		return fmt.Errorf("error getting image dimensions from file '%s': %w", imageFileName, err)
	}

	pdf.ImageOptions(imageFileName, 0, 0, width, height, false, gofpdf.ImageOptions{}, 0, "")
	text := h.ToString()

	pdf.MoveTo(0, 0)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(width, 10, text, "", 0, "CM", false, 0, "")
	if err := pdf.OutputFileAndClose(outFileName); err != nil {
		return fmt.Errorf("error writing to output PDF file '%s': %w", outFileName, err)
	}
	return nil
}

// getImageDimensions retrieves the dimensions of the image file.
func (h *HocrConverter) getImageDimensions(imageFileName string) (float64, float64, error) {
	file, err := os.Open(imageFileName)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode image config: %w", err)
	}

	return float64(img.Width), float64(img.Height), nil
}
