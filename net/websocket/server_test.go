package websocket

import (
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func newServer() *Server {
	return NewServer(NewProcessor, DefaultOptions())
}

func TestHandleFunc(t *testing.T) {
	server := newServer()
	server.HandleFunc("test", func(conn *Conn, message *Message) {
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
	server.Handle("test.other", BIND(testOtherHandler))
	server.Handle("test.struct", BIND(testStructHandler))
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
