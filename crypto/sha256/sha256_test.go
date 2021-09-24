package sha256

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
	"time"
)

func Test(t *testing.T) {
	startTime := time.Now()
	data := []byte("123qweasdzxc")
	var err error
	for i := 0; i < 2048; i++ {
		data, err = Hash(data)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Log("duration:", time.Since(startTime), data)
}

func TestHash(t *testing.T) {
	hashed, err := Hash([]byte("123qweasdzxc"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashed)
	t.Log(string(hashed))
}

func TestHex(t *testing.T) {
	hashed, err := Hex([]byte("123qweasdzxc"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashed)
}

func TestBase64(t *testing.T) {
	hashed, err := Base64([]byte("123qweasdzxc"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashed)
}

func TestHashReader(t *testing.T) {
	go func() {
		for {
			traceMemStats(t)
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		startTime := time.Now()
		file, err := os.Open("/Users/wyb/backup/software/os/ubuntu-20.04.2-live-server-amd64.iso")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		t.Log("Open duration:", time.Since(startTime))

		hash, err := HashReader(file)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("sha256:", hash.Hex(), "duration:", time.Since(startTime))

		//复原文件指针位置
		ret, err := file.Seek(0, io.SeekStart)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Seek:", ret, "duration:", time.Since(startTime))
	}()

	select {}
}

func TestFile1(t *testing.T) {
	go func() {
		for {
			traceMemStats(t)
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		startTime := time.Now()
		bs, err := ioutil.ReadFile("/Users/wyb/backup/software/os/ubuntu-20.04.2-live-server-amd64.iso")
		if err != nil {
			t.Fatal(err)
		}
		t.Log("ReadFile duration:", time.Since(startTime))

		hash := sha256.New()
		if _, err := hash.Write(bs); err != nil {
			t.Fatal(err)
		}
		sum := hash.Sum(nil)
		t.Log("sha256:", hex.EncodeToString(sum), "duration:", time.Since(startTime))
	}()

	time.Sleep(time.Hour)
}

func TestFile2(t *testing.T) {
	go func() {
		for {
			traceMemStats(t)
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		startTime := time.Now()
		file, err := os.Open("/Users/wyb/backup/software/os/ubuntu-20.04.2-live-server-amd64.iso")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		t.Log("Open duration:", time.Since(startTime))

		hash := sha256.New()
		//buff := make([]byte, 1024*1024)
		buff := make([]byte, 4096*32)
		for {
			_, err := file.Read(buff)
			if err != nil && err != io.EOF {
				t.Fatal(err)
			}
			if err == io.EOF {
				break
			}

			if _, err := hash.Write(buff); err != nil {
				t.Fatal(err)
			}
		}

		sum := hash.Sum(nil)
		t.Log("sha256:", hex.EncodeToString(sum), "duration:", time.Since(startTime))
	}()

	select {}
}

func TestFile3(t *testing.T) {
	go func() {
		for {
			traceMemStats(t)
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		startTime := time.Now()
		file, err := os.Open("/Users/wyb/backup/software/os/ubuntu-20.04.2-live-server-amd64.iso")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		t.Log("Open duration:", time.Since(startTime))

		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			t.Fatal(err)
		}

		sum := hash.Sum(nil)
		t.Log("sha256:", hex.EncodeToString(sum), "duration:", time.Since(startTime))
	}()

	select {}
}

func traceMemStats(t *testing.T) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	mb := 1024 * 1024.0
	logstr := fmt.Sprintf("Alloc=%.2fMB\tTotalAlloc=%.2fMB\tSys=%.2fMB\tNumGC=%v",
		float64(ms.Alloc)/mb, float64(ms.TotalAlloc)/mb, float64(ms.Sys)/mb, ms.NumGC)
	t.Log(logstr)
}
