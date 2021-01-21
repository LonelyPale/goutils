package hmac

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	// 对sha256算法进行hash加密,key随便设置
	hash := hmac.New(sha256.New, []byte("abc123")) // 创建对应的sha256哈希加密算法
	hash.Write([]byte("test"))                     // 写入加密数据
	fmt.Printf("%x\n", hash.Sum(nil))              // c10a04b78bcbcc1c4cba37f6afe0fa60cbf08f6e0a1d93b09387f7069be1aeff

	// 对md5算法进行hash加密，key随便设置
	hash = hmac.New(md5.New, []byte("abc123")) // 创建对应的md5哈希加密算法
	hash.Write([]byte("test"))                 // 写入加密数据
	fmt.Printf("%x\n", hash.Sum(nil))          // 0eee86e484505ec4ab48c18095e6a8ac
}

func TestHash(t *testing.T) {
	hash, err := Hash(sha256.New, []byte("test"), []byte("abc123"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%x\n", hash)
}
