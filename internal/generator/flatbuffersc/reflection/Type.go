// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package reflection

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Type struct {
	_tab flatbuffers.Table
}

func GetRootAsType(buf []byte, offset flatbuffers.UOffsetT) *Type {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Type{}
	x.Init(buf, n+offset)
	return x
}

func FinishTypeBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsType(buf []byte, offset flatbuffers.UOffsetT) *Type {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Type{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedTypeBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Type) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Type) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Type) BaseType() BaseType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return BaseType(rcv._tab.GetInt8(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *Type) MutateBaseType(n BaseType) bool {
	return rcv._tab.MutateInt8Slot(4, int8(n))
}

func (rcv *Type) Element() BaseType {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return BaseType(rcv._tab.GetInt8(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *Type) MutateElement(n BaseType) bool {
	return rcv._tab.MutateInt8Slot(6, int8(n))
}

func (rcv *Type) Index() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return -1
}

func (rcv *Type) MutateIndex(n int32) bool {
	return rcv._tab.MutateInt32Slot(8, n)
}

func (rcv *Type) FixedLength() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Type) MutateFixedLength(n uint16) bool {
	return rcv._tab.MutateUint16Slot(10, n)
}

/// The size (octets) of the `base_type` field.
func (rcv *Type) BaseSize() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 4
}

/// The size (octets) of the `base_type` field.
func (rcv *Type) MutateBaseSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(12, n)
}

/// The size (octets) of the `element` field, if present.
func (rcv *Type) ElementSize() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

/// The size (octets) of the `element` field, if present.
func (rcv *Type) MutateElementSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(14, n)
}

func TypeStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func TypeAddBaseType(builder *flatbuffers.Builder, baseType BaseType) {
	builder.PrependInt8Slot(0, int8(baseType), 0)
}
func TypeAddElement(builder *flatbuffers.Builder, element BaseType) {
	builder.PrependInt8Slot(1, int8(element), 0)
}
func TypeAddIndex(builder *flatbuffers.Builder, index int32) {
	builder.PrependInt32Slot(2, index, -1)
}
func TypeAddFixedLength(builder *flatbuffers.Builder, fixedLength uint16) {
	builder.PrependUint16Slot(3, fixedLength, 0)
}
func TypeAddBaseSize(builder *flatbuffers.Builder, baseSize uint32) {
	builder.PrependUint32Slot(4, baseSize, 4)
}
func TypeAddElementSize(builder *flatbuffers.Builder, elementSize uint32) {
	builder.PrependUint32Slot(5, elementSize, 0)
}
func TypeEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
