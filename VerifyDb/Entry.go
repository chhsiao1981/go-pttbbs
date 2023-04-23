// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package VerifyDb

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Entry struct {
	_tab flatbuffers.Table
}

func GetRootAsEntry(buf []byte, offset flatbuffers.UOffsetT) *Entry {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Entry{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsEntry(buf []byte, offset flatbuffers.UOffsetT) *Entry {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Entry{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Entry) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Entry) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Entry) Userid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Entry) Generation() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Entry) MutateGeneration(n int64) bool {
	return rcv._tab.MutateInt64Slot(6, n)
}

func (rcv *Entry) Vmethod() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Entry) MutateVmethod(n int32) bool {
	return rcv._tab.MutateInt32Slot(8, n)
}

func (rcv *Entry) Vkey() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Entry) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Entry) MutateTimestamp(n int64) bool {
	return rcv._tab.MutateInt64Slot(12, n)
}

func EntryStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func EntryAddUserid(builder *flatbuffers.Builder, userid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(userid), 0)
}
func EntryAddGeneration(builder *flatbuffers.Builder, generation int64) {
	builder.PrependInt64Slot(1, generation, 0)
}
func EntryAddVmethod(builder *flatbuffers.Builder, vmethod int32) {
	builder.PrependInt32Slot(2, vmethod, 0)
}
func EntryAddVkey(builder *flatbuffers.Builder, vkey flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(vkey), 0)
}
func EntryAddTimestamp(builder *flatbuffers.Builder, timestamp int64) {
	builder.PrependInt64Slot(4, timestamp, 0)
}
func EntryEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
