package rrors

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTest(t *testing.T, str string, idx int) {
	require.NoError(t, os.WriteFile(fmt.Sprintf("./tests/test%d", idx), []byte(str), os.ModePerm))
}

func readTest(t *testing.T, idx int) string {
	bytes, fdErr := os.ReadFile(fmt.Sprintf("./tests/test%d", idx))
	require.NoError(t, fdErr)
	return string(bytes)
}

func TestErr(t *testing.T) {
	var errs error
	errs = errors.Join(Errorf("hello"), Errorf("world\n-happened"), Errorf("woowowow %w", errors.Join(Errorf("more text"), Errorf("hw some error happened %w", Errorf("here %w", Errorf("here %w", fmt.Errorf("wow")))), Errorf("aaah"))))
	err := Errorf("my error: %w", errs)
	firstTest := readTest(t, 0)
	require.Equal(t, firstTest, err.Error(), err.Error())

	errs = Errorf("%w", fmt.Errorf("hello2: %w", errors.Join(fmt.Errorf("hello"), fmt.Errorf("world: %w", errors.Join(fmt.Errorf("more stuff"), fmt.Errorf("more stuff"))))))
	err = Errorf("my errors: %w", errs)
	secondTest := readTest(t, 1)
	require.Equal(t, secondTest, err.Error(), err.Error())
}

func init() {
	Initialise()
}
