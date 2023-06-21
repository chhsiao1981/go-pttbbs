package bbs

import "github.com/Ptt-official-app/go-pttbbs/ptt"

func LoadVerifyDB(offset int, limit int) (ents []*VerifyDBEntry, err error) {
	entRaws, err := ptt.LoadVerifyDB(offset, limit)
	if err != nil {
		return nil, err
	}

	ents = make([]*VerifyDBEntry, len(entRaws))
	for idx, each := range entRaws {
		ents[idx] = NewVerifyDBEntry(each)
	}

	return ents, nil
}
