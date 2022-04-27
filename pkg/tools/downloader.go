package tools

import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type httpDownloader struct {
	urls   []string
	md5    string
	sha256 string
}

func newHTTPDownloader(md5 string, sha256 string, urls ...string) *httpDownloader {
	return &httpDownloader{urls, md5, sha256}
}

func (d *httpDownloader) Download(dstFolder string) error {
	var allErrors error

	for _, url := range d.urls {
		dstFile := path.Join(dstFolder, path.Base(url))
		err := d.downloadFrom(url, dstFile)
		if err == nil {
			break
		}

		allErrors = multierr.Append(allErrors, err)
	}

	return allErrors
}

func (d *httpDownloader) downloadFrom(url string, outputFile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("download error: %s", string(bytes))
	}

	if err := checkMD5(bytes, d.md5); err != nil {
		return err
	}

	if err := checkSHA256(bytes, d.sha256); err != nil {
		return err
	}

	if err := ioutil.WriteFile(outputFile, bytes, fs.FileMode(int(0777))); err != nil {
		return err
	}

	return nil
}

func (d *httpDownloader) ensurePath() (string, error) {
	return "", nil
}

func checkMD5(bytes []byte, md5sum string) error {
	// "" means skip
	if len(md5sum) == 0 {
		return nil
	}

	hex := fmt.Sprintf("%x", md5.Sum(bytes))
	if hex != md5sum {
		return errors.New("md5sum mismatch")
	}

	return nil
}

func checkSHA256(bytes []byte, sha1sum string) error {
	// "" means skip
	if len(sha1sum) == 0 {
		return nil
	}

	hex := fmt.Sprintf("%x", sha256.Sum256(bytes))
	if hex != sha1sum {
		return errors.New("sha1sum mismatch")
	}

	return nil
}

func (d *httpDownloader) isUptodate(dstFile string) bool {
	if len(d.md5) == 0 && len(d.sha256) == 0 {
		// No way to check, always outdated
		return false
	}

	if _, err := os.Stat(dstFile); errors.Is(err, os.ErrNotExist) {
		return false
	}

	file, err := os.Open(dstFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: Failed to check %s: %v", dstFile, err)
		return false
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: Failed to check %s: %v", dstFile, err)
		return false
	}

	return checkMD5(bytes, d.md5) == nil && checkSHA256(bytes, d.sha256) == nil
}
