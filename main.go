package grmon

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var (
	newline   = byte(10)
	statusRe  = regexp.MustCompile("^goroutine\\s(\\d+)\\s\\[(.*)\\]:")
	createdRe = regexp.MustCompile("^created by (.*)")
	threadRe  = regexp.MustCompile("^threadcreate\\sprofile:\\stotal\\s(\\d+)")
)

type Routine struct {
	Num       int      `json:"no"`
	State     string   `json:"state"`
	CreatedBy string   `json:"created_by"`
	Trace     []string `json:"trace"`
}

func ReadRoutines(buf bytes.Buffer) (routines []*Routine) {
	var p *Routine

	for {
		line, err := buf.ReadString(newline)
		if err != nil {
			break
		}

		mg := statusRe.FindStringSubmatch(line)
		if len(mg) > 2 {
			// new routine block
			p = &Routine{}

			i, err := strconv.Atoi(mg[1])
			if err != nil {
				panic(err)
			}
			p.Num = i

			p.State = mg[2]
			routines = append(routines, p)
			continue
		}

		mg = createdRe.FindStringSubmatch(line)
		if len(mg) > 1 {
			p.CreatedBy = mg[1]
		}

		line = strings.Trim(line, "\n")
		if line != "" {
			p.Trace = append(p.Trace, line)
		}
	}

	return
}

type ThreadCreate struct {
	Count int      `json:"count"`
	Trace []string `json:"trace"`
}

func ReadThreads(buf bytes.Buffer) *ThreadCreate {
	t := &ThreadCreate{}

	for {
		line, err := buf.ReadString(newline)
		if err != nil {
			break
		}

		mg := threadRe.FindStringSubmatch(line)
		if len(mg) > 1 {
			i, err := strconv.Atoi(mg[1])
			if err != nil {
				panic(err)
			}
			t.Count = i
			continue
		}

		line = strings.Trim(line, "\n")
		if line != "" {
			t.Trace = append(t.Trace, line)
		}
	}

	return t
}
