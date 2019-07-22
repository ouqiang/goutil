package ip

import (
        "github.com/stretchr/testify/require"
        "testing"
)

func TestInternalIP(t *testing.T) {
        v := InternalIP()
        t.Logf("internal ip: %s", v)
        require.NotEmpty(t, v)
}