package source_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/errorsx"
	sourcerepo "github.com/asl-open/asl-core/services/api/internal/repository/source"
	"github.com/asl-open/asl-core/services/api/internal/services/source"
)

func newService(repo sourcerepo.Repo) source.Service {
	return source.New(source.Params{Repo: repo})
}

func validInput() *source.CreateInput {
	return &source.CreateInput{
		Title:         "Sahih al-Bukhari",
		Author:        "Muhammad al-Bukhari",
		Type:          sourcerepo.TypeHadithCollection,
		Edition:       "Dar Tawq al-Najat, 2001",
		Language:      "ar",
		LocatorScheme: "collection + hadith number",
	}
}

func TestService_Create(t *testing.T) {
	t.Run("valid input persists a trimmed source", func(t *testing.T) {
		repo := &sourcerepo.MockRepo{}
		svc := newService(repo)

		in := validInput()
		in.Title = "  Sahih al-Bukhari  "

		want := sourcerepo.Source{ID: uuid.New(), Title: "Sahih al-Bukhari"}
		repo.On("Create", mock.Anything, mock.MatchedBy(func(s *sourcerepo.Source) bool {
			return s.Title == "Sahih al-Bukhari" && s.Type == sourcerepo.TypeHadithCollection
		})).Return(want, nil)

		got, err := svc.Create(t.Context(), in)
		require.NoError(t, err)
		require.Equal(t, want, got)
		repo.AssertExpectations(t)
	})

	t.Run("rejects invalid input before touching the repository", func(t *testing.T) {
		cases := map[string]struct {
			mutate      func(in *source.CreateInput)
			wantMessage string
		}{
			"missing title":          {func(in *source.CreateInput) { in.Title = "   " }, "source title is required"},
			"missing author":         {func(in *source.CreateInput) { in.Author = "" }, "source author is required"},
			"missing type":           {func(in *source.CreateInput) { in.Type = "" }, "source type is required"},
			"invalid type":           {func(in *source.CreateInput) { in.Type = "newspaper" }, "source type is invalid"},
			"missing edition":        {func(in *source.CreateInput) { in.Edition = "" }, "source edition is required"},
			"missing language":       {func(in *source.CreateInput) { in.Language = "" }, "source language is required"},
			"missing locator scheme": {func(in *source.CreateInput) { in.LocatorScheme = "" }, "source locator scheme is required"},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				repo := &sourcerepo.MockRepo{}
				svc := newService(repo)

				in := validInput()
				tc.mutate(in)

				_, err := svc.Create(t.Context(), in)
				appErr, ok := errorsx.As(err)
				require.True(t, ok)
				require.Equal(t, http.StatusBadRequest, appErr.StatusCode())
				require.Equal(t, tc.wantMessage, appErr.Message())
				repo.AssertNotCalled(t, "Create")
			})
		}
	})

	t.Run("propagates a repository error", func(t *testing.T) {
		repo := &sourcerepo.MockRepo{}
		svc := newService(repo)

		wantErr := errors.New("db down")
		repo.On("Create", mock.Anything, mock.Anything).Return(sourcerepo.Source{}, wantErr)

		_, err := svc.Create(t.Context(), validInput())
		require.ErrorIs(t, err, wantErr)
	})
}

func TestService_Get(t *testing.T) {
	repo := &sourcerepo.MockRepo{}
	svc := newService(repo)

	id := uuid.New()
	want := sourcerepo.Source{ID: id, Title: "Sahih al-Bukhari"}
	repo.On("Get", mock.Anything, id).Return(want, nil)

	got, err := svc.Get(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, want, got)
	repo.AssertExpectations(t)
}

func TestService_List(t *testing.T) {
	repo := &sourcerepo.MockRepo{}
	svc := newService(repo)

	want := []sourcerepo.Source{{ID: uuid.New()}, {ID: uuid.New()}}
	repo.On("List", mock.Anything).Return(want, nil)

	got, err := svc.List(t.Context())
	require.NoError(t, err)
	require.Equal(t, want, got)
	repo.AssertExpectations(t)
}
