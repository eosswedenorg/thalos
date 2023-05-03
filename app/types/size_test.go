package types

import "testing"

func TestSize_Parse(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected int64
		wantErr  bool
	}{
		{"Empty", "", 0, false},
		{"NoDigit", "abcdefg", 0, true},
		{"Negative", "-10MB", 0, true},
		{"Invalid prefix", "100WAX", 0, true},
		{"Multiple spaces between prefix and value", "100  gb", 0, true},
		{"100kb", "100kb", 100 * 1000, false},
		{"10MB", "10 MB", 10 * 1000 * 1000, false},
		{"2gb", "2gb", 2 * 1000 * 1000 * 1000, false},
		{"4Tb", "4 Tb", 4 * 1000 * 1000 * 1000 * 1000, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Size(0)
			if err := s.Parse(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Size.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if int64(s) != tt.expected {
				t.Errorf("Size = %v, expected %v", s, tt.expected)
			}
		})
	}
}
