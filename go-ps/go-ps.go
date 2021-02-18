package go_ps

import (
	"errors"
	"github.com/mitchellh/go-ps"
)

func PS() string {
	ps, _ := ps.Processes()
	str := ": "
	for pp := range ps {
		str += ps[pp].Executable()+", "
	}
	return str
}

// FindProcess( key string ) ( int, string, error )
func FindProcess(key string) (int, string, error) {
	pname := ""
	pid := 0
	err := errors.New("not found")
	ps, _ := ps.Processes()

	for i, _ := range ps {
		if ps[i].Executable() == key {
			pid = ps[i].Pid()
			pname = ps[i].Executable()
			err = nil
			break
			break
		}
	}
	return pid, pname, err
}
