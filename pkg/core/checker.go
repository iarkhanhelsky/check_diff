package core

type Checker interface {
	Downloads() []Downloader
	Setup()
	Check(ranges []LineRange) ([]Issue, error)
}
