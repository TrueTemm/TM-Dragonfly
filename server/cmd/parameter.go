/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package cmd

import (
	"reflect"
	"strings"

	"github.com/go-gl/mathgl/mgl64"
)

type Parameter interface {
	Parse(line *Line, v reflect.Value) error

	Type() string
}

type Enum interface {
	Type() string

	Options(source Source) []string
}

type Suggestion interface {
	Enum

	Suggestion()
}

type ParamDescriber interface {
	DescribeParams(src Source) []ParamInfo
}

type SubCommand struct{}

type Varargs string

type Optional[T any] struct {
	val T
	set bool
}

func (o Optional[T]) Load() (T, bool) {
	return o.val, o.set
}

func (o Optional[T]) LoadOr(or T) T {
	if o.set {
		return o.val
	}
	return or
}

func (o Optional[T]) with(val any) any {
	return Optional[T]{val: val.(T), set: true}
}

type optionalT interface {
	with(val any) any
}

func typeNameOf(i any, name string) string {
	switch i.(type) {
	case int, int8, int16, int32, int64:
		return "int"
	case uint, uint8, uint16, uint32, uint64:
		return "uint"
	case float32, float64:
		return "float"
	case string:
		return "string"
	case Varargs:
		return "text"
	case bool:
		return "bool"
	case mgl64.Vec3:
		return "x y z"
	case []Target:
		return "target"
	case SubCommand:
		return name
	}
	if param, ok := i.(Parameter); ok {
		return param.Type()
	}
	if enum, ok := i.(Enum); ok {
		return enum.Type()
	}
	return "value"
}

func unwrap(v reflect.Value) reflect.Value {
	if _, ok := v.Interface().(optionalT); ok {
		return reflect.New(v.Field(0).Type()).Elem()
	}
	return v
}

func optional(v reflect.Value) bool {
	_, ok := v.Interface().(optionalT)
	return ok
}

func suffix(v reflect.StructField) string {
	_, str := tag(v)
	return str
}

func name(v reflect.StructField) string {
	str, _ := tag(v)
	if str == "" {
		return v.Name
	}
	return str
}

func tag(v reflect.StructField) (name string, suffix string) {
	t, _ := v.Tag.Lookup("cmd")
	a, b, _ := strings.Cut(t, ",")
	return a, b
}
