package slug

import (
	"fmt"
	"regexp"
	"strings"
	"time"
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
	slug = reg.ReplaceAllString(slug, "-")

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

// Mmebuat slug unique dengan menambah timestamp
func GenerateUnique(input string) string {
	baseSlug := Generate(input)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", baseSlug, timestamp)
}

func Validate(slug string) bool {
	if slug == "" || len(slug) > 100 {
		return false
	}

	reg := regexp.MustCompile("^[a-z0-9-]+$")
	return reg.MatchString(slug)
}
