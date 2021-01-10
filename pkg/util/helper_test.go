package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStrPtr(t *testing.T) {
	assert.Equal(t, "test123", *StrPtr("test123"))
}

func TestStrPtrSlice(t *testing.T) {
	assert.Equal(t, []*string{StrPtr("test1"), StrPtr("test2")}, StrPtrSlice([]string{"test1", "test2"}))
}

func TestStrVal(t *testing.T) {
	assert.Equal(t, "", StrVal(nil))
	assert.Equal(t, "test123", StrVal(StrPtr("test123")))
}

func TestStrValSlice(t *testing.T) {
	assert.Equal(t, []string{"test1", "test2", ""}, StrValSlice([]*string{StrPtr("test1"), StrPtr("test2"), nil}))
}

func TestBoolPtr(t *testing.T) {
	assert.True(t, *BoolPtr(true))
	assert.False(t, *BoolPtr(false))
}

func TestBoolPtrSlice(t *testing.T) {
	assert.Equal(t, []*bool{BoolPtr(true), BoolPtr(false)}, BoolPtrSlice([]bool{true, false}))
}

func TestBoolVal(t *testing.T) {
	assert.False(t, BoolVal(nil))
	assert.True(t, BoolVal(BoolPtr(true)))
	assert.False(t, BoolVal(BoolPtr(false)))
}

func TestBoolValSlice(t *testing.T) {
	assert.Equal(t, []bool{true, false, false}, BoolValSlice([]*bool{BoolPtr(true), BoolPtr(false), nil}))
}

func TestIntPtr(t *testing.T) {
	assert.Equal(t, 42, *IntPtr(42))
}

func TestIncrementIntPtr(t *testing.T) {
	assert.Equal(t, 43, *IncrementIntPtr(IntPtr(42)))
}

func TestIntPtrSlice(t *testing.T) {
	assert.Equal(t, []*int{IntPtr(42), IntPtr(63)}, IntPtrSlice([]int{42, 63}))
}

func TestIntVal(t *testing.T) {
	assert.Equal(t, 0, IntVal(nil))
	assert.Equal(t, 42, IntVal(IntPtr(42)))
}

func TestIntValSlice(t *testing.T) {
	assert.Equal(t, []int{42, 63, 0}, IntValSlice([]*int{IntPtr(42), IntPtr(63), nil}))
}

func TestInt8Ptr(t *testing.T) {
	assert.Equal(t, int8(42), *Int8Ptr(42))
}

func TestInt8Val(t *testing.T) {
	assert.Equal(t, int8(0), Int8Val(nil))
	assert.Equal(t, int8(42), Int8Val(Int8Ptr(42)))
}

func TestInt16Ptr(t *testing.T) {
	assert.Equal(t, int16(42), *Int16Ptr(42))
}

func TestInt16Val(t *testing.T) {
	assert.Equal(t, int16(0), Int16Val(nil))
	assert.Equal(t, int16(42), Int16Val(Int16Ptr(42)))
}

func TestInt32Ptr(t *testing.T) {
	assert.Equal(t, int32(42), *Int32Ptr(42))
}

func TestInt32Val(t *testing.T) {
	assert.Equal(t, int32(0), Int32Val(nil))
	assert.Equal(t, int32(42), Int32Val(Int32Ptr(42)))
}

func TestInt64Ptr(t *testing.T) {
	assert.Equal(t, int64(42), *Int64Ptr(42))
}

func TestInt64Val(t *testing.T) {
	assert.Equal(t, int64(0), Int64Val(nil))
	assert.Equal(t, int64(42), Int64Val(Int64Ptr(42)))
}

func TestUInt8Ptr(t *testing.T) {
	assert.Equal(t, uint8(42), *UInt8Ptr(42))
}

func TestUInt8Val(t *testing.T) {
	assert.Equal(t, uint8(0), UInt8Val(nil))
	assert.Equal(t, uint8(42), UInt8Val(UInt8Ptr(42)))
}

func TestUInt16Ptr(t *testing.T) {
	assert.Equal(t, uint16(42), *UInt16Ptr(42))
}

func TestUInt16Val(t *testing.T) {
	assert.Equal(t, uint16(0), UInt16Val(nil))
	assert.Equal(t, uint16(42), UInt16Val(UInt16Ptr(42)))
}

func TestUInt32Ptr(t *testing.T) {
	assert.Equal(t, uint32(42), *UInt32Ptr(42))
}

func TestUInt32Val(t *testing.T) {
	assert.Equal(t, uint32(0), UInt32Val(nil))
	assert.Equal(t, uint32(42), UInt32Val(UInt32Ptr(42)))
}

func TestUInt64Ptr(t *testing.T) {
	assert.Equal(t, uint64(42), *UInt64Ptr(42))
}

func TestUInt64Val(t *testing.T) {
	assert.Equal(t, uint64(0), UInt64Val(nil))
	assert.Equal(t, uint64(42), UInt64Val(UInt64Ptr(42)))
}

func TestFloat64Ptr(t *testing.T) {
	assert.Equal(t, 42.0, *Float64Ptr(42.0))
}

func TestFloat64PtrSlice(t *testing.T) {
	assert.Equal(t, []*float64{Float64Ptr(42.0), Float64Ptr(63.0)}, Float64PtrSlice([]float64{42.0, 63.0}))
}

func TestFloat64Val(t *testing.T) {
	assert.Equal(t, 0.0, Float64Val(nil))
	assert.Equal(t, 42.0, Float64Val(Float64Ptr(42.0)))
}

func TestFloat64ValSlice(t *testing.T) {
	assert.Equal(t, []float64{42.0, 63.0, 0.0}, Float64ValSlice([]*float64{Float64Ptr(42.0), Float64Ptr(63.0), nil}))
}

func TestParseDateTypeSuccess(t *testing.T) {
	v, err := ParseDateType("24102021")
	if err != nil {
		t.Fatalf("Expected ParseDateTime not to give an error on a valid date string, got %s", err)
	}
	if v.Format("2006-01-02") != "2021-10-24" {
		t.Fatalf("Expected ParseDateTime to return date 2021-10-24, got %s", v)
	}
}

func TestParseDateTypeFailure(t *testing.T) {
	if _, err := ParseDateType("01192021"); err == nil {
		t.Fatal("Expected ParseDateTime to give an error on a valid date string, got none")
	}
}

func TestTrimLeftChar(t *testing.T) {
	a := assert.New(t)

	a.Equal("tes", SubStr("test", 0, 3))
	a.Equal("te", SubStr("te", 0, 3))
	a.Equal("", SubStr("", 0, 3))
	a.Equal("", SubStr("test", 0, 0))
	a.Equal("", SubStr("test", 10, 11))
}
