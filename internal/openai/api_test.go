package openai

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, `{
  "error": {
    "message": "That model does not exist",
    "type": "invalid_request_error",
    "param": null,
    "code": null
  }
}`)
	}))
	defer ts.Close()

	apiClient := NewApiClient(ts.URL, "key")
	err, text := apiClient.Completion("any")
	if err == nil {
		t.Error("Error: ", fmt.Sprintf("%s", err))
	}
	if text != "" {
		t.Error("Text expected: ", "", " != ", text)
	}
	expected := "Response from: " + ts.URL + ": 404: That model does not exist"
	if err.Error() != expected {
		t.Error("Error expected: ", expected, " != ", err.Error())
	}
}

func Test_message(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, `{
  "id": "cmpl-6QE2gKc8rObBwj6M3f1tQ5JLmpk9f",
  "object": "text_completion",
  "created": 1671708526,
  "model": "text-davinci-003",
  "choices": [
    {
      "text": "\"Ну тогда давай будем ты и дети вместе вставать раньше утром, а то гамаляться будет некому!\"",
      "index": 0,
      "logprobs": null,
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 179,
    "completion_tokens": 96,
    "total_tokens": 275
  }
}`)
	}))
	defer ts.Close()

	apiClient := NewApiClient(ts.URL, "key")
	err, text := apiClient.Completion("any")
	if err != nil {
		t.Error("Error: ", fmt.Sprintf("%s", err))
	}
	expected := "\"Ну тогда давай будем ты и дети вместе вставать раньше утром, а то гамаляться будет некому!\""
	if text != expected {
		t.Error("Text expected: ", expected, " != ", text)
	}
}
