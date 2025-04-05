package main

import (
	"fmt"
	"image"
	"os"
	"os/exec"

	_ "image/png"

	"github.com/alexflint/go-arg"
	extensions "github.com/moyen-blog/goldmark-extensions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"jo-m.ch/go/gocatprint/pkg/printer"
	"jo-m.ch/go/gocatprint/pkg/simple"
)

type flags struct {
	Title        string `arg:"--title" default:"" help:"title at top of document, ignored if empty" placeholder:"STR"`
	MarkdownFile string `arg:"positional,required" help:"markdown file to process" placeholder:"MD-FILE"`
}

const header = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<style>
			body {
				margin: 0;
				font-family: sans-serif;
			}
			ul li.task {
				margin-left: -20px;
				list-style-type: none;
			}
		</style>
	</head>
	<body>
`

const footer = `
	</body>
</html>
`

var mdConverter = goldmark.New(
	goldmark.WithExtensions(extensions.TasklistExtension, extension.Strikethrough, extension.Table, extension.Linkify, extension.Typographer),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
)

func mdToHTML(mdPath, title string) (string, error) {
	src, err := os.ReadFile(mdPath)
	if err != nil {
		return "", err
	}

	if title != "" {
		src = append([]byte(fmt.Sprintf("# %s\n\n", title)), src...)
	}

	htmlFile, err := os.CreateTemp("", "*.html")
	if err != nil {
		return "", err
	}
	defer htmlFile.Close()

	htmlFile.Write([]byte(header))
	if err := mdConverter.Convert(src, htmlFile); err != nil {
		os.Remove(htmlFile.Name())
		return "", err
	}
	_, err = htmlFile.Write([]byte(footer))
	if err != nil {
		os.Remove(htmlFile.Name())
		return "", err
	}

	return htmlFile.Name(), nil
}

func htmlToImg(htmlPath string, width int) (image.Image, error) {
	pngFile, err := os.CreateTemp("", "*.html")
	if err != nil {
		return nil, err
	}
	defer os.Remove(pngFile.Name())
	pngFile.Close()

	cmd := exec.Command("wkhtmltoimage",
		"--width", fmt.Sprint(width),
		"--disable-smart-width",
		"--disable-local-file-access",
		"--no-images",
		"--format", "png", htmlPath, pngFile.Name())
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(pngFile.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func main() {
	f := flags{}
	arg.MustParse(&f)

	log.Logger = log.Logger.Level(zerolog.Disabled)

	htmlPath, err := mdToHTML(f.MarkdownFile, f.Title)
	if err != nil {
		panic(err)
	}
	defer os.Remove(htmlPath)

	img, err := htmlToImg(htmlPath, printer.PrintWidth)
	if err != nil {
		panic(err)
	}

	err = simple.Print(img, false)
	if err != nil {
		panic(err)
	}
}
