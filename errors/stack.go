package errors

import (
	"bytes"
	"fmt"
)

type stack []string

func New(message string, values ...any) *stack {
	result := stack([]string{fmt.Sprintf(message, values)})

	return &result
}

func (s *stack) Cause(err error) *stack {
	*s = append(*s, ":", err.Error())

	return s
}

func (s *stack) Error() string {
	result := bytes.NewBufferString("")
	stack := []string(*s)
	for i := range stack {
		result.WriteString(stack[i])
	}

	return result.String()
}
