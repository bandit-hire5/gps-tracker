package workers

import (
	"bufio"
)

type WorkerInterface interface {
	Work(input string, w *bufio.Writer)
}
