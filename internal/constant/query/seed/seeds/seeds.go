package seeds

import (
	"Eccomerce-website/internal/constant"
	"Eccomerce-website/internal/constant/model/dto"

	"Eccomerce-website/internal/constant/query/seed/seed"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func All() []seed.Seed {
	userID := uuid.New()
	productCategoryID := uuid.New()
	productID := uuid.New()
	sizeID := uuid.New()
	colorID := uuid.New()
	productItemID := uuid.New()

	return []seed.Seed{
		{
			// Seed user
			Name: "Create User",
			Run: func(c *pgxpool.Pool) error {
				return CreateUser(c, &dto.User{
					ID:             userID,
					Username:       "eagle",
					Email:          "eagle@gmail.com",
					PhoneNumber:    "0912241771",
					Password:       "Eagle@1224",
					FirstName:      "Ady",
					LastName:       "Yismaw",
					ProfilePicture: "eagle.jpg",
					EmailVerified:  true,
					Role:           constant.UserRoleCustomer,
				})

			},
		},
		{
			// Seed product category
			Name: "Create Product Category",
			Run: func(c *pgxpool.Pool) error {
				return CreateProductCategory(c, &dto.ProductCategory{
					ID:          productCategoryID,
					Name:        "Electronics",
					Description: "Eagle Ecommerce Electronics Category",
					Image:       "electronics.jpg",
				})

			},
		},

		{
			// Seed product
			Name: "Create Product",
			Run: func(c *pgxpool.Pool) error {
				return CreateProduct(c, &dto.Product{
					ID:            productID,
					CategoryID:    productCategoryID,
					Brand:         "Samsung",
					Name:          "Samsung Galaxy S21",
					Status:        constant.ProductStatusActive,
					Description:   "Samsung Galaxy S21 5G | Factory Unlocked Android Cell Phone | US Version 5G Smartphone | Pro-Grade Camera, 8K Video, 64MP High Res | 128GB, Phantom Gray (SM-G991UZAAXAA)",
					AverageRating: 4.5,
					TotalReviews:  100,
				})

			},
		},

		{
			// Seed size
			Name: "Create Size",
			Run: func(c *pgxpool.Pool) error {
				return CreateSize(c, &dto.Size{
					ID:   sizeID,
					Name: "S21",
				})

			},
		},

		{
			// Seed Color
			Name: "Create Color",
			Run: func(c *pgxpool.Pool) error {
				return CreateColor(c, &dto.Color{
					ID:   colorID,
					Name: "BLACK",
				})

			},
		},

		{
			// Seed product item
			Name: "Create Product Item",
			Run: func(c *pgxpool.Pool) error {
				return CreateProductItem(c, &dto.ProductItem{
					ID:        productItemID,
					ProductID: productID,
					ColorID:   colorID,
					SizeID:    sizeID,
					SKU:       "SUMSUNG-S21-BLACK",
					Status:    constant.ProductStatusActive,
					ImageURL:  "sumsung-galaxys21.png",
					Price:     10000,
					Discount:  120000,
				})

			},
		},
	}
}
