package util

func MaskPasswordVars(vars ...*string) {
	for _, v := range vars {
		*v = MaskPassword(*v)
	}
}

// MaskPassword masks password if not empty.
func MaskPassword(password string) string {
	if password == "" {
		return ""
	}

	return "*****"
}
