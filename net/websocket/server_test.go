package websocket

import (
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

var testServer = NewServer()

func newTestContainer() *Container {
	return NewContainer(NewProcessor, DefaultConfig())
}

func TestWebSocket(t *testing.T) {
	container1 := testContainer1()
	container2 := testContainer2()
	testServer.AddContainer(container1, container2).Start()

	engine := gin.Default()
	engine.Handle(http.MethodGet, "/ws1", container1.OpenConnGin)
	engine.Handle(http.MethodGet, "/ws2", container2.OpenConnGin)

	//go func() {
	//	for {
	//		time.Sleep(time.Second * 3)
	//		msg := &Message{
	//			Type: "test.push",
	//			SN:   0,
	//			Code: 0,
	//			Msg:  "ok",
	//			Data: fmt.Sprintf("你好，杭州。%s", time.Now().Format(time.DateTime)),
	//		}
	//		if err := container1.Send("test.push", msg); err != nil {
	//			t.Error(err)
	//		}
	//	}
	//}()

	if err := engine.Run(); err != nil {
		t.Fatal(err)
	}
}

func testContainer1() *Container {
	container := newTestContainer()
	container.HandleFunc("test.other", func(conn *Conn, message *Message) {
		sendWSMessage(conn)
	})
	return container
}

func testContainer2() *Container {
	container := newTestContainer()
	container.Handle("test.other", BIND(testOtherHandler))
	container.Handle("test.struct", BIND(testStructHandler))
	return container
}

func testOtherHandler(conn *Conn, message *Message, num string) {
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
