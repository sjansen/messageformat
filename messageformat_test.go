package messageformat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var messages = map[string]string{
	"en": `There {n, plural, one{is 1 item} other{are # items}} in your inbox.`,
	"pt": `{n, plural, one{Existe 1 item} other{Existem # itens}} na sua caixa de entrada.`,
}

func TestParse(t *testing.T) {
	require := require.New(t)

	for _, message := range messages {
		msg, err := Parse(message)
		require.NoError(err)
		require.NotNil(msg)
	}
}
