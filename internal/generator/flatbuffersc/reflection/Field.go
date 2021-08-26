// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package reflection

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Field struct {
	_tab flatbuffers.Table
}

func GetRootAsField(buf []byte, offset flatbuffers.UOffsetT) *Field {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Field{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsField(buf []byte, offset flatbuffers.UOffsetT) *Field {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Field{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Field) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Field) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Field) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Field) Type(obj *Type) *Type {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Type)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *Field) Id() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Field) MutateId(n uint16) bool {
	return rcv._tab.MutateUint16Slot(8, n)
}

func (rcv *Field) Offset() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Field) MutateOffset(n uint16) bool {
	return rcv._tab.MutateUint16Slot(10, n)
}

func (rcv *Field) DefaultInteger() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Field) MutateDefaultInteger(n int64) bool {
	return rcv._tab.MutateInt64Slot(12, n)
}

func (rcv *Field) DefaultReal() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Field) MutateDefaultReal(n float64) bool {
	return rcv._tab.MutateFloat64Slot(14, n)
}

func (rcv *Field) Deprecated() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Field) MutateDeprecated(n bool) bool {
	return rcv._tab.MutateBoolSlot(16, n)
}

func (rcv *Field) Required() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Field) MutateRequired(n bool) bool {
	return rcv._tab.MutateBoolSlot(18, n)
}

func (rcv *Field) Key() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Field) MutateKey(n bool) bool {
	return rcv._tab.MutateBoolSlot(20, n)
}

func (rcv *Field) Attributes(obj *KeyValue, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Field) AttributesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Field) Documentation(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *Field) DocumentationLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Field) Optional() bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.GetBool(o + rcv._tab.Pos)
	}
	return false
}

func (rcv *Field) MutateOptional(n bool) bool {
	return rcv._tab.MutateBoolSlot(26, n)
}

func FieldStart(builder *flatbuffers.Builder) {
	builder.StartObject(12)
}
func FieldAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func FieldAddType(builder *flatbuffers.Builder, type_ flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(type_), 0)
}
func FieldAddId(builder *flatbuffers.Builder, id uint16) {
	builder.PrependUint16Slot(2, id, 0)
}
func FieldAddOffset(builder *flatbuffers.Builder, offset uint16) {
	builder.PrependUint16Slot(3, offset, 0)
}
func FieldAddDefaultInteger(builder *flatbuffers.Builder, defaultInteger int64) {
	builder.PrependInt64Slot(4, defaultInteger, 0)
}
func FieldAddDefaultReal(builder *flatbuffers.Builder, defaultReal float64) {
	builder.PrependFloat64Slot(5, defaultReal, 0.0)
}
func FieldAddDeprecated(builder *flatbuffers.Builder, deprecated bool) {
	builder.PrependBoolSlot(6, deprecated, false)
}
func FieldAddRequired(builder *flatbuffers.Builder, required bool) {
	builder.PrependBoolSlot(7, required, false)
}
func FieldAddKey(builder *flatbuffers.Builder, key bool) {
	builder.PrependBoolSlot(8, key, false)
}
func FieldAddAttributes(builder *flatbuffers.Builder, attributes flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(9, flatbuffers.UOffsetT(attributes), 0)
}
func FieldStartAttributesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func FieldAddDocumentation(builder *flatbuffers.Builder, documentation flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(10, flatbuffers.UOffsetT(documentation), 0)
}
func FieldStartDocumentationVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func FieldAddOptional(builder *flatbuffers.Builder, optional bool) {
	builder.PrependBoolSlot(11, optional, false)
}
func FieldEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
