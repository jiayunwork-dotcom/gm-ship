package stability

// gmSink buffers an intact-stability Result and flushes on Close.
// A second Close is treated as a rewind of the destination so a
// caller that defers Close after an explicit Close empties the
// committed GM / GZ fields.
type gmSink struct {
	dst    *Result
	buf    Result
	nclose int
}

func (s *gmSink) commit(r Result) {
	s.buf = r
}

func (s *gmSink) Close() error {
	s.nclose++
	if s.nclose == 1 {
		*s.dst = s.buf
		return nil
	}
	*s.dst = Result{}
	return nil
}
