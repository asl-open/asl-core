package contributor

import (
	"context"
	"strings"

	"github.com/asl-open/asl-core/services/api/internal/http/apierrors"
	contributorrepo "github.com/asl-open/asl-core/services/api/internal/repository/contributor"
)

type CreateInput struct {
	DisplayName string
	Bio         *string
	Handle      *string
	Roles       []contributorrepo.Role
}

func (s *service) Create(ctx context.Context, in *CreateInput) (contributorrepo.Contributor, error) {
	if err := in.validate(); err != nil {
		return contributorrepo.Contributor{}, err
	}

	return s.repo.Create(ctx, &contributorrepo.Contributor{
		DisplayName: strings.TrimSpace(in.DisplayName),
		Roles:       in.Roles,
		Bio:         in.Bio,
		Handle:      in.Handle,
	})
}

func (in *CreateInput) validate() error {
	if strings.TrimSpace(in.DisplayName) == "" {
		return apierrors.ErrBadRequest.WithMessage("contributor display name is required")
	}
	if len(in.Roles) == 0 {
		return apierrors.ErrBadRequest.WithMessage("contributor requires at least one role")
	}
	for _, role := range in.Roles {
		if !role.Valid() {
			return apierrors.ErrBadRequest.WithMessage("contributor role is invalid")
		}
	}

	return nil
}
