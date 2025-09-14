package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"unicode/utf8"
)

const (
	Success = int32(iota)
	Failure
)

const (
	TypeRequest = int32(iota + 1)
	TypeResponse
)

const maxPaddingByte = 4

type Request struct {
	Type int32
	Str  string
}

type Response struct {
	Type    int32
	ErrorNo int32
	Buffer  []byte
}

func ReceiveRequest(r io.Reader) (Request, error) {
	requestType, err := readInt32(r)
	if err != nil {
		return Request{}, fmt.Errorf("reading request type: %w", err)
	}

	if requestType != TypeRequest {
		return Request{}, errors.New("invalid request type")
	}

	str, err := readString(r)
	if err != nil {
		return Request{}, fmt.Errorf("reading string: %w", err)
	}

	return Request{
		Type: requestType,
		Str:  str,
	}, nil
}

func SendResponse(w io.Writer, response Response) error {
	if err := writeInt32(w, response.Type); err != nil {
		return fmt.Errorf("writing response type: %w", err)
	}

	if err := writeInt32(w, response.ErrorNo); err != nil {
		return fmt.Errorf("writing error number: %w", err)
	}

	if err := writeBuffer(w, response.Buffer); err != nil {
		return fmt.Errorf("writing buffer: %w", err)
	}

	return nil
}

func SendRequest(w io.Writer, request Request) error {
	if err := writeInt32(w, request.Type); err != nil {
		return fmt.Errorf("writing request type: %w", err)
	}

	if err := writeString(w, request.Str); err != nil {
		return fmt.Errorf("writing string: %w", err)
	}

	return nil
}

func ReceiveResponse(r io.Reader) (Response, error) {
	ResponseType, err := readInt32(r)
	if err != nil {
		return Response{}, fmt.Errorf("reading response type: %w", err)
	}

	if ResponseType != TypeResponse {
		return Response{}, errors.New("invalid response type")
	}

	errorNo, err := readInt32(r)
	if err != nil {
		return Response{}, fmt.Errorf("reading error number: %w", err)
	}

	buffer, err := readBuffer(r)
	if err != nil {
		return Response{}, fmt.Errorf("reading buffer: %w", err)
	}

	return Response{
		Type:    ResponseType,
		ErrorNo: errorNo,
		Buffer:  buffer,
	}, nil
}

// readInt32 читает int32 из r, как 8 байт, кодируется, как Big Endian.
func readInt32(r io.Reader) (int32, error) {
	var buf [8]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}

	value := binary.BigEndian.Uint64(buf[:])
	if value > math.MaxInt32 {
		return 0, errors.New("int32 value too large")
	}

	return int32(value), nil
}

// writeInt32 записывает int32 в w, как 8 байт, кодируется, как Big Endian.
func writeInt32(w io.Writer, value int32) error {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(value))
	_, err := w.Write(buf[:])
	return err
}

// readString читает string из r, как: int32 - длина строки, bytes - данные
// строки, кодировка строки - UTF 8, выравнивается на границу 4 байта.
func readString(r io.Reader) (string, error) {
	length, err := readInt32(r)
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", errors.New("negative string length")
	}

	data := make([]byte, length)
	if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
		return "", err
	}

	if !utf8.Valid(data) {
		return "", errors.New("invalid UTF-8 string")
	}

	padding := (maxPaddingByte - (length % maxPaddingByte)) % maxPaddingByte
	if padding > 0 {
		if err := binary.Read(r, binary.LittleEndian, make([]byte, padding)); err != nil {
			return "", err
		}
	}

	return string(data), nil
}

// writeString записывает string в w, как: int32 - длина строки, bytes - данные
// строки, кодировка строки - UTF 8, выравнивается на границу 4 байта.
func writeString(w io.Writer, s string) error {
	data := []byte(s)
	length := int32(len(data))

	if !utf8.Valid(data) {
		return errors.New("invalid UTF-8 string")
	}

	if err := writeInt32(w, length); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, data); err != nil {
		return err
	}

	padding := (maxPaddingByte - (length % maxPaddingByte)) % maxPaddingByte
	if padding > 0 {
		if err := binary.Write(w, binary.LittleEndian, make([]byte, padding)); err != nil {
			return err
		}
	}

	return nil
}

// readBuffer читает []byte из r, как: buffer - кодируется, как: int32 - длина
// буфера, bytes - данные буфера выравнивается на границу 4 байта.
func readBuffer(r io.Reader) ([]byte, error) {
	length, err := readInt32(r)
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, errors.New("negative buffer length")
	}

	data := make([]byte, length)
	if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
		return nil, err
	}

	padding := (maxPaddingByte - (length % maxPaddingByte)) % maxPaddingByte
	if padding > 0 {
		if err := binary.Read(r, binary.LittleEndian, make([]byte, padding)); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// writeBuffer записывает []byte в w, как: buffer - кодируется, как: int32 - длина
// буфера, bytes - данные буфера выравнивается на границу 4 байта.
func writeBuffer(w io.Writer, data []byte) error {
	length := int32(len(data))
	if err := writeInt32(w, length); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, data); err != nil {
		return err
	}

	padding := (maxPaddingByte - (length % maxPaddingByte)) % maxPaddingByte
	if padding > 0 {
		if err := binary.Write(w, binary.LittleEndian, make([]byte, padding)); err != nil {
			return err
		}
	}

	return nil
}
