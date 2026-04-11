package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kiarrobino/divelog/internal/model"
)

type mockRepo struct {
	dives []*model.Dive
}

func (m *mockRepo) Create(ctx context.Context, dive *model.Dive) error {
	m.dives = append(m.dives, dive)
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*model.Dive, error) {
	return nil, model.ErrDiveNotFound
}

func (m *mockRepo) List(ctx context.Context, limit, offset int) ([]*model.Dive, error) {
	return m.dives, nil
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockRepo) NextDiveNumber(ctx context.Context) (int, error) {
	return len(m.dives) + 1, nil
}

func TestCreateDive(t *testing.T) {
	s := NewDiveService(&mockRepo{})

	tests := []struct {
		name    string
		input   model.CreateDiveInput
		wantErr error
	}{
		{
			name: "valid input",
			input: model.CreateDiveInput{
				Date:       "2006-01-01",
				SiteName:   "test",
				Location:   "test",
				MaxDepth:   1,
				AvgDepth:   1,
				Duration:   1,
				WaterTemp:  1,
				Visibility: 1,
				TankStart:  1,
				TankEnd:    1,
				O2Percent:  21,
				WaterType:  "test",
				DiveType:   "test",
				Notes:      "",
				Rating:     1,
			},
			wantErr: nil,
		},
		{
			name: "invalid input: date",
			input: model.CreateDiveInput{
				Date:       time.Now().String(),
				SiteName:   "test",
				Location:   "test",
				MaxDepth:   1,
				AvgDepth:   1,
				Duration:   1,
				WaterTemp:  1,
				Visibility: 1,
				TankStart:  1,
				TankEnd:    1,
				O2Percent:  21,
				WaterType:  "test",
				DiveType:   "test",
				Notes:      "",
				Rating:     1,
			}, wantErr: model.ErrInvalidDate,
		},
		{
			name: "invalid input: max depth",
			input: model.CreateDiveInput{
				Date:       "2006-01-01",
				SiteName:   "test",
				Location:   "test",
				MaxDepth:   0,
				AvgDepth:   1,
				Duration:   1,
				WaterTemp:  1,
				Visibility: 1,
				TankStart:  1,
				TankEnd:    1,
				O2Percent:  21,
				WaterType:  "test",
				DiveType:   "test",
				Notes:      "",
				Rating:     1,
			}, wantErr: model.ErrInvalidDepth,
		},
		{
			name: "invalid input: duration",
			input: model.CreateDiveInput{
				Date:       "2006-01-01",
				SiteName:   "test",
				Location:   "test",
				MaxDepth:   1,
				AvgDepth:   1,
				Duration:   0,
				WaterTemp:  1,
				Visibility: 1,
				TankStart:  1,
				TankEnd:    1,
				O2Percent:  21,
				WaterType:  "test",
				DiveType:   "test",
				Notes:      "",
				Rating:     1,
			}, wantErr: model.ErrInvalidDuration,
		},
		{
			name: "invalid input: o2 percent",
			input: model.CreateDiveInput{
				Date:       "2006-01-01",
				SiteName:   "test",
				Location:   "test",
				MaxDepth:   1,
				AvgDepth:   1,
				Duration:   1,
				WaterTemp:  1,
				Visibility: 1,
				TankStart:  1,
				TankEnd:    1,
				O2Percent:  1,
				WaterType:  "test",
				DiveType:   "test",
				Notes:      "",
				Rating:     1,
			}, wantErr: model.ErrInvalidO2Percent,
		},
		{
			name: "invalid input: rating",
			input: model.CreateDiveInput{
				Date:       "2006-01-01",
				SiteName:   "test",
				Location:   "test",
				MaxDepth:   1,
				AvgDepth:   1,
				Duration:   1,
				WaterTemp:  1,
				Visibility: 1,
				TankStart:  1,
				TankEnd:    1,
				O2Percent:  21,
				WaterType:  "test",
				DiveType:   "test",
				Notes:      "",
				Rating:     10,
			}, wantErr: model.ErrInvalidRating,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreateDive(context.Background(), tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Error("expected a dive back, got nil")
			}
		})
	}
}

func TestGetDive(t *testing.T) {
	s := NewDiveService(&mockRepo{})

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "invalid id",
			id:      "test",
			wantErr: model.ErrDiveNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetDive(context.Background(), tt.id)
			if !errors.Is(err, model.ErrDiveNotFound) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestListDives(t *testing.T) {
	s := NewDiveService(&mockRepo{
		dives: []*model.Dive{
			{ID: "1", SiteName: "test1", Location: "test1"},
			{ID: "2", SiteName: "test2", Location: "test2"},
			{ID: "3", SiteName: "test3", Location: "test3"},
		},
	})

	tests := []struct {
		name          string
		limit, offset int
		wantCount     int
	}{
		{name: "get all", limit: 10, offset: 0, wantCount: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListDives(context.Background(), tt.limit, tt.offset)
			if err != nil {
				t.Errorf("got error %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("got %d dives, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestCalculateNDL(t *testing.T) {
	s := NewDiveService(&mockRepo{})

	tests := []struct {
		name    string
		depth   float64
		want    int
		wantErr error
	}{
		{"10m", 10, 219, nil},
		{"12m", 12, 147, nil},
		{"14m", 14, 98, nil},
		{"16m", 16, 72, nil},
		{"18m", 18, 56, nil},
		{"20m", 20, 45, nil},
		{"25m", 25, 29, nil},
		{"30m", 30, 20, nil},
		{"35m", 35, 14, nil},
		{"40m", 40, 9, nil},
		{"45m", 45, 0, nil},
		{"32m", 32, 14, nil},
		{"0m", 0, 0, model.ErrInvalidDepth},
		{"-5m", -5, 0, model.ErrInvalidDepth},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CalculateNDL(tt.depth)
			if err != nil && tt.wantErr == nil {
				t.Errorf("got error %v", err)
			}
			if err != nil && !errors.Is(err, model.ErrInvalidDepth) {
				t.Errorf("got error %v, want error %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("got %d mins, want %d mins", got, tt.want)
			}
		})
	}
}
