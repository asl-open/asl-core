package source

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/database"
)

var Module = fx.Provide(New)

type Type string

const (
	TypeQuran            Type = "quran"
	TypeHadithCollection Type = "hadith_collection"
	TypeFiqhManual       Type = "fiqh_manual"
	TypeTafsir           Type = "tafsir"
	TypeArticle          Type = "article"
	TypeOther            Type = "other"
)

func (t Type) Valid() bool {
	switch t {
	case TypeQuran, TypeHadithCollection, TypeFiqhManual, TypeTafsir, TypeArticle, TypeOther:
		return true
	default:
		return false
	}
}

type Source struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Author        string
	Edition       string
	Language      string
	LocatorScheme string
	Type          Type
	ID            uuid.UUID
}

type Repo interface {
	Create(ctx context.Context, s *Source) (Source, error)
	Get(ctx context.Context, id uuid.UUID) (Source, error)
	List(ctx context.Context) ([]Source, error)
}

type Params struct {
	fx.In

	Conn database.Conn
}

type repo struct {
	conn database.Conn
}

func New(params Params) Repo {
	return &repo{conn: params.Conn}
}
