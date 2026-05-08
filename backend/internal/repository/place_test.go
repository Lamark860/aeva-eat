package repository

import "testing"

func TestEscapeLike(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain text untouched", "кафе", "кафе"},
		{"percent escaped", "100%", `100\%`},
		{"underscore escaped", "ka_e", `ka\_e`},
		{"backslash escaped first", `c:\foo`, `c:\\foo`},
		{"all three at once", `a_b%c\d`, `a\_b\%c\\d`},
		{"empty string", "", ""},
		{"only wildcards", `%_`, `\%\_`},
		// Ensures backslash escape happens first; otherwise `\` would become `\\` after `%` → `\%`
		// is processed, which would corrupt the escape sequences.
		{"backslash before percent", `\%`, `\\\%`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeLike(tt.in)
			if got != tt.want {
				t.Errorf("escapeLike(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
