package ptttype

import "unsafe"

type VerifyDBStatus int8

const (
	VERIFYDB_OK    VerifyDBStatus = 0
	VERIFYDB_ERROR VerifyDBStatus = -1
)

type VerifyDBVMethod int32

const (
	VMETHOD_UNSET VerifyDBVMethod = 0
	VMETHOD_EMAIL VerifyDBVMethod = 1
	VMETHOD_SMS   VerifyDBVMethod = 2
)

type VerifyDBMessage struct {
	RegMailDBReqHeader
	Message []byte
}

type VerifyDBMessage2 struct {
	RegMailDBReqHeader
	Message [0]byte
}

var EMPTY_VERIFY_DB_MESSAGE2 = VerifyDBMessage2{}

const VERIFY_DB_MESSAGE_SZ = unsafe.Sizeof(EMPTY_VERIFY_DB_MESSAGE2)

const VERIFYDB_VKEY_SIZE = 160
