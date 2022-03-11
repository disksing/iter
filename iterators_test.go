package iter_test

// func TestChanIterator(t *testing.T) {
// 	assert := assert.New(t)
// 	ch := make(chan int)
// 	go func() {
// 		CopyN(IotaReader(1), 10, ChanWriter(ch))
// 		close(ch)
// 	}()
// 	assert.Equal(Accumulate(ChanReader(ch), ChanEOF, 0), 55)

// 	ch = make(chan int)
// 	close(ch)
// 	assert.True(ChanReader(ch).Eq(ChanEOF))

// 	ch = make(chan int)
// 	close(ch)
// 	assert.True(ChanEOF.Eq(ChanReader(ch)))

// 	assert.Nil(ChanEOF.Read())
// 	assert.Equal(ChanEOF.Next(), ChanEOF)
// 	assert.True(ChanEOF.Eq(ChanEOF))

// 	ch = make(chan int, 1)
// 	ch <- 100
// 	assert.Equal(ChanReader(ch).Read(), 100)

// 	ch = make(chan int, 2)
// 	ch <- 42
// 	ch <- 43
// 	it := ChanReader(ch)
// 	it.Next()
// 	assert.Equal(it.Read(), 43)

// 	ch = make(chan int, 10)
// 	CopyN(IotaReader(1), 5, ChanWriter(ch))
// 	close(ch)
// 	it = ChanReader(ch)
// 	it = AdvanceN(it, 3).(InputIter)
// 	assert.Equal(it.Read(), 4)
// 	assert.Equal(Distance(it, ChanEOF), 2)

// 	ch = make(chan int)
// 	w := ChanWriter(ch)
// 	assert.Panics(func() { AdvanceN(w, 1) })
// 	assert.Panics(func() { Distance(w, ChanEOF) })
// }

// type dummyWriter struct {
// 	tick int
// }

// func (w *dummyWriter) Write(b []byte) (int, error) {
// 	w.tick--
// 	if w.tick < 0 {
// 		return 0, errors.New("boom!")
// 	}
// 	return len(b), nil
// }

// func TestIOWriterPanics(t *testing.T) {
// 	assert := assert.New(t)
// 	assert.Panics(func() {
// 		CopyN(IotaReader(0), 10, IOWriter(&dummyWriter{}, ","))
// 	})
// 	assert.Panics(func() {
// 		CopyN(IotaReader(0), 10, IOWriter(&dummyWriter{tick: 1}, ","))
// 	})
// }
