package kubernetes

import (
	"errors"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"ops-api/kubernetes"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var PodTerminal podTerminal

type podTerminal struct{}

type streamHandler struct {
	ws          *websocket.Conn
	inputMsg    chan []byte
	resizeEvent chan remotecommand.TerminalSize
}

func (p *podTerminal) Init(namespace, podName, containerName string, cols, rows int, ws *websocket.Conn, client *kubernetes.ClientList) error {
	req := client.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(client.RestConfig, "POST", req.URL())
	if err != nil {
		return err
	}

	handler := &streamHandler{
		ws:          ws,
		inputMsg:    make(chan []byte, 1024),
		resizeEvent: make(chan remotecommand.TerminalSize, 1),
	}
	handler.resizeEvent <- remotecommand.TerminalSize{Width: uint16(cols), Height: uint16(rows)}

	go executeTask(handler)

	go func() {
		<-time.After(30 * time.Minute)
		_ = ws.WriteMessage(websocket.TextMessage, []byte("会话超时，退出..."))
		_ = ws.Close()
	}()

	if err := exec.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		Tty:               true,
		TerminalSizeQueue: handler,
	}); err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("启动终端失败："+err.Error()))
		return err
	}

	go func() {
		<-time.After(30 * time.Minute)
		_ = ws.WriteMessage(websocket.TextMessage, []byte("会话超时，退出..."))
		_ = ws.Close()
	}()
	return nil
}

func (h *streamHandler) Next() *remotecommand.TerminalSize {
	resize := <-h.resizeEvent
	return &resize
}

func (h *streamHandler) Read(p []byte) (int, error) {
	data, ok := <-h.inputMsg
	if !ok {
		return 0, errors.New("I/O data reading failed")
	}
	copy(p, data)
	return len(data), nil
}

func (h *streamHandler) Write(p []byte) (int, error) {
	if !utf8.Valid(p) {
		bufStr := string(p)
		buf := make([]rune, 0, len(bufStr))
		for _, r := range bufStr {
			if r == utf8.RuneError {
				buf = append(buf, []rune("@")...)
			} else {
				buf = append(buf, r)
			}
		}
		p = []byte(string(buf))
	}
	err := h.ws.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}

func executeTask(h *streamHandler) {
	for {
		_, msg, err := h.ws.ReadMessage()
		if err != nil {
			return
		}
		if string(msg) == "ping" {
			continue
		}
		if strings.Contains(string(msg), "resize") {
			resizeSlice := strings.Split(string(msg), ":")
			rows, _ := strconv.Atoi(resizeSlice[1])
			cols, _ := strconv.Atoi(resizeSlice[2])
			h.resizeEvent <- remotecommand.TerminalSize{
				Width:  uint16(cols),
				Height: uint16(rows),
			}
			continue
		}
		h.inputMsg <- msg
	}
}
