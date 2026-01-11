package services

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
)

type ContributorService struct {
	contribRepo interfaces.ContributorRepository
}

func NewContributorService(contribRepo interfaces.ContributorRepository) *ContributorService {
	return &ContributorService{
		contribRepo: contribRepo,
	}
}

func (s *ContributorService) CreateContributor(ctx context.Context, contributor models.CreateContributorRequest) error {
	if len(contributor.Fullname) < 4 {
		return errors.New("fullname must be at least 4 characters long")
	}
	if len(contributor.Email) < 5 {
		return errors.New("email must be at least 5 characters long")
	}
	_, err := mail.ParseAddress(strings.TrimSpace(contributor.Email))
	if err != nil {
		return errors.New("invalid email format")
	}
	err = s.contribRepo.CreateContributor(ctx, contributor)
	if err != nil {
		return err
	}
	return nil
}

func (s *ContributorService) GetContributorsByUserID(ctx context.Context, userID string) ([]models.ContributorResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	contributors, err := s.contribRepo.GetContributorsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return contributors, nil
}

func (s *ContributorService) GetContributorByID(ctx context.Context, contributorID string) (*models.ContributorResponse, error) {
	if contributorID == "" {
		return nil, errors.New("contributor ID is required")
	}
	contributor, err := s.contribRepo.GetContributorByID(ctx, contributorID)
	if err != nil {
		return nil, err
	}
	return contributor, nil
}
