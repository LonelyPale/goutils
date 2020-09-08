package md5

import (
	"crypto/md5"
	"fmt"
	"io"
	"testing"
)

func Test(t *testing.T) {
	str := "123456"

	//方法一
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	fmt.Println(md5str1)

	//方法二
	w := md5.New()
	if _, err := io.WriteString(w, str); err != nil {
		t.Fatal(err)
	}
	//将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))

	fmt.Println(md5str2)

	//结果
	//e10adc3949ba59abbe56e057f20f883e
	//e10adc3949ba59abbe56e057f20f883e
}
