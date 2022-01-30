package kubelinter

import (
	"github.com/iarkhanhelsky/check_diff/pkg/downloader"
)

func (linter *KubeLinter) Downloads() []downloader.Interface {
	return []downloader.Interface{
		downloader.NewZipDownloader(linter.handleDownload, "kube-linter.exe",
			"877db4eb2825d3237e4e0aec40ca15c0",
			"c767e757109c79df9f51b0691d43da4f80e5dbf5b956929389af0face19719f4",
			"https://github.com/stackrox/kube-linter/releases/download/0.2.5/kube-linter-windows.zip"),
	}
}
