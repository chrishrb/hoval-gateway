package datapoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval/datapoint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDatapointConfigFromCSV(t *testing.T) {
	t.Run("successful load with valid CSV", func(t *testing.T) {
		// Create a temporary CSV file for testing
		testCSV := createTestCSVFile(t)
		//nolint:errcheck
		defer os.Remove(testCSV)

		maps, err := datapoint.NewCsvDatapointProvider(testCSV)
		require.NoError(t, err)
		require.NotNil(t, maps)

		// Test ByUnitID mapping
		assert.Contains(t, maps.ByUnitID, "65")
		assert.Len(t, maps.ByUnitID["65"], 3)
		assert.Contains(t, maps.ByUnitID, "66")
		assert.Len(t, maps.ByUnitID["66"], 1)

		// Test ByFunction mapping
		assert.Contains(t, maps.ByFunction, "22:0:14")
		assert.Contains(t, maps.ByFunction, "22:1:14")
		assert.Contains(t, maps.ByFunction, "23:2:15")

		// Test ByName mapping
		assert.Contains(t, maps.ByName, "TKO1 Kollektor Temperatur")
		assert.Contains(t, maps.ByName, "TKO2 Kollektor Temperatur")
		assert.Contains(t, maps.ByName, "Test Datapoint")

		// Verify specific datapoint values
		dp1 := maps.ByFunction["22:0:14"]
		assert.Equal(t, uint16(1), dp1.RegisterAddress)
		assert.Equal(t, "SOL", dp1.UnitName)
		assert.Equal(t, "65", dp1.UnitID)
		assert.Equal(t, uint8(22), dp1.FunctionGroup)
		assert.Equal(t, uint8(0), dp1.FunctionNumber)
		assert.Equal(t, uint16(14), dp1.DatapointID)
		assert.Equal(t, "TKO1 Kollektor Temperatur", dp1.DatapointName)
		assert.Equal(t, "1", dp1.Type)
		assert.Equal(t, "S16", dp1.TypeName)
		assert.Equal(t, uint8(1), dp1.Decimal)
		assert.Equal(t, "Kollektor", dp1.FunctionGroupName)
		assert.Equal(t, "Kollektor 1", dp1.FunctionName)
	})

	t.Run("file not found error", func(t *testing.T) {
		maps, err := datapoint.NewCsvDatapointProvider("nonexistent.csv")
		assert.Error(t, err)
		assert.Nil(t, maps)
		assert.Contains(t, err.Error(), "failed to open CSV file")
	})

	t.Run("invalid CSV format", func(t *testing.T) {
		// Create a temporary file with invalid CSV content
		testFile := filepath.Join(t.TempDir(), "invalid.csv")
		err := os.WriteFile(testFile, []byte("invalid,csv\ndata,with,too,many,fields"), 0644)
		require.NoError(t, err)

		maps, err := datapoint.NewCsvDatapointProvider(testFile)
		assert.Error(t, err)
		assert.Nil(t, maps)
		assert.Contains(t, err.Error(), "failed to parse CSV file")
	})

	t.Run("empty CSV file", func(t *testing.T) {
		// Create a temporary empty CSV file
		testFile := filepath.Join(t.TempDir(), "empty.csv")
		err := os.WriteFile(testFile, []byte("Register Address,UnitName,UnitId,FunctionGroup,FunctionNumber,DatapointId,DatapointName,Type,TypeName,Decimal,FunctionGroup name,Function name,Steps,Min. value,Max. value,Writable,unit,Commentary,Text 0,Text 1,Text 2,Text 3,Text 4,Text 5,Text 6,Text 7,Text 8,Text 9,Text 10,Text 11,Text 12,Text 13,Text 14,Text 15,Text 16,Text 17,Text 18,Text 19,Text 20,Text 21,Text 22,Text 23,Text 24,Text 25,Text 26,Text 27,Text 28,Text 29,Text 30,Text 31\n"), 0644)
		require.NoError(t, err)

		maps, err := datapoint.NewCsvDatapointProvider(testFile)
		require.NoError(t, err)
		require.NotNil(t, maps)

		// All maps should be empty
		assert.Empty(t, maps.ByUnitID)
		assert.Empty(t, maps.ByFunction)
		assert.Empty(t, maps.ByName)
	})
}

