package v8plugin

import (
	"errors"

	v8 "github.com/herb-go/v8go"
)

type Releaser interface {
	Release()
}
type Manager struct {
	Context *v8.Context
	managed []Releaser
}

func MustStringArray(ctx *v8.Context, args []string) *v8.Value {
	arr := make([]v8.Valuer, len(args))

	for i, v := range args {
		var item = mustNewValue(ctx, v)
		arr[i] = item
		defer item.Release()
	}
	return MustNewArray(ctx, arr)
}
func MustNewArray(ctx *v8.Context, args []v8.Valuer) *v8.Value {
	global := ctx.Global()
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
	return MustCall(array, fnargs...)
}
func MustCall(fn *v8.Value, args ...v8.Valuer) *v8.Value {
	f, err := fn.AsFunction()
	if err != nil {
		panic(err)
	}
	result, err := f.Call(fn, args...)
	if err != nil {
		panic(err)
	}
	return result
}
func MustReturn(call *v8.FunctionCallbackInfo, value interface{}) *v8.Value {
	call.Release()
	if value == nil {
		return nil
	}
	return mustNewValue(call.Context(), value)
}
func mustNewValue(ctx *v8.Context, value interface{}) *v8.Value {
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
	val, err := v8.NewValue(ctx.Isolate(), value)
	if err != nil {
		panic(err)
	}
	return val
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
func (m *Manager) manage(v *v8.Value) *v8.Value {
	m.managed = append(m.managed, v)
	return v

}
func (m *Manager) newValue(value interface{}) *v8.Value {
	return mustNewValue(m.Context, value)
}
func (m *Manager) NewValue(value interface{}) *v8.Value {
	return m.manage(m.newValue(value))
}
func (m *Manager) NewArray(arags []v8.Valuer) *v8.Value {
	return m.manage(MustNewArray(m.Context, arags))
}
func (m *Manager) NewStringArray(arags []string) *v8.Value {
	return m.manage(MustStringArray(m.Context, arags))
}
func (m *Manager) Call(fn *v8.Value, args ...v8.Valuer) *v8.Value {
	return m.manage(MustCall(fn, args...))
}

func (m *Manager) GetItem(o *v8.Object, key string) *v8.Value {
	item, err := o.Get(key)
	if err != nil {
		panic(err)
	}
	if item == nil || item.IsNullOrUndefined() {
		return nil
	}
	return m.manage(item)
}

func (m *Manager) ConvertToArray(ctx *v8.Context, val *v8.Value) []*v8.Value {
	if val.IsNull() || val.IsUndefined() {
		return []*v8.Value{}
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
	result := make([]*v8.Value, l)
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
