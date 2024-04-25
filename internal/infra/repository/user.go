package repository

import (
	"Eccomerce-website/internal/core/common/utils"
	"Eccomerce-website/internal/core/dto"
	"Eccomerce-website/internal/core/port/repository"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// import (
// 	"Eccomerce-website/internal/core/dto"
// 	"Eccomerce-website/internal/core/port/repository"
// 	dbmodels "Eccomerce-website/internal/infra/db_models"
// 	"errors"
// 	"fmt"

//	"golang.org/x/crypto/bcrypt"
//	"gorm.io/gorm"
//
// )

type userRepository struct {
	db repository.Database
}

func NewUserRepository(db repository.Database) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func matchPassword(hashPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	isMatch := err == nil
	return isMatch, err
}

func (u userRepository) InsertUser(user dto.User) (int, error) {
	DB := u.db.GetDB()

	var count int
	if err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", user.Username, user.Email).Scan(&count); err != nil {
		return 0, err
	}
	if count > 0 {
		err := errors.New("user already exists")
		return 0, err
	}

	query := `INSERT INTO users(username, email, password, first_name, last_name, phone_number, address, role, email_verified, profile_picture)
	          VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(query, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.PhoneNumber, user.Address, user.Role, user.EmailVerified, user.ProfilePicture)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// result := DB.Create(user)
	// if result.Error != nil {
	// 	return result.Error
	// }

	return int(id), nil
}

func (u userRepository) Authentication(email string, password string) (utils.User, error) {
	var user utils.User
	DB := u.db.GetDB()

	query := `
	 SELECT id, username, email, password, first_name,
	 last_name, phone_number, address, profile_picture, email_verified, 
	 role, created_at, updated_at FROM users WHERE email = ?
	 `

	if err := DB.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &user.PhoneNumber, &user.Address, &user.ProfilePicture,
		&user.EmailVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
		err = fmt.Errorf("user with email %s is not found: %s", email, err)
		return utils.User{}, err
	}
	// if err := DB.Where("email=?", email).First(&user).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		err = fmt.Errorf("user with email %s is not found", email)
	// 		return dbmodels.User{}, err
	// 	}

	// 	return dbmodels.User{}, err
	// }

	match, err := matchPassword(user.Password, password)
	if err != nil || !match {
		err = errors.New("invalid password")
		return utils.User{}, err
	}

	return user, nil

}

func (u userRepository) ListUsers() ([]utils.User, error) {
	var users []utils.User
	DB := u.db.GetDB()

	query := `SELECT id, username, email, password, first_name,
	last_name, phone_number, address, profile_picture, email_verified, 
	role, created_at, updated_at FROM users`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user utils.User
		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName,
			&user.LastName, &user.PhoneNumber, &user.Address, &user.ProfilePicture,
			&user.EmailVerified, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	// result := DB.Find(&users)
	// if result.Error != nil {
	// 	return []dbmodels.User{}, result.Error
	// }

	return users, nil
}
