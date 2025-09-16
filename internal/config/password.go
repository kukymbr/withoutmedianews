package config

func maskPasswordVars(vars ...*string) {
	for _, v := range vars {
		*v = maskPassword(*v)
	}
}

func maskPassword(password string) string {
	if password == "" {
		return ""
	}

	return "*****"
}
