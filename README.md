# RetryGO

![Retrygo](./docs/assets/retrygo.png)

---
[![codecov](https://codecov.io/gh/ic-it/retrygo/graph/badge.svg?token=HXT5N3O452)](https://codecov.io/gh/ic-it/retrygo)
[![Go Reference](https://pkg.go.dev/badge/github.com/ic-it/retrygo.svg)](https://pkg.go.dev/github.com/ic-it/retrygo)

**RetryGO** is a simple and flexible library for retrying functions in Go. 

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
- **Backoff:** The user can create its own retry strategy or use one
of the predefined backoff strategies.

## Features
- **Predefined Backoff Strategies:** The library provides some predefined
backoff strategies, such as constant, exponential, and linear backoff.
- **Combine Backoff Strategies:** The user can combine multiple backoff
strategies to create a custom backoff strategy.
- **Recover:** The user can enable the recover feature, which will recover
panics and behave as if the function returned an error.
- **Context:** The user can use the context to cancel the retry process.

## Installation
```bash
go get github.com/ic-it/retrygo
```

## Usage
[Examples](./examples/)

### Simple Example With Predefined Backoff
```go
package main

import (
	"context"
	"time"

	"github.com/ic-it/retrygo"
)

func main() {
	type ReturnType struct{}
	retry, err := retrygo.New[ReturnType](
		retrygo.Combine(
			retrygo.Constant(1*time.Second),
			retrygo.LimitCount(5),
		),
	)

	if err != nil {
		// Handle error
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := retry.Do(ctx, func(context.Context) (ReturnType, error) {
		// Do something
		return ReturnType{}, nil
	})

    if err != nil {
        // Handle error
    }

    // Continue with val
}
```

**Note:** If you don't want to get any value from a function, NewZero and DoZero should be a solution for you:

```go
func main() {
	retry, err := retrygo.NewZero(
		retrygo.Combine(
			retrygo.Constant(1*time.Second),
			retrygo.LimitCount(5),
		),
	)

	// Do something

	err := retry.DoZero(ctx, func(context.Context) error {
		// Do something
		return nil
	})

	// Handle error if needed
}
```

**Its also possible to use DoZero with New function:**
```go
func main() {
	type ReturnType struct{}
	retry, err := retrygo.New[ReturnType](
		retrygo.Combine(
			retrygo.Constant(1*time.Second),
			retrygo.LimitCount(5),
		),
	)

	// Do something

	err := retry.DoZero(ctx, func(context.Context) error {
		// Do something
		return nil
	})

	// Handle error if needed
}
```

**Important:** You can't use Do with NewZero.

### Example with Custom Backoff
```go
package main

import (
	"context"
	"time"

	"github.com/ic-it/retrygo"
)

func main() {
	type ReturnType struct{}
	retry, err := retrygo.New[ReturnType](
		func(ri retrygo.RetryInfo) (continueRetry bool, sleep time.Duration) {
			// Custom logic
			return false, 0
		},
	)

	if err != nil {
		// Handle error
	}

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

## Benchmarks
See benchmarks [here](./benchmarks/).

Results [gist](https://gist.github.com/ic-it/99a569a99772c38fafb447ba12baa19a).

## License
[MIT](./LICENSE.txt)
