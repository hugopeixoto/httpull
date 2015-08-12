package httpull

import "fmt"

type Job struct {
	ResponseChan chan JobResponse
	Request      JobRequest
}

type Queue struct {
	in  chan Job
  in2 chan JobRequest
	out chan JobResponse

  jobs map[string]chan JobResponse
}

func NewQueue() *Queue {
  q := Queue{
    make(chan Job),
    make(chan JobRequest),
    make(chan JobResponse),
    make(map[string]chan JobResponse),
  }

  go q.Loop()

  return &q
}

func (q *Queue) Loop() {
  for {
    select {
    case job := <-q.in:
      q.jobs[job.Request.JobIdentifier] = job.ResponseChan
      go func() { q.in2 <- job.Request }()

    case response := <-q.out:
      ch := q.jobs[response.JobIdentifier]
      delete(q.jobs, response.JobIdentifier)
      go func() { ch <- response }()
    }

    fmt.Printf("{\"pending_jobs\": %v}\n", len(q.jobs))
  }
}

func (q *Queue) HandleJob(req JobRequest) JobResponse {
	ch := make(chan JobResponse)

  q.in <- Job{ch, req}

	return <-ch
}

func (q *Queue) AskJob() JobRequest {
	return <-q.in2
}

func (q *Queue) FinishJob(resp JobResponse) {
  q.out <- resp
}
