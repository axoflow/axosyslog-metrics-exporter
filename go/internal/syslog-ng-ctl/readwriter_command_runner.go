package syslogngctl

import (
	"bytes"
	"fmt"
	"io"

	iox "github.com/axoflow/axo-edge/x/io"
	"go.uber.org/multierr"
)

func NewReadWriterCommandRunner(rwCtor func() (io.ReadWriter, error)) *ReadWriterCommandRunner {
	return &ReadWriterCommandRunner{
		rwCtor: rwCtor,
	}
}

type ReadWriterCommandRunner struct {
	rwCtor func() (io.ReadWriter, error)
}

func (r ReadWriterCommandRunner) RunCommand(cmd string) (rsp string, err error) {
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
		err = multierr.Append(err, MissingResponseTerminator{
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
