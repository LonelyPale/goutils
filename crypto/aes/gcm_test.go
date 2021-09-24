package aes

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testAESKey            = []byte("testAPIv3Key0000")
	testAESPlaintext      = []byte("Hello World")
	testAESNonce          = []byte("wCq51qZv4Yfg")
	testAESAssociatedData = []byte("Dl56FrqWQJF1t9LtC3vEUsniZKvbqdR8")
	testAESCiphertext     = "FsdXzxryWfKwvLJKf8LG/ToRPTRh8RN9wROC"
)

func TestDecryptAes256Gcm(t *testing.T) {
	type args struct {
		apiv3Key       []byte
		associatedData []byte
		nonce          []byte
		ciphertext     string
	}
	tests := []struct {
		name      string
		args      args
		plaintext []byte
		wantErr   bool
	}{
		{
			name: "decrypt success",
			args: args{
				apiv3Key:       testAESKey,
				associatedData: testAESAssociatedData,
				nonce:          testAESNonce,
				ciphertext:     testAESCiphertext,
			},
			wantErr:   false,
			plaintext: testAESPlaintext,
		},
		{
			name: "invalid base64 ciphertext",
			args: args{
				apiv3Key:       testAESKey,
				associatedData: testAESAssociatedData,
				nonce:          testAESNonce,
				ciphertext:     "invalid cipher",
			},
			wantErr: true,
		},
		{
			name: "invalid ciphertext",
			args: args{
				apiv3Key:       testAESKey,
				associatedData: testAESAssociatedData,
				nonce:          testAESNonce,
				ciphertext:     "SGVsbG8gV29ybGQK",
			},
			wantErr: true,
		},
		{
			name: "invalid aes key",
			args: args{
				apiv3Key:       []byte("not a aes key"),
				associatedData: testAESAssociatedData,
				nonce:          testAESNonce,
				ciphertext:     testAESCiphertext,
			},
			wantErr: true,
		},
		{
			name: "wrong aes key",
			args: args{
				apiv3Key:       []byte("testAPIv3Key1111"),
				associatedData: testAESAssociatedData,
				nonce:          testAESNonce,
				ciphertext:     testAESCiphertext,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				decodedCiphertext, err := base64.StdEncoding.DecodeString(tt.args.ciphertext)
				if err != nil {
					t.Log(err)
				}

				plaintext, err := GCMDecrypt(tt.args.apiv3Key, tt.args.nonce, decodedCiphertext, tt.args.associatedData)
				require.Equal(t, tt.wantErr, err != nil)
				assert.Equal(t, tt.plaintext, plaintext)
			},
		)
	}
}
