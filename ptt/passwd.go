package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func InitCurrentUser(userID *ptttype.UserID_t) (uid ptttype.Uid, user *ptttype.UserecRaw, err error) {
	uid, user, err = cmbbs.PasswdLoadUser(userID)
	if err != nil {
		return uid, user, err
	}

	if types.Cstrcmp(user.UserID[:], []byte(ptttype.STR_GUEST)) == 0 {
		pwcuInitGuestPerm(user)
	}

	return uid, user, nil
}

func InitCurrentUserByUid(uid ptttype.Uid) (user *ptttype.UserecRaw, err error) {
	user, err = cmbbs.PasswdQuery(uid)
	if err != nil {
		return nil, err
	}
	if types.Cstrcmp(user.UserID[:], []byte(ptttype.STR_GUEST)) == 0 {
		pwcuInitGuestPerm(user)
	}

	return user, nil
}

func passwdSyncUpdate(uid ptttype.Uid, user *ptttype.UserecRaw) error {
	if !uid.IsValid() {
		return cache.ErrInvalidUID
	}

	user.Money = cache.MoneyOf(uid)

	err := cmbbs.PasswdUpdate(uid, user)
	if err != nil {
		return err
	}

	return nil
}

func passwdSyncQuery(uid ptttype.Uid) (*ptttype.UserecRaw, error) {
	user, err := cmbbs.PasswdQuery(uid)
	if err != nil {
		return nil, err
	}

	user.Money = cache.MoneyOf(uid)

	return user, nil
}
