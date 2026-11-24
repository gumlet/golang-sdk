// File generated from our OpenAPI spec by Scalar. See README.md for details.

package requestconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gumlet/golang-sdk/internal/apierror"
)

// Uploads to pre-signed URLs for the SDK's signed-upload methods.
//
// A signed-upload method calls an operation that answers with a pre-signed URL, sends a file to
// that URL, then optionally confirms. The upload is deliberately NOT sent through the SDK's request
// pipeline: the client's requests carry its credentials (auth headers, query auth, cookies), its
// default headers and query, and its middleware. None of that may reach the upload host, which is
// usually a different host with the authorization already signed into the URL. So each upload
// builds its request from nothing and sends exactly the file, `Content-Type`, and the headers the
// initiating response itself names. Only the client's per-attempt timeout and retry budget are
// reused, with the network path below, and the timeout counts time without progress rather than the
// whole transfer (see [sendSignedUpload]).
//
// The one piece of the client's HTTP client an upload does take is its network setup (see
// [SignedUploadTransport]). A plain [*http.Transport] is cloned, so its proxy, custom CA, client
// certificate and dialers apply. Any other [http.RoundTripper] and a custom [HTTPDoer] may add
// headers or credentials of their own (a tracing or OAuth2 wrapper, say), so the upload never goes
// through one unless the caller marked it with option.SignedUploadTransport, declaring that it sends
// a request as it is handed over; otherwise the upload falls back to a transport of its own and the
// proxy settings in the environment (HTTPS_PROXY, NO_PROXY). On every path the client's cookie jar,
// redirect policy and timeout, the SDK's middleware and every header and credential option stay
// behind: those are added by the request pipeline and the [*http.Client], not the transport.

// signedUploadDefaultContentType is the `Content-Type` an upload falls back to when neither the
// initiating response, the SDK configuration, nor the file names one.
const signedUploadDefaultContentType = "application/octet-stream"

// signedUploadResponseHeaderTimeout bounds the wait for the upload host's response headers, the
// same bound the SDK's default HTTP client applies, so a host that accepts the connection and never
// answers cannot hang an upload that has no per-attempt timeout. On the upload's own clone of a
// transport it starts once the file is sent; through a marked RoundTripper or HTTPDoer of the
// caller's, whose wait the upload cannot set, it is each attempt's idle timeout instead (see
// [signedUploadAttemptTimeout]).
const signedUploadResponseHeaderTimeout = 10 * time.Minute

// signedUploadMaxRetryAfter is the longest wait a `Retry-After`/`retry-after-ms` header can ask for
// and be honoured; anything longer (or not positive) falls back to the client's own backoff.
const signedUploadMaxRetryAfter = time.Minute

// SignedUploadSpec describes how a generated signed-upload method sends its file. It is fixed when
// the SDK is generated.
type SignedUploadSpec struct {
	// URLField is the wire-name path to the upload URL in the initiating response.
	URLField []string
	// Method is the HTTP verb the file is sent with.
	Method string
	// Multipart sends `multipart/form-data` with the file as one part instead of the raw bytes.
	Multipart bool
	// MultipartField is the form field the file is sent under, for a multipart upload.
	MultipartField string
	// ContentType is the configured `Content-Type` of a binary upload, or empty.
	ContentType string
	// HeadersField is the wire-name path to a string map of headers the upload must carry, or nil.
	HeadersField []string
}

// SignedUploadFile is the file a signed upload sends, read into memory once before the initiating
// request so that every retried attempt sends the same bytes.
type SignedUploadFile struct {
	// Name is the file's own base name, or empty when the reader does not carry one.
	Name string
	// ContentType is the file's own media type, or empty when the reader does not carry one.
	ContentType string
	// Content is the whole file.
	Content []byte
}

// ReadSignedUploadFile reads the file a signed upload sends, before the initiating request is made.
//
// A file that cannot be read then fails the call before the API has created an upload that would
// never receive its file. The file's name comes from a `Filename() string` or `Name() string` method
// (an [*os.File] has the latter), and its media type from a `ContentType() string` method — the same
// conventions the SDK's multipart encoder reads.
func ReadSignedUploadFile(file io.Reader) (SignedUploadFile, error) {
	if file == nil {
		return SignedUploadFile{}, errors.New("signed upload: the file to upload is nil")
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return SignedUploadFile{}, fmt.Errorf("signed upload: reading the file to upload: %w", err)
	}
	upload := SignedUploadFile{Content: content}
	if named, ok := file.(interface{ Filename() string }); ok {
		upload.Name = named.Filename()
	} else if named, ok := file.(interface{ Name() string }); ok && named.Name() != "" {
		upload.Name = filepath.Base(named.Name())
	}
	if typed, ok := file.(interface{ ContentType() string }); ok {
		upload.ContentType = typed.ContentType()
	}
	return upload, nil
}

// ReadSignedUploadField reads a value out of a raw JSON response body by wire-name path, or nil when
// a step is missing or the body is not JSON.
//
// It reads the raw body rather than a decoded struct, because the path is spelled in wire names and
// because an API that issues upload URLs does not always label its JSON as JSON: a body sent as
// `text/plain` is still read.
func ReadSignedUploadField(body []byte, path []string) interface{} {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	var current interface{}
	if err := decoder.Decode(&current); err != nil {
		return nil
	}
	for _, segment := range path {
		object, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = object[segment]
		if !ok {
			return nil
		}
	}
	return current
}

// signedUploadIntegerPattern matches a JSON number that is a plain integer.
var signedUploadIntegerPattern = regexp.MustCompile(`^-?[0-9]+$`)

// SignedUploadPathParam reads a confirm-step path parameter off the initiating response body.
//
// It fails rather than filling the path with an empty or zero value: a confirm request sent to
// `/files//complete` would at best 404 and at worst act on a resource nobody named.
func SignedUploadPathParam(initiated []byte, path []string) (string, error) {
	switch value := ReadSignedUploadField(initiated, path).(type) {
	case string:
		// A `.` or `..` segment would send the confirm request to another resource. NewRequestConfig
		// refuses one too, but only once the confirm is built, after the file has been sent; refusing it
		// here keeps the upload from going out with a confirm that can never be made.
		if segment, found := firstDotSegment(value); found {
			return "", fmt.Errorf("signed upload: the upload response's `%s` contains a %q segment, so the upload cannot be confirmed", strings.Join(path, "."), segment)
		}
		if value != "" {
			return value, nil
		}
	case json.Number:
		if signedUploadIntegerPattern.MatchString(value.String()) {
			return value.String(), nil
		}
	}
	return "", fmt.Errorf("signed upload: expected the upload response to carry a non-empty `%s` to confirm the upload with", strings.Join(path, "."))
}

