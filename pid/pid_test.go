package pid

import (
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadWriteToFile(t *testing.T) {
	f, err := ioutil.TempFile("", "pid")
	require.NoError(t, err)
	_ = f.Close()
	defer func() {
		_ = os.Remove(f.Name())
	}()
	t.Logf("temp file: %s, pid: %d", f.Name(), os.Getpid())
	err = WriteToFile(f.Name())
	require.NoError(t, err)

	pid, err := ReadFromFile(f.Name())
	require.NoError(t, err)

	require.Equal(t, os.Getpid(), pid)
}

func TestIsRunning(t *testing.T) {
	require.Equal(t, true, IsRunning(os.Getpid()))
	require.Equal(t, false, IsRunning(math.MaxInt32))
}
