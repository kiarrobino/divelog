package exporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/kiarrobino/divelog/internal/model"
)

func WriteCSV(w io.Writer, dives []*model.Dive) error {
	writer := csv.NewWriter(w)

	err := writer.Write([]string{
		"id", "dive_number", "date", "site_name", "location",
		"max_depth_m", "avg_depth_m", "duration_min", "water_temp_c",
		"visibility_m", "tank_start_bar", "tank_end_bar", "o2_percent",
		"water_type", "dive_type", "notes", "rating",
		"created_at", "updated_at",
	})
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, dive := range dives {
		err := writer.Write([]string{
			dive.ID,
			strconv.Itoa(dive.DiveNumber),
			dive.Date.Format("2006-01-02"),
			dive.SiteName,
			dive.Location,
			strconv.FormatFloat(dive.MaxDepth, 'f', 2, 64),
			strconv.FormatFloat(dive.AvgDepth, 'f', 2, 64),
			strconv.Itoa(dive.Duration),
			strconv.FormatFloat(dive.WaterTemp, 'f', 2, 64),
			strconv.FormatFloat(dive.Visibility, 'f', 2, 64),
			strconv.Itoa(dive.TankStart),
			strconv.Itoa(dive.TankEnd),
			strconv.FormatFloat(dive.O2Percent, 'f', 2, 64),
			dive.WaterType,
			dive.DiveType,
			dive.Notes,
			strconv.Itoa(dive.Rating),
			dive.CreatedAt.Format("2006-01-02"),
			dive.UpdatedAt.Format("2006-01-02"),
		})
		if err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush: %w", err)
	}

	return nil
}