// ExecuteSignedUploadInitiation sends a signed upload's initiating request, built by
// [NewRequestConfig] with no destination, and returns its raw body: the upload URL, the headers the
// upload must carry and any confirm path parameter are read off it.
//
// The body is captured here whatever the caller's options asked for. A caller's
// option.WithResponseBodyInto replaces the destination of every request its options reach, and the
// upload cannot go on without these bytes, so it never takes them from this request. What happens to
// the caller's response options depends on isResult, whether this response is the method's result:
//
//   - With a confirm step (isResult false), option.WithResponseBodyInto and option.WithResponseInto
//     are for the confirm response, which the confirm request hands them itself, so this request
//     leaves them unset.
//   - Without one (isResult true), both are honoured once the response is in: WithResponseInto
//     receives the [*http.Response], and a WithResponseBodyInto destination receives the body as the
//     client would deliver it (see [deliverSignedUploadResponseBody]). The method still decodes its
//     own return value from the same bytes.
func ExecuteSignedUploadInitiation(cfg *RequestConfig, isResult bool) ([]byte, error) {
	bodyInto := cfg.ResponseBodyInto
	responseInto := cfg.ResponseInto
	var contents []byte
	var res *http.Response
	cfg.ResponseBodyInto = &contents
	cfg.ResponseInto = &res
	err := cfg.Execute()
	if !isResult {
		return contents, err
	}
	// Set even when the request failed, as [RequestConfig.Execute] does, so a caller can inspect it.
	if responseInto != nil {
		*responseInto = res
	}
	if err != nil {
		return nil, err
	}
	if bodyInto != nil {
		if err := deliverSignedUploadResponseBody(bodyInto, res, contents); err != nil {
			return nil, err
		}
	}
	return contents, nil
}

// deliverSignedUploadResponseBody hands an initiating response body to the destination a caller gave
// option.WithResponseBodyInto, as [RequestConfig.Execute] would have: the bytes for a []byte, the
// text for a string, the response itself (its body readable again) for an [*http.Response], and the
// decoded JSON for anything else.
func deliverSignedUploadResponseBody(dst interface{}, res *http.Response, contents []byte) error {
	switch dst := dst.(type) {
	case **http.Response:
		if res != nil {
			res.Body = io.NopCloser(bytes.NewReader(contents))
		}
		*dst = res
	case *[]byte:
		*dst = contents
	case **[]byte:
		tmp := contents
		*dst = &tmp
	case *string:
		*dst = string(contents)
	case **string:
		tmp := string(contents)
		*dst = &tmp
	default:
		if len(contents) == 0 {
			return nil
		}
		if err := json.Unmarshal(contents, dst); err != nil {
			return fmt.Errorf("error parsing response json: %w", err)
		}
	}
	return nil
}

// ExecuteSignedUploadConfirm sends a signed upload's confirm request, an ordinary bodyless API call
// through the client pipeline, decoding its response into dst.
//
// opts are the client's options, the caller's per-call ones and the confirm's own; clientOpts are the
// client's alone. The per-call options are the ones the caller passed to the upload method, so they
// were written for the initiating request too. Those that shape how a request is sent still apply —
// headers, the timeout, retries, the base URL, middleware and the response destinations — but the
// confirm sends no body and has a query of its own. So whatever a per-call option.WithJSONSet,
// option.WithRequestBody or option.WithQuery added is dropped: the body is cleared, and the query,
// `Content-Type` and idempotency key are the ones the client's options alone give the confirm
// request, except that a query parameter a client setting writes (an API key sent in the query) is
// taken from the call's options, since a per-call one is a credential for the whole call (see
// [signedUploadConfirmQuery]).
func ExecuteSignedUploadConfirm(ctx context.Context, method string, path string, dst interface{}, clientOpts []RequestOption, opts ...RequestOption) error {
	client, err := NewRequestConfig(ctx, method, path, nil, nil, clientOpts...)
	if err != nil {
		return err
	}
	cfg, err := NewRequestConfig(ctx, method, path, nil, dst, opts...)
	if err != nil {
		return err
	}
	cfg.Body = nil
	cfg.Request.Body = nil
	cfg.Request.GetBody = nil
	cfg.Request.ContentLength = 0
	cfg.Request.Header.Del("Content-Type")
	if contentType := client.Request.Header.Values("Content-Type"); len(contentType) > 0 {
		cfg.Request.Header["Content-Type"] = contentType
	}
	cfg.Request.URL.RawQuery = signedUploadConfirmQuery(client.Request.URL, cfg.Request.URL)
	// A per-call option.WithHeader for the idempotency header was written for the initiating request:
	// replayed on the confirm, an API that honours the key would answer the confirm with the initiating
	// response, or refuse the key as reused. So the confirm keeps the key the client's options give it,
	// which NewRequestConfig generates fresh unless the client pins one of its own.
	if signedUploadIdempotencyHeader != "" {
		cfg.Request.Header.Del(signedUploadIdempotencyHeader)
		for _, key := range client.Request.Header.Values(signedUploadIdempotencyHeader) {
			cfg.Request.Header.Add(signedUploadIdempotencyHeader, key)
		}
	}
	return cfg.Execute()
}

// signedUploadIdempotencyHeader is the SDK's idempotency header, or empty when it sends none; see
// [ExecuteSignedUploadConfirm].
const signedUploadIdempotencyHeader = "Idempotency-Key"

// signedUploadClientQueryKeys are the query parameters the SDK's own client settings write, such as an
// API key sent in the query. A per-call option for one of them (option.WithAPIKey on a single call)
// is a credential for the whole call, so the confirm request keeps it.
var signedUploadClientQueryKeys = []string{}

// signedUploadConfirmQuery is the confirm request's query: the one the client's options alone give it,
// with each client-setting parameter (see [signedUploadClientQueryKeys]) taken from the merged
// options instead, so a credential passed for this call reaches the confirm and anything else a
// per-call option.WithQuery added for the initiating request does not.
func signedUploadConfirmQuery(client *url.URL, merged *url.URL) string {
	if len(signedUploadClientQueryKeys) == 0 {
		return client.RawQuery
	}
	query := client.Query()
	mergedQuery := merged.Query()
	for _, key := range signedUploadClientQueryKeys {
		if values, ok := mergedQuery[key]; ok {
			query[key] = values
		}
	}
	return query.Encode()
}

