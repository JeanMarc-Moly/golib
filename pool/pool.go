package pool

import "sync"

type Pool struct {
	count, queue int
	jobs         chan func()
}

func New(workers, queue int) *Pool {
	return &Pool{count: workers, queue: queue}
}

func (p *Pool) Run() {
	count, jobs, wg := p.count, make(chan func()), new(sync.WaitGroup)

	p.jobs = jobs
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			for w := range jobs {
				w()
			}
		}()
	}

	wg.Wait()
}

func (p *Pool) Hand(job func()) {
	p.jobs <- job
}

func (p *Pool) Feed(jobs <-chan func()) {
	for job := range jobs {
		p.Hand(job)
	}
}

func (p *Pool) Close() {
	if p.jobs != nil {
		close(p.jobs)
	}
}
