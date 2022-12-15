package cache

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"

	log "github.com/sirupsen/logrus"
)

// AddToUHash
func AddToUHash(uidInCache ptttype.UIDInStore, userID *ptttype.UserID_t) error {
	h := cmsys.StringHashWithHashBits(userID[:])

	// line: 166
	Shm.Shm.Userid[uidInCache] = *userID

	// init vars
	p := h
	val := ptttype.UIDInStore(0)
	pval := &val
	valptr := unsafe.Pointer(pval)
	offset := unsafe.Offsetof(Shm.Raw.HashHead)

	// line: 168
	Shm.ReadAt(
		offset+ptttype.UID_IN_STORE_SZ*uintptr(p),
		ptttype.UID_IN_STORE_SZ,
		valptr,
	)

	times := 0
	isNext := false
	offsetNextInHash := unsafe.Offsetof(Shm.Raw.NextInHash)
	for ; times < ptttype.MAX_USERS && val != -1; times++ {
		offset = offsetNextInHash
		isNext = true
		p = cmsys.Fnv32_t(val)
		Shm.ReadAt(
			offset+ptttype.UID_IN_STORE_SZ*uintptr(p),
			ptttype.UID_IN_STORE_SZ,
			valptr,
		)
	}
	if times >= ptttype.MAX_USERS {
		log.Errorf("Unable to add-to-uhash! uid-in-cache: %v userID: %v", uidInCache, string(userID[:]))
		return ErrAddToUHash
	}

	// set current ptr
	if !isNext {
		Shm.Shm.HashHead[p] = uidInCache
	} else {
		Shm.Shm.NextInHash[p] = uidInCache
	}

	// set next as -1
	Shm.Shm.NextInHash[uidInCache] = -1

	return nil
}

// RemoveFromUHash
func RemoveFromUHash(uidInCache ptttype.UIDInStore) error {
	userID := &ptttype.UserID_t{}

	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Userid)+ptttype.USER_ID_SZ*uintptr(uidInCache),
		ptttype.USER_ID_SZ,
		unsafe.Pointer(userID),
	)

	h := cmsys.StringHashWithHashBits(userID[:])

	// line: 191
	p := h
	val := ptttype.UIDInStore(0)
	pval := &val
	valptr := unsafe.Pointer(pval)
	offset := unsafe.Offsetof(Shm.Raw.HashHead)
	Shm.ReadAt(
		offset+types.INT32_SZ*uintptr(p),
		types.INT32_SZ,
		valptr,
	)

	// line: 194
	times := 0
	isNext := false
	for ; times < ptttype.MAX_USERS && val != -1 && val != uidInCache; times++ {
		p = cmsys.Fnv32_t(val)
		isNext = true
		offset = unsafe.Offsetof(Shm.Raw.NextInHash)
		Shm.ReadAt(
			offset+types.INT32_SZ*uintptr(p),
			types.INT32_SZ,
			valptr,
		)
	}
	if times >= ptttype.MAX_USERS {
		log.Errorf("Unable to remove-from-uhash! uid-in-cache: %v userID: %v", uidInCache, string(userID[:]))
		return ErrRemoveFromUHash
	}

	if val == uidInCache {
		nextNum := ptttype.UIDInStore(0)
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.NextInHash)+ptttype.UID_IN_STORE_SZ*uintptr(uidInCache),
			ptttype.UID_IN_STORE_SZ,
			unsafe.Pointer(&nextNum),
		)

		if !isNext {
			Shm.Shm.HashHead[p] = nextNum
		} else {
			Shm.Shm.NextInHash[p] = nextNum
		}
	}
	return nil
}

// SearchUser
// Params
//
//	userID: querying user-id.
//	isReturn: is return the user-id in the shm.
//
// Return
//
//	uid:
//	err:
func SearchUserRaw(userID *ptttype.UserID_t, rightID *ptttype.UserID_t) (uid ptttype.UID, err error) {
	if userID[0] == 0 {
		return 0, nil
	}
	return DoSearchUserRaw(userID, rightID)
}

