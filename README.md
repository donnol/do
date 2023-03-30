# do

Do something interesting.

Base Go1.18 with generic.

## Must

Panic if error is not nil, otherwise return some result except the error.

```go
package main

import (
	"fmt"

	"github.com/donnol/do"
)

func main() {
	do.Must(retErr()) // without result

	// specify result type with type parameter
	_ = do.Must1[int](retErrAndOneResult()) // with one result

	_, _ = do.Must2[int, int](retErrAndTwoResult()) // with two result
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
```
