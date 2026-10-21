package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *Store) SaveVote(ctx context.Context, memberID, movieID int64, score int) error {
	if score < 1 || score > 10 {
		return errors.New("score must be between 1 and 10")
	}

	_, err := s.Pool.Exec(ctx,
		`INSERT INTO votes (member_id, movie_id, score)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (member_id, movie_id)
		 DO UPDATE SET score = EXCLUDED.score, updated_at = now()`,
		memberID, movieID, score,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) RemoveVote(ctx context.Context, memberID, movieID int64) error {
	tag, err := s.Pool.Exec(ctx,
		`DELETE FROM votes
		 WHERE member_id = $1 AND movie_id = $2`,
		memberID, movieID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return nil
	}
	return nil
}

func (s *Store) GetVote(ctx context.Context, memberID, movieID int64) (int, error) {
	var score int
	err := s.Pool.QueryRow(ctx,
		`SELECT score FROM votes
		 WHERE member_id = $1 AND movie_id = $2`,
		memberID, movieID,
	).Scan(&score)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNotFound
	}
	return score, err
}
