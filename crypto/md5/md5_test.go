package md5

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"time"
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

func TestBigFile(t *testing.T) {
	t1 := time.Now()
	bs, err := ioutil.ReadFile("/Users/wyb/backup/software/os/ubuntu-20.04.2-live-server-amd64.iso")
	if err != nil {
		t.Fatal(err)
	}
	elapsed := time.Since(t1)
	t.Log(len(bs), elapsed)

	t1 = time.Now()
	w := md5.New()
	_, err = w.Write(bs)
	if err != nil {
		t.Fatal(err)
	}

	md5str := fmt.Sprintf("%x", w.Sum(nil))
	fmt.Println(md5str, time.Since(t1))
}
