package redis

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserRepository struct {
	RD *redis.Client
	DB *sql.DB
}

type Destination struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Country           string    `json:"country"`
	Description       string    `json:"description"`
	BestTimeToVisit   string    `json:"best_time_to_visit"`
	AverageCostPerDay float64   `json:"average_cost_per_day"`
	Currency          string    `json:"currency"`
	Language          string    `json:"language"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

var logger *zap.Logger

func NewUserRepository(rd *redis.Client, db *sql.DB) *UserRepository {
	return &UserRepository{
		RD: rd,
		DB: db,
	}
}

func (repo *UserRepository) UpdateTopDestinations(ctx context.Context, destinations []Destination) error {
	destinationsJSON, err := json.Marshal(destinations)
	if err != nil {
		logger.Error("Failed to marshal destinations to JSON", zap.Error(err))
		return err
	}

	err = repo.RD.Set(ctx, "top_destinations", destinationsJSON, 1*time.Hour).Err()
	if err != nil {
		logger.Error("Failed to set top destinations in Redis", zap.Error(err))
		return err
	}

	return nil
}

func (repo *UserRepository) GetTopDestinations(ctx context.Context) ([]Destination, error) {
	destinationsJSON, err := repo.RD.Get(ctx, "top_destinations").Bytes()
	if err != nil {
		if err == redis.Nil {
			logger.Info("Top destinations not found in Redis, updating from database")
			return repo.updateAndGetTopDestinations(ctx)
		}
		logger.Error("Failed to fetch top destinations from Redis", zap.Error(err))
		return nil, err
	}

	var destinations []Destination
	err = json.Unmarshal(destinationsJSON, &destinations)
	if err != nil {
		logger.Error("Failed to unmarshal destinations JSON", zap.Error(err))
		return nil, err
	}

	return destinations, nil
}

func (repo *UserRepository) updateAndGetTopDestinations(ctx context.Context) ([]Destination, error) {
	// Implement your logic to fetch top destinations from the database
	destinations := []Destination{} // Replace with your actual database query

	// Update destinations in Redis
	err := repo.UpdateTopDestinations(ctx, destinations)
	if err != nil {
		logger.Error("Failed to update top destinations in Redis", zap.Error(err))
		return nil, err
	}

	return destinations, nil
}
