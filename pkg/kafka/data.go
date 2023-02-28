package kafka

import (
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"
)

// ProtoBuffer -
type ProtoBuffer struct {
	srv      *Client
	sensorID uint64
	max      int32
	lasttime int64
	mutex    sync.Mutex
	count    int32
	rows     []ProtoData
}

// ProtoBuffer insert insert numerical data to database
func (n *ProtoBuffer) insert() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	defer atomic.StoreInt32(&n.count, 0)

	// send to kafka
	if err := n.srv.SendResultToKafka(n); err != nil {
		return err
	}

	return nil
}

// flush -
func (n *ProtoBuffer) flush() error {
	if n.count > n.max {
		return fmt.Errorf("numerical count more than max")
	}
	if atomic.LoadInt32(&n.count) > 0 {
		return n.insert()
	}
	return nil
}

// append -
func (n *ProtoBuffer) append(data ...ProtoData) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for i, v := range data {
		if i == 0 && v.Ts == n.lasttime {
			continue
		}

		// 缓存
		count := atomic.LoadInt32(&n.count)
		if count >= n.max {
			return nil
		}

		n.lasttime = v.Ts
		n.rows[count] = v
		atomic.AddInt32(&n.count, 1)

	}
	return nil
}

// sendToKafka -
func (s *Server) sendToKafka(id uint64, mode uint16, timestamp int64) error {
	if !s.status {
		return nil
	}
	key := uint64(mode) + (id << 2)
	buff, ok := s.numericalBuffer.Load(key)
	if !ok {
		buff, _ = s.numericalBuffer.LoadOrStore(key,
			&ProtoBuffer{
				sensorID: id,
				max:      8192,
				rows:     make([]ProtoData, 8192),
				srv:      s.srv,
			},
		)
	}

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	sensorID := fmt.Sprintf("%X", b[2:])

	return buff.(*ProtoBuffer).append(ProtoData{
		Ts:       timestamp / 1e3,
		SensorID: sensorID,
	})
}