// UploadToSignedURL sends file to the pre-signed URL in the initiating response body and returns
// once the upload host has answered with a success status.
//
// maxRetries and timeout are the client's retry budget and per-attempt timeout, and network is the
// network path [SignedUploadTransport] picked: the only settings the upload takes from the client.
// The timeout bounds how long an attempt may go without progress, not how long it takes: a large
// file on a slow link is not cut off while its bytes keep going out (see [sendSignedUpload]). It
// retries connection errors, timeouts, 408, 409, 429 and 5xx (or whatever `x-should-retry` says),
// waiting as `Retry-After`/`retry-after-ms` asks when that is under a minute and otherwise backing
// off as the client does. Retrying is safe because the file was read once (see
// [ReadSignedUploadFile]) and every attempt sends the same bytes. A redirect is not followed: it
// would re-send the file and the named headers to a host the API's response did not name. A custom
// [HTTPDoer] follows redirects however it was built to, so the upload fails, unretried, once it has
// followed one or failed on another URL (see [signedUploadDoer]). Any other failure status is
// returned as an [*apierror.Error].
func UploadToSignedURL(ctx context.Context, maxRetries int, timeout time.Duration, network SignedUploadNetwork, initiated []byte, spec SignedUploadSpec, file SignedUploadFile) error {
	target, err := signedUploadURL(initiated, spec.URLField)
	if err != nil {
		return err
	}
	headers, body, err := signedUploadRequestParts(initiated, spec, file)
	if err != nil {
		return err
	}
	method := spec.Method
	if method == "" {
		method = http.MethodPut
	}

	client, release := network.sender()
	defer release()
	timeout = signedUploadAttemptTimeout(timeout, network)

	for attempt := 0; ; attempt++ {
		res, err := sendSignedUpload(ctx, client, timeout, method, target, headers, body)
		if ctx.Err() != nil {
			if res != nil {
				_ = res.Body.Close()
			}
			return ctx.Err()
		}
		// Not retried: the HTTPDoer may already have sent the file wherever the redirect pointed.
		if isSignedUploadRedirectRefusal(err) {
			return err
		}
		if err != nil {
			if attempt < maxRetries {
				if waitErr := waitForSignedUploadRetry(ctx, signedUploadRetryDelay(nil, attempt)); waitErr != nil {
					return waitErr
				}
				continue
			}
			return redactSignedUploadError(err, target)
		}
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			_, _ = io.Copy(io.Discard, res.Body)
			_ = res.Body.Close()
			return nil
		}
		if res.StatusCode >= 300 && res.StatusCode < 400 {
			_ = res.Body.Close()
			return fmt.Errorf("signed upload: the upload host answered with a redirect (%d), which is not followed", res.StatusCode)
		}
		// Asked of the same predicate the client's own retries use, so the two agree on
		// `x-should-retry` and on which statuses are worth another try. Its request only answers
		// whether the body can be replayed, which every upload attempt's can.
		if attempt < maxRetries && shouldRetry(&http.Request{GetBody: replayableSignedUploadBody(body)}, res) {
			_ = res.Body.Close()
			if waitErr := waitForSignedUploadRetry(ctx, signedUploadRetryDelay(res, attempt)); waitErr != nil {
				return waitErr
			}
			continue
		}
		return signedUploadStatusError(res)
	}
}

// SignedUploadNetwork is the network path an upload goes out over, as [SignedUploadTransport] picks
// it. Exactly one field is set.
type SignedUploadNetwork struct {
	// transport is a plain [*http.Transport], cloned with the upload's own settings.
	transport *http.Transport
	// roundTripper is a caller's [http.RoundTripper] marked as a signed-upload transport, used as it is.
	roundTripper http.RoundTripper
	// doer is a caller's [HTTPDoer] marked as a signed-upload transport, used as it is.
	doer HTTPDoer
}

// signedUploadMarker is option.SignedUploadTransport, spelled here so the runtime can read it off a
// caller's RoundTripper or HTTPDoer without importing the option package, which imports this one.
type signedUploadMarker interface {
	SignedUploadTransport()
}

// SignedUploadTransport picks the network path of cfg's HTTP client that an upload may go out over,
// and nothing else of cfg.
//
// A plain [*http.Transport] (the client's, or [http.DefaultTransport] when the client sets none and
// that is still the standard library's) is network setup alone, and a clone of it adds no header to
// a request, so the upload clones it (see [newSignedUploadClient]). Any other [http.RoundTripper],
// and a custom [HTTPDoer], may add headers or credentials the upload host must never see (an OAuth2
// transport sets `Authorization` on every request, whatever the host), and the SDK cannot tell one
// that does from a proxy or TLS wrapper. So one is used only when the caller marked it with
// option.SignedUploadTransport, declaring that it sends a request as it is handed over: a caller's
// declaration, never a guess, and off by default. An unmarked one, a nil transport and a typed-nil
// one fall back to the upload's own transport (see [signedUploadFallback]).
//
// The marker is read off what the caller handed to option.WithHTTPClient: the HTTPDoer itself, or
// the [*http.Client]'s Transport. A wrapper around a marked one is marked only if it says so itself.
// Whichever path is picked, the upload request is still built from nothing (see
// [signedUploadRequestParts]), and the client's cookie jar, timeout and redirect policy stay behind.
//
// A cloned transport is used exactly as configured, for the upload host as well as the API host: its
// RootCAs must trust the upload host, a nil Proxy means no proxy (not the environment's), and a
// ServerName, dialer or InsecureSkipVerify set for the API applies to the upload too. Adjusting any
// of them here would second-guess the caller's own network setup. A marked RoundTripper or HTTPDoer
// keeps its own settings too, and the upload cannot set its header wait, so each attempt through one
// is bounded by the idle timeout even when the client sets no timeout (see
// [signedUploadAttemptTimeout]). One that reads the whole file before sending it (a retrying wrapper
// that buffers bodies) or has a large WriteBufferSize gives the idle timeout no sign of progress while
// the bytes go out (see [sendSignedUpload]), so that bound then covers its whole transfer.
//
// Go's embedding promotes the marker's method: a wrapper type that embeds a marked one, by value or by
// pointer, is marked too, and whatever it adds reaches the upload host. A wrapper that holds the
// marked one in a named field is not.
func SignedUploadTransport(cfg *RequestConfig) SignedUploadNetwork {
	if cfg.CustomHTTPDoer != nil {
		if _, marked := cfg.CustomHTTPDoer.(signedUploadMarker); marked {
			return SignedUploadNetwork{doer: cfg.CustomHTTPDoer}
		}
		return signedUploadFallback()
	}
	var roundTripper http.RoundTripper
	if cfg.HTTPClient != nil {
		roundTripper = cfg.HTTPClient.Transport
	}
	// No transport, or http.DefaultTransport itself (the SDK's default client sets it when the process
	// wrapped it), is the process's, not something the caller handed this client: it is cloned when it
	// is the standard library's and never lent otherwise, marked or not.
	if roundTripper == nil || isSignedUploadDefaultTransport(roundTripper) {
		return signedUploadFallback()
	}
	if transport, ok := roundTripper.(*http.Transport); ok && transport != nil {
		return SignedUploadNetwork{transport: transport}
	}
	if _, marked := roundTripper.(signedUploadMarker); marked {
		return SignedUploadNetwork{roundTripper: roundTripper}
	}
	return signedUploadFallback()
}

// isSignedUploadDefaultTransport reports whether roundTripper is [http.DefaultTransport] itself.
//
// Comparing two interfaces panics when they hold the same type and the value cannot be compared: a
// RoundTripper that is a func, or a struct wrapper holding one in an interface field, whose type
// alone still reads as comparable. So the types and both values are checked first, and a value of
// the default's own type that cannot be compared counts as the default: it may well be, and lending
// the process's transport is the mistake to avoid, not falling back for a caller's own of that type.
func isSignedUploadDefaultTransport(roundTripper http.RoundTripper) bool {
	if roundTripper == nil || http.DefaultTransport == nil || reflect.TypeOf(roundTripper) != reflect.TypeOf(http.DefaultTransport) {
		return false
	}
	if !reflect.ValueOf(roundTripper).Comparable() || !reflect.ValueOf(http.DefaultTransport).Comparable() {
		return true
	}
	return roundTripper == http.DefaultTransport
}

