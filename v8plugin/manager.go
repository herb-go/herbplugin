package v8plugin

import (
	"errors"

	v8 "github.com/herb-go/v8go"
)

type Releaser interface {
	Release()
}
type Value struct {
	*v8.Value
}
type Manager struct {
	Context *v8.Context
	managed []Releaser
}

func NewManager(ctx *v8.Context) *Manager {
	return &Manager{
		Context: ctx,
		managed: []Releaser{},
	}
}
func (m *Manager) Global() *v8.Object {
	global := m.Context.Global()
	m.manage(global.Value)
	return global
}
func (m *Manager) Release() {
	for _, r := range m.managed {
		r.Release()
	}
}
func (m *Manager) manage(v *v8.Value) *Value {
	m.managed = append(m.managed, v)
	return &Value{
		Value: v,
	}
}
func (m *Manager) newValue(value interface{}) *v8.Value {
	switch v := value.(type) {
	case *v8.Object:
		return v.Value
	case *v8.Value:
		return v
	// case []*v8.Value:
	// 	return MustNewArray(m.Context, v)
	// case []string:
	// 	arr := make([]*v8.Value, len(v))
	// 	for k, val := range v {
	// 		arr[k] = MustNewValue(m.Context, val)
	// 	}
	// 	result := MustNewArray(m.Context, arr)
	// 	for _, val := range arr {
	// 		val.Release()
	// 	}
	// 	return result
	case int:
		value = int64(v)
	}
	val, err := v8.NewValue(m.Context.Isolate(), value)
	if err != nil {
		panic(err)
	}
	return val

}
func (m *Manager) Return(value interface{}) *v8.Value {
	return m.newValue(value)
}
func (m *Manager) NewValue(value interface{}) *Value {
	return m.manage(m.newValue(value))
}
func (m *Manager) NewStringArray(arags []string) *Value {
	arr := make([]v8.Valuer, len(arags))
	for i, v := range arags {
		arr[i] = m.NewValue(v).Value
	}
	result := m.NewArray(arr)
	return result
}
func (m *Manager) NewArray(args []v8.Valuer) *Value {
	global := m.Context.Global()
	defer global.Release()
	array, err := global.Get("Array")
	defer array.Release()
	if err != nil {
		panic(err)
	}
	fnargs := make([]v8.Valuer, len(args))
	for i, v := range args {
		fnargs[i] = v
	}
	return m.Call(array, fnargs...)
}
func (m *Manager) Call(fn *v8.Value, args ...v8.Valuer) *Value {
	f, err := fn.AsFunction()
	if err != nil {
		panic(err)
	}
	result, err := f.Call(fn, args...)
	if err != nil {
		panic(err)
	}
	return m.manage(result)
}

func (m *Manager) GetItem(o *v8.Object, key string) *Value {
	item, err := o.Get(key)
	if err != nil {
		panic(err)
	}
	if item == nil || item.IsNullOrUndefined() {
		return nil
	}
	return m.manage(item)
}

func (m *Manager) ConvertToArray(ctx *v8.Context, val *v8.Value) []*Value {
	if val.IsNull() || val.IsUndefined() {
		return []*Value{}
	}
	if !val.IsArray() {
		panic(errors.New("value is not an array"))
	}

	obj, err := val.AsObject()
	if err != nil {
		panic(err)
	}
	length, err := obj.Get("length")
	defer length.Release()
	if err != nil {
		panic(err)
	}
	l := length.Int32()
	if l < 0 {
		panic(errors.New("array length is negative"))
	}
	result := make([]*Value, l)
	for i := int32(0); i < l; i++ {

		item, err := obj.GetIdx(uint32(i))
		if err != nil {
			panic(err)
		}
		result[int(i)] = m.manage(item)
	}
	return result
}
func (m *Manager) ConvertToStringArray(ctx *v8.Context, val *v8.Value) []string {
	values := m.ConvertToArray(ctx, val)
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = v.String()
	}
	return result
}
