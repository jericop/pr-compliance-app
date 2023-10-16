package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jericop/pr-compliance-app/storage/postgres"
)

var validPullRequests []postgres.PullRequest

func TestGetPullRequests(t *testing.T) {
	expected := []postgres.PullRequest{
		{ID: 1, PrID: 991, RepoID: 444, PrNumber: 1, OpenedBy: 651, IsMerged: false},
	}

	// Test the happy path first
	fakeStore.GetPullRequestsCall.Returns.PullRequestSlice = expected

	urlPath, err := router.Get("GetPullRequests").URL()
	if err != nil {
		t.Fatalf("expected 'err' (%v) be nil", err)
	}

	resp, err := http.Get(server.URL + urlPath.String())
	if err != nil {
		t.Fatalf("expected 'err' (%v) be nil", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'http.StatusOK' (%v)", resp.StatusCode, http.StatusOK)
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	result := []postgres.PullRequest{}

	err = decoder.Decode(&result)
	if err != nil {
		t.Fatalf("expected 'err' (%v) be nil", err)
	}

	indentedJson, err := PrettyStruct(result)
	if err != nil {
		t.Fatalf("expected 'err' (%v) to be nil", err)
	}
	fmt.Println("PrettyStruct(pullRequests)\n", indentedJson)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected 'result' (%v) to equal 'expected' (%v)", result, expected)
	}

	if fakeStore.GetPullRequestsCall.CallCount != 1 {
		t.Error("unexpected call count")
	}

	// Set apiServer.jsonMarshal to an anonymous function with the same signature as json.Marshal in order to force an error.
	apiServer.jsonMarshal = func(v interface{}) ([]byte, error) {
		return []byte{}, fmt.Errorf("Marshalling failed")
	}
	resp, err = http.Get(server.URL + urlPath.String())
	if err != nil {
		t.Fatalf("expected 'err' (%v) to be nil", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'http.StatusInternalServerError' (%v)", resp.StatusCode, http.StatusInternalServerError)
	}

	if fakeStore.GetPullRequestsCall.CallCount != 2 {
		t.Error("unexpected call count")
	}

	// Set the
	fakeStore.GetPullRequestsCall.Returns.PullRequestSlice = []postgres.PullRequest{}
	fakeStore.GetPullRequestsCall.Returns.Error = fmt.Errorf("db error")

	resp, err = http.Get(server.URL + urlPath.String())
	if err != nil {
		t.Fatalf("expected 'err' (%v) to be nil", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'http.StatusInternalServerError' (%v)", resp.StatusCode, http.StatusInternalServerError)
	}

	if fakeStore.GetPullRequestsCall.CallCount != 3 {
		t.Error("unexpected call count")
	}
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

// func TestGetFoo(t *testing.T) {
// 	key := "ABeeSee"
// 	urlPath, err := router.Get("GetFoo").URL("key", key)
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) be nil", err)
// 	}

// 	resp, err := http.Get(server.URL + urlPath.String())
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) be nil", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'http.StatusOK' (%v)", resp.StatusCode, http.StatusOK)
// 	}

// 	decoder := json.NewDecoder(resp.Body)
// 	defer resp.Body.Close()

// 	var foo Foo
// 	err = decoder.Decode(&foo)
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) be nil", err)
// 	}

// 	if foo.B != AllPullRequests[key].B {
// 		t.Fatalf("expected 'foo.B' (%v) to equal 'AllPullRequests[key].B' (%v)", foo.B, AllPullRequests[key].B)
// 	}
// 	if foo.A != AllPullRequests[key].A {
// 		t.Fatalf("expected 'foo.A' (%v) to equal 'AllPullRequests[key].A' (%v)", foo.A, AllPullRequests[key].A)
// 	}
// 	if foo.R != AllPullRequests[key].R {
// 		t.Fatalf("expected 'foo.R' (%v) to equal 'AllPullRequests[key].R' (%v)", foo.R, AllPullRequests[key].R)
// 	}
// }
