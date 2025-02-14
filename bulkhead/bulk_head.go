package bulkhead

type Bulkhead struct {
	sem chan struct{}
}

func New(limit int) *Bulkhead {
	return &Bulkhead{sem: make(chan struct{}, limit)}
}

func (b *Bulkhead) Execute(task func()) {
	b.sem <- struct{}{} // Acquire slot
	go func() {
		defer func() { <-b.sem }() // Release slot
		task()
	}()
}