package names

import "github.com/z7zmey/php-parser/errors"

type ParserErrors []*errors.Error

func (E ParserErrors) Error() string {
	if len(E) == 0 {
		return ""
	}

	msg := E[0].String()

	for _, e := range E[1:] {
		msg += "\n" + e.String()
	}

	return msg
}
