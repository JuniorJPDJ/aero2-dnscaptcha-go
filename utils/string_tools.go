package utils

import (
	"strings"
	"io"
	"bytes"
	"encoding/base32"
)

func StringBetween(s string, start string, end string) string{
	str := s[strings.Index(s, start) + len(start):]
	return str[:strings.Index(str, end)]
}

func Base32(in io.Reader, remove_padding bool) (string, error){
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(in)
	if err != nil{
		return "", err
	}
	b32 := bytes.Buffer{}
	b32w := base32.NewEncoder(base32.StdEncoding, &b32)
	_, err = buf.WriteTo(b32w)
	if err != nil{
		return "", err
	}
	err = b32w.Close()
	if err != nil{
		return "", err
	}

	if remove_padding{
		b32.Truncate(b32.Len() - strings.Count(b32.String(), "="))
	}
	return b32.String(), nil
}

func UnBase32(in string) (*bytes.Buffer, error){
	inbuf := bytes.NewBufferString(in)
	inbuf.WriteString(strings.Repeat("=", (8 - inbuf.Len() % 8) % 8))  // restore padding if needed

	unb32 := base32.NewDecoder(base32.StdEncoding, inbuf)
	outbuf := bytes.Buffer{}
	_, err := outbuf.ReadFrom(unb32)
	if err != nil{
		return nil, err
	}
	return &outbuf, nil
}
