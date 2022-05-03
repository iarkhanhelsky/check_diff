package tools

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestRegistry_Install(t *testing.T) {
	assert := assert.New(t)
	r := newRegistry(t.TempDir(), &noopUnpacker{}, zap.NewNop().Sugar())

	var handler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("abc"))
	}
	ts := httptest.NewServer(handler)
	defer ts.Close()

	b := &Binary{
		Name:    "foo",
		Path:    "x",
		DstFile: "x",
		Targets: map[TargetTuple]TargetSource{
			Any: {
				Urls: []string{ts.URL},
				MD5:  fmt.Sprintf("%x", md5.Sum([]byte("abc"))),
			},
		}}
	err := r.Install(b)
	assert.NoError(err)
}

func TestRegistry_install(t *testing.T) {
	assert := assert.New(t)
	r := newRegistry("test", &noopUnpacker{}, zap.NewNop().Sugar())
	b := &Binary{Name: "foo", Path: "bar"}
	r.install(b)
	assert.Equal(b.executable, "test/foo/bar")
}

func TestRegistry_binHome(t *testing.T) {
	assert := assert.New(t)
	r := newRegistry("test", &noopUnpacker{}, zap.NewNop().Sugar())
	assert.Equal(r.binHome(&Binary{Name: "foo"}), "test/foo")
	assert.Equal(r.binHome(&Binary{Name: "foo"}, "bar"), "test/foo/bar")
}

func TestRegistry_ensureBinary(t *testing.T) {
	// Should not be called
	httpFail := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - internal error"))
	}

	writeBytes := func(bytes []byte) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			writer.Write(bytes)
		}
	}

	testCases := map[string]struct {
		binary        *Binary
		remoteHandler http.HandlerFunc
		vendorSetup   func(bin *Binary, vendorDir string) error
		err           string
		errf          string
	}{
		"http error": {
			binary: &Binary{
				Name: "test-binary",
				Targets: map[TargetTuple]TargetSource{
					Any: {
						Urls: []string{"http://test.org/a"},
						MD5:  fmt.Sprintf("%x", md5.Sum([]byte("abc"))),
					},
				},
			},
			remoteHandler: httpFail,
			vendorSetup: func(bin *Binary, vendorDir string) error {
				return nil
			},
			err: "fetch failed; binary = test-binary: download failed: download error: 500 - internal error",
		},
		"mkdir failed": {
			binary: &Binary{
				Name: "test-binary",
				Targets: map[TargetTuple]TargetSource{
					Any: {
						Urls: []string{"http://test.org/a"},
						MD5:  fmt.Sprintf("%x", md5.Sum([]byte("abc"))),
					},
				},
				Path:    "x",
				DstFile: "x",
			},
			vendorSetup: func(bin *Binary, vendorDir string) error {
				data := []byte("abc")
				if err := os.MkdirAll(path.Join(vendorDir), 0755); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path.Join(vendorDir, bin.Name), data, 0755); err != nil {
					return err
				}
				return nil
			},
			remoteHandler: httpFail,
			errf:          "mkdir %s/test-binary: not a directory",
		},
		"successful download dist": {
			binary: &Binary{
				Name: "test-binary",
				Targets: map[TargetTuple]TargetSource{
					Any: {
						Urls: []string{"http://test.org/a"},
						MD5:  fmt.Sprintf("%x", md5.Sum([]byte("abc"))),
					},
				},
				Path:    "x",
				DstFile: "x",
			},
			remoteHandler: writeBytes([]byte("abc")),
			vendorSetup: func(bin *Binary, vendorDir string) error {
				data := []byte("abc")
				if err := os.MkdirAll(path.Join(vendorDir, bin.Name), 0755); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path.Join(vendorDir, bin.Name, bin.Path), data, 0755); err != nil {
					return err
				}
				return nil
			},
		},
		"uptodate": {
			binary: &Binary{
				Name: "test-binary",
				Targets: map[TargetTuple]TargetSource{
					Any: {
						Urls: []string{"http://test.org/a"},
						MD5:  fmt.Sprintf("%x", md5.Sum([]byte("abc"))),
					},
				},
				DstFile: "x",
				Path:    "x",
			},
			remoteHandler: httpFail,
			vendorSetup: func(bin *Binary, vendorDir string) error {
				data := []byte("abc")
				if err := os.MkdirAll(path.Join(vendorDir, bin.Name), 0755); err != nil {
					return err
				}
				if err := os.MkdirAll(path.Join(vendorDir, bin.Name, "_dist"), 0755); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path.Join(vendorDir, bin.Name, bin.Path), data, 0755); err != nil {
					return err
				}
				if err := ioutil.WriteFile(path.Join(vendorDir, bin.Name, "_dist", bin.DstFile), data, 0755); err != nil {
					return err
				}

				digest, _ := bin.digest()
				m := manifest{
					dir:          path.Join(vendorDir, bin.Name),
					binaryDigest: digest,
					logger:       zap.NewNop().Sugar(),
				}
				return m.saveManifest()
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			vendorDir := t.TempDir()
			registry := newRegistry(vendorDir, &noopUnpacker{}, zap.NewNop().Sugar())

			ts := httptest.NewServer(tc.remoteHandler)
			defer ts.Close()

			binary := tc.binary
			target, _ := binary.Targets[Any]
			target.Urls = []string{ts.URL + "/a"}
			binary.Targets[Any] = target

			assert.NoError(tc.vendorSetup(binary, vendorDir), "vendor setup failed")

			err := registry.ensureBinary(tc.binary)
			if tc.err != "" {
				assert.EqualError(err, tc.err)
			} else if tc.errf != "" {
				assert.EqualError(err, fmt.Sprintf(tc.errf, vendorDir))
			} else {
				assert.NoError(err)
				assert.Equal(registry.binHome(binary, binary.Path), binary.executable)
			}
		})
	}
}
