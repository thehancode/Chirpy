package database

import (
	"context"

	"github.com/google/uuid"
)

func (q *Queries) GetAllChirps(ctx context.Context, sortdirection string) ([]Chirp, error) {
	var chirps []Chirp
	var err error
	if sortdirection == "desc" {
		chirps, err = q.GetAllChirpsDesc(ctx)
	} else {
		chirps, err = q.GetAllChirpsAsc(ctx)
	}

	if err != nil {
		return nil, err
	}
	return chirps, nil
}

func (q *Queries) GetChirpsByAuthorID(ctx context.Context, userID uuid.UUID, sortdirection string) ([]Chirp, error) {
	var chirps []Chirp
	var err error
	if sortdirection == "desc" {
		chirps, err = q.GetChirpsByAuthorIDDesc(ctx, userID)
	} else {
		chirps, err = q.GetChirpsByAuthorIDAsc(ctx, userID)
	}

	if err != nil {
		return nil, err
	}
	return chirps, nil
}
