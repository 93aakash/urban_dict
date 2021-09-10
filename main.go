package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/93aakash/urban_dict/models"
)

const (
	URL        = "https://api.urbandictionary.com/v0/define?term="
	LINE_WIDTH = 49
)

type Response struct {
	Entries []models.Def `json:"List"`
}

func main() {
	l := log.New(os.Stderr, "", 0)

	DB_PATH := os.Getenv("DB_PATH")
	_, err := os.Stat(DB_PATH)
	if os.IsNotExist(err) {
		l.Fatalln("The specified database doesn't exist")
	}

	db, err := models.InitDB("sqlite3", DB_PATH)
	if err != nil {
		l.Fatalln(err)
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	var query string
	switch os.Args[1] {
	case "-d", "--delete":
		query = strings.ToLower(strings.Join(os.Args[2:], " "))
		err := db.DeleteDef(query)
		if err != nil {
			l.Fatalln(err)
		}
		return
	case "-h", "--help":
		printUsage()
		return
	default:
		query = strings.ToLower(strings.Join(os.Args[1:], " "))
	}

	results := []models.Def{}
	if db.IfExists(query) {
		results, err = db.FetchDef(query)
		if err != nil {
			l.Fatalln(err)
		}
	} else {
		results, err = getDef(query)
		if err != nil {
			l.Fatalln(err)
		}
		err = db.InsertDef(results)
		if err != nil {
			l.Fatalln(err)
		}
	}
	InfoColor := "\033[1;34m%s\033[0m\n\n"
	fmt.Printf(InfoColor, results[0].Word)

	for i, r := range results {
		padding := "  "
		fmt.Printf("%d.", i+1)
		for _, line := range strings.Split(r.Definition, "\n") {
			fmt.Println(padding + wordWrap(line, LINE_WIDTH))
			padding = "    "
		}
		fmt.Println()
	}
}

func printUsage() {
	fmt.Printf("Usage:\n  urban_dict [OPTIONS] word\n\n")
	fmt.Println("Options:\n  -d, --delete\tdelete word from the database")
	fmt.Println("  -h, --help\tprint this help and exit\n")
}

func getDef(query string) ([]models.Def, error) {
	res, err := http.Get(URL + url.QueryEscape(query))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	r := &Response{}
	err = json.NewDecoder(res.Body).Decode(r)
	if len(r.Entries) == 0 {
		return nil, fmt.Errorf("No results found for \"%s\"", query)
	}

	replacer := strings.NewReplacer("T", " ", "Z", "")
	for i := range r.Entries {
		r.Entries[i].WrittenOn = replacer.Replace(r.Entries[i].WrittenOn)
	}
	return r.Entries, nil
}

func wordWrap(text string, lineWidth int) string {
	wrap := make([]byte, 0, len(text)+2*len(text)/lineWidth)
	eoLine := lineWidth
	inWord := false
	for i, j := 0, 0; ; {
		r, size := utf8.DecodeRuneInString(text[i:])
		if size == 0 && r == utf8.RuneError {
			r = ' '
		}
		if unicode.IsSpace(r) {
			if inWord {
				if i >= eoLine {
					wrap = append(wrap, '\n')
					wrap = append(wrap, []byte{' ', ' ', ' ', ' '}...)
					eoLine = len(wrap) + lineWidth
				} else if len(wrap) > 0 {
					wrap = append(wrap, ' ')
				}
				wrap = append(wrap, text[j:i]...)
			}
			inWord = false
		} else if !inWord {
			inWord = true
			j = i
		}
		if size == 0 && r == ' ' {
			break
		}
		i += size
	}
	return string(wrap)
}
