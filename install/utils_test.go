package install

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLineInFileWindowsPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bashrc")
	line := `complete -C C:\Users\foo\bar\gocomplete.exe mycmd`
	require.NoError(t, os.WriteFile(path, []byte("# comment\n"+line+"\nexport FOO=bar\n"), 0644))

	assert.NotPanics(t, func() {
		assert.True(t, lineInFile(path, line))
	})
}

func TestRemoveFromFileWindowsPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bashrc")
	line := `complete -C C:\Users\foo\bar\gocomplete.exe mycmd`
	require.NoError(t, os.WriteFile(path, []byte("export FOO=bar\n"+line+"\n"), 0644))

	assert.NotPanics(t, func() {
		require.NoError(t, removeFromFile(path, line))
	})
	assert.False(t, lineInFile(path, line))
}