// signedUploadAttemptTimeout is the idle timeout an upload's attempts get: the client's per-attempt
// timeout, or, with none, [signedUploadResponseHeaderTimeout] on a marked RoundTripper or HTTPDoer.
//
// The upload's own clone already bounds the wait for response headers, but a caller's code cannot be
// told to, and the upload host is one the API's response named. Without a bound, a marked transport
// against a host that never answers would hang the upload for good.
func signedUploadAttemptTimeout(timeout time.Duration, network SignedUploadNetwork) time.Duration {
	if timeout == 0 && network.transport == nil {
		return signedUploadResponseHeaderTimeout
	}
	return timeout
}

// signedUploadFallback is the upload's own transport, for a client whose network setup it may not
// borrow: a clone of [http.DefaultTransport] when that is still the standard library's, and otherwise
// a new one with the standard library's defaults. A process that replaced http.DefaultTransport with
// a wrapping RoundTripper may have it add headers (a tracing or auth wrapper, say), so the upload never
// runs through one; leaving the transport unset would send it through exactly that.
func signedUploadFallback() SignedUploadNetwork {
	if transport, ok := http.DefaultTransport.(*http.Transport); ok && transport != nil {
		return SignedUploadNetwork{transport: transport}
	}
	return SignedUploadNetwork{transport: newSignedUploadTransport()}
}

// newSignedUploadTransport builds the transport an upload uses when [http.DefaultTransport] is not
// the standard library's own, with the same settings that one is declared with.
func newSignedUploadTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// sender returns the client an upload's attempts are sent with, and what to call once the upload is
// done: never the SDK's own client, and never the caller's [*http.Client] itself, whose cookie jar,
// timeout and redirect policy would apply. A marked HTTPDoer is the one exception, since it is the
// only way to send through it.
func (n SignedUploadNetwork) sender() (HTTPDoer, func()) {
	switch {
	case n.transport != nil:
		client := newSignedUploadClient(n.transport)
		return client, client.CloseIdleConnections
	case n.doer != nil:
		return signedUploadBoundedDoer{doer: signedUploadDoer{doer: n.doer}}, func() {}
	case n.roundTripper != nil:
		// The caller's RoundTripper keeps its connections pooled for the client's other requests, so
		// the upload leaves its idle connections alone rather than close them under the client.
		client := &http.Client{
			Transport:     n.roundTripper,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		}
		return signedUploadBoundedDoer{doer: client}, func() {}
	default:
		// A zero network, which SignedUploadTransport never returns, takes the safe fallback rather
		// than an http.Client with no Transport, which would send through http.DefaultTransport.
		return signedUploadFallback().sender()
	}
}

// errSignedUploadFollowedRedirect is what an upload through a caller's own [HTTPDoer] fails with once
// the doer has followed a redirect: refused, never retried.
var errSignedUploadFollowedRedirect = errors.New("signed upload: the HTTP client followed a redirect from the upload host, which an upload must not do")

// errSignedUploadOtherURL is what an upload through a caller's own [HTTPDoer] fails with when the doer
// failed on another URL than the upload's, wrapped around that failure: never retried, since a retry
// could send the file there again.
var errSignedUploadOtherURL = errors.New("signed upload: the HTTP client failed on another URL than the upload's (a redirect it followed or refused, or a URL it rewrote), so the upload is not retried")

// isSignedUploadRedirectRefusal reports whether err is one of the redirect refusals above, which end
// the upload rather than being retried.
func isSignedUploadRedirectRefusal(err error) bool {
	return errors.Is(err, errSignedUploadFollowedRedirect) || errors.Is(err, errSignedUploadOtherURL)
}

// signedUploadOtherURLError is [errSignedUploadOtherURL] around the doer's redacted failure. It is a
// type of its own so that [signedUploadWatchdog.timeoutError] can swap the cause for a timeout and
// keep the refusal.
type signedUploadOtherURLError struct {
	cause error
}

// Error is the refusal followed by its cause.
func (e *signedUploadOtherURLError) Error() string {
	return errSignedUploadOtherURL.Error() + ": " + e.cause.Error()
}

// Unwrap matches both the refusal and its cause.
func (e *signedUploadOtherURLError) Unwrap() []error {
	return []error{errSignedUploadOtherURL, e.cause}
}

// signedUploadDoer sends an upload through a caller's own [HTTPDoer].
//
// Only an [*http.Client] can be told not to follow a redirect, and an HTTPDoer follows them however it
// was built to, so a response it followed one to is refused (see [UploadToSignedURL]). net/http sets
// a request's Response only on a request made to follow a redirect, so a doer that rewrites the URL
// is not mistaken for one, and a redirect to the same URL is still caught.
type signedUploadDoer struct {
	doer HTTPDoer
}

// Do sends req through the caller's doer, refusing a response it followed a redirect to.
//
// A redirected request that fails instead of answering leaves no response to look at, only the
// [*url.Error] naming the URL it failed on (net/http names a redirect's Location there when the
// client's CheckRedirect refuses it, too). One naming another URL than req's is refused, its cause
// redacted and kept: a retry could send the file there again. A doer that rewrites URLs itself then
// gets no retry on a failure, the safe side to err on.
func (d signedUploadDoer) Do(req *http.Request) (*http.Response, error) {
	// Without GetBody, net/http declines to follow a 307 or 308, which would send the file again, and
	// hands the redirect back for the upload to refuse; a 301, 302 or 303 it follows sends no body.
	req.GetBody = nil
	res, err := d.doer.Do(req)
	if err == nil && res != nil && res.Request != nil && res.Request.Response != nil {
		// A response a doer built itself may have no Body.
		if res.Body != nil {
			_ = res.Body.Close()
		}
		return nil, errSignedUploadFollowedRedirect
	}
	var urlErr *url.Error
	if err != nil && errors.As(err, &urlErr) {
		if !isSignedUploadRequestURL(urlErr.URL, req.URL) {
			return nil, &signedUploadOtherURLError{cause: redactSignedUploadCause(err, req.URL.String())}
		}
	}
	return res, err
}

// signedUploadBoundedDoer sends an attempt through a marked RoundTripper or HTTPDoer of the caller's,
// returning once the attempt's context ends even when that code ignores it.
//
// The upload's own transport stops when the request's context is done, but caller code can hand the
// request on under a context of its own (a rate limiter's wait, a detached tracing context), and the
// upload's idle timeout and the caller's cancel would then never reach it. So Do runs on a goroutine
// of its own; once the context ends the attempt is given up, while the caller's code runs on until
// it returns, and a response it hands back late is closed. A body it returns is closed when the
// context ends too, so a read the code does not end cannot hold the upload either.
type signedUploadBoundedDoer struct {
	doer HTTPDoer
}

