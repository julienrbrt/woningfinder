package util

import "time"

func StrPtr(v string) *string {
	return &v
}

func StrPtrSlice(v []string) []*string {
	result := make([]*string, len(v))
	for i, str := range v {
		result[i] = StrPtr(str)
	}
	return result
}

func StrVal(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func StrValSlice(v []*string) []string {
	result := make([]string, len(v))
	for i, str := range v {
		result[i] = StrVal(str)
	}
	return result
}

func IntPtr(v int) *int {
	return &v
}

func IncrementIntPtr(v *int) *int {
	return IntPtr(IntVal(v) + 1)
}

func IntPtrSlice(v []int) []*int {
	result := make([]*int, len(v))
	for i, str := range v {
		result[i] = IntPtr(str)
	}
	return result
}

func IntVal(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func IntValSlice(v []*int) []int {
	result := make([]int, len(v))
	for i, str := range v {
		result[i] = IntVal(str)
	}
	return result
}

func Int8Ptr(v int8) *int8 {
	return &v
}

func Int8Val(v *int8) int8 {
	if v == nil {
		return 0
	}
	return *v
}

func Int16Ptr(v int16) *int16 {
	return &v
}

func Int16Val(v *int16) int16 {
	if v == nil {
		return 0
	}
	return *v
}

func Int32Ptr(v int32) *int32 {
	return &v
}

func Int32Val(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

func Int64Ptr(v int64) *int64 {
	return &v
}

func Int64Val(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func UInt8Ptr(v uint8) *uint8 {
	return &v
}

func UInt8Val(v *uint8) uint8 {
	if v == nil {
		return 0
	}
	return *v
}

func UInt16Ptr(v uint16) *uint16 {
	return &v
}

func UInt16Val(v *uint16) uint16 {
	if v == nil {
		return 0
	}
	return *v
}

func UInt32Ptr(v uint32) *uint32 {
	return &v
}

func UInt32Val(v *uint32) uint32 {
	if v == nil {
		return 0
	}
	return *v
}

func UInt64Ptr(v uint64) *uint64 {
	return &v
}

func UInt64Val(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func Float64Ptr(v float64) *float64 {
	return &v
}

func Float64PtrSlice(v []float64) []*float64 {
	result := make([]*float64, len(v))
	for i, str := range v {
		result[i] = Float64Ptr(str)
	}
	return result
}

func Float64Val(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func Float64ValSlice(v []*float64) []float64 {
	result := make([]float64, len(v))
	for i, str := range v {
		result[i] = Float64Val(str)
	}
	return result
}

func BoolPtr(v bool) *bool {
	return &v
}

func BoolPtrSlice(v []bool) []*bool {
	result := make([]*bool, len(v))
	for i, str := range v {
		result[i] = BoolPtr(str)
	}
	return result
}

func BoolVal(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func BoolValSlice(v []*bool) []bool {
	result := make([]bool, len(v))
	for i, str := range v {
		result[i] = BoolVal(str)
	}
	return result
}

func ParseDateType(date string) (time.Time, error) {
	return time.Parse("02012006", date)
}

func SubStr(s string, start int, length int) string {
	end := start + length
	if start > len(s) {
		return ""
	}
	if end > len(s) {
		end = len(s)
	}
	return string([]rune(s)[start:end])
}
