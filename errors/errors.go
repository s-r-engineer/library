package libraryErrors

import "fmt"

func Panicer(err any) {
	if err != nil {
		panic(err)
	}
}

func WrapError(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s -> %w", msg, err)
	}
	return nil
}
