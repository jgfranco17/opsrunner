package runner

import (
	"context"
	"fmt"
)

type UnexpectedCommandError struct {
	Want string
	Got  string
}

type MockExpectedInput struct {
	Step       Step
	ReturnCode int
	Output     string
	Error      error
}

func (e *UnexpectedCommandError) Error() string {
	return fmt.Sprintf("Unexpeced command: wanted %s but got %s", e.Want, e.Got)
}

type MockExecutor struct {
	ExpectedSteps []MockExpectedInput
	called        bool
}

func (m *MockExecutor) WithStep(command string, args string, mockReturnCode int, mockOutput string, mockError error) *MockExecutor {
	newStep := Step{
		Command: command,
		Args:    args,
	}
	m.ExpectedSteps = append(m.ExpectedSteps, MockExpectedInput{
		Step:       newStep,
		ReturnCode: mockReturnCode,
		Output:     mockOutput,
		Error:      mockError,
	})
	return m
}

func (m *MockExecutor) Exec(ctx context.Context, name string, args string) (int, string, error) {
	m.called = true
	for _, entry := range m.ExpectedSteps {
		if name == entry.Step.Command && args == entry.Step.Args {
			return entry.ReturnCode, entry.Output, entry.Error
		}
	}
	return -1, "", fmt.Errorf("Undefined call: %s %s", name, args)
}

func (m *MockExecutor) IsCalled() bool {
	return m.called
}

func NewMockExecutor(expectedCallCount int) *MockExecutor {
	return &MockExecutor{
		ExpectedSteps: make([]MockExpectedInput, expectedCallCount),
	}
}
