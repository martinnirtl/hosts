package helpers

import "fmt"

func PrintFile(filename string, file interface{}) string {
	return filename + ":\n" + fmt.Sprint(file)
}

func PrintFileWithSpacer(filename string, file interface{}) string {
	return PrintFile(filename, file) + "\n---\n"
}
