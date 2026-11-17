package lbaas

import (
	"context"
	"errors"
	"fmt"

	"github.com/anexia/go-anxsdk/paging"
	"github.com/anexia/go-anxsdk/v2/common"
)

var (
	errEmptyReferenceFilter      = errors.New("lbaas reference filter is empty")
	errUnsupportedReferenceField = errors.New("unsupported lbaas reference field")
)

// ReferenceListParams selects one page of an LBaaS collection by a single relation.
// RefField is "frontend", "backend" or "default_backend".
type ReferenceListParams struct {
	Resource   string
	RefField   string
	RefID      string
	Attributes string
}

// ReferenceListItem is an LBaaS list row reduced to its identity and the relations that can filter it.
type ReferenceListItem struct {
	Identifier     string          `json:"identifier"`
	Name           string          `json:"name"`
	Frontend       common.Resource `json:"frontend"`
	Backend        common.Resource `json:"backend"`
	DefaultBackend common.Resource `json:"default_backend"`
}

// RelationID returns the identifier of the frontend, backend or default_backend relation.
func (item ReferenceListItem) RelationID(field string) string {
	switch field {
	case "frontend":
		return item.Frontend.Identifier
	case "backend":
		return item.Backend.Identifier
	case "default_backend":
		return item.DefaultBackend.Identifier
	default:
		return ""
	}
}

type referenceListQuery struct {
	Frontend       string `url:"frontend,omitempty"`
	Backend        string `url:"backend,omitempty"`
	DefaultBackend string `url:"default_backend,omitempty"`
	Attributes     string `url:"attributes,omitempty"`
}

func (p ReferenceListParams) query() (referenceListQuery, error) {
	if p.RefID == "" {
		return referenceListQuery{}, fmt.Errorf("%w: %s %s", errEmptyReferenceFilter, p.Resource, p.RefField)
	}
	q := referenceListQuery{Attributes: p.Attributes}
	switch p.RefField {
	case "frontend":
		q.Frontend = p.RefID
	case "backend":
		q.Backend = p.RefID
	case "default_backend":
		q.DefaultBackend = p.RefID
	default:
		return referenceListQuery{}, fmt.Errorf("%w: %q", errUnsupportedReferenceField, p.RefField)
	}
	return q, nil
}

// ListByReference lists one page of an LBaaS resource filtered by frontend, backend or default_backend.
func (c *Client) ListByReference(ctx context.Context, page paging.Params, params ReferenceListParams) (paging.PagedResponse[ReferenceListItem], error) {
	query, err := params.query()
	if err != nil {
		return paging.PagedResponse[ReferenceListItem]{}, err
	}
	resp := paging.PagedResponse[ReferenceListItem]{}
	err = c.transport.Get(ctx, fmt.Sprintf("/api/LBaaS/v2/%s", params.Resource), &resp, page, query)
	return resp, common.MapTransportError(err)
}
