package command

type None struct {
}

var _ Command = &None{}

func (*None) Run() error {
	return nil
}
