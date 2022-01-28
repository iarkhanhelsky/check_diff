package kubelinter

import "github.com/iarkhanhelsky/check_diff/pkg/core"

func (linter *KubeLinter) Downloads() []core.Downloader {
	return []core.Downloader{
		core.NewZipDownloader(linter.handleDownload, "kube-linter.exe",
			"877db4eb2825d3237e4e0aec40ca15c0",
			"c767e757109c79df9f51b0691d43da4f80e5dbf5b956929389af0face19719f4",
			"https://github.com/stackrox/kube-linter/releases/download/0.2.5/kube-linter-windows.zip"),
	}
}
