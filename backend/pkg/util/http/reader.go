package util_http

import (
	"context"
	"github.com/go-playground/validator/v10"
)

type Reader struct {
	validation *validator.Validate
}

func NewReader() *Reader {
	return &Reader{
		validation: validator.New(),
	}
}

func (r Reader) Read(ctx context.Context, parser func(out interface{}) error, request interface{}) error {
	if err := parser(request); err != nil {
		return err
	}

	return r.validation.StructCtx(ctx, request)
}
