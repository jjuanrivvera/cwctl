package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// Resource[T] is the generic CRUD handle. Every resource reuses it; the only per-resource
// code is the struct T and a Client accessor. This is the "generic core, thin resources"
// guarantee: adding a resource never edits this file.
type Resource[T any] struct {
	client       *Client
	path         string // full collection path relative to the instance root
	updateMethod string // PATCH by default; contacts/profile use PUT (DECISIONS.md #10)
}

// NewResource builds a typed handle to a collection path (already account-scoped by the
// Client accessor when applicable).
func NewResource[T any](c *Client, path string) *Resource[T] {
	return &Resource[T]{client: c, path: path, updateMethod: http.MethodPatch}
}

// WithUpdateMethod overrides the HTTP verb Update uses — the generic-core knob for
// Chatwoot's PUT endpoints, never a per-resource fork.
func (r *Resource[T]) WithUpdateMethod(m string) *Resource[T] {
	r.updateMethod = m
	return r
}

// Client returns the underlying client (used by custom verbs that need Raw/multipart).
func (r *Resource[T]) Client() *Client { return r.client }

// Path returns the collection path (used by Extra commands to build sub-paths).
func (r *Resource[T]) Path() string { return r.path }

// List fetches one page, normalizing whichever envelope the endpoint uses.
func (r *Resource[T]) List(ctx context.Context, p ListParams) ([]T, error) {
	items, _, err := r.listPage(ctx, p)
	return items, err
}

func (r *Resource[T]) listPage(ctx context.Context, p ListParams) ([]T, *ListMeta, error) {
	_, _, data, err := r.client.Raw(ctx, http.MethodGet, r.path, p.values(), nil)
	if err != nil || data == nil { // error or dry-run
		return nil, nil, err
	}
	return decodeList[T](data)
}

// ListAll walks pages until the API stops producing new items, honoring ctx cancellation
// between pages. Chatwoot page sizes differ per endpoint and some lists ignore `page`
// entirely, so the stop conditions are (in order): an empty page, a page identical to the
// previous one (unpaginated endpoints echo page 1 forever), or the envelope's advertised
// total reached (DECISIONS.md #7). A short page is deliberately NOT a stop signal.
func (r *Resource[T]) ListAll(ctx context.Context, p ListParams) ([]T, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	var all []T
	var prevSig []byte
	// A hard page cap so a pathological server can't loop us forever; 10k pages of 25 is
	// far beyond any real account.
	for range 10000 {
		items, meta, err := r.listPage(ctx, p)
		if err != nil {
			return all, err
		}
		if len(items) == 0 {
			break
		}
		sig, err := json.Marshal(items)
		if err != nil {
			return all, err
		}
		if prevSig != nil && bytes.Equal(sig, prevSig) {
			break
		}
		prevSig = sig
		all = append(all, items...)

		if total := meta.total(); total >= 0 && len(all) >= total {
			break
		}
		p.Page++
		if err := ctx.Err(); err != nil {
			return all, err
		}
	}
	return all, nil
}

// Get fetches a single record by id, unwrapping {payload:{…}}/{data:{…}} envelopes.
func (r *Resource[T]) Get(ctx context.Context, id string) (*T, error) {
	_, _, data, err := r.client.Raw(ctx, http.MethodGet, r.path+"/"+url.PathEscape(id), nil, nil)
	if err != nil || data == nil {
		return nil, err
	}
	var out T
	if err := decodeOne(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create POSTs a body and decodes the response. POST is never auto-retried (see retry.go).
func (r *Resource[T]) Create(ctx context.Context, body any, out any) error {
	return r.send(ctx, http.MethodPost, r.path, body, out)
}

// Update modifies a record using the resource's configured verb (PATCH or PUT).
func (r *Resource[T]) Update(ctx context.Context, id string, body any, out any) error {
	return r.send(ctx, r.updateMethod, r.path+"/"+url.PathEscape(id), body, out)
}

// Delete removes a record by id.
func (r *Resource[T]) Delete(ctx context.Context, id string) error {
	_, err := r.client.doJSON(ctx, http.MethodDelete, r.path+"/"+url.PathEscape(id), nil, nil, nil)
	return err
}

// Action performs a custom verb on the collection or a member (e.g.
// POST conversations/{id}/toggle_status, GET contacts/search), decoding the (possibly
// enveloped) response into out. Method defaults to GET when empty.
func (r *Resource[T]) Action(ctx context.Context, method, subPath string, query url.Values, body any, out any) error {
	if method == "" {
		method = http.MethodGet
	}
	path := r.path
	if subPath != "" {
		path = r.path + "/" + subPath
	}
	return r.client.Send(ctx, method, path, query, body, out)
}

// send marshals body (when non-nil) and decodes the enveloped response into out.
func (r *Resource[T]) send(ctx context.Context, method, path string, body, out any) error {
	return r.client.Send(ctx, method, path, nil, body, out)
}

// Send is the client-level JSON round-trip used by resources and irregular commands
// alike: marshal body → request → unwrap single-object envelopes into out.
func (c *Client) Send(ctx context.Context, method, path string, query url.Values, body, out any) error {
	var reader *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}
	var (
		data []byte
		err  error
	)
	if reader == nil {
		_, _, data, err = c.Raw(ctx, method, path, query, nil)
	} else {
		_, _, data, err = c.Raw(ctx, method, path, query, reader)
	}
	if err != nil || data == nil || out == nil {
		return err
	}
	switch o := out.(type) {
	case *json.RawMessage:
		*o = json.RawMessage(data)
		return nil
	default:
		return decodeOne(data, o)
	}
}
