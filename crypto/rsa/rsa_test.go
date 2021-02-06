package rsa

import (
	"bytes"
	"testing"
)

var data = []byte("123qwe你好杭州")

func TestXRsa(t *testing.T) {
	pubKey := bytes.NewBuffer(make([]byte, 0))
	priKey := bytes.NewBuffer(make([]byte, 0))

	if err := CreateKeys(pubKey, priKey, Bits); err != nil {
		t.Fatal(err)
	}

	xrsa, err := NewXRsa(pubKey.Bytes(), priKey.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	encode, err := xrsa.EncryptToBase64(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("encode:", encode)

	decode, err := xrsa.DecryptFromBase64(encode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("decode:", string(decode))

	signcode, err := xrsa.SignToBase64(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Sign:", signcode)

	if err := xrsa.VerifyFromBase64(data, signcode); err != nil {
		t.Fatal(err)
	}
	t.Log("Verify:", err)
}

func TestXRsaKeys(t *testing.T) {
	pub, pri, err := CreateKeysToBase64()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", pub)
	t.Log("\n", pri)
}

func TestXRsa2(t *testing.T) {
	pubKey := `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+1RGNdS9MYgvDr+9crtr
irWw9yuxdP/St0sq9knqK2UdzC8/qWPStbT3TFouLIUTveKND9X1VylKsW2ut+L5
St6YWVUJQFDYIMV/j79jEGOxs2+0w7uoi2PB3/4hXYjGQwAYLIBH8yA06yaX5/UU
mUZqsMFAlg+QSBo0sbKcgZSw6kIqBsW3PEUC7iiigz1o3wCeJNkEABG6ur/GvRcS
9yObrWnBDOuJDjH+ONRqUvk0v4NRq9ythUjQ0/6adK4fJB3EYJutGgB0xDsXL2cm
NeW2p0+DpxWi3Hr/RXvKRZyj3pnipsUlUHJarm8agjwt2tZfmNDg5+ymNaAX6DjI
7wIDAQAB
-----END PUBLIC KEY-----
`
	priKey := `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA+1RGNdS9MYgvDr+9crtrirWw9yuxdP/St0sq9knqK2UdzC8/
qWPStbT3TFouLIUTveKND9X1VylKsW2ut+L5St6YWVUJQFDYIMV/j79jEGOxs2+0
w7uoi2PB3/4hXYjGQwAYLIBH8yA06yaX5/UUmUZqsMFAlg+QSBo0sbKcgZSw6kIq
BsW3PEUC7iiigz1o3wCeJNkEABG6ur/GvRcS9yObrWnBDOuJDjH+ONRqUvk0v4NR
q9ythUjQ0/6adK4fJB3EYJutGgB0xDsXL2cmNeW2p0+DpxWi3Hr/RXvKRZyj3pni
psUlUHJarm8agjwt2tZfmNDg5+ymNaAX6DjI7wIDAQABAoIBAG5Tq1J8JeU1c/rU
frR7w5Srd5i5LHhAyN4/eAePoOUVyUPVxet0741meFyjBbvzWxwy9FtsP/vYG3rR
vC3qzCZqmpQ0eRArrQSNzhZrHHCYYuxy7/YwTaeKzhOm+jyWCvpkczgtw+fdTn2f
fnWLp1Q1jiYzX0fAY0TThCgxUPSsU13posnof47t7nWjsv2qYY6fshsdGPUJS+08
xTsYhOvo5mIfAXNVNYublNLbHe+nL80lAwH7xgegE+xlhawPNvNIIkhHvg5KsEng
fN0mHFWjMgobyrSsMQmGbrtEiMnF2Q+Clzr8UMGO2CmkskHRGoMB2//Hx32sCVaB
bV7UHnECgYEA/tIF9ZsJFjKlLQDegHTPBoB0S3smRpApVo1Eh3sAU3uz9otlMZl3
JlnEX4AB9wgvQT7GWRgRj8M55tpyNMhVBWSsnOJXQtxJbxC7iV/B6Mjm/i/hBK8f
dlJhjhOEFMb2SRXizhZ0ZWEB1jesFt+KcBi47agjk0qy6MIV1e48S9cCgYEA/H4d
G3b20iT5UK9hrDqAginbqf5/sFinYSN25oYaIRXsipRlLwh4rPtf4sYHJiBuWu+L
Jap3+k2RRYkqbKdCVpUMiMbH16W9g1Me9qEf39M0/H4bYs0CdyEBGkKEMo0RjlIe
lvBemRqAtMvJsqVg4MFEMF7cGQm6kCT3/F7UCKkCgYBlRZb0u37rAYm/zv6e0s2M
afTOIs1dceHb8hzwMyQ4CYvGSjQXeERwS3DN+5PMV8ZgCdDOi9A+8HnMk7ib3Zpc
oFwxpYrEmcPdjiraN+Ja361eDC1DrU21upvm2T6++yvadAZFnYr414rhVhLgrEra
rhig6xfoa+Gau7ft49a21QKBgQCJ+NAnBebyBkpGkM+qsX0vo3fpeKyFzKwKJLsA
VR9KHRTY1SZFgTeQLvzCiru1Vdt3zZYXywMsv942RTHtlahmb6QdyaHCcUsRzAYL
dxhX0q4Nm0uTvbsvJdXYZ6idhwCk6LLWgBrxRs41/XYGLOC3cGS2md9jvzE3OzxX
p+ntoQKBgAlMnSHHI2AUY3JBbEpqz8wyXY9mGkprhVqZ7vnyDH6AexzIQShXuYHy
YCApzxi6CwLtveNWIVpI6971QaXaQhTekcTCvJfpwB2LE9L7jJRfyX0T3Xgw++LU
k7oCULOB7/tnPrtSWnpT7Tx6u6hxVMuE6D/b5o+88K+MqJuWAh5E
-----END RSA PRIVATE KEY-----
`

	xrsa, err := NewXRsa([]byte(pubKey), []byte(priKey))
	if err != nil {
		t.Fatal(err)
	}

	encode, err := xrsa.EncryptToBase64(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("encode:", encode)

	decode, err := xrsa.DecryptFromBase64(encode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("decode:", string(decode))

	signcode, err := xrsa.SignToBase64(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Sign:", signcode)

	if err := xrsa.VerifyFromBase64(data, signcode); err != nil {
		t.Fatal(err)
	}
	t.Log("Verify:", err)
}
