package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/cli/go-gh/pkg/api"

	graphql "github.com/cli/shurcooL-graphql"
)

// Borrowed from https://github.com/cli/shurcooL-graphql/blob/trunk/graphql_test.go

// localRoundTripper is an http.RoundTripper that executes HTTP transactions
// by using handler directly, instead of going over an HTTP connection.
type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)
	return w.Result(), nil
}

func NewLocalHTTP(resp func(w http.ResponseWriter, req *http.Request)) *http.Client {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", resp)
	return &http.Client{Transport: localRoundTripper{handler: mux}}
}

func NewFakeGQLClient(host string, resp func(w http.ResponseWriter, req *http.Request)) api.GQLClient {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", resp)
	endpoint := fmt.Sprintf("%s/graphql", host)
	g := graphql.NewClient("/graphql", NewLocalHTTP(resp))
	httpClient := NewLocalHTTP(resp)
	return gqlClient{
		client:     g,
		host:       endpoint,
		httpClient: httpClient,
	}
}

// Implements api.GQLClient interface.
type gqlClient struct {
	client     *graphql.Client
	host       string
	httpClient *http.Client
}

// DoWithContext executes a single GraphQL query request and populates the response into the data argument.
func (c gqlClient) DoWithContext(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
	reqBody, err := json.Marshal(map[string]interface{}{"query": query, "variables": variables})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.host, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		return api.HandleHTTPError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	gr := gqlResponse{Data: response}
	err = json.Unmarshal(body, &gr)
	if err != nil {
		return err
	}

	if len(gr.Errors) > 0 {
		return api.GQLError{Errors: gr.Errors}
	}

	return nil
}

// Do wraps DoWithContext using context.Background.
func (c gqlClient) Do(query string, variables map[string]interface{}, response interface{}) error {
	return c.DoWithContext(context.Background(), query, variables, response)
}

// MutateWithContext executes a single GraphQL mutation request,
// with a mutation derived from m, populating the response into it.
// "m" should be a pointer to struct that corresponds to the GitHub GraphQL schema.
func (c gqlClient) MutateWithContext(ctx context.Context, name string, m interface{}, variables map[string]interface{}) error {
	return c.client.MutateNamed(ctx, name, m, variables)
}

// Mutate wraps MutateWithContext using context.Background.
func (c gqlClient) Mutate(name string, m interface{}, variables map[string]interface{}) error {
	return c.MutateWithContext(context.Background(), name, m, variables)
}

// QueryWithContext executes a single GraphQL query request,
// with a query derived from q, populating the response into it.
// "q" should be a pointer to struct that corresponds to the GitHub GraphQL schema.
func (c gqlClient) QueryWithContext(ctx context.Context, name string, q interface{}, variables map[string]interface{}) error {
	return c.client.QueryNamed(ctx, name, q, variables)
}

// Query wraps QueryWithContext using context.Background.
func (c gqlClient) Query(name string, q interface{}, variables map[string]interface{}) error {
	return c.QueryWithContext(context.Background(), name, q, variables)
}

type gqlResponse struct {
	Data   interface{}
	Errors []api.GQLErrorItem
}
