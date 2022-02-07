package stdout

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestContext(t *testing.T) {
	ctx := newCachedFileContext(3)

	type testcase struct {
		file string
		line int

		expect       []string
		expectOffset int
	}
	testcases := []testcase{
		{file: "testdata/a.txt", line: 5, expect: []string{"a02", "a03", "a04", "a05", "a06", "a07", "a08"}, expectOffset: 3},
		{file: "testdata/a.txt", line: 1, expect: []string{"a01", "a02", "a03", "a04"}, expectOffset: 0},
		{file: "testdata/b.txt", line: 1, expect: []string{"b01", "b02", "b03", "b04"}, expectOffset: 0},
	}

	for _, tc := range testcases {

		t.Run(fmt.Sprintf("%s-%d", tc.file, tc.line), func(t *testing.T) {
			assert := assert.New(t)
			lines, offset, err := ctx.readContext(tc.file, tc.line)
			assert.NoError(err)
			assert.Equal(tc.expect, lines)
			assert.Equal(tc.expectOffset, offset)
		})
	}
}
