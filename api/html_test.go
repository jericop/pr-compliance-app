package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jericop/pr-compliance-app/fakes"
)

func TestGetHtml(t *testing.T) {
	t.Skip()
	api := getApiServer(&fakes.Querier{})
	api.AddHtmlRoutes()

	urlPath := getRouteUrlPath(t, api.router, "GetHtml")

	server := httptest.NewServer(api.router)
	defer server.Close()

	t.Run(fmt.Sprintf("StatusOK"), func(t *testing.T) {

		resp := makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
			// makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
			return http.Get(server.URL + urlPath)
		})

		fmt.Println("---html returned is:")
		io.Copy(os.Stdout, resp.Body)

		// decoder := json.NewDecoder(resp.Body)
		// defer resp.Body.Close()

		// var result postgres.Approval
		// if err := decoder.Decode(&result); err != nil {
		// 	t.Fatalf("expected 'err' (%v) be nil", err)
		// }

		// if !reflect.DeepEqual(result, expected) {
		// 	t.Fatalf("expected 'result' (%v) to equal 'expected' (%v)", result, expected)
		// }
	})

}
