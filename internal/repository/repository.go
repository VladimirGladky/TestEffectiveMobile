package repository

import (
	"TestEffectiveMobile/internal/models"
	"TestEffectiveMobile/pkg/suberrors"
	"TestEffectiveMobile/pkg/timeparser"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type SubscriptionRepositoryInterface interface {
	Create(sub *models.Subscription) error
	Read(id string) (*models.Subscription, error)
	Update(id string, sub *models.UpdateSubscription) error
	Delete(id string) error
	ListSubscriptions(userId string) ([]*models.Subscription, error)
	CalculateSumSubscriptions(userId string, startDate string, endDate string, serviceName string) (int, error)
}

type SubscriptionRepository struct {
	db  *pgx.Conn
	ctx context.Context
}

func NewSubscriptionRepository(db *pgx.Conn, ctx context.Context) *SubscriptionRepository {
	return &SubscriptionRepository{
		db:  db,
		ctx: ctx,
	}
}

func (s *SubscriptionRepository) Create(sub *models.Subscription) error {
	stD, err := timeparser.ParseMonthYear(sub.StartDate)
	if err != nil {
		return fmt.Errorf("error creating subscription: %w", err)
	}
	endD, err := timeparser.ParseMonthYear(sub.EndDate)
	if err != nil {
		return fmt.Errorf("error creating subscription: %w", err)
	}
	_, err = s.db.Exec(s.ctx,
		"INSERT INTO subscriptions (id,service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)",
		sub.Id,
		sub.ServiceName,
		sub.Price,
		sub.UserId,
		stD,
		endD)
	if err != nil {
		return fmt.Errorf("error creating subscription: %w", err)
	}
	return err
}

func (s *SubscriptionRepository) Read(id string) (*models.Subscription, error) {
	var sub models.Subscription
	var startDate, endDate time.Time
	err := s.db.QueryRow(s.ctx,
		"SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1",
		id).Scan(
		&sub.Id,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserId,
		&startDate,
		&endDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, suberrors.ErrIdSubscriptionNotFound
		}
		return nil, fmt.Errorf("error reading subscription: %w", err)
	}
	sub.StartDate = startDate.Format("01-2006")
	sub.EndDate = endDate.Format("01-2006")
	sub.Id = id
	return &sub, nil
}

func (s *SubscriptionRepository) Update(id string, sub *models.UpdateSubscription) error {
	const query = `
        UPDATE subscriptions 
        SET 
            service_name = COALESCE($1, service_name),
            price = COALESCE($2, price),
            start_date = COALESCE($3, start_date),
            end_date = COALESCE($4, end_date)
        WHERE id = $5
        RETURNING id
    `
	stD, err := timeparser.ParseMonthYear(sub.StartDate)
	if err != nil {
		return fmt.Errorf("error updating subscription: %w", err)
	}
	endD, err := timeparser.ParseMonthYear(sub.EndDate)
	if err != nil {
		return fmt.Errorf("error updating subscription: %w", err)
	}

	var updatedID string
	err = s.db.QueryRow(s.ctx, query,
		sub.ServiceName,
		sub.Price,
		stD,
		endD,
		id,
	).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return suberrors.ErrIdSubscriptionNotFound
		}
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

func (s *SubscriptionRepository) Delete(id string) error {
	res, err := s.db.Exec(s.ctx,
		"DELETE FROM subscriptions WHERE id = $1",
		id)
	if err != nil {
		return fmt.Errorf("error deleting subscription: %w", err)
	}
	if res.RowsAffected() == 0 {
		return suberrors.ErrIdSubscriptionNotFound
	}
	return nil
}

func (s *SubscriptionRepository) ListSubscriptions(userId string) ([]*models.Subscription, error) {
	var subscriptions []*models.Subscription

	rows, err := s.db.Query(s.ctx,
		"SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE user_id = $1",
		userId)
	if err != nil {
		return nil, fmt.Errorf("error listing subscriptions: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var startDate, endDate time.Time
		var sub models.Subscription
		err := rows.Scan(&sub.Id,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserId,
			&startDate,
			&endDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning subscription: %w", err)
		}
		sub.StartDate = startDate.Format("01-2006")
		sub.EndDate = endDate.Format("01-2006")
		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	if len(subscriptions) == 0 {
		return nil, suberrors.ErrUserIdNotFound
	}
	return subscriptions, nil
}

func (s *SubscriptionRepository) CalculateSumSubscriptions(userId string, startDate string, endDate string, serviceName string) (int, error) {
	var sum int
	stD, err := timeparser.ParseMonthYear(startDate)
	if err != nil {
		return 0, fmt.Errorf("error calculating sum subscriptions: %w", err)
	}
	endD, err := timeparser.ParseMonthYear(endDate)
	if err != nil {
		return 0, fmt.Errorf("error calculating sum subscriptions: %w", err)
	}
	sql := "SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE start_date <= $1 AND end_date >= $2"
	args := []interface{}{
		endD.AddDate(0, 1, -1),
		stD}
	paramIndex := 3
	if serviceName != "" {
		sql += fmt.Sprintf(" AND service_name = $%d", paramIndex)
		args = append(args, serviceName)
		paramIndex++
	}
	if userId != "" {
		sql += fmt.Sprintf(" AND user_id = $%d", paramIndex)
		args = append(args, userId)
		exists, err := s.userExists(userId)
		if err != nil {
			return 0, fmt.Errorf("error checking user existence: %w", err)
		}
		if !exists {
			return 0, suberrors.ErrUserIdNotFound
		}
	}
	err = s.db.QueryRow(s.ctx, sql, args...).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("error calculating sum subscriptions: %w", err)
	}
	return sum, nil
}

func (s *SubscriptionRepository) userExists(userId string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(s.ctx,
		"SELECT EXISTS(SELECT 1 FROM subscriptions WHERE user_id = $1)",
		userId).Scan(&exists)
	return exists, err
}
