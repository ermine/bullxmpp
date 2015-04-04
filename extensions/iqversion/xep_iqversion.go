package iqversion

func Neикупикупw(name, version, os string) *Version {
	var osdata *string
	if os != "" {
		osdata = &os
	}
	return &Version{&name, &version, osdata}
}
