package repository

import (
	"Eccomerce-website/internals/core/common/utils"
	"Eccomerce-website/internals/core/dto"
	"Eccomerce-website/internals/core/entity"
	"Eccomerce-website/internals/core/port/repository"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type reviewRepository struct {
	db       repository.Database
	dbLogger *zap.Logger
}

func NewReviewRepository(db repository.Database, dbLogger *zap.Logger) repository.ReviewRepository {
	return &reviewRepository{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (r reviewRepository) InsertReview(ctx context.Context, review dto.Review, requestID string) (*int, error) {
	DB := r.db.GetDB()
	userId := review.UserID
	productId := review.ProductID
	rating := review.Rating
	comment := review.Comment

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM review WHERE product_id = ? AND user_id = ?", productId, userId).Scan(&count); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "unable to read COUNT in the query").WithProperty(entity.StatusCode, 500)
		r.dbLogger.Error("failed to read 'COUNT' in the query",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertReview"),
			zap.String("requestID", requestID),
			zap.String("query", "SELECT COUNT(*) FROM review WHERE product_id = ? AND user_id = ?"),
			zap.Int("productID", productId),
			zap.Int("userID", userId),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	if count > 0 {
		err := fmt.Errorf("review with user_id '%d' and product_id '%d' already exists", userId, productId)
		errorResponse := entity.DuplicateEntry.Wrap(err, "conflict error").WithProperty(entity.StatusCode, 409)
		r.dbLogger.Error("duplicate entry",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertReview"),
			zap.String("requestID", requestID),
			zap.Any("requestData", review),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	query := `INSERT INTO review(user_id, product_id, rating, comment) VALUES (?, ?, ?, ?)`
	result, err := DB.ExecContext(ctx, query, userId, productId, rating, comment)
	if err != nil {
		errorResponse := entity.UnableToSave.Wrap(err, "failed to insert product review").WithProperty(entity.StatusCode, 500)
		r.dbLogger.Error("failed to create product review",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertReview"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Any("requestData", review),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	id64, err := result.LastInsertId()
	if err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "failed to get the inserted id").WithProperty(entity.StatusCode, 500)
		r.dbLogger.Error("unable to get lastInserted id",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "InsertReview"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	id := int(id64)

	return &id, nil
}

func (r reviewRepository) ListReviews(ctx context.Context, offset, limit int, requestID string) ([]utils.Review, error) {
	DB := r.db.GetDB()
	var reviews []utils.Review

	query := `SELECT review_id, user_id, product_id, rating, comment, created_at, updated_at, deleted_at FROM review WHERE deleted_at IS NULL ORDER BY review_id LIMIT ? OFFSET ?`
	rows, err := DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		errorResponse := entity.UnableToFindResource.Wrap(err, "failed to get list of reviews").WithProperty(entity.StatusCode, 404)
		r.dbLogger.Error("reviews not found",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListReviews"),
			zap.String("requestID", requestID),
			zap.String("query", query),
			zap.Int("offset", offset),
			zap.Int("limit", limit),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)

		return nil, errorResponse
	}

	defer rows.Close()

	for rows.Next() {
		var review utils.Review

		if err := rows.Scan(&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt, &review.DeletedAt); err != nil {
			errorResponse := entity.UnableToRead.Wrap(err, "failed to scan review data").WithProperty(entity.StatusCode, 500)
			r.dbLogger.Error("unable to scan the review data",
				zap.String("timestamp", time.Now().Format(time.RFC3339)),
				zap.String("layer", "databaseLayer"),
				zap.String("function", "ListReviews"),
				zap.String("requestID", requestID),
				zap.Error(errorResponse),
				zap.Stack("stacktrace"),
			)

			return nil, errorResponse
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		errorResponse := entity.UnableToRead.Wrap(err, "db rows error").WithProperty(entity.StatusCode, 500)
		r.dbLogger.Error("db rows error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "databaseLayer"),
			zap.String("function", "ListReviews"),
			zap.String("requestID", requestID),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		return nil, errorResponse
	}

	return reviews, nil
}
