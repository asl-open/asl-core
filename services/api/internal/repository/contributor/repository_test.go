package contributor_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/pkg/database"
	"github.com/asl-open/asl-core/pkg/errorsx"
	"github.com/asl-open/asl-core/pkg/logger"
	"github.com/asl-open/asl-core/services/api/internal/repository/contributor"
)

func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_DSN"); v != "" {
		return v
	}
	return "postgres://postgres:postgres@localhost:5439/asl_core_test?sslmode=disable"
}

func newRepo(t *testing.T) contributor.Repo {
	t.Helper()
	t.Setenv("DATABASE_DSN", testDSN())

	cfg, err := config.New()
	require.NoError(t, err)

	conn, err := database.New(database.Params{
		Lifecycle: fxtest.NewLifecycle(t),
		Config:    cfg,
		Logger:    &logger.MockLogger{},
	})
	require.NoError(t, err)

	_, err = conn.Exec(t.Context(), "TRUNCATE TABLE contributors")
	require.NoError(t, err)

	return contributor.New(contributor.Params{Conn: conn})
}

func ptr(s string) *string {
	return &s
}

func sample(name string) *contributor.Contributor {
	return &contributor.Contributor{
		DisplayName: name,
		Roles:       []contributor.Role{contributor.RoleContributor, contributor.RoleReviewer},
		Bio:         ptr("Hadith scholar"),
		Handle:      ptr("@" + name),
	}
}

func TestRepository_CreateGetList(t *testing.T) {
	repo := newRepo(t)
	ctx := t.Context()

	created, err := repo.Create(ctx, sample("bukhari"))
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, created.ID)
	require.Equal(t, "bukhari", created.DisplayName)
	require.Equal(t, []contributor.Role{contributor.RoleContributor, contributor.RoleReviewer}, created.Roles)
	require.Equal(t, "Hadith scholar", *created.Bio)
	require.Equal(t, "@bukhari", *created.Handle)
	require.False(t, created.CreatedAt.IsZero())

	got, err := repo.Get(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created, got)

	list, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, created, list[0])
}

func TestRepository_Create_nullableFields(t *testing.T) {
	repo := newRepo(t)

	created, err := repo.Create(t.Context(), &contributor.Contributor{
		DisplayName: "anon",
		Roles:       []contributor.Role{contributor.RoleTranslator},
	})
	require.NoError(t, err)
	require.Nil(t, created.Bio)
	require.Nil(t, created.Handle)
}

func TestRepository_Get_notFound(t *testing.T) {
	repo := newRepo(t)

	_, err := repo.Get(t.Context(), uuid.New())
	appErr, ok := errorsx.As(err)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, appErr.StatusCode())
}

func TestRepository_List_orderedByCreatedAtDesc(t *testing.T) {
	repo := newRepo(t)
	ctx := t.Context()

	first, err := repo.Create(ctx, sample("first"))
	require.NoError(t, err)

	time.Sleep(2 * time.Millisecond)

	second, err := repo.Create(ctx, sample("second"))
	require.NoError(t, err)

	list, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, second.ID, list[0].ID, "newest contributor first")
	require.Equal(t, first.ID, list[1].ID)
}
