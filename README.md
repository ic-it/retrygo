# RetryGO
[![codecov](https://codecov.io/gh/ic-it/retrygo/graph/badge.svg?token=HXT5N3O452)](https://codecov.io/gh/ic-it/retrygo)

RetryGO is a simple library for retrying functions in Go. 

**RetryGO** is based on giving the user control over the logic responsible for 
deciding whether to continue retry attempts. 

## Purpose
The purpose of this library is to provide a simple but flexible way to retry
functions in Go.

Unlike other libraries, **RetryGO** does not provide a way to specify the number
of retry attempts. Instead, it allows the user to specify a function that will
be called after each retry attempt. This function will be responsible for
deciding whether to continue retrying or not.

**Main ideas:**
- **Retry Policy:** The user can specify a function that will be called after
each retry attempt. This function will be responsible for deciding whether to
continue retrying or not.
- **Prepared Backoff:** The user can create its own backoff strategy or use one
of the predefined strategies.
- **Context:** The user can specify a context to cancel the retry attempts.

## Installation
```bash
go get github.com/ic-it/retrygo
```

## Usage
[Examples](./examples/)
```go
package main

import (
	"context"
	"time"

	"github.com/ic-it/retrygo"
)

func main() {
	type ReturnType struct{}
	retry := retrygo.New[ReturnType](
		retrygo.Combine(
			retrygo.Constant(1*time.Second),
			retrygo.LimitCount(5),
			retrygo.LimitTime(10*time.Second),
		),
	)

	val, err := retry.Do(context.TODO(), func(context.Context) (ReturnType, error) {
		// Do something
		return ReturnType{}, nil
	})

    if err != nil {
        // Handle error
    }

    // Continue with val
}
```

## Documentation
[Doumentation](./docs/) is available in the `docs` folder. This documentation
was generated using [gomarkdoc](https://github.com/princjef/gomarkdoc).

## License
[MIT](./LICENSE)
