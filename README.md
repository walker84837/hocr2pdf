# hocr2pdf

> Convert HOCR data into PDFs with integrated image support

hocr2pdf is a tool for converting HOCR (HTML-based OCR) documents into PDF format, integrating text with associated images. This tool is ideal for users needing to create searchable PDFs from OCR data and images, such as scanned documents or annotated text.

## Installing / Getting started

To get started with `hocr2pdf`, you'll need to have Go installed on your machine. The following instructions assume you have Go set up.

1. **Clone the repository**:
   ``` console
   $ git clone https://winlogon.ddns.net/winlogon/hocr2pdf.git
   $ cd hocr2pdf/
   ```

2. **Build the project**:
   ``` console
   $ make
   ```

3. **Run the application**:
   ``` console
   $ ./hocr2pdf -hocr path/to/your.hocr -image path/to/your-image.png -pdf output.pdf
   ```

   This command generates a PDF named `output.pdf` from the HOCR file and image provided.

### Initial Configuration

No additional initial configuration is required beyond the standard Go setup and dependencies.

## Developing

To contribute to `hocr2pdf`, clone the repository:
``` console
$ git clone https://winlogon.ddns.net/winlogon/hocr2pdf.git
$ cd hocr2pdf/
```

### Building

After making code changes, you can build the project with:

``` console
$ make
```

This command compiles the source code into an executable named `hocr2pdf`.

### Deploying / Publishing

To deploy or distribute the project, simply distribute the built binary. For publishing on a server, ensure the executable is included in your deployment package.

## Features

- **Convert HOCR to PDF**: Takes HOCR data and an image file to produce a PDF.
- **Bounding box parsing**: Extracts text coordinates from HOCR data for accurate placement.
- **Text extraction**: Converts HOCR document text into a plain text string for use in PDFs.

## Configuration

The application uses command-line arguments for configuration:

| Argument   | Type    | Default | Description                                              | Example                                                 |
|------------|---------|---------|----------------------------------------------------------|---------------------------------------------------------|
| `-hocr`    | String  | `""`    | Path to the HOCR file to process.                       | `./hocr2pdf -hocr myfile.hocr -image myimage.png -pdf output.pdf` |
| `-image`   | String  | `""`    | Path to the image file to be included in the PDF.        | `./hocr2pdf -hocr myfile.hocr -image myimage.png -pdf output.pdf` |
| `-pdf`     | String  | `""`    | Path to the output PDF file.                             | `./hocr2pdf -hocr myfile.hocr -image myimage.png -pdf output.pdf` |
| `-overwrite` | Boolean | `false` | If `true`, will overwrite the output PDF file if it already exists. | `./hocr2pdf -hocr myfile.hocr -image myimage.png -pdf output.pdf -overwrite` |

## Contributing

We welcome contributions to improve hocr2pdf. Please fork the repository, make your changes, and submit a pull request.

## Links

- Repository: [https://winlogon.ddns.net/winlogon/hocr2pdf/](https://winlogon.ddns.net/winlogon/hocr2pdf/)
- Issue tracker: [https://winlogon.ddns.net/winlogon/hocr2pdf/issues](https://winlogon.ddns.net/winlogon/hocr2pdf/issues)
  - For sensitive bugs or security vulnerabilities, please contact me at `@winlogon.exe:matrix.org` directly.

## Licensing

The code in this project is licensed under the [BSD 3-Clause](LICENSE.md).
