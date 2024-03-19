package problem_test

import (
	"encoding/json"
	"errors"
	"goprojects/httptordle/internal/problem"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProblem(t *testing.T) {
	p := problem.New(problem.Title("title string"), problem.Custom("x", "value"))

	assert.JSONEq(t, p.JSONString(), `{"title":"title string", "x": "value"}`)

	b, err := json.Marshal(p)
	require.NoError(t, err)
	assert.JSONEq(t, string(b), `{"title":"title string", "x": "value"}`)

	p = problem.New(
		problem.Title("title string"),
		problem.Status(404),
		problem.Custom("x", "value"),
	)
	str := p.JSONString()
	assert.JSONEq(t, str, `{"title":"title string", "x": "value", "status":404}`)

	p.Append(problem.Detail("some more details"), problem.Instance("https://example.com/details"))
	str = p.JSONString()
	assert.JSONEq(
		t,
		str,
		`{"title":"title string", "x": "value", "status":404, "detail":"some more details", "instance":"https://example.com/details"}`,
	)

	p = problem.Of(http.StatusAccepted)
	str = p.JSONString()
	assert.JSONEq(
		t,
		str,
		`{"status":202, "title":"Accepted", "type":"https://tools.ietf.org/html/rfc7231#section-6.3.3"}`,
	)

}

func TestProblemHTTP(t *testing.T) {
	p := problem.New(
		problem.Title("title string"),
		problem.Status(404),
		problem.Custom("x", "value"),
	)
	p.Append(problem.Detail("some more details"), problem.Instance("https://example.com/details"))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.Append(problem.Type("https://example.com/404"))
		if r.Method == "HEAD" {
			p.WriteHeaderTo(w)
		} else {
			p.WriteTo(w)
		}
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	require.NoError(t, err)

	bodyBytes, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusNotFound)
	assert.Equal(t, res.Header.Get("Content-Type"), problem.ContentTypeJson)
	assert.Equal(
		t,
		string(bodyBytes),
		`{"detail":"some more details","instance":"https://example.com/details","status":404,"title":"title string","type":"https://example.com/404","x":"value"}`,
	)

	// Try HEAD request
	res, err = http.Head(ts.URL)
	require.NoError(t, err)

	bodyBytes, err = io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	assert.Zero(t, len(bodyBytes))
	assert.Equal(t, res.StatusCode, http.StatusNotFound)
	assert.Equal(t, res.Header.Get("Content-Type"), problem.ContentTypeJson)
}

func TestMarshalUnmarshal(t *testing.T) {
	p := problem.New(problem.Status(500), problem.Title("Strange"))

	newProblem := problem.New()
	err := json.Unmarshal(p.JSON(), &newProblem)
	require.NoError(t, err)
	assert.Equal(t, p.Error(), newProblem.Error())
}

func TestErrors(t *testing.T) {
	knownProblem := problem.New(problem.Status(404), problem.Title("not found"))

	responseFromExternalService := http.Response{
		StatusCode: 404,
		Header: map[string][]string{
			"Content-Type": {problem.ContentTypeJson},
		},
		Body: io.NopCloser(strings.NewReader(`{"status":404, "title":"not found"}`)),
	}
	defer responseFromExternalService.Body.Close()

	if responseFromExternalService.Header.Get("Content-Type") == problem.ContentTypeJson {
		problemDecoder := json.NewDecoder(responseFromExternalService.Body)

		problemFromExternalService := problem.New()
		problemDecoder.Decode(&problemFromExternalService)

		assert.ErrorIs(t, problemFromExternalService, knownProblem)
	}
}

func TestNestedErrors(t *testing.T) {
	rootProblem := problem.New(problem.Status(404), problem.Title("not found"))
	p := problem.New(problem.Wrap(rootProblem), problem.Title("high level error msg"))

	unwrappedProblem := errors.Unwrap(p)
	assert.ErrorIs(t, unwrappedProblem, rootProblem)
	assert.Empty(t, errors.Unwrap(unwrappedProblem))

	// See wrapped error in 'reason'
	assert.JSONEq(
		t,
		p.JSONString(),
		`{"reason":"{\"status\":404,\"title\":\"not found\"}", "title":"high level error msg"}`,
	)

	p = problem.New(problem.WrapSilent(rootProblem), problem.Title("high level error msg"))
	assert.JSONEq(t, p.JSONString(), `{"title":"high level error msg"}`)
}

func TestOSErrorInProblem(t *testing.T) {
	_, err := os.ReadFile("non-existing")
	if err != nil {
		p := problem.New(problem.Wrap(err), problem.Title("Internal Error"), problem.Status(404))
		assert.ErrorIs(t, p, os.ErrNotExist)
		assert.NotErrorIs(t, p, os.ErrPermission)

		var o *os.PathError
		assert.ErrorAs(t, p, &o)

		newErr := errors.New("new error")
		p = problem.New(problem.Wrap(newErr), problem.Title("new problem"))

		assert.ErrorIs(t, p, newErr)
	}
}

func TestTitlef(t *testing.T) {
	expected := `{"title":"this is a test"}`
	toTest := problem.New(problem.Titlef("this is a %s", "test")).JSONString()
	assert.Contains(t, expected, toTest)
}

func TestDetailf(t *testing.T) {
	expected := `{"detail":"this is a test"}`
	toTest := problem.New(problem.Detailf("this is a %s", "test")).JSONString()
	assert.Contains(t, expected, toTest)
}

func TestInstancef(t *testing.T) {
	expected := `{"instance":"this is a test"}`
	toTest := problem.New(problem.Instancef("this is a %s", "test")).JSONString()
	assert.Contains(t, expected, toTest)
}