func DoSearchUserRaw(userID *ptttype.UserID_t, rightID *ptttype.UserID_t) (ptttype.UID, error) {
	// XXX we should have 0 as non-exists.
	//     currently the reason why it's ok is because the probability of collision on 0 is low.

	_ = StatInc(ptttype.STAT_SEARCHUSER)
	h := cmsys.StringHashWithHashBits(userID[:])

	// p = SHM->hash_head[h]  //line: 219
	p := ptttype.UIDInStore(0)
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.HashHead)+types.INT32_SZ*uintptr(h),
		types.INT32_SZ,
		unsafe.Pointer(&p),
	)

	shmUserID := ptttype.UserID_t{}
	for times := 0; times < ptttype.MAX_USERS && p != -1 && p < ptttype.MAX_USERS; times++ {
		// if (strcasecmp(SHM->userid[p], userid) == 0)  //line: 222
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.Userid)+ptttype.USER_ID_SZ*uintptr(p),
			ptttype.USER_ID_SZ,
			unsafe.Pointer(&shmUserID),
		)
		if types.Cstrcasecmp(userID[:], shmUserID[:]) == 0 {
			if userID[0] != 0 && rightID != nil {
				copy(rightID[:], shmUserID[:])
			}
			return p.ToUID(), nil
		}
		Shm.ReadAt(
			unsafe.Offsetof(Shm.Raw.NextInHash)+types.INT32_SZ*uintptr(p),
			types.INT32_SZ,
			unsafe.Pointer(&p),
		)
	}

	return 0, nil
}

// GetUserID
//
// XXX uid = uid-in-cache + 1
func GetUserID(uid ptttype.UID) (*ptttype.UserID_t, error) {
	uidInCache := uid.ToUIDInStore()
	if uidInCache < 0 || uidInCache >= ptttype.MAX_USERS {
		return nil, ErrInvalidUID
	}

	userID := &ptttype.UserID_t{}
	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.Userid)+ptttype.USER_ID_SZ*uintptr(uidInCache),
		ptttype.USER_ID_SZ,
		unsafe.Pointer(userID),
	)

	return userID, nil
}

// SetUserID
//
// XXX uid = uid-in-cache + 1
func SetUserID(uid ptttype.UID, userID *ptttype.UserID_t) (err error) {
	if uid <= 0 || uid > ptttype.MAX_USERS {
		return ErrInvalidUID
	}

	uidInCache := uid.ToUIDInStore()
	errRemove := RemoveFromUHash(uidInCache)
	errAdd := AddToUHash(uidInCache, userID)
	if errRemove != nil {
		return errRemove
	}
	if errAdd != nil {
		return errAdd
	}

	return nil
}

// CooldownTimeOf
// https://github.com/ptt/pttbbs/blob/master/include/cmbbs.h#L97
func CooldownTimeOf(uid ptttype.UID) (cooldowntime types.Time4) {
	uidInCache := uid.ToUIDInStore()

	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.CooldownTime)+types.TIME4_SZ*uintptr(uidInCache),
		types.TIME4_SZ,
		unsafe.Pointer(&cooldowntime),
	)

	// types.Time4 is int32, not uint32
	// we use 0x7FFFFFF0 instead of 0xFFFFFFF0
	return cooldowntime & 0x7FFFFFF0
}

func AddCooldownTime(uid ptttype.UID, minutes int) (err error) {
	cooldowntime := CooldownTimeOf(uid)
	base := types.NowTS()
	if base < cooldowntime {
		base = cooldowntime
	}

	base += types.Time4(minutes) * 60
	base &= 0x7FFFFFF0

	uidInCache := uid.ToUIDInStore()
	Shm.Shm.CooldownTime[uidInCache] = base

	return nil
}

// PosttimesOf
// https://github.com/ptt/pttbbs/blob/master/include/cmbbs.h#L98
func PosttimesOf(uid ptttype.UID) (posttimes types.Time4) {
	uidInCache := uid.ToUIDInStore()

	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.CooldownTime)+types.TIME4_SZ*uintptr(uidInCache),
		types.TIME4_SZ,
		unsafe.Pointer(&posttimes),
	)

	return posttimes & 0xF
}

func AddPosttimes(uid ptttype.UID, times int) (err error) {
	posttimes := PosttimesOf(uid)
	newPosttimes := posttimes + types.Time4(times)

	uidInCache := uid.ToUIDInStore()
	if newPosttimes < 0xf {
		Shm.Shm.CooldownTime[uidInCache] += types.Time4(times)
	} else {
		Shm.Shm.CooldownTime[uidInCache] |= 0xf
	}

	return nil
}
