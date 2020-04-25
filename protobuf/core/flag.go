package core

type Flag struct {
	value  bool
	parent Flagged
}

func (f Flag) Get() bool {
	return f.value
}

func (f *Flag) Set(value bool) error {
	old := f.value
	f.value = value
	if err := f.parent.validateFlag(f); err != nil {
		f.value = old
		return err
	}
	return nil
}

func (f Flag) Parent() Flagged {
	return f.parent
}

type Flagged interface {
	validateFlag(*Flag) error
}
