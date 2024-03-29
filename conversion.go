package osutil

// AnySliceToTypeSlice returns a slice of the designated type with the values from the given "any" slice that match the type.
func AnySliceToTypeSlice[T any](anySlice []any) (typeSlice []T) {
	for _, v := range anySlice {
		if x, ok := v.(T); ok {
			typeSlice = append(typeSlice, x)
		}
	}

	return typeSlice
}
