package tools

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

func TestHttpDownloader_Download(t *testing.T) {
	testCases := map[string]struct {
		md5    string
		sha256 string
		data   []byte
		cached []byte

		handler http.HandlerFunc
		error   string
	}{
		"md5": {
			md5:  fmt.Sprintf("%x", md5.Sum([]byte("abc"))),
			data: []byte("abc"),
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte("abc"))
			},
		},

		"wrong md5": {
			md5:  fmt.Sprintf("%x", md5.Sum([]byte("abx"))),
			data: []byte("abc"),
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte("abc"))
			},
			error: "md5sum mismatch",
		},

		"sha256": {
			sha256: fmt.Sprintf("%x", sha256.Sum256([]byte("abc"))),
			data:   []byte("abc"),
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte("abc"))
			},
		},

		"wrong sha256": {
			sha256: fmt.Sprintf("%x", sha256.Sum256([]byte("abx"))),
			data:   []byte("abc"),
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte("abc"))
			},
			error: "sha256sum mismatch",
		},

		"cached uptodate": {
			sha256: fmt.Sprintf("%x", sha256.Sum256([]byte("abc"))),
			data:   []byte("abc"),
			cached: []byte("abc"),
			handler: func(writer http.ResponseWriter, request *http.Request) {
				// should not be called
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("500 - internal error"))
			},
		},

		"cached old": {
			sha256: fmt.Sprintf("%x", sha256.Sum256([]byte("abc"))),
			data:   []byte("abc"),
			cached: []byte("abx"),
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte("abc"))
			},
		},

		"cached no sums": {
			data:   []byte("abc"),
			cached: []byte("abx"),
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte("abc"))
			},
		},

		"http error": {
			sha256: fmt.Sprintf("%x", sha256.Sum256([]byte("abc"))),
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("500 - internal error"))
			},

			error: "download error: 500 - internal error",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(tc.handler)
			defer ts.Close()

			assert := assert.New(t)
			vendorDir := t.TempDir()
			downloader := newHTTPDownloader(tc.md5, tc.sha256, ts.URL+"/a.bin")
			dstFile := path.Join(vendorDir, "a.bin")

			if len(tc.cached) > 0 {
				err := ioutil.WriteFile(dstFile, tc.data, fs.FileMode(int(0777)))
				assert.NoError(err)
			}

			err := downloader.Download(vendorDir)
			if tc.error != "" {
				assert.EqualError(err, tc.error)
				return
			}
			assert.NoError(err)

			bytes, err := ioutil.ReadFile(dstFile)
			assert.NoError(err)
			assert.Equal(tc.data, bytes)
		})
	}
}
