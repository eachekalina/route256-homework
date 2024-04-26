package core

type ValidationError struct {
	Err string
}

func (v ValidationError) Error() string {
	return v.Err
}
