package aurapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInfo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(
			`{"type": "multiinfo", "results": [
				{"Name": "foo", "Version": "1-0.1", "Maintainer": "neer", "LastModified": 1700000000},
				{"Name": "orphan-pkg", "Version": "3.4.1", "Maintainer": null, "LastModified": 1600000000}
			]}`))
	}))
	defer srv.Close()
	c := &Client{BaseUrl: srv.URL, HTTP: srv.Client()}
	got, err := c.Info(context.Background(), []string("foo", "orphan-pkg", "gone"))

	if err != nil {
		t.Fatal(err)
	}
	if got["foo"].Maintainer == nil || *got["foo"].Maintainer != "neer" {
		t.Errorf("foo maintainer wrong: &+v", got["foo"])
	}
	if got["orphan-pkg"].Maintainer != nil {
		t.Errorf("orphan-pkg should have nil maintainer.")
	}
	if _, ok := got["gone"]; ok {
		t.Errorf("missing package shouldn't be in the map")
	}
}
