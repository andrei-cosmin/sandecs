package marker

import "reflect"

type Marker reflect.Type

func Symbol[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().Name()
}

func Of[T any](component T) string {
	return reflect.TypeOf(component).Name()
}
