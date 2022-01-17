package k8s_kubelint

import "github.com/iarkhanhelsky/check_diff/pkg/core"

func (linter *K8sKubeLint) Downloads() []core.Downloader {
	return []core.Downloader{
		core.NewZipDownloader(linter.handleDownload, "kube-linter",
			"58b4a9b8d55c1997c866471c14bbcb3a",
			"dd75ba0a35db6ee12f36e8e36dac0e3e361e9a43103196962da86458092f9ab7",
			"https://github.com/stackrox/kube-linter/releases/download/0.2.5/kube-linter-darwin.zip"),
	}
}
