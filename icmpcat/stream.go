package icmpcat

import "io"

type stream struct {
	done   chan bool
	isDone bool
	more   chan []byte
	unread []byte
}

func newStream() *stream {
	return &stream{
		isDone: false,
		done:   make(chan bool, 1),
		more:   make(chan []byte, 1024),
	}
}

func (s *stream) Read(bytes []byte) (int, error) {
	if s.isDone {
		return 0, io.EOF
	}
	if s.unread != nil && len(s.unread) > 0 {
		n := 0
		for i := 0; i < len(bytes); i++ {
			if i > len(s.unread)-1 {
				break
			}
			bytes[i] = s.unread[i]
			n++
		}
		s.unread = s.unread[n:]
		return n, nil
	}
	select {
	case <-s.done:
		return 0, io.EOF
	case b := <-s.more:
		n := 0
		for i := 0; i < len(bytes); i++ {
			if i > len(b)-1 {
				break
			}
			bytes[i] = b[i]
			n++
		}
		s.unread = b[n:]
		return n, nil
	}
}

func (s *stream) Write(b []byte) error {
	s.more <- b
	return nil
}

func (s *stream) Close() error {
	s.isDone = true
	s.done <- true
	return nil
}
