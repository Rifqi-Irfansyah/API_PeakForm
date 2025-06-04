package repository

import (
	"api-peak-form/domain"
	"api-peak-form/dto"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	logrus.Info("Creating new UserRepository")
	return &userRepository{db: db}
}

func (u userRepository) FindByEmail(ctx context.Context, email string) (user domain.User, err error) {
	err = u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("User not found for email: %s", email)
			return user, nil
		}
		logrus.Errorf("Error querying user: %v", err)
		return user, err
	}
	logrus.Infof("User found: %s", user.Email)
	return user, nil
}

func (u userRepository) FindByID(ctx context.Context, id string) (user domain.User, err error) {
	err = u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warnf("User not found for id: %s", id)
			return user, err
		}
		logrus.Errorf("Error querying user: %v", err)
		return user, err
	}
	logrus.Infof("User found: %s", user.Email)
	return user, nil
}

func (u userRepository) Save(ctx context.Context, user domain.User) error {
	err := u.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		logrus.Errorf("Error while saving user: %v", err)
		return err
	}
	logrus.Infof("User successfully saved: %s", user.Email)
	return nil
}

func (u userRepository) UpdatePassword(ctx context.Context, email string, password string) error {
	err := u.db.WithContext(ctx).Model(&domain.User{}).Where("email = ?", email).Update("password", password).Error
	if err != nil {
		logrus.Errorf("Error updating password for email: %s, error: %v", email, err)
		return err
	}
	logrus.Infof("Password successfully updated for email: %s", email)
	return nil
}

func (u userRepository) UpdatePoint(ctx context.Context, id string, point int) error {
	err := u.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("point", point).Error
	if err != nil {
		logrus.Errorf("Error updating point for user ID: %s, error: %v", id, err)
		return err
	}
	logrus.Infof("Point successfully updated for user ID: %s", id)
	return nil
}

func (u userRepository) UpdateStreak(ctx context.Context, id string, streak int) error {
	err := u.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("streak", streak).Error
	if err != nil {
		logrus.Errorf("Error updating streak for user ID: %s, error: %v", id, err)
		return err
	}
	logrus.Infof("Streak successfully updated for user ID: %s", id)
	return nil
}

func (u userRepository) GetAllUsersDesc(ctx context.Context) ([]dto.UserLeaderboardResponse, error) {
	var users []domain.User
	err := u.db.WithContext(ctx).Order("point DESC").Find(&users).Error
	if err != nil {
		logrus.Errorf("Error retrieving all users: %v", err)
		return nil, err
	}
	logrus.Infof("Retrieved %d users", len(users))

	var userResponses []dto.UserLeaderboardResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserLeaderboardResponse{
			Name:  user.Name,
			Point: user.Point,
		})
	}

	return userResponses, nil
}

func (u *userRepository) UpdatePhoto(ctx context.Context, id string, photoURL string) error {
	return u.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("photo_url", photoURL).Error
}