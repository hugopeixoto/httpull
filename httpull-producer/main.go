package main

import "net/http"
import "flag"
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
	var (
		frontend = flag.String("frontend-address", ":8080", "frontend address")
		backend  = flag.String("backend-address", ":9090", "backend address")
	)

	flag.Parse()

	queue := QueueHandler{httpull.NewQueue()}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", queue.HandleJob)
		http.ListenAndServe(*frontend, mux)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/jobs/ask", queue.AskJob)
	mux.HandleFunc("/jobs/finish", queue.FinishJob)
	http.ListenAndServe(*backend, mux)
}
