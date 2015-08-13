package main

import "time"
import "net/http"
import "encoding/json"
import "os/exec"
import "bytes"
import "fmt"
import "flag"

import "github.com/mattn/go-shellwords"
import "github.com/hugopeixoto/httpull"

func run(cmd string) {
	args, _ := shellwords.Parse(cmd)

	x := exec.Command(args[0], args[1:]...)

	x.Start()
	x.Wait()
}

func GetJob(server_url string) (httpull.JobRequest, error) {
	job := httpull.JobRequest{}

	resp, err := http.Get(server_url + "/jobs/ask")
	if err != nil {
		return job, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&job)
	if err != nil {
		return job, err
	}

	return job, nil
}

func FinishJob(server_url string, job httpull.JobResponse) error {
	b, _ := json.Marshal(job)
	http.Post(server_url+"/jobs/finish", "application/json", bytes.NewBuffer(b))

	return nil
}

func ExecuteJob(worker_url string, job httpull.JobRequest) (httpull.JobResponse, error) {
	req, err := http.NewRequest(job.Method, worker_url+job.RequestURI, bytes.NewBuffer(job.Body))

	if err != nil {
		return httpull.JobResponse{}, err
	}

	req.Header = job.Headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return httpull.JobResponse{}, err
	}

	return httpull.SerializeResponse(job, resp), nil
}

func handle(server_url string, worker_url string) error {
	fmt.Println("getting job at " + server_url)
	job, err := GetJob(server_url)
	if err != nil {
		return err
	}

	fmt.Println("got job")

	response, err := ExecuteJob(worker_url, job)
	if err != nil {
		return err
	}

	fmt.Println("finished job")
	return FinishJob(server_url, response)
}

func worker(server_url string, worker_url string) {
	fmt.Println("launching worker")
	for {
		handle(server_url, worker_url)
	}
}

func main() {
	var (
		workers    = flag.Int("workers", 16, "number of requests the worker can handle simultaneously")
		server_url = flag.String("server-url", "http://localhost:8080", "address of the httpull-producer")
		worker_url = flag.String("worker-url", "http://localhost:3000", "address of the worker")
		command    = flag.String("worker-cmd", "", "worker command")
	)

	flag.Parse()

	for i := 0; i < *workers; i++ {
		go func() {
			time.Sleep(3 * time.Second)
			worker(*server_url, *worker_url)
		}()
	}

	run(*command)
}
