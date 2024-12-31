package fix

import (
	"fmt"
)

func Setup() {
	fmt.Printf(
		"just printing from '%s' with '%s'", // want "explicit single-quoted '%s' should be replaced by `%q` in fmt.Printf"
		"function",
		"style",
	)
}
