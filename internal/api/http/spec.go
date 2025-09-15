package apihttp

func GetOpenAPIContent() []byte {
	spec, err := rawSpec()
	if err != nil {
		panic(err)
	}

	return spec
}
