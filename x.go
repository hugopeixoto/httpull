package httpull

import "net/http"
import "io/ioutil"
import "crypto/rand"
import "encoding/base64"

type JobRequest struct {
	JobIdentifier string
	RemoteAddress string
	Method        string
	RequestURI    string
	Headers       map[string][]string
	Body          []byte
}

type JobResponse struct {
	JobIdentifier string
	StatusCode    int
	Headers       map[string][]string
	Body          []byte
}

func generate_job_id() string {
  buf := make([]byte, 16)
  rand.Read(buf)
  return base64.StdEncoding.EncodeToString(buf)
}

func SerializeRequest(r *http.Request) JobRequest {
	body, _ := ioutil.ReadAll(r.Body)

	return JobRequest{
		generate_job_id(),
		r.RemoteAddr,
		r.Method,
		r.RequestURI,
		r.Header,
		body,
	}
}

func DeserializeResponse(response JobResponse, w http.ResponseWriter) {
	w.WriteHeader(response.StatusCode)

	for key, values := range response.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Write(response.Body)
}
