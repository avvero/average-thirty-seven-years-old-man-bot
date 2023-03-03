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
  "id": "chatcmpl-6pqeGmrLomEyv9nCElEakuX9jXaja",
  "object": "chat.completion",
  "created": 1677815128,
  "model": "gpt-3.5-turbo-0301",
  "usage": {
    "prompt_tokens": 209,
    "completion_tokens": 97,
    "total_tokens": 306
  },
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Ну тогда давай будем ты и дети вместе вставать раньше утром, а то гамаляться будет некому!"
      },
      "finish_reason": "stop",
      "index": 0
    }
  ]
}`)
	}))
	defer ts.Close()

	apiClient := NewApiClient(ts.URL, "key")
	err, text := apiClient.Completion("any")
	if err != nil {
		t.Error("Error: ", fmt.Sprintf("%s", err))
	}
	expected := "Ну тогда давай будем ты и дети вместе вставать раньше утром, а то гамаляться будет некому!"
	if text != expected {
		t.Error("Text expected: ", expected, " != ", text)
	}
}
