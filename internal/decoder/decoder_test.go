package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	require := require.New(t)

	d := New("cão")
	expected := []rune{'c', 'ã', 'o'}
	for _, ch := range expected {
		require.True(d.Decode())
		require.Equal(ch, d.Decoded())
	}
	require.False(d.Decode())
}
