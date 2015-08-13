package main

import "net/http"
import "encoding/json"
import "github.com/hugopeixoto/httpull"

type QueueHandler struct {
	Queue *httpull.Queue
}

func (h *QueueHandler) HandleJob(w http.ResponseWriter, r *http.Request) {
	req := httpull.SerializeRequest(r)

	response := h.Queue.HandleJob(req)
	httpull.DeserializeResponse(response, w)
}

func (h *QueueHandler) AskJob(w http.ResponseWriter, req *http.Request) {
	job := h.Queue.AskJob()

  json.NewEncoder(w).Encode(job)
}

func (h *QueueHandler) FinishJob(w http.ResponseWriter, req *http.Request) {
  job := httpull.JobResponse{}
  json.NewDecoder(req.Body).Decode(&job)

	h.Queue.FinishJob(job)
}

func main() {
	queue := QueueHandler{httpull.NewQueue()}

	go func() {
    mux := http.NewServeMux()
		mux.HandleFunc("/", queue.HandleJob)
		http.ListenAndServe(":8080", mux)
	}()

  mux := http.NewServeMux()
	mux.HandleFunc("/jobs/ask", queue.AskJob)
	mux.HandleFunc("/jobs/finish", queue.FinishJob)
  http.ListenAndServe(":9090", mux)
}
