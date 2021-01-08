package websocket

import "github.com/go-spring/spring-utils"

type Server struct {
	containers []*Container //WS容器列表
}

func NewServer() *Server {
	return &Server{
		containers: make([]*Container, 0),
	}
}

// Containers 返回 Container 实例列表
func (s *Server) Containers() []*Container {
	return s.containers
}

// AddContainer 添加 Container 实例
func (s *Server) AddContainer(container ...*Container) *Server {
	s.containers = append(s.containers, container...)
	return s
}

// Start 启动 WS 容器，非阻塞调用
func (s *Server) Start() {
	for _, c := range s.containers {
		c.Start()
	}
}

// Stop 停止 WS 容器，阻塞调用
func (s *Server) Stop() {
	var wg SpringUtils.WaitGroup
	for _, container := range s.containers {
		c := container //避免延迟绑定
		wg.Add(func() {
			if err := c.Stop(); err != nil {
				DefaultLogger.Error(err)
			}
		})
	}
	wg.Wait()
}
