package slug

import (
	"regexp"
	"strings"
)

// Generate membuat slug dari string input
func Generate(input string) string {
	// Membuat lowercase semua karakter
	slug := strings.ToLower(input)

	// Mengganti spasi dengan dash (-)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Menghilangkan karakter yang tidak di inginkan
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Menghilangkan multiple dash (--- menjadi -)
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "")

	// Trim dash awal dan akhir
	slug = strings.Trim(slug, "-")

	// Limiit panjang slug max 100 karakter
	if len(slug) > 100 {
		slug = slug[:100]
		// Trim dash di akhir kalau ada
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}
