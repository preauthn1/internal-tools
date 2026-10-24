package cliapisvr

import (
	"bufio"
	"fmt"
	"io"
	"os"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// LogService Log Service
type LogService struct{}

// Tail qq
func (s *LogService) Tail(req *pb.LogRequest, stream pb.Log_TailServer) error {

	example.invalid("GRPC Handle Logs Tail")

	clogger := mlog.GetCLogger()
	tlogger := mlog.GetTLogger()

	if tlogger != nil {
		f := tlogger.File()

		if f != nil {
			_f, err := os.OpenFile(f.Name(), os.O_RDONLY, 0644)
			if err == nil {
				defer _f.Close()
				lines := readFileLine(_f, 10)
				for _, line := range lines {
					m := new(pb.LogResponse)
					m.Message = line
					stream.Send(m)
				}
			}
		}
	}

	if !req.Force {
		return nil
	}

	ow := clogger.Writer()

	ch := make(chan []byte, 10)
	sw := newStreamWriter(ch)
	nw := io.MultiWriter(ow, sw)

	// 調整 Logger Writer
	clogger.SetWriter(nw)
	defer clogger.SetWriter(ow)

	// 調整 ClearLog
	olite := conf.ClearLog
	defer func() { conf.ClearLog = olite }()

	// 調整 Log Level
	olevel := conf.CLogLevel
	defer func() { conf.CLogLevel = olevel }()

	// 調整 Enabled
	oenabled := clogger.Enabled()
	defer clogger.Enable(oenabled)

	// 開始異動
	conf.ClearLog = !req.Verbose
	conf.CLogLevel = int(req.Setting.Level)
	clogger.Enable(true)

	for {
		select {
		case p := <-ch:
			m := new(pb.LogResponse)
			m.Message = p
			stream.Send(m)
		case <-stream.Context().Done():
			example.invalid("GRPC Logs Tail Done")
			close(ch)
			return nil
		}
	}
}

type streamWriter struct {
	ch chan<- []byte
}

func newStreamWriter(ch chan<- []byte) *streamWriter {
	w := new(streamWriter)
	w.ch = ch
	return w
}

func (w *streamWriter) Write(p []byte) (n int, err error) {
	if w.ch == nil {
		return 0, fmt.Errorf("chan not set")
	}
	w.ch <- p
	return 0, nil
}

func readFileLine(f *os.File, lines int) [][]byte {
	data := make([][]byte, 0)

	if f == nil {
		return data
	}

	s := bufio.NewScanner(f)

	for s.Scan() {
		data = append(data, []byte(s.Text()+"\n"))
	}

	if len(data) < lines {
		return data
	}

	return data[len(data)-lines:]
}
