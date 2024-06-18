package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"github.com/shopspring/decimal"
	"testing"
)

func TestWeb_Execute(t *testing.T) {
	web, err := NewWeb(
		&model.Connector{
			ID:   decimal.NewFromInt(1),
			Name: "web test",
			Type: "web",
			ConnectorSpecificConfig: model.JSONMap{
				"url1": "https://help.collaboard.app/",
				"url2": "https://apidog.com/blog/openapi-specification/",
				"url":  "https://developer.apple.com/documentation/visionos/improving-accessibility-support-in-your-app",
				"url3": "https://developer.mozilla.org/en-US/docs/Learn/HTML/Introduction_to_HTML",
			},
			DocsMap: make(map[string]*model.Document),
		})
	if err != nil {
		t.Log(err.Error())
		t.Fatal(err)
	}
	conn := web.Execute(context.Background(), nil)

	for res := range conn {
		t.Log(res.URL)
	}
}
