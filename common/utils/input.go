package utils

import (
	"bufio"
	"fmt"
	"os"
)

func RegisterInputWithMessage(message string) (input string) {
	fmt.Print(message)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}
