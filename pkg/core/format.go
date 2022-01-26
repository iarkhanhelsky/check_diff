package core

type Format string

const (
	STDOUT      Format = "stdout"
	Codeclimate Format = "codeclimate"
	Gitlab      Format = "gitlab"
	Phabricator Format = "phabricator"
)
