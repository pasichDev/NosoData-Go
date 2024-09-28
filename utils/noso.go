package utils

import "fmt"

func ToNoso(n int64) string {
	var noso = float32(n) / 10e7
	return fmt.Sprintf("%.8f Noso", noso)
}
