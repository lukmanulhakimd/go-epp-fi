package registry

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/xml"
	"github.com/pkg/errors"
	"io"
	"net"
	"time"
)

const APIVersion = "1.0"
const APILanguage = "en"

func (s *Client) Connect() error {
	dialConn, err := tls.Dial("tcp", s.registryServer, &s.tlsConfig)
	if err != nil {
		return err
	}
	s.conn = dialConn

	greet, err := s.Read()
	if err != nil {
		return err
	}

	s.Greeting, err = unmarshalGreeting(greet)
	if err != nil {
		return err
	}

	if s.Greeting.SvcMenu.Version != APIVersion {
		return errors.New("Unexpected version: " + s.Greeting.SvcMenu.Version)
	}

	return nil
}

func (s *Client) Read() ([]byte, error) {
	var rawResponse int32

	if s.readTimeout > 0 {
		s.conn.SetReadDeadline(time.Now().Add(s.readTimeout))
	}

	err := binary.Read(s.conn, binary.BigEndian, &rawResponse)
	if err != nil {
		return nil, err
	}

	rawResponse -= 4
	if rawResponse < 0 {
		return nil, io.ErrUnexpectedEOF
	}

	bytesResponse, err := readStreamToBytes(s.conn, rawResponse)
	if err != nil {
		return nil, err
	}

	return bytesResponse, nil
}

func (s *Client) Write(payload []byte) error {
	payload = []byte(xml.Header + string(payload))

	sendBytesLength := uint32(4 + len(payload))

	if s.writeTimeout > 0 {
		s.conn.SetWriteDeadline(time.Now().Add(s.writeTimeout))
	}

	err := binary.Write(s.conn, binary.BigEndian, sendBytesLength)
	if err != nil {
		return err
	}
	if _, err = s.conn.Write(payload); err != nil {
		return err
		// TODO log first param (amount of bytes written) if error
	}

	return nil
}

func (s *Client) Send(payload []byte) ([]byte, error) {
	s.log.Debug("Sending message:\n" + string(payload))
	err := s.Write(payload)
	if err != nil {
		// TODO log error
		return nil, err
	}

	time.Sleep(s.sendWaitTime)

	apiResp, err := s.Read()
	if err != nil {
		// TODO log error
		return nil, err
	}

	s.log.Debug("Received response:\n" + string(apiResp))

	return apiResp, nil
}

func (s *Client) Close() error {
	if err := s.conn.Close(); err != nil {
		return err
	}

	s.conn = nil
	return nil
}

func readStreamToBytes(conn net.Conn, rawResponse int32) ([]byte, error) {
	lr := io.LimitedReader{R: conn, N: int64(rawResponse)}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(&lr); err != nil {
		return nil, err
		// TODO log first param (amount of bytes read) if error
	}
	return buf.Bytes(), nil
}
