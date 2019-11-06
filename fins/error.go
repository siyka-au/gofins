package fins

import (
	"fmt"
	"time"
)

// Client errors

// ResponseTimeoutError An error when the other end of a command code hasn't responded within the set timeout period
type ResponseTimeoutError struct {
	duration time.Duration
}

func (e ResponseTimeoutError) Error() string {
	return fmt.Sprintf("Response timeout of %s has been reached", e.duration)
}

// IncompatibleMemoryAreaError An error for when an imcompatible memory area for the given operation is given
type IncompatibleMemoryAreaError struct {
	memoryArea MemoryArea
}

func (e IncompatibleMemoryAreaError) Error() string {
	return fmt.Sprintf("The memory area is incompatible with the data type to be read: 0x%X", e.memoryArea)
}
