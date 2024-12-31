package a

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	something = fmt.Sprintf("starting '%s'", "program") // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Sprintf"

	errSomething = fmt.Errorf("unexpected error from '%s'", "program") // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Errorf"
)

func setup() {
	_ = fmt.Sprintf("starting '%s'", "function") // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Sprintf"

	_ = fmt.Errorf("unexpected error from '%s'", "function") // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Errorf"

	fmt.Printf("just printing from '%s'", "function") // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Printf"

	fmt.Printf(
		"just printing from '%s' with '%s'", // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Printf"
		"function",
		"style",
	)

	fmt.Fprintf(os.Stderr, "just printing from '%s' on stderr", "function") // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Fprintf"

	_ = func(origin string) string {
		return fmt.Sprintf("starting '%s'", origin) // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Sprintf"
	}("callback")

	for i := 0; i < 3; i++ {
		fmt.Printf("for loop iteraration '%s'", strconv.Itoa(i)) // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Printf"
	}

	defer log.Printf("ending '%s'", "function/defer") // want "explicit single-quoted '%s' should be replaced by `%q` in log.Printf"

	logger := log.Default()

	defer func() {
		logger.Printf("ending '%s'", "callback/defer") // want "explicit single-quoted '%s' should be replaced by `%q` in logger.Printf"
	}()
}
