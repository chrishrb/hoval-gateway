package config

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

type DatapointProvider interface {
	GetByUnit(unit string) []Datapoint
	GetByFunction(fg, fn uint8, dpID uint16) *Datapoint
	GetByName(name string) *Datapoint
}

// Datapoint represents the configuration for a single datapoint in the CSV
type Datapoint struct {
	RegisterAddress   uint16 `csv:"Register Address"`
	UnitName          string `csv:"UnitName"`
	UnitID            string `csv:"UnitId"`
	FunctionGroup     uint8  `csv:"FunctionGroup"`
	FunctionNumber    uint8  `csv:"FunctionNumber"`
	DatapointID       uint16 `csv:"DatapointId"`
	DatapointName     string `csv:"DatapointName"`
	Type              string `csv:"Type"`
	TypeName          string `csv:"TypeName"`
	Decimal           uint8  `csv:"Decimal"`
	FunctionGroupName string `csv:"FunctionGroup name"`
	FunctionName      string `csv:"Function name"`
}

// CsvDatapointProvider holds the different indexed maps for datapoint configuration
type CsvDatapointProvider struct {
	ByUnitID   map[string][]Datapoint
	ByFunction map[string]Datapoint // key: "FunctionGroup:FunctionNumber:DatapointID"
	ByName     map[string]Datapoint
}

func (d *CsvDatapointProvider) GetByUnit(unit string) []Datapoint {
	dp, exists := d.ByUnitID[unit]
	if !exists {
		return nil
	}
	return dp
}

func (d *CsvDatapointProvider) GetByFunction(fg, fn uint8, dpID uint16) *Datapoint {
	key := fmt.Sprintf("%d:%d:%d", fg, fn, dpID)
	dp, exists := d.ByFunction[key]
	if !exists {
		return nil
	}
	return &dp
}

func (d *CsvDatapointProvider) GetByName(name string) *Datapoint {
	dp, exists := d.ByName[name]
	if !exists {
		return nil
	}
	return &dp
}

// NewCsvDatapointProvider reads a CSV file and returns datapoint configurations
// indexed by different keys for efficient lookup
func NewCsvDatapointProvider(filename string) (*CsvDatapointProvider, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	var datapoints []Datapoint
	if err := gocsv.UnmarshalFile(file, &datapoints); err != nil {
		return nil, fmt.Errorf("failed to parse CSV file: %w", err)
	}

	maps := &CsvDatapointProvider{
		ByUnitID:   make(map[string][]Datapoint),
		ByFunction: make(map[string]Datapoint),
		ByName:     make(map[string]Datapoint),
	}

	for _, dp := range datapoints {
		// Index by UnitID (one-to-many relationship)
		maps.ByUnitID[dp.UnitID] = append(maps.ByUnitID[dp.UnitID], dp)

		// Index by FunctionGroup:FunctionNumber:DatapointID
		functionKey := fmt.Sprintf("%d:%d:%d", dp.FunctionGroup, dp.FunctionNumber, dp.DatapointID)
		maps.ByFunction[functionKey] = dp

		// Index by DatapointName
		maps.ByName[dp.DatapointName] = dp
	}

	return maps, nil
}
