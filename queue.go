package httpull

import "fmt"
import "time"
import "net/http"

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
    make(chan Job, 100),
    make(chan JobRequest, 100),
    make(chan JobResponse, 100),
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
      go func() { fmt.Println("request ready to be served"); q.in2 <- job.Request }()
      go func(job JobRequest) {
        time.Sleep(30 * time.Second)
        q.FinishJob(JobResponse{
          job.JobIdentifier,
          http.StatusGatewayTimeout,
          map[string][]string{},
          []byte(``),
        })
      }(job.Request)

    case response := <-q.out:
      ch, ok := q.jobs[response.JobIdentifier]
      if ok {
        delete(q.jobs, response.JobIdentifier)
        go func() { ch <- response }()
      }
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
  fmt.Println("worker waiting for job")
	return <-q.in2
}

func (q *Queue) FinishJob(resp JobResponse) {
  q.out <- resp
}
