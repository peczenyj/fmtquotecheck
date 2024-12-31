package b

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	something = fmt.Sprintf("starting '%s'", "program")

	errSomething = fmt.Errorf("unexpected error from '%s'", "program")
)

func setup() {
	_ = fmt.Sprintf("starting '%s'", "function")

	_ = fmt.Errorf("unexpected error from '%s'", "function")

	fmt.Printf("just printing from '%s'", "function") // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Printf"

	fmt.Printf(
		"just printing from '%s' with '%s'", // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Printf"
		"function",
		"style",
	)

	fmt.Fprintf(os.Stderr, "just printing from '%s' on stderr", "function")

	_ = func(origin string) string {
		return fmt.Sprintf("starting '%s'", origin)
	}("callback")

	for i := 0; i < 3; i++ {
		fmt.Printf("for loop iteraration '%s'", strconv.Itoa(i)) // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Printf"
	}

	defer log.Printf("ending '%s'", "function/defer")

	defer func() {
		logger := log.Default()
		logger.Printf("ending '%s'", "callback/defer")
	}()
}
