package openai

import "testing"

func TestHasSuffix(t *testing.T) {
	cases := []struct {
		s      string
		suffix string
		want   bool
	}{
		{"hello.jpg", ".jpg", true},
		{"hello.jpeg", ".jpg", false},
		{"image.png", ".png", true},
		{"image.png?query=1", ".png", false},
	}
	for _, c := range cases {
		if got := hasSuffix(c.s, c.suffix); got != c.want {
			t.Errorf("hasSuffix(%q, %q)=%v want %v", c.s, c.suffix, got, c.want)
		}
	}
}

func TestExtFromContentType(t *testing.T) {
	cases := []struct {
		ct   string
		want string
	}{
		{"image/jpeg", ".jpg"},
		{"image/png", ".png"},
		{"image/webp", ".webp"},
		{"application/json", ""},
	}
	for _, c := range cases {
		if got := extFromContentType(c.ct); got != c.want {
			t.Errorf("extFromContentType(%q)=%q want %q", c.ct, got, c.want)
		}
	}
}
