package internal

import (
	"bytes"
	"path/filepath"
)

const (
	Scheme = "etcd"
)

func BuildPath(scheme string, paths ...string) string {
	var nPath = filepath.Join(paths...)

	if len(nPath) > 0 && nPath[0] == '/' {
		nPath = nPath[1:]
	}

	var buf = bytes.NewBufferString(scheme)
	buf.WriteString(":///")
	buf.WriteString(nPath)

	if len(nPath) > 0 && nPath[len(nPath)-1] != '/' {
		buf.WriteString("/")
	}
	return buf.String()
}
