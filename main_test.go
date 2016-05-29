package main

import (
	"crypto/rand"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const (
	testSz       = 1000
	testFilename = "testout.mkv"
)

var (
	testData = make([]byte, testSz)
)

func TestNewVideoStream(t *testing.T) {
	os.Remove(testFilename)
	_, err := rand.Read(testData)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Length", strconv.Itoa(testSz))
		_, err := w.Write(testData)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	vs, err := NewVideoStream(ts.URL, time.Second, "testout.mkv", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if vs.Size != testSz {
		t.Fatalf("VideoStream created with wrong size, got %v wanted %v\n", vs.Size, testSz)
	}
	if vs.Duration != time.Second {
		t.Fatal("VideoStream did not set duration")
	}
	received := make([]byte, testSz)
	if _, err := io.ReadFull(vs.fs, received); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(received, testData) {
		t.Fatal("VideoStream did not set the expected underlying reader interface")
	}
	if _, err := os.Stat(testFilename); os.IsNotExist(err) {
		t.Fatal("VideoStream did not create outfile")
	}
	if err := os.Remove(testFilename); err != nil {
		t.Fatal(err)
	}
}
