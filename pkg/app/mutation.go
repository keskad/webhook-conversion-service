package app

import (
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type StreamMutator struct {
	ParentStream io.ReadCloser
	Replacements []Replacement
}

func (m StreamMutator) Read(p []byte) (int, error) {
	tmpBuff := make([]byte, len(p))
	n, err := m.ParentStream.Read(tmpBuff)

	// strip out last null-bytes
	buff := make([]byte, n)
	copy(buff, tmpBuff)

	// // todo: create a more advanced replacing that would consider replacing line-by-line
	buff, n = m.replace(buff, n)

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

func (m StreamMutator) replace(p []byte, readUpTo int) ([]byte, int) {
	asStr := string(p)[0:readUpTo]

	for _, replacement := range m.Replacements {
		asStr = strings.Replace(asStr, replacement.From, replacement.To, -1)
	}

	asBytes := make([]byte, len(asStr))
	copy(asBytes, asStr)

	return asBytes, len(asBytes)
}
