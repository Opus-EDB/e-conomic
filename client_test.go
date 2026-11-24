package economic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// fakeRestAPI routes http.DefaultClient to handler for the duration of the test,
// so nothing goes over the network.
func fakeRestAPI(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	old := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		rec := httptest.NewRecorder()
		handler(rec, r)
		return rec.Result(), nil
	})
	t.Cleanup(func() { http.DefaultClient.Transport = old })
}

type pagedEntity struct {
	N    int    `json:"n"`
	Note string `json:"note,omitempty"`
}

// The filtered collection is n = 1..5, served two per page; the unfiltered one
// (what a request without the filter gets) is n = 101..120. Entity 1 carries a
// Note that no later entity has, to catch fields leaking between pages.
func TestGetEntitiesPagination(t *testing.T) {
	const filter = "date$gt$2026-11-01"
	const pageSize = 2
	filtered := []pagedEntity{{N: 1, Note: "first"}, {N: 2}, {N: 3}, {N: 4}, {N: 5}}
	var unfiltered []pagedEntity
	for n := 101; n <= 120; n++ {
		unfiltered = append(unfiltered, pagedEntity{N: n})
	}

	var requests []string
	fakeRestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RawQuery)
		q := r.URL.Query()
		all := unfiltered
		if q.Get("filter") == filter {
			all = filtered
		}
		skip, _ := strconv.Atoi(q.Get("skippages"))
		size, _ := strconv.Atoi(q.Get("pagesize"))
		from, to := min(skip*size, len(all)), min((skip+1)*size, len(all))
		body, _ := json.Marshal(map[string]any{
			"collection": all[from:to],
			"pagination": map[string]any{"results": len(all), "pageSize": size, "skipPages": skip},
		})
		io.WriteString(w, string(body))
	})

	client := &Client{AgreementGrant: "grant", AppSecretToken: "secret"}
	tc := &TypedClient[pagedEntity]{client: client}
	got, err := tc.getEntities("invoices/paid", pageSize, filter)
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprint(got) != fmt.Sprint(filtered) {
		t.Errorf("got %v, want %v", got, filtered)
	}
	if len(requests) != 3 {
		t.Errorf("made %d requests, want 3 (5 results at 2 per page): %v", len(requests), requests)
	}
	for _, raw := range requests {
		q, _ := url.ParseQuery(raw)
		if q.Get("filter") != filter {
			t.Errorf("request %q lost the filter", raw)
		}
	}
}

// An exact multiple of the page size must not cost an extra, empty request.
func TestGetEntitiesExactMultiple(t *testing.T) {
	var requests int
	fakeRestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		skip, _ := strconv.Atoi(r.URL.Query().Get("skippages"))
		body, _ := json.Marshal(map[string]any{
			"collection": []pagedEntity{{N: 2*skip + 1}, {N: 2*skip + 2}},
			"pagination": map[string]any{"results": 4},
		})
		io.WriteString(w, string(body))
	})

	client := &Client{AgreementGrant: "grant", AppSecretToken: "secret"}
	got, err := (&TypedClient[pagedEntity]{client: client}).getEntities("invoices/booked", 2, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 || requests != 2 {
		t.Errorf("got %v in %d requests, want 4 entities in 2", got, requests)
	}
}
