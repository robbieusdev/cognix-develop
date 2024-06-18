package oauth

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MicrosoftGetURL(t *testing.T) {
	cfg := MicrosoftConfig{}
	err := utils.ReadConfig(&cfg)
	assert.NoError(t, err)
	oauth := NewMicrosoft(&cfg)
	url, err := oauth.GetAuthURL(context.Background(), "http://localhost:8080/", "")
	t.Log(url)
}
