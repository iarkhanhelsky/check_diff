package kubelinter

import "github.com/iarkhanhelsky/check_diff/pkg/core"

func (linter *KubeLinter) Downloads() []core.Downloader {
	return []core.Downloader{
		core.NewZipDownloader(linter.handleDownload, "kube-linter",
			"05c8a6c57cb6d84ebae6a09efc9f46c2",
			"a858572d7b673574855ce8cb84476b6d4690e79905d3fbaf303fe0a70eb8798e",
			"https://github.com/stackrox/kube-linter/releases/download/0.2.5/kube-linter-linux.zip"),
	}
}
