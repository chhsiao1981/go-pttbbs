package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const LOAD_VERIFYDB_R = "/verifydb"

type LoadVerifyDBParams struct {
	Offset   int  `json:"offset" form:"offset" url:"offset"`
	Limit    int  `json:"limit" form:"limit" url:"limit"`
	IsSystem bool `json:"system,omitempty" form:"system,omitempty" url:"system"`
}

type LoadVerifyDBResult struct {
	Entry *bbs.VerifyDBEntry `json:"content"`
}

func LoadVerifyDBWrapper(c *gin.Context) {
	params := &LoadVerifyDBParams{}
	LoginRequiredQuery(LoadVerifyDB, params, c)
}

func LoadVerifyDB(remoteAddr string, uuserID bbs.UUserID, params interface{}) (result interface{}, err error) {
	theParams, ok := params.(*LoadVerifyDBParams)
	if !ok {
		return nil, ErrInvalidParams
	}

	if theParams.IsSystem {
		uuserID = bbs.UUserID(string(ptttype.STR_SYSOP))
	}

	if uuserID != bbs.UUserID(string(ptttype.STR_SYSOP)) {
		return nil, ErrInvalidUser
	}

	ents, err := bbs.LoadVerifyDB(theParams.Offset, theParams.Limit)
	if err != nil {
		return nil, err
	}

	return ents, nil
}
