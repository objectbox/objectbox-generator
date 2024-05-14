package object

import "time"

// Tests all available GO & ObjectBox types
type Typeful struct {
	Id           uint64 `objectbox:"id(assignable)"` // NOTE ID is currently required
	Int          int
	Int8         int8
	Int16        int16
	Int32        int32
	Int64        int64
	Uint         uint
	Uint8        uint8
	Uint16       uint16
	Uint32       uint32
	Uint64       uint64
	Bool         bool
	String       string
	StringVector []string
	Byte         byte
	ByteVector   []byte
	Rune         rune
	Float32      float32
	FloatVector  []float32
	Float64      float64
	Date         int64     `objectbox:"date index"`
	Time         time.Time `objectbox:"date,index"`
	Time2        time.Time // prints a warning, otherwise the same as with an annotation
	TimeNano     time.Time `objectbox:"date-nano,index"`
}

type TSDate struct {
	Id        uint64
	timestamp int64 `objectbox:"id-companion,date"`
}

type TSDateNano struct {
	Id        uint64
	timestamp int64 `objectbox:"id-companion,date-nano"`
}
