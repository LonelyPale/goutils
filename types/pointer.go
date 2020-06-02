package types

func String(s string) *string {
	return &s
}

func Bool(b bool) *bool {
	return &b
}

func True() *bool {
	b := true
	return &b
}

func False() *bool {
	b := false
	return &b
}

func Int(n int) *int {
	return &n
}

func Uint(u uint) *uint {
	return &u
}
