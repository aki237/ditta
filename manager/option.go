package manager

import "errors"

type Option struct {
	some string
	none bool
}

func NewOption() Option {
	return Option{some: "", none: true}
}

func (o Option) IsSome() bool {
	return !o.none
}

func (o Option) Unexpect() (string, error) {
	if o.IsSome() {
		return o.some, nil
	}
	return "", errors.New("Nothing set")
}

func (o *Option) Set(opt string) error {
	if opt == "" {
		return errors.New("Cannot hold a empty string")
	}
	o.some = opt
	o.none = false
	return nil
}
