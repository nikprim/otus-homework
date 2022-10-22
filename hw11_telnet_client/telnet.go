package main

import (
	"errors"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *telnetClient) Connect() error {
	if t.conn != nil {
		return nil
	}

	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}

	t.conn = conn

	return nil
}

func (t *telnetClient) Close() error {
	return t.conn.Close()
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return errors.New("need to connect first")
	}

	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return err
	}

	if _, err = os.Stderr.Write([]byte("...EOF\n")); err != nil {
		return err
	}

	return nil
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return errors.New("need to connect first")
	}

	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return err
	}

	if _, err = os.Stderr.Write([]byte("...Connection was closed by peer\n")); err != nil {
		return err
	}

	return nil
}
