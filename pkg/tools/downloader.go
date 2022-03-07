package tools

type Downloader interface {
	Download(dstFolder string) error
}

var Empty = make([]Downloader, 0)

type Handler func(path string) error