func TestDatapointMapsIndexing(t *testing.T) {
	testCSV := createTestCSVFile(t)
	//nolint:errcheck
	defer os.Remove(testCSV)

	maps, err := datapoint.NewCsvDatapointProvider(testCSV)
	require.NoError(t, err)

	t.Run("ByUnitID indexing", func(t *testing.T) {
		// Test that UnitID "65" has multiple datapoints
		datapoints := maps.ByUnitID["65"]
		assert.Len(t, datapoints, 3)

		// Verify all datapoints have the correct UnitID
		for _, dp := range datapoints {
			assert.Equal(t, "65", dp.UnitID)
		}
	})

	t.Run("ByFunction indexing", func(t *testing.T) {
		// Test function key format: "FunctionGroup:FunctionNumber:DatapointID"
		dp, exists := maps.ByFunction["22:0:14"]
		assert.True(t, exists)
		assert.Equal(t, uint8(22), dp.FunctionGroup)
		assert.Equal(t, uint8(0), dp.FunctionNumber)
		assert.Equal(t, uint16(14), dp.DatapointID)

		dp2, exists := maps.ByFunction["23:2:15"]
		assert.True(t, exists)
		assert.Equal(t, uint8(23), dp2.FunctionGroup)
		assert.Equal(t, uint8(2), dp2.FunctionNumber)
		assert.Equal(t, uint16(15), dp2.DatapointID)
	})

	t.Run("ByName indexing", func(t *testing.T) {
		dp, exists := maps.ByName["TKO1 Kollektor Temperatur"]
		assert.True(t, exists)
		assert.Equal(t, "TKO1 Kollektor Temperatur", dp.DatapointName)

		dp2, exists := maps.ByName["Test Datapoint"]
		assert.True(t, exists)
		assert.Equal(t, "Test Datapoint", dp2.DatapointName)
	})
}

func TestDatapointConfigStruct(t *testing.T) {
	t.Run("struct field mapping", func(t *testing.T) {
		testCSV := createTestCSVFile(t)
		//nolint:errcheck
		defer os.Remove(testCSV)

		maps, err := datapoint.NewCsvDatapointProvider(testCSV)
		require.NoError(t, err)

		dp := maps.ByFunction["22:0:14"]

		// Test all fields are properly mapped
		assert.Equal(t, uint16(1), dp.RegisterAddress)
		assert.Equal(t, "SOL", dp.UnitName)
		assert.Equal(t, "65", dp.UnitID)
		assert.Equal(t, uint8(22), dp.FunctionGroup)
		assert.Equal(t, uint8(0), dp.FunctionNumber)
		assert.Equal(t, uint16(14), dp.DatapointID)
		assert.Equal(t, "TKO1 Kollektor Temperatur", dp.DatapointName)
		assert.Equal(t, "1", dp.Type)
		assert.Equal(t, "S16", dp.TypeName)
		assert.Equal(t, uint8(1), dp.Decimal)
		assert.Equal(t, "Kollektor", dp.FunctionGroupName)
		assert.Equal(t, "Kollektor 1", dp.FunctionName)
	})
}

// createTestCSVFile creates a temporary CSV file with test data
func createTestCSVFile(t *testing.T) string {
	content := `Register Address,UnitName,UnitId,FunctionGroup,FunctionNumber,DatapointId,DatapointName,Type,TypeName,Decimal,FunctionGroup name,Function name,Steps,Min. value,Max. value,Writable,unit,Commentary,Text 0,Text 1,Text 2,Text 3,Text 4,Text 5,Text 6,Text 7,Text 8,Text 9,Text 10,Text 11,Text 12,Text 13,Text 14,Text 15,Text 16,Text 17,Text 18,Text 19,Text 20,Text 21,Text 22,Text 23,Text 24,Text 25,Text 26,Text 27,Text 28,Text 29,Text 30
1,SOL,65,22,0,14,TKO1 Kollektor Temperatur,1,S16,1,Kollektor,Kollektor 1,1,-300,3000,No,°C,TKO1 Kollektor Temperatur,text0,text1,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
2,SOL,65,22,1,14,TKO2 Kollektor Temperatur,1,S16,1,Kollektor,Kollektor 2,1,-300,3000,No,°C,TKO2 Kollektor Temperatur,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
3,SOL,65,22,1,2034,Gesamtertrag Kollektor,2,S32,0,Kollektor,Kollektor 2,1,0,0,No,kWh,Gesamtertrag Kollektor,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
4,TEST,66,23,2,15,Test Datapoint,3,U16,2,Test Group,Test Function,5,10,100,Yes,V,Test Commentary,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,`
	testFile := filepath.Join(t.TempDir(), "test_datapoints.csv")
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)
	return testFile
}
