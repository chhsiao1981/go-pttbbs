package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/VerifyDb"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	flatbuffers "github.com/google/flatbuffers/go"
)

func LoadVerifyDB(offset int, limit int) (ents []*VerifyDb.Entry, err error) {
	fbb := flatbuffers.NewBuilder(VERIFYDB_BUF_SIZE)
	VerifyDb.ListRequestStart(fbb)
	VerifyDb.ListRequestAddOffset(fbb, int32(offset))
	VerifyDb.ListRequestAddLimit(fbb, int32(limit))
	listreq := VerifyDb.ListRequestEnd(fbb)

	VerifyDb.MessageStart(fbb)
	VerifyDb.MessageAddMessageType(fbb, VerifyDb.AnyMessageListRequest)
	VerifyDb.MessageAddMessage(fbb, listreq)
	req := VerifyDb.MessageEnd(fbb)

	fbb.Finish(req)
	reqBytes := fbb.FinishedBytes()

	msgReplyBytes, err := verifydbTransact(reqBytes)
	if err != nil {
		return nil, err
	}
	msgReply := VerifyDb.GetRootAsMessage(msgReplyBytes, 0)
	if msgReply.MessageType() != VerifyDb.AnyMessageGetReply {
		return nil, ErrInvalidVerifyDBReply
	}

	replyTab := &flatbuffers.Table{}
	if !msgReply.Message(replyTab) {
		return nil, ErrInvalidVerifyDBReply
	}
	reply := &VerifyDb.GetReply{}
	reply.Init(replyTab.Bytes, replyTab.Pos)

	ok := reply.Ok()
	if !ok {
		return nil, ErrInvalidVerifyDBReply
	}
	nEntry := reply.EntryLength()
	ents = make([]*VerifyDb.Entry, nEntry)
	for idx := 0; idx < nEntry; idx++ {
		reply.Entry(ents[idx], idx)
	}
	return ents, nil
}

func verifydbTransact(theBytes []byte) (outBytes []byte, err error) {
	conn, err := InitConn(ptttype.REGMAILD_ADDR)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rep := &ptttype.VerifyDBMessage{}
	rep.Cb = types.Size_t(ptttype.REG_MAILDB_REQ_HEADER_SZ) + types.Size_t(len(theBytes))
	rep.Operation = ptttype.VERIFYDB_MESSAGE

	nWrite, err := conn.Write(&rep.RegMailDBReqHeader, ptttype.REG_MAILDB_REQ_HEADER_SZ)
	if err != nil {
		return nil, err
	}
	if nWrite != ptttype.REG_MAILDB_REQ_HEADER_SZ {
		return nil, ErrInvalidVerifyDBTransact
	}
	nWrite, err = conn.Write(theBytes, len(theBytes))
	if err != nil {
		return nil, err
	}
	if nWrite != len(theBytes) {
		return nil, ErrInvalidVerifyDBTransact
	}

	nRead, err := conn.Read(&rep.RegMailDBReqHeader, ptttype.REG_MAILDB_REQ_HEADER_SZ)
	if err != nil {
		return nil, err
	}
	if nRead != ptttype.REG_MAILDB_REQ_HEADER_SZ {
		return nil, ErrInvalidVerifyDBTransact
	}

	nBytes := rep.Cb - types.Size_t(ptttype.REG_MAILDB_REQ_HEADER_SZ)
	outBytes = make([]byte, nBytes)
	nRead, err = conn.Read(outBytes, nBytes)
	if err != nil {
		return nil, err
	}
	if nRead != nBytes {
		return nil, ErrInvalidVerifyDBTransact
	}

	return outBytes, nil
}
