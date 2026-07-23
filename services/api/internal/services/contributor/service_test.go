package contributor_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/errorsx"
	contributorrepo "github.com/asl-open/asl-core/services/api/internal/repository/contributor"
	"github.com/asl-open/asl-core/services/api/internal/services/contributor"
)

func newService(repo contributorrepo.Repo) contributor.Service {
	return contributor.New(contributor.Params{Repo: repo})
}

func validInput() *contributor.CreateInput {
	return &contributor.CreateInput{
		DisplayName: "Muhammad al-Bukhari",
		Roles:       []contributorrepo.Role{contributorrepo.RoleContributor, contributorrepo.RoleReviewer},
	}
}

func TestService_Create(t *testing.T) {
	t.Run("valid input persists a trimmed contributor", func(t *testing.T) {
		repo := &contributorrepo.MockRepo{}
		svc := newService(repo)

		in := validInput()
		in.DisplayName = "  Muhammad al-Bukhari  "

		want := contributorrepo.Contributor{ID: uuid.New(), DisplayName: "Muhammad al-Bukhari"}
		repo.On("Create", mock.Anything, mock.MatchedBy(func(c *contributorrepo.Contributor) bool {
			return c.DisplayName == "Muhammad al-Bukhari" && len(c.Roles) == 2
		})).Return(want, nil)

		got, err := svc.Create(t.Context(), in)
		require.NoError(t, err)
		require.Equal(t, want, got)
		repo.AssertExpectations(t)
	})

	t.Run("rejects invalid input before touching the repository", func(t *testing.T) {
		cases := map[string]struct {
			mutate      func(in *contributor.CreateInput)
			wantMessage string
		}{
			"missing display name": {func(in *contributor.CreateInput) { in.DisplayName = "   " }, "contributor display name is required"},
			"no roles":             {func(in *contributor.CreateInput) { in.Roles = nil }, "contributor requires at least one role"},
			"invalid role":         {func(in *contributor.CreateInput) { in.Roles = []contributorrepo.Role{"scholar"} }, "contributor role is invalid"},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				repo := &contributorrepo.MockRepo{}
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
		repo := &contributorrepo.MockRepo{}
		svc := newService(repo)

		wantErr := errors.New("db down")
		repo.On("Create", mock.Anything, mock.Anything).Return(contributorrepo.Contributor{}, wantErr)

		_, err := svc.Create(t.Context(), validInput())
		require.ErrorIs(t, err, wantErr)
	})
}

func TestService_Get(t *testing.T) {
	repo := &contributorrepo.MockRepo{}
	svc := newService(repo)

	id := uuid.New()
	want := contributorrepo.Contributor{ID: id, DisplayName: "Muhammad al-Bukhari"}
	repo.On("Get", mock.Anything, id).Return(want, nil)

	got, err := svc.Get(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, want, got)
	repo.AssertExpectations(t)
}

func TestService_List(t *testing.T) {
	repo := &contributorrepo.MockRepo{}
	svc := newService(repo)

	want := []contributorrepo.Contributor{{ID: uuid.New()}, {ID: uuid.New()}}
	repo.On("List", mock.Anything).Return(want, nil)

	got, err := svc.List(t.Context())
	require.NoError(t, err)
	require.Equal(t, want, got)
	repo.AssertExpectations(t)
}
