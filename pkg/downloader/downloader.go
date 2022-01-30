package downloader

type Interface interface {
	Download(dstFolder string) error
}

type Handler func(path string) error
