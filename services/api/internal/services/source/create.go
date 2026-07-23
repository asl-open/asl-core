package source

import (
	"context"
	"strings"

	"github.com/asl-open/asl-core/services/api/internal/http/apierrors"
	sourcerepo "github.com/asl-open/asl-core/services/api/internal/repository/source"
)

type CreateInput struct {
	Title         string
	Author        string
	Edition       string
	Language      string
	LocatorScheme string
	Type          sourcerepo.Type
}

func (s *service) Create(ctx context.Context, in *CreateInput) (sourcerepo.Source, error) {
	if err := in.validate(); err != nil {
		return sourcerepo.Source{}, err
	}

	return s.repo.Create(ctx, &sourcerepo.Source{
		Title:         strings.TrimSpace(in.Title),
		Author:        strings.TrimSpace(in.Author),
		Type:          in.Type,
		Edition:       strings.TrimSpace(in.Edition),
		Language:      strings.TrimSpace(in.Language),
		LocatorScheme: strings.TrimSpace(in.LocatorScheme),
	})
}

func (in *CreateInput) validate() error {
	switch {
	case strings.TrimSpace(in.Title) == "":
		return apierrors.ErrBadRequest.WithMessage("source title is required")
	case strings.TrimSpace(in.Author) == "":
		return apierrors.ErrBadRequest.WithMessage("source author is required")
	case in.Type == "":
		return apierrors.ErrBadRequest.WithMessage("source type is required")
	case !in.Type.Valid():
		return apierrors.ErrBadRequest.WithMessage("source type is invalid")
	case strings.TrimSpace(in.Edition) == "":
		return apierrors.ErrBadRequest.WithMessage("source edition is required")
	case strings.TrimSpace(in.Language) == "":
		return apierrors.ErrBadRequest.WithMessage("source language is required")
	case strings.TrimSpace(in.LocatorScheme) == "":
		return apierrors.ErrBadRequest.WithMessage("source locator scheme is required")
	}

	return nil
}
