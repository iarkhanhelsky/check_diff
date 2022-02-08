package downloader

type Interface interface {
	Download(dstFolder string) error
}

var Empty = make([]Interface, 0)

type Handler func(path string) error
