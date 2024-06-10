package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/port/repository"
	"fmt"
)

type reviewRepository struct {
	db repository.Database
}

func NewReviewRepository(db repository.Database) repository.ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}

func (r reviewRepository) InsertReview(review dto.Review) (*int, error) {
	DB := r.db.GetDB()
	userId := review.UserID
	productId := review.ProductID
	rating := review.Rating
	comment := review.Comment

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM review WHERE product_id = ? AND user_id = ?", productId, userId).Scan(&count); err != nil {
		return nil, err
	}

	if count > 0 {
		err := fmt.Errorf("review with user_id '%d' and product_id '%d' already exists", userId, productId)
		return nil, err
	}

	query := `INSERT INTO review(user_id, product_id, rating, comment) VALUES (?, ?, ?, ?)`
	result, err := DB.Exec(query, userId, productId, rating, comment)
	if err != nil {
		return nil, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	id := int(id64)

	return &id, nil
}

func (r reviewRepository) ListReviews() ([]utils.Review, error) {
	DB := r.db.GetDB()
	var reviews []utils.Review

	query := `SELECT review_id, user_id, product_id, rating, comment, created_at, updated_at, deleted_at FROM review WHERE deleted_at IS NULL`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var review utils.Review

		if err := rows.Scan(&review.ID, &review.UserID, &review.ProductID, &review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt, &review.DeletedAt); err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}
