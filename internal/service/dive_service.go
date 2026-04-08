package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kiarrobino/divelog/internal/calculator"
	"github.com/kiarrobino/divelog/internal/model"
	"github.com/kiarrobino/divelog/internal/repository"
)

type DiveService struct {
	repo repository.DiveRepository
}

func NewDiveService(repo repository.DiveRepository) *DiveService {
	return &DiveService{repo: repo}
}

func (s *DiveService) CreateDive(ctx context.Context, input model.CreateDiveInput) (*model.Dive, error) {
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, model.ErrInvalidDate
	}
	if input.MaxDepth <= 0 {
		return nil, model.ErrInvalidDepth
	}
	if input.Duration <= 0 {
		return nil, model.ErrInvalidDuration
	}
	if input.Rating != 0 && (input.Rating < 1 || input.Rating > 5) {
		return nil, model.ErrInvalidRating
	}

	id := uuid.New().String()

	newDiveNum, err := s.repo.NextDiveNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("get next dive number: %w", err)
	}

	newDive := &model.Dive{
		ID:         id,
		DiveNumber: newDiveNum,
		Date:       date,
		SiteName:   input.SiteName,
		Location:   input.Location,
		MaxDepth:   input.MaxDepth,
		AvgDepth:   input.AvgDepth,
		Duration:   input.Duration,
		WaterTemp:  input.WaterTemp,
		Visibility: input.Visibility,
		TankStart:  input.TankStart,
		TankEnd:    input.TankEnd,
		O2Percent:  input.O2Percent,
		WaterType:  input.WaterType,
		DiveType:   input.DiveType,
		Notes:      input.Notes,
		Rating:     input.Rating,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.repo.Create(ctx, newDive)
	if err != nil {
		return nil, fmt.Errorf("create new dive: %w", err)
	}

	return newDive, nil
}

func (s *DiveService) GetDive(ctx context.Context, id string) (*model.Dive, error) {
	dive, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get dive by id: %w", err)
	}
	return dive, nil
}

func (s *DiveService) ListDives(ctx context.Context, limit, offset int) ([]*model.Dive, error) {
	dives, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list dives: %w", err)
	}
	return dives, nil
}

func (s *DiveService) CalculateNDL(depth, o2Percent float64) (int, error) {
	if depth <= 0 {
		return 0, model.ErrInvalidDepth
	}
	if o2Percent < 21 || o2Percent > 100 {
		return 0, model.ErrInvalidO2Percent
	}

	ndl := calculator.Calculate(depth, o2Percent)
	return ndl, nil
}
