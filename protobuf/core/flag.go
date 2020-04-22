package core

type flag struct {
	value  bool
	parent Flagged
}

func (f flag) Get() bool {
	return f.value
}

func (f *flag) Set(value bool) error {
	old := f.value
	f.value = value
	if err := f.parent.validateFlag(f); err != nil {
		f.value = old
		return err
	}
	return nil
}

func (f flag) Parent() Flagged {
	return f.parent
}

type Flagged interface {
	validateFlag(*flag) error
}
