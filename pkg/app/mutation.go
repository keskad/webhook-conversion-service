package app

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
)

type StreamMutator struct {
	ParentStream io.ReadCloser
	Replacements []Replacement
}

func (m StreamMutator) Read(p []byte) (int, error) {
	buff := make([]byte, len(p))
	n, err := m.ParentStream.Read(buff)

	// todo: create a more advanced replacing that would consider replacing line-by-line
	buff = m.replace(buff)

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.Debugln("read()")
		logrus.Debugln(string(buff))
	}

	copy(p, buff)

	return n, err
}

func (m StreamMutator) Close() error {
	return m.ParentStream.Close()
}

func (m StreamMutator) replace(p []byte) []byte {
	for _, replacement := range m.Replacements {
		p = bytes.Replace(p, []byte(replacement.From), []byte(replacement.To), -1)
	}
	return p
}