// signedUploadAbandonGrace is how long a marked attempt whose context ended may still return before
// it is given up on.
const signedUploadAbandonGrace = 250 * time.Millisecond

// signedUploadAttemptResult is one Do call's outcome, handed from its goroutine.
type signedUploadAttemptResult struct {
	res *http.Response
	err error
}

// Do sends req, bounded by req's context.
func (d signedUploadBoundedDoer) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	done := make(chan signedUploadAttemptResult, 1)
	go func() {
		res, err := d.doer.Do(req)
		done <- signedUploadAttemptResult{res: res, err: err}
	}()
	select {
	case result := <-done:
		if result.res != nil && result.res.Body != nil {
			body := result.res.Body
			stop := context.AfterFunc(ctx, func() { _ = body.Close() })
			result.res.Body = &signedUploadStoppingBody{ReadCloser: body, stop: stop}
		}
		return result.res, result.err
	case <-ctx.Done():
		// Code that honours the context returns at once, and its answer (a redirect refusal, say) is
		// kept; only code that ignores it is given up on after the grace.
		grace := time.NewTimer(signedUploadAbandonGrace)
		defer grace.Stop()
		select {
		case result := <-done:
			return result.res, result.err
		case <-grace.C:
		}
		go func() {
			if late := <-done; late.res != nil && late.res.Body != nil {
				_ = late.res.Body.Close()
			}
		}()
		method := req.Method
		if method != "" {
			method = method[:1] + strings.ToLower(method[1:])
		}
		return nil, &url.Error{Op: method, URL: req.URL.String(), Err: ctx.Err()}
	}
}

// signedUploadStoppingBody is a response body that, once closed, no longer needs closing when its
// attempt's context ends (see [signedUploadBoundedDoer]).
type signedUploadStoppingBody struct {
	io.ReadCloser
	stop func() bool
}

// Close stops the context hook and closes the body.
func (b *signedUploadStoppingBody) Close() error {
	b.stop()
	return b.ReadCloser.Close()
}

// isSignedUploadRequestURL reports whether named, the URL a [*url.Error] names, is req's own: as
// written, as net/http writes it with the password masked, or parsing to the same scheme, host and
// path. net/http writes the masked user name decoded, so a `/` in it would parse as another host;
// comparing net/http's own spelling first keeps that from reading as another URL.
func isSignedUploadRequestURL(named string, req *url.URL) bool {
	if named == req.String() || named == signedUploadMaskedURL(req) {
		return true
	}
	parsed, err := url.Parse(named)
	return err == nil && strings.EqualFold(parsed.Scheme, req.Scheme) && strings.EqualFold(parsed.Host, req.Host) && parsed.EscapedPath() == req.EscapedPath()
}

// signedUploadMaskedURL is u as net/http names it in a [*url.Error]: with a password masked as `***`
// after the decoded user name.
func signedUploadMaskedURL(u *url.URL) string {
	if _, set := u.User.Password(); set {
		return strings.Replace(u.String(), u.User.String()+"@", u.User.Username()+":***@", 1)
	}
	return u.String()
}

// newSignedUploadClient returns the fresh client an upload through a clone of base, a plain
// transport [SignedUploadTransport] picked, is sent with.
//
// Redirects are handed back instead of followed (see [UploadToSignedURL]), and the client has no
// cookie jar or overall timeout. The clone keeps base's proxy, TLS configuration and dialers, so they
// apply to the upload as they do to the API. The wait for response headers is bounded unless base
// already bounds it, the write buffer is always the default size, no connection is kept alive past
// its request, and HTTP/2 runs on the clone's own connections rather than on a pool an upgrade copied
// from base.
func newSignedUploadClient(base *http.Transport) *http.Client {
	client := &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	transport := base.Clone()
	if transport.ResponseHeaderTimeout == 0 {
		transport.ResponseHeaderTimeout = signedUploadResponseHeaderTimeout
	}
	// The idle timeout sees progress only when the transport reads the next piece of the file, which
	// it does once the connection took the last. A large write buffer a caller chose for throughput
	// would take megabytes in one read and then drain for longer than the timeout, so a steady upload
	// on a slow link would be cut off as idle. The default size keeps every read a step of progress.
	transport.WriteBufferSize = 0
	// Each upload's transport is thrown away when it returns, so no connection may outlive it. The
	// CloseIdleConnections when the upload returns misses an HTTP/2 connection whose stream is not
	// yet marked done, and a caller transport often has no IdleConnTimeout, so that connection would
	// stay open for good. Without keep-alives every connection closes once its request is done; a
	// retry opens a new one.
	transport.DisableKeepAlives = true
	// A transport set up with golang.org/x/net/http2's ConfigureTransport carries an "h2" upgrade
	// bound to that transport's own HTTP/2 connection pool, and Clone copies it: the upload's HTTP/2
	// connections would land in the caller's pool, outlive the upload and ignore DisableKeepAlives.
	// Dropping the copied upgrades makes the clone set up HTTP/2 of its own; ForceAttemptHTTP2 keeps
	// it from skipping that for a transport with custom TLS or dialers, which would otherwise still
	// offer "h2" in its TLS handshake with nothing to speak it. Any other upgrade is dropped with it.
	// An empty, non-nil map is how a caller turns HTTP/2 off, and is kept.
	if len(transport.TLSNextProto) > 0 {
		transport.TLSNextProto = nil
		transport.ForceAttemptHTTP2 = true
	}
	client.Transport = transport
	return client
}

