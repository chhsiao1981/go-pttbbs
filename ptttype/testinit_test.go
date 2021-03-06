package ptttype

import (
	"os"

	"github.com/Ptt-official-app/go-pttbbs/types"
	jww "github.com/spf13/jwalterweatherman"
)

func setupTest() {
	jww.SetLogOutput(os.Stderr)
	jww.SetLogThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelDebug)

	types.SetIsTest()

	SetIsTest()
}

func teardownTest() {
	UnsetIsTest()

	types.UnsetIsTest()
}
