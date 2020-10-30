package types

func String(s string) *string {
	return &s
}

func Bool(b bool) *bool {
	return &b
}

func True() *bool {
	return Bool(true)
}

func False() *bool {
	return Bool(false)
}

func Int(n int) *int {
	return &n
}

func Uint(u uint) *uint {
	return &u
}
