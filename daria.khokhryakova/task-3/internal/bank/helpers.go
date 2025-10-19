package bank

import "os"

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func FileClose(file *os.File) {
	if file != nil {
		err := file.Close()
		panicErr(err)
	}
}
