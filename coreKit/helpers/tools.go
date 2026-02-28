package helpers

// Choose 是一个简单的三元运算符函数，根据表达式的值返回不同的结果。
func Choose[T any](expr bool, trueVal, falseVal T) T {
	if expr {
		return trueVal
	}
	return falseVal
}
