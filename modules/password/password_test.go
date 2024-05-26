package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashValidate(t *testing.T) {
	passwordLen := 32
	password := TemporaryPassword(passwordLen)
	require.Equal(t, passwordLen, len(password))
	hash, err := Hash(password)
	require.NoError(t, err)
	require.True(t, Validate(hash, password))
}
