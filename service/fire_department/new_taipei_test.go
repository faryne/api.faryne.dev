package fire_department

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTaipei(t *testing.T) {
	resp, err := NewTaipei()
	require.NoError(t, err)
	t.Logf("%+v", resp)
}
