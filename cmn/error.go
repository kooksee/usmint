package cmn

import "strings"

var ErrDelim = "\n"

func Errs(errs ... string) string {
	return strings.Join(errs, ErrDelim)
}

func MustNotErr(errs ... error) {
	for _, err := range errs {
		if err != nil {
			logger.Error(err.Error())
			panic("")
		}
	}
}
