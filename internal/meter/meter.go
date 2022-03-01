package meter

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

type Meter struct {
	startedAt time.Time
	total     uint64
}

func (m *Meter) Write(p []byte) (n int, err error) {
	written := len(p)
	m.AddBytes(written)
	return written, nil
}

func NewMeter() *Meter {
	return &Meter{
		startedAt: time.Now(),
		total:     0,
	}
}

func (m *Meter) AddBytes(bytes int) {
	atomic.AddUint64(&m.total, uint64(bytes))
}

func (m *Meter) CalculateSpeed() ByteSize {
	totalDiff := time.Now().Sub(m.startedAt).Seconds()
	speed := float64(m.total) / totalDiff
	return ByteSize(speed)
}

func CopyWithMeter(dest io.Writer, src io.Reader, meter *Meter) (int64, error) {
	if meter != nil {
		return io.Copy(dest, io.TeeReader(src, meter))
	}

	return io.Copy(dest, src)
}