// sendSignedUpload sends one upload attempt, bounded by the per-attempt timeout when there is one.
//
// The timeout is an idle timeout: it counts time without progress, and more of the file handed to
// the connection, the response headers arriving and each piece of the response body all restart it
// (see [signedUploadWatchdog]). A deadline over the whole attempt would cut a large file off on a
// slow link however steadily it was going out, and the retry would send all of it again, only to be
// cut off at the same point. A host that stops reading or never answers still fails the attempt
// after timeout, with the [context.DeadlineExceeded] a deadline would have given.
//
// The response body is read in full before the attempt's context is released, so a timeout that
// fires while it is being read is reported by this attempt rather than by a later read.
func sendSignedUpload(ctx context.Context, client HTTPDoer, timeout time.Duration, method string, target string, headers http.Header, body []byte) (*http.Response, error) {
	attemptCtx := ctx
	var watchdog *signedUploadWatchdog
	// A caller's deadline that falls before the first moment the watchdog could fire needs no
	// watchdog: the deadline ends the attempt first, and is reported as the caller's own.
	if timeout > 0 && isBeforeContextDeadline(time.Now().Add(timeout), ctx) {
		var cancel context.CancelFunc
		attemptCtx, cancel = context.WithCancel(ctx)
		defer cancel()
		watchdog = newSignedUploadWatchdog(timeout, cancel)
		defer watchdog.stop()
	}
	// The length is set by hand because the body is no longer a type net/http knows the size of, and
	// without one the file would go out with chunked encoding, which object stores refuse.
	getBody := func() (io.ReadCloser, error) {
		return io.NopCloser(&signedUploadProgressReader{reader: bytes.NewReader(body), watchdog: watchdog}), nil
	}
	requestBody, _ := getBody()
	req, err := http.NewRequestWithContext(attemptCtx, method, target, requestBody)
	if err != nil {
		return nil, err
	}
	req.ContentLength = int64(len(body))
	req.GetBody = getBody
	if len(body) == 0 {
		// A non-nil body of length zero would read as "length unknown" and be sent chunked.
		req.Body = http.NoBody
		req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
	}
	req.Header = headers.Clone()
	res, err := client.Do(req)
	if err != nil {
		return nil, watchdog.timeoutError(ctx, err)
	}
	// Only an *http.Client promises a response and a body; a marked HTTPDoer of the caller's can return
	// neither, which would otherwise panic below.
	if res == nil {
		return nil, errors.New("the HTTP client returned no response and no error")
	}
	if res.Body == nil {
		res.Body = http.NoBody
	}
	// The response headers arriving is progress too.
	watchdog.progress()
	// Only net/http's own transport promises a response its request. A caller's RoundTripper or
	// HTTPDoer can answer with a response of its own making (a circuit breaker's 503, say), and the
	// error built from it (see [signedUploadStatusError]) prints and redacts that request.
	if res.Request == nil {
		res.Request = req
	}
	res.Body = &signedUploadProgressBody{ReadCloser: res.Body, watchdog: watchdog}
	// Bounded like every other response the SDK reads into memory: an error answer is cut at
	// maxErrorBodySize, and a success larger than maxResponseBodySize fails the attempt.
	read := readBoundedBody
	if res.StatusCode >= 400 {
		read = readErrorBody
	}
	// Without a per-attempt timeout, the upload host (which the API's response named) could
	// send headers and then stall: the read gets defaultBodyReadTimeout, as in Execute, cut
	// short by closing the body under it.
	var timedOut atomic.Bool
	if watchdog == nil {
		stalled := res.Body
		timer := time.AfterFunc(defaultBodyReadTimeout, func() {
			timedOut.Store(true)
			_ = stalled.Close()
		})
		defer timer.Stop()
	}
	contents, err := read(res.Body)
	_ = res.Body.Close()
	if err != nil {
		if timedOut.Load() {
			return nil, fmt.Errorf("the response body was not read within %s: %w", defaultBodyReadTimeout, err)
		}
		return nil, watchdog.timeoutError(ctx, err)
	}
	res.Body = io.NopCloser(bytes.NewReader(contents))
	return res, nil
}

// signedUploadWatchdog ends an upload attempt that has made no progress for timeout, by cancelling
// the attempt's context. Every sign of progress pushes that moment back by the full timeout.
//
// A nil watchdog is an attempt with no idle timeout, and every method on it does nothing.
type signedUploadWatchdog struct {
	timeout time.Duration
	// fired records that the watchdog, and not the caller, cancelled the attempt, so the failure can
	// be reported as a timeout instead of a cancellation.
	fired atomic.Bool

	mu      sync.Mutex
	timer   *time.Timer
	stopped bool
}

// newSignedUploadWatchdog starts a watchdog that calls cancel after timeout without progress.
func newSignedUploadWatchdog(timeout time.Duration, cancel context.CancelFunc) *signedUploadWatchdog {
	watchdog := &signedUploadWatchdog{timeout: timeout}
	watchdog.timer = time.AfterFunc(timeout, func() {
		watchdog.fired.Store(true)
		cancel()
	})
	return watchdog
}

// progress restarts the idle clock.
func (w *signedUploadWatchdog) progress() {
	if w == nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	// The transport can still be reading the request body after the attempt has returned (a host
	// that answers before it has the whole file), and a reset then would re-arm a timer nobody stops.
	if !w.stopped {
		w.timer.Reset(w.timeout)
	}
}

// stop ends the watchdog once its attempt is over.
func (w *signedUploadWatchdog) stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.stopped = true
	w.timer.Stop()
}

// timeoutError turns the failure of an attempt the watchdog cancelled into the error a deadline on
// the attempt's context would have produced; any other failure is returned as it is.
//
// A cancelled context fails the request with [context.Canceled], which would tell the caller that
// they cancelled the upload. Callers match a timeout with errors.Is(err, context.DeadlineExceeded)
// or [net.Error]'s Timeout, and [UploadToSignedURL] retries it as a failed attempt, so the error
// keeps the shape net/http gives a deadline. When the caller's own context is done as well, that
// is the reason the caller is told, so the error is left alone.
func (w *signedUploadWatchdog) timeoutError(ctx context.Context, err error) error {
	if w == nil || !w.fired.Load() || ctx.Err() != nil {
		return err
	}
	// A redirect refusal stays one, so the upload ends instead of being retried as a timeout; only the
	// cause the watchdog's cancellation gave it becomes the timeout.
	var otherURL *signedUploadOtherURLError
	if errors.As(err, &otherURL) {
		return &signedUploadOtherURLError{cause: w.timeoutError(ctx, otherURL.cause)}
	}
	if errors.Is(err, errSignedUploadFollowedRedirect) {
		return err
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return &url.Error{Op: urlErr.Op, URL: urlErr.URL, Err: context.DeadlineExceeded}
	}
	return context.DeadlineExceeded
}

// signedUploadProgressReader is the request body of an upload attempt, reporting each read to the
// attempt's watchdog.
//
// The transport reads the next piece of the file only once the connection has accepted the one
// before, so a read is the sign that bytes are going out; when the host stops reading, the reads
// stop and the idle clock runs out. It deliberately has no WriteTo: with one, the whole file would
// be handed over in a single call that reports nothing until it ends.
type signedUploadProgressReader struct {
	reader   io.Reader
	watchdog *signedUploadWatchdog
}

// Read reads more of the file, counting it as progress.
func (r *signedUploadProgressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		r.watchdog.progress()
	}
	return n, err
}

// signedUploadProgressBody is the upload host's response body, reporting each piece that arrives to
// the attempt's watchdog, so a host answering slowly but steadily is not cut off either.
type signedUploadProgressBody struct {
	io.ReadCloser
	watchdog *signedUploadWatchdog
}

// Read reads more of the response, counting it as progress.
func (b *signedUploadProgressBody) Read(p []byte) (int, error) {
	n, err := b.ReadCloser.Read(p)
	if n > 0 {
		b.watchdog.progress()
	}
	return n, err
}

// replayableSignedUploadBody returns a GetBody that hands out the same bytes on every call.
func replayableSignedUploadBody(body []byte) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
}

// waitForSignedUploadRetry waits out a retry delay, or returns early when the caller cancels.
func waitForSignedUploadRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// signedUploadRetryDelay is how long to wait before retrying: what `retry-after-ms`/`Retry-After`
// asks for when that is a sane wait (positive and under a minute), else the client's exponential
// backoff with jitter.
func signedUploadRetryDelay(res *http.Response, attempt int) time.Duration {
	if wait, ok := parseRetryAfterHeader(res); ok && wait > 0 && wait < signedUploadMaxRetryAfter {
		return wait
	}
	return retryDelay(nil, attempt)
}

