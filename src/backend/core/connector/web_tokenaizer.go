package connector

import (
	"bytes"
	"fmt"
	"github.com/gocolly/colly/v2"
	"golang.org/x/net/html"
	"io"
)

var skipTag = map[string]bool{
	"script": true,
	"style":  true,
	"meta":   true,
	"link":   true,
	"a":      true,
	"li":     true,
	"ui":     true,
}

type webTokenizer struct {
}

func (c *webTokenizer) tokenizer(r *colly.Response) {
	tokenizer := html.NewTokenizer(bytes.NewBuffer(r.Body))

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()
		if tokenizer.Err() == io.EOF {
			break
		}
		if _, ok := skipTag[token.Data]; ok {
			if tokenType == html.StartTagToken {

			}
			continue
		}
		if tokenType.String() == "Text" {
			//		fmt.Println(tokenType.String(), "", token.Data)
			fmt.Println(token.String())
		}
	}
	fmt.Println("--")
}
