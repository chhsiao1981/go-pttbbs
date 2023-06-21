package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/VerifyDb"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type VerifyDBEntry struct {
	UserID     UUserID
	Generation types.Time4
	Method     ptttype.VerifyDBVMethod
	VKey       string
	Timestamp  types.Time4
}

func NewVerifyDBEntry(ent *VerifyDb.Entry) (newEnt *VerifyDBEntry) {
	userID := ent.Userid()
	userIDRaw := ptttype.UserID_t{}
	copy(userIDRaw[:], userID[:])

	return &VerifyDBEntry{
		UserID:     ToUUserID(&userIDRaw),
		Generation: types.Time4(ent.Generation()),
		Method:     ptttype.VerifyDBVMethod(ent.Vmethod()),
		VKey:       types.CstrToString(ent.Vkey()),
		Timestamp:  types.Time4(ent.Timestamp()),
	}
}
