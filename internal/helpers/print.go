package helpers

import "fmt"

func PrintFile(filename string, file interface{}) (output string) {
	output = filename + "\n"
	output = output + fmt.Sprint(file)

	return
}

func PrintFileWithSpacer(filename string, file interface{}) (output string) {
	output = "\n--\n" + filename + ":\n"
	output = output + fmt.Sprint(file)

	return
}
