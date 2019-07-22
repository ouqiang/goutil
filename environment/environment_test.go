package environment

import (
        "github.com/stretchr/testify/require"
        "testing"
)

func TestEnvironment(t *testing.T)  {
        require.Equal(t, "", Mode())

        Init(DevMode)
        require.True(t, IsDev())


        Init(TestMode)
        require.True(t, IsTest())

        Init(PreMode)
        require.True(t, IsPre())

        Init(ProdMode)
        require.True(t, IsProd())
}
