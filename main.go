package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alexflint/go-arg"
	extensions "github.com/moyen-blog/goldmark-extensions"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type flags struct {
	Title        string `arg:"--title" default:"" help:"title at top of document, ignored if empty" placeholder:"STR"`
	Width        uint   `arg:"--width" default:"384" help:"output image width" placeholder:"PX"`
	MarkdownFile string `arg:"positional,required" help:"markdown file to process" placeholder:"MD-FILE"`
	PNGFile      string `arg:"positional,required" help:"PNG file name to write to" placeholder:"OUT-FILE"`
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

func htmlToImg(htmlPath, pngPath string, width int) error {
	cmd := exec.Command("wkhtmltoimage",
		"--width", fmt.Sprint(width),
		"--disable-smart-width",
		"--disable-local-file-access",
		"--no-images",
		"--format", "png", htmlPath, pngPath)
	return cmd.Run()
}

func main() {
	f := flags{}
	arg.MustParse(&f)

	htmlPath, err := mdToHTML(f.MarkdownFile, f.Title)
	if err != nil {
		panic(err)
	}
	defer os.Remove(htmlPath)

	err = htmlToImg(htmlPath, f.PNGFile, int(f.Width))
	if err != nil {
		panic(err)
	}
}
