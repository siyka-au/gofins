package main

func FormatBoolAs(b bool, trueVal string, falseVal string) string {
	if b {
		return trueVal
	}
	return falseVal
}
