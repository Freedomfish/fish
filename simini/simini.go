package simini

import (
	"fmt"
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type StrMap map[string]string
type SimIni struct {
	sess_map_ map[string]StrMap
	loaded_   bool
	errmsg_   string
}

func (p *SimIni) IsLoaded() bool {
	return p.loaded_
}

func (p *SimIni) ErrMsg() string {
	return p.errmsg_
}

func (p *SimIni) LoadFile(filename string) int {
	fd, err := os.Open(filename)
	if nil != err {
		p.errmsg_ = err.Error()
		return 1
	}
	defer fd.Close()
	p.sess_map_ = make(map[string]StrMap)
	buf := bufio.NewReader(fd)
	curkey := ""
	flag := true
	num := 0
	for flag {
		num++
		line, err := buf.ReadString('\n')
		if io.EOF == err {
			flag = false
		}

		line = strings.TrimSpace(line)
		if len(line) < 2 || '#' == line[0] {
			continue
		}

		length := len(line)

		if '[' == line[0] && ']' == line[length-1] {
			curkey = line[1 : length-1]
			p.sess_map_[curkey] = make(StrMap)
			continue
		}
		if curkey == "" {
			p.errmsg_ = fmt.Sprintf("lack of [] on line %d", num)
			return 1
		}
		val := strings.SplitN(line, "=", 2)
		if 2 != len(val) || 0 == len(val[0]) {
			continue
		}
		val[0] = strings.TrimSpace(val[0])
		val[1] = strings.TrimSpace(val[1])
		fmt.Println(val[0], val[1])
		p.sess_map_[curkey][val[0]] = val[1]
	}
	p.loaded_ = true
	return 0
}

const (
	extern_head_label = "<begin>"
	extern_end_label  = "<end>"
)

func (p *SimIni) LoadFileExtern(filename string) int {
	fd, err := os.Open(filename)
	if nil != err {
		p.errmsg_ = err.Error()
		return 1
	}
	defer fd.Close()
	p.sess_map_ = make(map[string]StrMap)
	buf := bufio.NewReader(fd)
	curkey := ""  //[curkey]
	datakey := "" //datakey=dataval
	dataval := ""
	dataflag := false
	for {
		line, err := buf.ReadString('\n')
		if io.EOF == err {
			break
		}
		line = strings.TrimLeft(line, " ")
		if len(line) < 2 || '#' == line[0] || '\n' == line[0] || '\r' == line[0] {
			continue
		}

		length := len(line)
		if line[length-2] == '\r' {
			length -= 1
			line = line[:length]
		}
		if '[' == line[0] && ']' == line[length-2] {
			curkey = line[1 : length-2]
			p.sess_map_[curkey] = make(StrMap)
			datakey = ""
			dataval = ""
			dataflag = false
			continue
		}
		if curkey == "" {
			p.errmsg_ = "lack of []|line=" + line
			return 1
		}
		if datakey == "" {
			val := strings.SplitN(line, "=", 2)
			if 2 != len(val) || 0 == len(val[0]) {
				continue
			}
			datakey = strings.TrimRight(val[0], " ")
			tmpval := strings.TrimLeft(val[1], " ")
			if tmpval[0:len(tmpval)-1] != extern_head_label {
				p.sess_map_[curkey][datakey] = tmpval[0 : len(tmpval)-1]
				datakey = ""
			} else {
				dataflag = true
			}
			continue
		}
		if dataflag {
			if line[0:len(line)-1] == extern_end_label {
				p.sess_map_[curkey][datakey] = dataval
				datakey = ""
				dataval = ""
				dataflag = false
			} else {
				dataval += line
			}
			continue
		}
	}
	p.loaded_ = true
	return 0

}

func (p *SimIni) GetStringVal(sess, key string) string {
	sv, sok := p.sess_map_[sess]
	if sok {
		v, ok := sv[key]
		if ok {
			return v
		}
	}
	return ""
}

func (p *SimIni) GetStringValWithDefault(sess, key, default_v string) string {
	str := p.GetStringVal(sess, key)
	if str == "" {
		str = default_v
	}
	return str
}

func (p *SimIni) GetIntVal(sess, key string) (int, error) {
	str := p.GetStringVal(sess, key)
	if str == "" {
		return 0, nil
	}
	return strconv.Atoi(str)
}

func (p *SimIni) GetIntValWithDefault(sess, key string, default_v int) (int, error) {
	str := p.GetStringVal(sess, key)
	if str == "" {
		return default_v, nil
	}
	return strconv.Atoi(str)
}

func (p *SimIni) GetSession(sess string) StrMap {
	strmap := make(StrMap)
	sv, sok := p.sess_map_[sess]
	if sok {
		for k, v := range sv {
			strmap[k] = v
		}
	}
	return strmap
}

func (p *SimIni) GetAllSession() map[string]StrMap {
	return p.sess_map_
}