// signedUploadStatusError turns a failed upload response into the SDK's own API error.
//
// The error's request carries the upload URL without its query string: a pre-signed URL's
// signature lives there, and it would otherwise reach every log the error is printed to. A body that
// is not JSON (object stores often answer in XML) still yields the API error, with no decoded fields.
func signedUploadStatusError(res *http.Response) error {
	contents, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return err
	}
	res.Body = io.NopCloser(bytes.NewReader(contents))
	request := res.Request
	if request != nil {
		request = request.Clone(context.Background())
		request.URL = redactedSignedUploadURL(request.URL)
		request.Body = nil
		request.GetBody = nil
		res.Request = request
	}
	apiErr := apierror.Error{Request: request, Response: res, StatusCode: res.StatusCode}
	if len(contents) > 0 {
		_ = apiErr.UnmarshalJSON(contents)
	}
	return &apiErr
}

// signedUploadMinRedactedLength is the shortest query, user info or user name the error scrub looks
// for on its own. A signature is far longer; a shorter one (`v=1`) is no secret, and taking it out
// would shred unrelated text.
const signedUploadMinRedactedLength = 8

// redactSignedUploadError strips the pre-signed URL's query and user info from a transport error, for
// the reason [signedUploadStatusError] gives.
//
// The [*url.Error] net/http wraps a failure in names the URL in a field of its own, which is
// redacted; the standard library's errors below it never name the URL. A caller's RoundTripper or
// HTTPDoer can, though, in its own text (a retrying wrapper's "PUT <url> giving up", say) or in an
// error it wraps (an inner client's [*url.Error]), so a cause whose chain names the URL's secrets
// anywhere is sealed (see [sealSignedUploadError]).
func redactSignedUploadError(err error, target string) error {
	return fmt.Errorf("signed upload: %w", redactSignedUploadCause(err, target))
}

// redactSignedUploadCause is err redacted as [redactSignedUploadError] describes, without its prefix.
func redactSignedUploadCause(err error, target string) error {
	parsed, parseErr := url.Parse(target)
	if parseErr != nil {
		return err
	}
	if urlErr, ok := err.(*url.Error); ok {
		failed := redactedSignedUploadURL(parsed).String()
		if named, namedErr := url.Parse(urlErr.URL); namedErr == nil {
			failed = redactedSignedUploadURL(named).String()
		}
		// A URL that does not parse as written (a user name with a `/`, which net/http writes decoded)
		// can leave a secret in what redaction kept, so the field is scrubbed as text too.
		for _, secret := range signedUploadSecrets(target, parsed) {
			// The whole URL is a secret only for its query and user info, which redaction already took
			// out; without either, the redacted URL is the URL itself and stays readable.
			if secret != target {
				failed = strings.ReplaceAll(failed, secret, "[redacted]")
			}
		}
		return &url.Error{Op: urlErr.Op, URL: failed, Err: sealSignedUploadError(urlErr.Err, target, parsed)}
	}
	return sealSignedUploadError(err, target, parsed)
}

// sealSignedUploadError returns err as it is when no error in its chain names the upload URL's
// secrets, and otherwise a [signedUploadSealedError] carrying only its scrubbed text.
//
// Scrubbing the text alone is not enough: errors.As and an error reporter walking the chain reach
// every link, and a caller's own error type prints whatever it holds. So a chain that names a secret
// is not handed out at all. The cost is that a caller's own error type in such a chain is no longer
// found with errors.As; errors.Is (a sentinel, a timeout, a cancellation) still matches, since it
// answers only yes or no.
func sealSignedUploadError(err error, target string, parsed *url.URL) error {
	secrets := signedUploadSecrets(target, parsed)
	if err == nil || !signedUploadChainNames(err, secrets) {
		return err
	}
	var timeout interface{ Timeout() bool }
	scrubbed := err.Error()
	for _, secret := range secrets {
		replacement := "[redacted]"
		if secret == target {
			replacement = redactedSignedUploadURL(parsed).String()
		}
		scrubbed = strings.ReplaceAll(scrubbed, secret, replacement)
	}
	return &signedUploadSealedError{
		message: scrubbed,
		timeout: errors.As(err, &timeout) && timeout.Timeout(),
		cause:   err,
	}
}

// signedUploadSecrets are the pieces of the upload URL an error must not carry, longest first so the
// query is not cut out of the middle of the whole URL and leave the rest behind: the URL itself, its
// query, its user info, its user name on its own, both as written and as decoded, since net/http
// writes a URL with only the password masked (`user:***@`), and each query value of
// [signedUploadMinRedactedLength] bytes or more. signedUploadURL refuses user info in an upload URL,
// so those entries only guard a caller of redactSignedUploadError that passes one.
func signedUploadSecrets(target string, parsed *url.URL) []string {
	secrets := []string{target}
	var name string
	if parsed.User != nil {
		name = url.User(parsed.User.Username()).String()
	}
	candidates := []string{parsed.RawQuery, parsed.User.String(), name, parsed.User.Username()}
	// Each query value on its own too, as written and decoded: code that rebuilds the query (Encode
	// sorts the keys) or prints it as a map never writes the whole query as it came.
	for _, part := range strings.Split(parsed.RawQuery, "&") {
		_, value, _ := strings.Cut(part, "=")
		candidates = append(candidates, value)
		if decoded, err := url.QueryUnescape(value); err == nil && decoded != value {
			candidates = append(candidates, decoded)
		}
	}
	for _, secret := range candidates {
		if len(secret) >= signedUploadMinRedactedLength {
			secrets = append(secrets, secret)
		}
	}
	return secrets
}

// signedUploadChainNames reports whether err, or any error it wraps (through Unwrap() error or
// Unwrap() []error), names one of secrets in its text.
func signedUploadChainNames(err error, secrets []string) bool {
	if err == nil {
		return false
	}
	text := err.Error()
	for _, secret := range secrets {
		if strings.Contains(text, secret) {
			return true
		}
	}
	switch wrapped := err.(type) {
	case interface{ Unwrap() error }:
		return signedUploadChainNames(wrapped.Unwrap(), secrets)
	case interface{ Unwrap() []error }:
		for _, inner := range wrapped.Unwrap() {
			if signedUploadChainNames(inner, secrets) {
				return true
			}
		}
	}
	return false
}

// signedUploadSealedError stands in for a transport failure whose chain named the upload URL's
// secrets (see [sealSignedUploadError]). It has no Unwrap or As, so nothing reaches the original
// text, and keeps only what callers match a failed attempt by: errors.Is against the original chain,
// and [net.Error]'s Timeout.
type signedUploadSealedError struct {
	message string
	timeout bool
	// cause is the original failure, consulted only by Is, which never hands it out.
	cause error
}

// Error is the original text with the secrets scrubbed out.
func (e *signedUploadSealedError) Error() string { return e.message }

// Timeout reports whether the original failure was a timeout.
func (e *signedUploadSealedError) Timeout() bool { return e.timeout }

