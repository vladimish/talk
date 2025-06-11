package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"talk/db/generated"
	"talk/internal/domain"
	"talk/internal/port/storage"
)

type PG struct {
	q *generated.Queries
}

func NewPg(q *generated.Queries) *PG {
	return &PG{q: q}
}

func (p *PG) GetUserByExternalUserID(ctx context.Context, id string) (*domain.User, error) {
	iid, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("can't convert id to int: %w", err)
	}

	u, err := p.q.GetUserByForeignID(ctx, int64(iid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("can't get user by foreign id: %w", err)
	}

	return &domain.User{
		ID:         strconv.FormatInt(u.ID, 10),
		ExternalID: strconv.FormatInt(u.ForeignID, 10),
		Language:   u.Language,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}, nil
}
