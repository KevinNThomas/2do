package main

import (
	"io"
	"log"
	"os"

	"github.com/namsral/flag"
	"github.com/prologic/bitcask"
)

var (
	db *bitcask.Bitcask
)

func main() {
	var (
		dbpath               string
		bind                 string
		colorTheme           string
		colorPageBackground  string
		colorInputBackground string
		colorForeground      string
		colorCheckMark       string
		colorXMark           string
		colorLabel           string
	)

	flag.StringVar(&dbpath, "dbpath", "todo.db", "Database path")
	flag.StringVar(&bind, "bind", "0.0.0.0:8000", "[int]:<port> to bind to")
	flag.Parse()

	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "COLOR", 0)
	fs.StringVar(&colorTheme, "theme", "dracula", "color theme of the todo list, or 'custom'")
	fs.StringVar(&colorPageBackground, "pagebackground", "282a36", "page background color")
	fs.StringVar(&colorInputBackground, "inputbackground", "44475a", "input boxes color")
	fs.StringVar(&colorForeground, "foreground", "f8f8f2", "text color")
	fs.StringVar(&colorCheckMark, "check", "50fa7b", "check mark color")
	fs.StringVar(&colorXMark, "x", "ff5555", "x mark color")
	fs.StringVar(&colorLabel, "label", "ff79c6", "label color")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	db, err = bitcask.Open(dbpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	selectColorTheme(colorTheme, colorPageBackground, colorInputBackground, colorForeground,
		colorCheckMark, colorXMark, colorLabel)

	newServer(bind).listenAndServe()
}

func selectColorTheme(colorTheme string, colorPageBackground string, colorInputBackground string,
	colorForeground string, colorCheckMark string, colorXMark string, colorLabel string) {
	if colorTheme == "custom" {
		customThemeFile, err := os.OpenFile("./static/css/color-theme.css", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer customThemeFile.Close()

		_, err = customThemeFile.WriteString(":root {" +
			"\n\t--page-background: #" + colorPageBackground + ";" +
			"\n\t--input-background: #" + colorInputBackground + ";" +
			"\n\t--foreground: #" + colorForeground + ";" +
			"\n\t--check: #" + colorCheckMark + ";" +
			"\n\t--x: #" + colorXMark + ";" +
			"\n\t--label: #" + colorLabel + ";" +
			"\n}")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		from, err := os.Open("./static/color-themes/" + colorTheme + ".css")
		if err != nil {
			log.Fatal(err)
		}
		defer from.Close()

		to, err := os.OpenFile("./static/css/color-theme.css", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			log.Fatal(err)
		}
	}
}