// Temporary completes [net.Error], so errors.As for one still finds a sealed failure and its Timeout;
// net.Error deprecates it, and nothing here can tell.
func (e *signedUploadSealedError) Temporary() bool { return false }

// Is matches what the original failure matched, such as [context.DeadlineExceeded].
func (e *signedUploadSealedError) Is(target error) bool { return errors.Is(e.cause, target) }

// redactedSignedUploadURL is u without its query, fragment or user info.
func redactedSignedUploadURL(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}
	redacted := *u
	redacted.User = nil
	redacted.RawQuery = ""
	redacted.ForceQuery = false
	redacted.Fragment = ""
	redacted.RawFragment = ""
	return &redacted
}

// signedUploadURL reads and validates the upload URL, refusing a missing value or one that is not
// https.
//
// The URL came from a response rather than from the caller, so a `file:` or `data:` URL is refused
// before any request is made, and so is plain `http:` to anything but a loopback host: it would send
// the file, the signature in the query and the named headers in cleartext, for anyone on the path to
// replay until the URL expires.
func signedUploadURL(initiated []byte, path []string) (string, error) {
	field := strings.Join(path, ".")
	value, ok := ReadSignedUploadField(initiated, path).(string)
	if !ok || value == "" {
		return "", fmt.Errorf("signed upload: expected the upload response to carry an upload URL in `%s`", field)
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return "", fmt.Errorf("signed upload: the upload URL in `%s` is not a valid URL", field)
	}
	scheme := strings.ToLower(parsed.Scheme)
	if (scheme != "https" && scheme != "http") || parsed.Hostname() == "" {
		return "", fmt.Errorf("signed upload: refusing to upload to `%s`: expected an absolute https: URL", field)
	}
	if scheme == "http" && !isSignedUploadLoopbackHost(parsed.Hostname()) {
		return "", fmt.Errorf("signed upload: refusing to upload to `%s` over plain http:; expected an https: URL", field)
	}
	// A user name or password in the URL would reach whatever sends the upload, and an *http.Client
	// turns it into an Authorization header before a marked RoundTripper of the caller's sees the
	// request, which that code may log or trace. A pre-signed URL carries its authorization in the
	// query, so one with user info is refused before any transport sees it, naming only the field.
	if parsed.User != nil {
		return "", fmt.Errorf("signed upload: refusing to upload to `%s`: the URL carries a user name or password", field)
	}
	return value, nil
}

// isSignedUploadLoopbackHost reports whether host is this machine: exactly `localhost`, `127.0.0.1`
// or `::1`. A `*.localhost` name is not one: not every resolver sends it to loopback.
func isSignedUploadLoopbackHost(host string) bool {
	switch strings.ToLower(host) {
	case "localhost", "127.0.0.1", "::1":
		return true
	}
	return false
}

// signedUploadRequestParts builds the headers and body of the upload request: nothing from the SDK
// client, by construction.
//
// The headers are the ones the initiating response names. A multipart upload drops a named
// `Content-Type`, since the multipart writer sets one carrying its boundary. A binary upload's
// `Content-Type` is a named one when the response gives it — a pre-signed URL can be signed for the
// type it was issued for — then the configured one, then the file's own, then
// `application/octet-stream`.
func signedUploadRequestParts(initiated []byte, spec SignedUploadSpec, file SignedUploadFile) (http.Header, []byte, error) {
	named := map[string]string{}
	if spec.HeadersField != nil {
		if object, ok := ReadSignedUploadField(initiated, spec.HeadersField).(map[string]interface{}); ok {
			for name, value := range object {
				if text, ok := value.(string); ok {
					named[name] = text
				}
			}
		}
	}
	names := make([]string, 0, len(named))
	for name := range named {
		names = append(names, name)
	}
	// Sorted so that two names differing only in case resolve the same way on every run.
	sort.Strings(names)

	headers := http.Header{}
	hasContentType := false
	for _, name := range names {
		if strings.EqualFold(name, "Content-Type") {
			if spec.Multipart {
				continue
			}
			hasContentType = true
		}
		if err := checkSignedUploadHeader(name, named[name]); err != nil {
			return nil, nil, err
		}
		headers.Set(name, named[name])
	}

	if spec.Multipart {
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)
		field := spec.MultipartField
		if field == "" {
			field = "file"
		}
		filename := file.Name
		if filename == "" {
			filename = "file"
		}
		partType := file.ContentType
		if partType == "" {
			partType = signedUploadDefaultContentType
		}
		if err := checkSignedUploadHeader("Content-Type", partType); err != nil {
			return nil, nil, err
		}
		part := textproto.MIMEHeader{}
		part.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeSignedUploadDispositionValue(field), escapeSignedUploadDispositionValue(filename)))
		part.Set("Content-Type", partType)
		partWriter, err := writer.CreatePart(part)
		if err != nil {
			return nil, nil, err
		}
		if _, err := partWriter.Write(file.Content); err != nil {
			return nil, nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, nil, err
		}
		headers.Set("Content-Type", writer.FormDataContentType())
		return headers, buffer.Bytes(), nil
	}

	if !hasContentType {
		contentType := spec.ContentType
		if contentType == "" {
			contentType = file.ContentType
		}
		if contentType == "" {
			contentType = signedUploadDefaultContentType
		}
		if err := checkSignedUploadHeader("Content-Type", contentType); err != nil {
			return nil, nil, err
		}
		headers.Set("Content-Type", contentType)
	}
	return headers, file.Content, nil
}

// signedUploadDispositionEscaper escapes the two characters a quoted `Content-Disposition`
// parameter cannot hold as they are, and drops the line breaks and NUL that would end the header.
var signedUploadDispositionEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"", "\r", "", "\n", "", "\x00", "")

// escapeSignedUploadDispositionValue makes a file or field name safe inside a quoted
// `Content-Disposition` parameter.
func escapeSignedUploadDispositionValue(value string) string {
	return signedUploadDispositionEscaper.Replace(value)
}

// checkSignedUploadHeader refuses a header no HTTP request can carry, once, rather than on every
// retried attempt.
//
// The transport would otherwise fail every attempt the same way, and that failure reads as a
// connection error, which is retried with backoff. The rules are the ones [http.Transport] enforces:
// a name must be a non-empty token (so no CR, LF, NUL or `:`), and a value may hold no control
// character other than a tab.
func checkSignedUploadHeader(name string, value string) error {
	valid := name != ""
	for i := 0; valid && i < len(name); i++ {
		valid = isSignedUploadTokenByte(name[i])
	}
	for i := 0; valid && i < len(value); i++ {
		c := value[i]
		valid = c == '\t' || (c >= 0x20 && c != 0x7f)
	}
	if !valid {
		return fmt.Errorf("signed upload: the upload header %q is not a valid HTTP header", name)
	}
	return nil
}

// isSignedUploadTokenByte reports whether c may appear in an HTTP header name (an RFC 7230 token).
func isSignedUploadTokenByte(c byte) bool {
	if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
		return true
	}
	return strings.IndexByte("!#$%&'*+-.^_`|~", c) >= 0
}
