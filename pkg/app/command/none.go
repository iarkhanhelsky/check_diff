package command

// None command is needed when we need to display help section.
//
// This is handled by flags package and we effectively do not need to do
// anything.
type None struct {
}

var _ Command = &None{}

func (*None) Run() error {
	return nil
}
