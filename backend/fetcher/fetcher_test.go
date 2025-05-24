package fetcher

import (
    "testing"
)
// Realizo el testing de la función ParseDollar
func TestParseDollar(t *testing.T) {
	// Defino estructura del test
    tests := []struct {
        in      string
        want    float64
        wantErr bool
    }{
        {in: "$1,234.56", want: 1234.56, wantErr: false},
        {in: "  $0.50  ", want: 0.5, wantErr: false},
        {in: "100", want: 100.0, wantErr: false},
        {in: "$bad", want: 0, wantErr: true},
        {in: "", want: 0, wantErr: true},
    }
		// Recorro los casos a probar
    for _, tc := range tests {
        got, err := parseDollar(tc.in)
				// En caso de deber dar error
        if tc.wantErr {
            if err == nil {
                t.Errorf("parseDollar(%q): se esperaba un error, se obtuvo err == nil", tc.in)
            }
            continue
        }
				// En caso de error inesperado informo
        if err != nil {
            t.Errorf("parseDollar(%q): error inesperado %v", tc.in, err)
            continue
        }
				// En caso de no coincidir lo buscado con lo obtenido
        if got != tc.want {
            t.Errorf("parseDollar(%q) = %v; buscado %v", tc.in, got, tc.want)
        }
    }
}
