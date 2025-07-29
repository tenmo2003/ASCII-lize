package main

import (
	"ascii-lize/internal/ascii"
	"ascii-lize/internal/config"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	targetedWidth := flag.Int("targeted-width", 100, "The width of the output that the program will try to fit the image in by characters")
	outputFileName := flag.String("output", "", "Name of the file the program will write the output to")
	characterSet := flag.String("charset", "default", fmt.Sprintf("Character set to use (%s)", strings.Join(ascii.GetAvailableCharacterSets(), ", ")))
	spaceDensity := flag.Int("space-density", 1, "Space density (default to 1)")
	showHelp := flag.Bool("help", false, "Show help information")

	flag.Parse()

	if *showHelp {
		showUsage()
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: Please specify the path to image using argument")
		showUsage()
		os.Exit(1)
	}

	cfg := config.NewConfig()
	cfg.ImagePath = args[0]
	cfg.TargetedWidth = *targetedWidth
	cfg.OutputPath = *outputFileName
	cfg.CharacterSet = *characterSet
	cfg.SpaceDensity = *spaceDensity

	if err := cfg.Validate(); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		os.Exit(1)
	}

	if ascii.GetCharacterSet(cfg.CharacterSet) == ascii.DefaultCharacterSet && cfg.CharacterSet != "default" {
		fmt.Printf("Warning: Unknown character set '%s', using default\n", cfg.CharacterSet)
		cfg.CharacterSet = "default"
	}

	converter := ascii.NewConverter()
	if err := converter.ConvertToASCII(cfg); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println("Image to ASCII Converter")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ascii-lize [options] <image-path>")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ascii-lize image.jpg")
	fmt.Println("  ascii-lize -targeted-width 80 -output result.txt image.jpg")
	fmt.Println("  ascii-lize -output stdout -charset blocks image.jpg")
	fmt.Println("  ascii-lize -charset detailed -targeted-width 120 image.jpg")
	fmt.Println()
	fmt.Printf("Available character sets: %s\n", strings.Join(ascii.GetAvailableCharacterSets(), ", "))
}
