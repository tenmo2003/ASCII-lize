# Image/Video to ASCII Converter

This is a simple command line tool that converts images and videos to ASCII art.

## Usage

```
Image to ASCII Converter

Usage:
  ascii-lize [options] <image-path>

Options:
  -targeted-width int
        The width of the output that the program will try to fit the image in by characters (default 100)
  -output string
        Name of the file the program will write the output to
  -charset string
        Character set to use (default, default-ascii, ascii-extended, ascii-block, ascii-detailed, ascii-classic, ascii-classic-extended, ascii-basic, ascii-basic-extended, ascii-full, ascii-full-extended, ascii-half, ascii-half-extended, ascii-double, ascii-double-extended, ascii-triple, ascii-triple-extended) (default "default")
  -space-density int
        Space density (default to 1)
  -help
        Show help information

Examples:
  ascii-lize image.jpg
  ascii-lize -targeted-width 80 -output result.txt ~/image.jpg
  ascii-lize -output out/output.txt -charset blocks ./image.jpg
  ascii-lize -charset detailed -targeted-width 120 image.jpg
```

## TODO

- [x] Add gif support
- [ ] Add video support
