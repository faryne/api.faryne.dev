package fire_department

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTaipei(t *testing.T) {
	resp, err := Taipei()
	require.NoError(t, err)
	t.Logf("resp:%v", resp)
}
