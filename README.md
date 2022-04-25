# Convert a Markdown file to PNG

`go run ./main.go example.md example.png`

Requires the `wkhtmltoimage` binary in PATH.

```
$ go run ./main.go --help
Usage: main [--title STR] [--width PX] MD-FILE OUT-FILE

Positional arguments:
  MD-FILE                markdown file to process
  OUT-FILE               PNG file name to write to

Options:
  --title STR            title at top of document, ignored if empty
  --width PX             output image width [default: 384]
  --help, -h             display this help and exit
```

Example [input](example.md) [output](example.png).
