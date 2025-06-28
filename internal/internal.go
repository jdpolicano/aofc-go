package internal

import "fmt"

func MapSlice[T any, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

func FilterSlice[T any](slice []T, f func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func InRange(i, j, k int) bool {
	return i <= k && k <= j
}

func DeepCopyBytes(original [][]byte) [][]byte {
	if original == nil {
		return nil
	}
	cp := make([][]byte, len(original))
	for i, inner := range original {
		cp[i] = make([]byte, len(inner))
		copy(cp[i], inner)
	}
	return cp
}

type Location [2]int

func (n Location) String() string {
	return fmt.Sprintf("{ row: %d, col: %d }", n.Row(), n.Col())
}

func (n Location) Row() int {
	return n[0]
}

func (n Location) Col() int {
	return n[1]
}
