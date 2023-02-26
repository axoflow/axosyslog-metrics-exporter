package syslogngctl

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	iox "github.com/axoflow/axo-edge/x/io"
)

func NewReadWriterControlChannel(rwCtor func() (io.ReadWriter, error)) *ReadWriterControlChannel {
	return &ReadWriterControlChannel{
		rwCtor: rwCtor,
	}
}

type ReadWriterControlChannel struct {
	rwCtor func() (io.ReadWriter, error)
}

func (r ReadWriterControlChannel) SendCommand(cmd string) (rsp string, err error) {
	rw, err := r.rwCtor()
	if err != nil {
		return
	}

	if closer, _ := rw.(io.Closer); closer != nil {
		defer closer.Close()
	}

	if _, err = io.WriteString(rw, cmd+"\n"); err != nil {
		return
	}

	// command is sent

	dat, rst, err := iox.ReadUntil(rw, []byte("\n"+responseTerminator))
	if len(rst) > 0 {
		dat, rst = append(dat, rst[0]), rst[1:] // re-add last new line removed by ReadUntil
		if err == io.EOF {
			err = nil // ignore EOF if terminator has been matched
		}
	}
	if !bytes.HasPrefix(rst, []byte(responseTerminator)) {
		err = errors.Join(err, MissingResponseTerminator{
			Response: dat,
		})
	}
	// TODO: check if there is something after the terminator
	rsp = string(bytes.ToValidUTF8(dat, []byte("�")))
	return
}

type MissingResponseTerminator struct {
	Response []byte
}

func (err MissingResponseTerminator) Error() string {
	return fmt.Sprintf("missing response terminator %q", responseTerminator)
}

const responseTerminator string = ".\n"
