package main

import (
	"fmt"

	"github.com/donnol/do"
)

func main() {
	do.Must(retErr()) // without result

	// specify result type with type parameter
	_ = do.Must1(retErrAndOneResult()) // with one result

	_, _ = do.Must2(retErrAndTwoResult()) // with two result
}

func retErr() error {
	return fmt.Errorf("a new error")
}

func retErrAndOneResult() (int, error) {
	return 1, fmt.Errorf("a new error")
}

func retErrAndTwoResult() (int, int, error) {
	return 0, 1, fmt.Errorf("a new error")
}
