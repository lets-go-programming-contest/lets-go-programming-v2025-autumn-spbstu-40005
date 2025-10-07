package must

import (
	"fmt"
	"io"
)

func Must(op string, err error) {
	if err != nil {
		panic(fmt.Sprintf("op: %s", err))
	}
}

func Close(path string, closer io.Closer) {
	Must(fmt.Sprintf("close %q", path), closer.Close())
}
