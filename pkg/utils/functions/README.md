### `GoSafe` Function in Go

The `GoSafe` function allows you to execute a provided function (`fn`) asynchronously in a goroutine with built-in error handling. If an error occurs (a panic), it is captured and logged by default. If an APM (Application Performance Management) service is available, the error will also be sent to the APM system for monitoring.

#### Usage

The `GoSafe` function is useful for safely executing functions concurrently. If an error occurs during execution, it will be logged, and optionally, it will be sent to an APM service if one is configured.

#### Function Definition

```go
func GoSafe(fn func(), logFunc func(string))
```

#### Parameters:

* `fn`: The function that you want to execute safely in a separate goroutine.
* `logFunc`: A function used to log errors. If no APM service is provided, this function will log the panic message. If an APM service is available, it can be used to send the error to the APM system.

#### How It Works:

* The provided `fn` function is executed in a separate goroutine.
* If a panic occurs during the execution of `fn`, the panic is caught by a `defer` block.
* If the APM service is available, the error is sent to it via `apm.CaptureError`. Otherwise, the error is logged using the provided `logFunc`.

#### Example Usage:

```go
package main

import (
	"fmt"
	"log"
)

// Default log function
func defaultLog(message string) {
	log.Println(message)
}

// Example function to run asynchronously
func exampleFunction() {
	panic("something went wrong!")
}

func main() {
	// Using GoSafe with default logging
	GoSafe(exampleFunction, defaultLog)

	// Continue with other tasks...
	fmt.Println("Program continues running...")
}
```

#### Example with APM:

```go
import "go.elastic.co/apm"

func GoSafe(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				apm.CaptureError(nil, fmt.Errorf("panic: %v", r)).Send()
			}
		}()
		fn()
	}()
}
```

In this example:

* If no APM service is provided, the error is logged using the default `defaultLog` function.
* If an APM service (e.g., Elastic APM) is available, it captures the error and sends it for monitoring.

#### Why Use `GoSafe`?

* **Error Handling:** Ensures that panics in goroutines are safely handled without crashing the application.
* **Optional APM Integration:** Logs errors to an APM service if available; otherwise, logs to standard output.
* **Concurrency Safety:** Ideal for functions that might panic and need to be run asynchronously, especially in production environments where reliability is crucial.
