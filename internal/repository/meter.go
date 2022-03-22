package repository

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

type Session struct {
	startedAt time.Time
	meters    []*Meter
}

func NewSession() *Session {
	return &Session{
		startedAt: time.Now(),
		meters:    []*Meter{},
	}
}

func (s *Session) NewMeter() *Meter {

	meter := &Meter{
		total: 0,
	}
	s.meters = append(s.meters, meter)
	return meter
}

func (s *Session) CalculateSpeed() ByteSize {
	totalDiff := time.Now().Sub(s.startedAt).Seconds()
	var total int64
	for _, meter := range s.meters {
		total = total + meter.total
	}

	speed := float64(total) / totalDiff
	return ByteSize(speed)
}

type Meter struct {
	total int64
}

func (m *Meter) Write(p []byte) (n int, err error) {
	written := len(p)
	m.AddBytes(written)
	return written, nil
}

func (m *Meter) AddBytes(bytes int) {
	atomic.AddInt64(&m.total, int64(bytes))
}

func (m *Meter) SetBytes(bytes int64) {
	atomic.StoreInt64(&m.total, bytes)
}

func CopyWithMeter(dest io.Writer, src io.Reader, meter *Meter) (int64, error) {
	if meter != nil {
		return io.Copy(dest, io.TeeReader(src, meter))
	}

	return io.Copy(dest, src)
}
