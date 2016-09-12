package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const (
	baseURL = "http://www.oursuperadventure.com/"
)

func main() {
	page := 265
	for i := page; i > 1; i-- { // 1 is a special case here, diff url
		fromURL := baseURL + "comic/page/" + strconv.Itoa(i) + "/"
		getImage(fromURL, i)
	}
	getImage(baseURL+"comic/", 1) // this handles just first page
}

func getImage(fromURL string, page int) {
	onPageCounter := 0
	resp, err := http.Get(fromURL)
	mustNotErr(err)

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.SelfClosingTagToken:
			t := z.Token()

			if t.Data == "img" {
				for _, v := range t.Attr {
					if v.Key == "src" {
						onPageCounter++
						URL := "http:" + v.Val

						filename := ""
						if strings.HasSuffix(v.Val, ".jpg") {
							filename = fmt.Sprintf("oursuperadventure.com-page%04d-image%04d.jpg", page, onPageCounter)
						} else if strings.HasSuffix(v.Val, ".png") {
							filename = fmt.Sprintf("oursuperadventure.com-page%04d-image%04d.png", page, onPageCounter)
						} else if strings.HasSuffix(v.Val, ".gif") {
							filename = fmt.Sprintf("oursuperadventure.com-page%04d-image%04d.gif", page, onPageCounter)
						}

						log.Println("downloading:", URL, "to:", filename)

						out, err := os.Create(filename)
						mustNotErr(err)
						defer out.Close()

						resp, err := http.Get(URL)
						mustNotErr(err)
						defer resp.Body.Close()

						_, err = io.Copy(out, resp.Body)
						mustNotErr(err)
					}
				}
			}
		}
	}

	resp.Body.Close()
}

func mustNotErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
