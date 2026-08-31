package paging

import (
	"context"
	"fmt"
	"iter"
)

// ItemWithID is an interface that ensures an item has an id.
type ItemWithID interface {
	GetID() string
}

// PagedResponse is a generic response type wrapping paging logic of the anexia engine.
type PagedResponse[T any] struct {
	Page       int `json:"page"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Limit      int `json:"limit"`
	Data       []T `json:"data"`
}

// Params represents the pagination params for the anexia engine.
type Params struct {
	Page  int
	Limit int
}

// ParamsError represents an error in the params.
type ParamsError struct {
	Page int
}

func (p *ParamsError) Error() string {
	return fmt.Sprintf("page must be at least %d, got: %d", EngineFirstPage, p.Page)
}

// Validate validates if the [Params] are valid.
func (p Params) Validate() error {
	if p.Page < EngineFirstPage {
		return &ParamsError{p.Page}
	}
	return nil
}

// DefaultParams creates a new instance with default values.
func DefaultParams() Params {
	return NewParams(EngineFirstPage, EngineMaxPageLimit)
}

// NewParams creates a new instance of the [Params].
func NewParams(page, limit int) Params {
	return Params{
		Page:  page,
		Limit: limit,
	}
}

// PageFetcher is a function that fetches the desired page from the api.
type PageFetcher[T any] func(
	ctx context.Context,
	pageParams Params,
) (PagedResponse[T], error)

// ItemFetcher is a function that fetches a single resource from the api via its id.
type ItemFetcher[T any] func(
	ctx context.Context,
	id string,
) (T, error)

const (
	// EngineFirstPage is the first page the anexia engine provides.
	EngineFirstPage = 1
	// EngineMaxPageLimit is the maximum paging limit for most of the anexia engine endpoints.
	EngineMaxPageLimit = 100
)

// Paginate iterates all resources from a paged endpoint using the provided PageFetcher.
func Paginate[T any]( //revive:disable:cognitive-complexity
	ctx context.Context,
	fetchPage PageFetcher[T],
) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var zero T
		page := EngineFirstPage

		for {
			err := ctx.Err()
			if err != nil {
				yield(zero, err)
				return
			}

			resp, err := fetchPage(ctx, NewParams(page, EngineMaxPageLimit))
			if err != nil {
				yield(zero, err)
				return
			}

			for _, v := range resp.Data {
				if !yield(v, nil) {
					return
				}
			}

			if page >= resp.TotalPages {
				return
			}

			page++
		}
	}
}

// PaginateAndLoad iterates all resources from a paged endpoint
// using the provided PageFetcher and loads each item using the provided ItemFetcher.
func PaginateAndLoad[T ItemWithID, TResult any](
	ctx context.Context,
	fetchPage PageFetcher[T],
	fetchItem ItemFetcher[TResult],
) iter.Seq2[TResult, error] {
	return func(yield func(TResult, error) bool) {
		var zero TResult

		for item, err := range Paginate(ctx, fetchPage) {
			if err != nil {
				yield(zero, err)
				return
			}

			engineResource, err := fetchItem(ctx, item.GetID())
			if err != nil {
				yield(zero, err)
				return
			}

			if !yield(engineResource, nil) {
				return
			}
		}
	}
}

// CollectAll is a helper method to collect all items into a slice.
//
// Returns the first error found, if any, otherwise the resulting slice and nil is returned.
func CollectAll[T any](seq iter.Seq2[T, error]) ([]T, error) {
	result := make([]T, 0, EngineMaxPageLimit)

	for item, itemErr := range seq {
		if itemErr != nil {
			return nil, itemErr
		}
		result = append(result, item)
	}

	return result, nil
}
