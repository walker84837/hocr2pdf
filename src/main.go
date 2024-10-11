package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// Args holds command line arguments for the program.
type Args struct {
	HocrFile      string
	InputImage    string
	OutputPdfFile string
	Overwrite     bool
}

func main() {
	var args Args

	flag.StringVar(&args.HocrFile, "hocr", "", "The HOCR file to process.")
	flag.StringVar(&args.InputImage, "image", "", "The image file to process.")
	flag.StringVar(&args.OutputPdfFile, "pdf", "", "The output PDF file.")
	flag.BoolVar(&args.Overwrite, "overwrite", false, "Overwrite the output PDF file if it already exists.")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if args.HocrFile == "" || args.InputImage == "" || args.OutputPdfFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	hocr, err := newHocrConverter(args.HocrFile)
	if err != nil {
		logger.WithError(err).Error("Failed to create HOCR object")
		os.Exit(1)
	}

	if err := hocr.ToPDF(args.InputImage, args.OutputPdfFile, args.Overwrite); err != nil {
		logger.WithError(err).Error("Failed to generate PDF")
		os.Exit(1)
	}

	logger.Info("PDF generated successfully")
}
