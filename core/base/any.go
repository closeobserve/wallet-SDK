package base

import (
	"encoding/json"
	"math/big"
)

// exchange `Aniable Object` & `Any`
type Aniable interface {
	AsAny() *Any

	// You need to implement the following methods if your class name is Xxx
	// func AsXxx(a *Any) *Xxx

	// ======= template
	// func (o *Xxx) AsAny() *base.Any {
	// 	return &base.Any{Value: o}
	// }
	//	func AsXxx(a *base.Any) *Xxx {
	//		if r, ok := a.Value.(*Xxx); ok {
	//			return r
	//		}
	//		if r, ok := a.Value.(Xxx); ok {
	//			return &r
	//		}
	//		return nil
	//	}
}

// 如果需要自定义类型支持 Any, 需要遵循协议 Aniable
type Any struct {
	Value any
}

func NewAny() *Any {
	return &Any{}
}

// `Any` only support Marshal
func (a Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Value)
}

func (a *Any) JsonString() (*OptionalString, error) {
	return JsonString(a)
}

func (a *Any) SetString(v string)  { a.Value = v }
func (a *Any) SetBool(v bool)      { a.Value = v }
func (a *Any) SetInt(v int)        { a.Value = v }
func (a *Any) SetInt8(v int8)      { a.Value = v }
func (a *Any) SetInt16(v int16)    { a.Value = v }
func (a *Any) SetInt32(v int32)    { a.Value = v }
func (a *Any) SetInt64(v int64)    { a.Value = v }
func (a *Any) SetUInt8(v *BigInt)  { n := uint8(v.bigint.Uint64()); a.Value = &n }
func (a *Any) SetUInt16(v *BigInt) { n := uint16(v.bigint.Uint64()); a.Value = &n }
func (a *Any) SetUInt32(v *BigInt) { n := uint32(v.bigint.Uint64()); a.Value = &n }
func (a *Any) SetUInt64(v *BigInt) { n := v.bigint.Uint64(); a.Value = &n }
func (a *Any) SetBigInt(v *BigInt) { a.Value = v }

func (a *Any) GetString() string  { return a.Value.(string) }
func (a *Any) GetBool() bool      { return a.Value.(bool) }
func (a *Any) GetInt() int        { return a.Value.(int) }
func (a *Any) GetInt8() int8      { return a.Value.(int8) }
func (a *Any) GetInt16() int16    { return a.Value.(int16) }
func (a *Any) GetInt32() int32    { return a.Value.(int32) }
func (a *Any) GetInt64() int64    { return a.Value.(int64) }
func (a *Any) GetUInt8() *BigInt  { return &BigInt{new(big.Int).SetUint64(uint64(*a.Value.(*uint8)))} }
func (a *Any) GetUInt16() *BigInt { return &BigInt{new(big.Int).SetUint64(uint64(*a.Value.(*uint16)))} }
func (a *Any) GetUInt32() *BigInt { return &BigInt{new(big.Int).SetUint64(uint64(*a.Value.(*uint32)))} }
func (a *Any) GetUInt64() *BigInt { return &BigInt{new(big.Int).SetUint64(*a.Value.(*uint64))} }
func (a *Any) GetBigInt() *BigInt { return a.Value.(*BigInt) }

type AnyArray struct {
	Values []any
}

func NewAnyArray() *AnyArray {
	return &AnyArray{Values: make([]any, 0)}
}

// `AnyArray` only support Marshal
func (a AnyArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Values)
}

func (a *AnyArray) JsonString() (*OptionalString, error) {
	return JsonString(a)
}

func (a *AnyArray) Count() int {
	return len(a.Values)
}

func (a *AnyArray) Append(any *Any) {
	a.Values = append(a.Values, any.Value)
}

func (a *AnyArray) Remove(index int) {
	a.Values = append(a.Values[:index], a.Values[index+1:]...)
}

func (a *AnyArray) SetValue(value *Any, index int) {
	a.Values[index] = value.Value
}

func (a *AnyArray) Contains(any *Any) bool {
	return a.IndexOf(any) != -1
}

// return -1 if not found
func (a *AnyArray) IndexOf(any *Any) int {
	for idx, item := range a.Values {
		if item == any.Value {
			return idx
		}
	}
	return -1
}

func (a *AnyArray) ValueOf(index int) *Any {
	return &Any{Value: a.Values[index]}
}

func (a *AnyArray) String() string {
	data, err := json.Marshal(a.Values)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func (a *AnyArray) AsAny() *Any {
	return &Any{a.Values}
}

func AsAnyArray(a *Any) *AnyArray {
	if res, ok := a.Value.([]any); ok {
		return &AnyArray{res}
	}
	return nil
}

type AnyMap struct {
	Values map[string]any
}

func NewAnyMap() *AnyMap {
	return &AnyMap{Values: make(map[string]any)}
}

// `AnyMap` only support Marshal
func (a AnyMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Values)
}

func (a *AnyMap) JsonString() (*OptionalString, error) {
	return JsonString(a)
}

func (a *AnyMap) ValueOf(key string) *Any {
	if v, ok := a.Values[key]; ok {
		return &Any{v}
	}
	return nil
}

func (a *AnyMap) SetValue(value *Any, key string) {
	a.Values[key] = value.Value
}

func (a *AnyMap) Remove(key string) *Any {
	if v, ok := a.Values[key]; ok {
		delete(a.Values, key)
		return &Any{v}
	}
	return nil
}

func (a *AnyMap) HasKey(key string) bool {
	_, ok := a.Values[key]
	return ok
}

func (a *AnyMap) Keys() *StringArray {
	keys := make([]string, len(a.Values))
	i := 0
	for k := range a.Values {
		keys[i] = k
		i++
	}
	return &StringArray{Values: keys}
}

func (a *AnyMap) String() string {
	data, err := json.Marshal(a.Values)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func (a *AnyMap) AsAny() *Any {
	return &Any{a.Values}
}

func AsAnyMap(a *Any) *AnyMap {
	if res, ok := a.Value.(map[string]any); ok {
		return &AnyMap{res}
	}
	return nil
}
