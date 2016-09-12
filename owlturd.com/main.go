package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	baseURL = "http://owlturd.com/search/COMIC"
	timeout = 15 * time.Second
)

func main() {
	page := 26
	for i := page; i > 1; i-- { // 1 is a special case here, diff url
		fromURL := baseURL + "/page/" + strconv.Itoa(i)
		getImage(fromURL, i)
	}
	getImage(baseURL, 1) // this handles just first page
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
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "img" {
				src := ""
				class := ""
				for _, v := range t.Attr {
					if v.Key == "class" {
						class = v.Val
					}
					if v.Key == "src" {
						src = v.Val
					}
				}
				if class == "photo-md" && src != "" {
					onPageCounter++
					URL := src
					//_, err := exec.Command("wget", URL).CombinedOutput()

					filename := ""
					if strings.HasSuffix(src, ".jpg") {
						filename = fmt.Sprintf("page%04d-image%04d.jpg", page, onPageCounter)
					} else if strings.HasSuffix(src, ".png") {
						filename = fmt.Sprintf("page%04d-image%04d.png", page, onPageCounter)
					} else if strings.HasSuffix(src, ".gif") {
						filename = fmt.Sprintf("page%04d-image%04d.gif", page, onPageCounter)
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
		case tt == html.SelfClosingTagToken:
			t := z.Token()

			if t.Data == "img" {
				log.Println(t.Data)
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
