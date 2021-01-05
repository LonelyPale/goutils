package websocket

import (
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func newServer() *Server {
	opts := &Options{
		Enable:          true,
		Origin:          true,
		ReadDeadline:    0,
		WriteDeadline:   10,
		ReadBufferSize:  20480,
		WriteBufferSize: 20480,
		ReadPoolSize:    10,
		WritePoolSize:   10,
		MaxMessageSize:  65535,
	}
	return NewServer(NewReader, NewWriter, opts)
}

func TestHandleFunc(t *testing.T) {
	server := newServer()
	server.ReaderHandleFunc("test", func(conn *Conn, message *Message) {
		sendWSMessage(conn)
	})
	go server.Run()

	engine := gin.Default()
	engine.Handle(http.MethodGet, "/", server.Open)

	if err := engine.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestBIND(t *testing.T) {
	server := newServer()
	server.ReaderHandle("test.other", BIND(testOtherHandler))
	server.ReaderHandle("test.struct", BIND(testStructHandler))
	//server.WriterHandle("test.other", BIND(writerHandler))
	//server.WriterHandle("test.struct", BIND(writerHandler))
	go server.Run()

	engine := gin.Default()
	engine.Handle(http.MethodGet, "/", server.Open)

	if err := engine.Run(); err != nil {
		t.Fatal(err)
	}
}

func testOtherHandler(conn *Conn, message *Message, num *string) {
	log.Println("test.other:", conn, message, num)

	sendWSMessage(conn)
}

func testStructHandler(user struct {
	Name string
	Age  int
}) interface{} {
	log.Println("test.struct", user)
	return user
	//return "user"
}

func sendWSMessage(conn *Conn) {
	if err := conn.Write(&WSMessage{
		Type: 1,
		Data: []byte("test: sendWSMessage: 你好中国🇨🇳"),
	}); err != nil {
		panic(err)
	}
}

func writerHandler(msg *Message) {
	log.Println("write:", msg)
}
