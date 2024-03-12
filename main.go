package main

import (
	"io"
	"log"
	"net/http"
)

var defaultConfig = http.DefaultTransport

func main(){

  log.Println("server running on port 8000")
  http.ListenAndServe(":8000", http.HandlerFunc(handleRequest))
}

func handleRequest(w http.ResponseWriter, r *http.Request){
  targetURL := r.URL
  proxyRequest, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
  if err != nil{
    http.Error(w, "you request can't be resolve", http.StatusInternalServerError)
    return
  }

  for name, values := range r.Header{
    for _, value := range values{
      r.Header.Set(name, value)
    }
  }

  resp, err := defaultConfig.RoundTrip(proxyRequest)
  if err != nil{
    http.Error(w, "error while sending request", http.StatusInternalServerError)
    return
  }
  defer resp.Body.Close()

  for name, values := range r.Header{
    for _, value := range values{
      w.Header().Add(name, value)
    }
  }

  io.Copy(w, resp.Body)
}
