package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

// Telnet .
type Telnet struct {
	addr    string
	timeout time.Duration
	conn    net.Conn
	notify  chan error
}

var (
	host    string
	port    string
	timeout time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")

	flag.Parse()
	host = flag.Arg(0)
	port = flag.Arg(1)
}

func main() {
	t := NewTelnet(host, port, timeout)
	t.Run()
}

// NewTelnet .
func NewTelnet(host, port string, timeout time.Duration) *Telnet {
	return &Telnet{
		addr:    fmt.Sprintf("%s:%s", host, port),
		timeout: timeout,
		notify:  make(chan error),
	}
}

// Run .
func (t *Telnet) Run() {
	if err := t.connect(); err != nil {
		fmt.Printf("telnet: %s\n", err.Error())
		return
	}

	go t.sender()
	go t.receiver()

	select {
	case err := <-t.notify:
		switch {
		case errors.Is(err, io.EOF):
			fmt.Printf("telnet: server has closed the connection %s\n", t.addr)
		default:
			fmt.Printf("telnet: %s\n", err.Error())
		}
	}

	err := t.close()
	if err != nil {
		fmt.Printf("telnet: %s\n", err.Error())
	}
}

func (t *Telnet) connect() error {
	conn, err := net.DialTimeout("tcp", t.addr, t.timeout)
	if err != nil {
		return fmt.Errorf("Run - connect - conn.Dial: %w", err)
	}

	t.conn = conn

	return nil
}

func (t *Telnet) close() error {
	err := t.conn.Close()
	if err != nil {
		return err
	}

	close(t.notify)

	return nil
}

func (t *Telnet) sender() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		data, err := inputReader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			t.notify <- fmt.Errorf("receive signal to close the connection")
		}
		if err != nil {
			t.notify <- err
		}

		if _, err := t.conn.Write(data); err != nil {
			t.notify <- err
		}
	}
}

func (t *Telnet) receiver() {
	for {
		buf := make([]byte, 1024)
		_, err := t.conn.Read(buf)
		if err != nil {
			t.notify <- err
		}

		if _, err := fmt.Fprint(os.Stdout, string(buf)); err != nil {
			t.notify <- err
		}
	}
}
