package source_test

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
	"github.com/asl-open/asl-core/services/api/internal/repository/source"
)

func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_DSN"); v != "" {
		return v
	}
	return "postgres://postgres:postgres@localhost:5439/asl_core_test?sslmode=disable"
}

// newRepo connects to the (migrated) test database and returns a Repo
// backed by an empty sources table, so each test starts from a known
// state. CI applies the migrations before running the tests.
func newRepo(t *testing.T) source.Repo {
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

	_, err = conn.Exec(t.Context(), "TRUNCATE TABLE sources")
	require.NoError(t, err)

	return source.New(source.Params{Conn: conn})
}

func sampleSource(title string) *source.Source {
	return &source.Source{
		Title:         title,
		Author:        "Muhammad al-Bukhari",
		Type:          source.TypeHadithCollection,
		Edition:       "Dar Tawq al-Najat, 2001",
		Language:      "ar",
		LocatorScheme: "collection + hadith number",
	}
}

func TestRepository_CreateGetList(t *testing.T) {
	repo := newRepo(t)
	ctx := t.Context()

	created, err := repo.Create(ctx, sampleSource("Sahih al-Bukhari"))
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, created.ID)
	require.Equal(t, "Sahih al-Bukhari", created.Title)
	require.Equal(t, source.TypeHadithCollection, created.Type)
	require.False(t, created.CreatedAt.IsZero())
	require.False(t, created.UpdatedAt.IsZero())

	got, err := repo.Get(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created, got)

	list, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, created, list[0])
}

func TestRepository_GetByID_notFound(t *testing.T) {
	repo := newRepo(t)

	_, err := repo.Get(t.Context(), uuid.New())
	appErr, ok := errorsx.As(err)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, appErr.StatusCode())
}

func TestRepository_List_orderedByCreatedAtDesc(t *testing.T) {
	repo := newRepo(t)
	ctx := t.Context()

	first, err := repo.Create(ctx, sampleSource("First"))
	require.NoError(t, err)

	// Separate the two inserts so their now() timestamps are distinct and
	// the descending order is deterministic.
	time.Sleep(2 * time.Millisecond)

	second, err := repo.Create(ctx, sampleSource("Second"))
	require.NoError(t, err)

	list, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, second.ID, list[0].ID, "newest source first")
	require.Equal(t, first.ID, list[1].ID)
}
