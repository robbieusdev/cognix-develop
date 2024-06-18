package ai

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStaticChunking struct {
	Source        string
	ExpectedCount int
	Expected      []string
}

var testData = []testStaticChunking{
	{
		Source:        "foo",
		ExpectedCount: 1,
		Expected:      []string{"foo"},
	},
	{
		Source:        "123456789012345678901234567890123456789012345",
		ExpectedCount: 7,
		Expected: []string{"1234567890",
			"8901234567",
			"5678901234",
			"2345678901",
			"9012345678",
			"6789012345",
			"345",
		},
	},
}

func Test_StaticChunkingText(t *testing.T) {

	service := NewStaticChunking(&ChunkingConfig{
		Strategy:           StrategyStatic,
		StaticChunkSize:    10,
		StaticChunkOverlap: 3,
	})
	for i, test := range testData {
		t.Run(fmt.Sprintf("%d", i),
			func(t *testing.T) {
				result, err := service.Split(context.Background(), test.Source)
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectedCount, len(result))
				for n, chunk := range result {
					assert.Equal(t, test.Expected[n], chunk)
				}
			})
	}
}
