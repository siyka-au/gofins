package main

// FormatBoolAs Format a boolean into one of two string values for display
func FormatBoolAs(b bool, trueVal string, falseVal string) string {
	if b {
		return trueVal
	}
	return falseVal
}
