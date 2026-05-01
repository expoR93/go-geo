package geo

import (
	"errors"
	"math"
	"testing"
)

// tolerance defines the allowed difference between expected and actual values
const tolerance = 1e-5

// assertAlmostEqual is a helper to check if two float64 values are almost equal
func assertAlmostEqual(t *testing.T, expected, actual float64, name string) {
	if math.Abs(expected-actual) > tolerance {
		t.Errorf("%s: expected %f, got %f", name, expected, actual)
	}
}

func TestNewLatLngDegrees(t *testing.T) {
	testCases := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{
			name: "Equator/Prime Meridian",
			lat:  0,
			lng:  0,
		},
		{
			name: "Madrid",
			lat:  40.4168,
			lng:  -3.7038,
		},
		{
			name: "Tokyo",
			lat:  35.6895,
			lng:  139.6917,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewLatLngDegrees(tc.lat, tc.lng)
			assertAlmostEqual(t, tc.lat, result.Lat, "Lat")
			assertAlmostEqual(t, tc.lng, result.Lng, "Lng")
		})
	}
}

func TestNewLatLngRadians(t *testing.T) {
	testCases := []struct {
		name    string
		latRad  float64
		lngRad  float64
		wantLat float64
		wantLng float64
	}{
		{
			name:    "Equator/Prime Meridian",
			latRad:  0,
			lngRad:  0,
			wantLat: 0,
			wantLng: 0,
		},
		{
			name:    "New York City",
			latRad:  0.7105724005118507,
			lngRad:  -1.291648375,
			wantLat: 40.7128,
			wantLng: -74.0060,
		},
		{
			name:    "Sydney",
			latRad:  -0.591122177,
			lngRad:  2.63910015,
			wantLat: -33.8688,
			wantLng: 151.2093,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewLatLngRadians(tc.latRad, tc.lngRad)
			assertAlmostEqual(t, tc.wantLat, result.Lat, "Lat")
			assertAlmostEqual(t, tc.wantLng, result.Lng, "Lng")
		})
	}
}

func TestNewDistance(t *testing.T) {
	testCases := []struct {
		name     string
		value    float64
		unit     Unit
		expected float64
	}{
		{
			name:     "Kilometers (Identity)",
			value:    10.0,
			unit:     UnitKilometers,
			expected: 10.0,
		},
		{
			name:     "Miles to Kilometers",
			value:    10.0,
			unit:     UnitMiles,
			expected: 16.0934,
		},
		{
			name:     "Marathon Miles to KM",
			value:    26.2188226,
			unit:     UnitMiles,
			expected: 42.195,
		},
		{
			name:     "Zero distance",
			value:    0.0,
			unit:     UnitMiles,
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewDistance(tc.value, tc.unit)
			assertAlmostEqual(t, float64(got), tc.expected, "Distance")
		})
	}
}

func TestConfigBuilder(t *testing.T) {
	testCases := []struct {
		name    string
		lat     float64
		lng     float64
		res     int
		wantErr error
	}{
		{"Valid Config", 45.0, 90.0, 9, nil},
		{"Invalid Latitude", 120.0, 90.0, 9, ErrInvalidLatitude},
		{"Invalid Longitude", 45.0, 182.0, 9, ErrInvalidLongitude},
		{"Invalid Resolution", 45.0, 90.0, 20, ErrInvalidResolution},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := NewConfigBuilder().
				WithCoordinate(tc.lat, tc.lng).
				WithResolution(tc.res)

			_, err := builder.Build()

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("Expected error %v, got %v", tc.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
