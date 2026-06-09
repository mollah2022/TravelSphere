package utils

import (
	"TravelSphere/models"
	"fmt"
	"strings"
	"unicode"
)

// FormatPopulation converts a population number into a human-readable format.
func FormatPopulation(pop int64) string {
	switch {
	case pop >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(pop)/1_000_000_000)
	case pop >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(pop)/1_000_000)
	case pop >= 1_000:
		return fmt.Sprintf("%.1fK", float64(pop)/1_000)
	default:
		return fmt.Sprintf("%d", pop)
	}
}

// FormatCurrencies converts a currency map into a readable string format.
func FormatCurrencies(currencies map[string]models.Currency) string {
	if len(currencies) == 0 {
		return "N/A"
	}
	parts := make([]string, 0, len(currencies))
	for code, curr := range currencies {
		if curr.Symbol != "" {
			parts = append(parts, fmt.Sprintf("%s (%s)", curr.Symbol, curr.Name))
		} else {
			parts = append(parts, fmt.Sprintf("%s (%s)", code, curr.Name))
		}
	}
	return strings.Join(parts, ", ")
}

// FormatCurrenciesWithCode converts currency codes into a readable format including the code.
func FormatCurrenciesWithCode(currencies map[string]models.Currency) string {
	if len(currencies) == 0 {
		return "N/A"
	}
	parts := make([]string, 0, len(currencies))
	for code, curr := range currencies {
		parts = append(parts, fmt.Sprintf("%s (%s)", code, curr.Name))
	}
	return strings.Join(parts, ", ")
}

// FormatLanguages converts a language map into a readable string.
func FormatLanguages(languages map[string]string) string {
	if len(languages) == 0 {
		return "N/A"
	}
	parts := make([]string, 0, len(languages))
	for _, lang := range languages {
		parts = append(parts, lang)
	}
	return strings.Join(parts, ", ")
}

// FormatCapital returns the first item from a capital slice.
func FormatCapital(capitals []string) string {
	if len(capitals) == 0 {
		return "N/A"
	}
	return capitals[0]
}

// FormatRegion combines region and subregion into a single readable string.
func FormatRegion(region, subregion string) string {
	if subregion == "" {
		return region
	}
	return fmt.Sprintf("%s - %s", region, subregion)
}

// FormatKinds converts OpenTripMap kinds string into a clean list.
func FormatKinds(kinds string) []string {
	if kinds == "" {
		return []string{}
	}
	raw := strings.Split(kinds, ",")
	result := make([]string, 0, len(raw))
	for _, k := range raw {
		k = strings.TrimSpace(k)
		k = strings.ReplaceAll(k, "_", " ")
		if k != "" {
			result = append(result, k)
		}
	}
	return result
}

// SlugToName converts a slug into a display name.
func SlugToName(slug string) string {
	words := strings.Split(slug, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// NameToSlug converts a country name into a URL-friendly slug.
func NameToSlug(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	var sb strings.Builder
	prevHyphen := false
	for _, r := range name {
		if r == ' ' || r == '_' || r == '-' {
			if !prevHyphen {
				sb.WriteRune('-')
				prevHyphen = true
			}
			continue
		}
		if r == '\'' {
			sb.WriteRune('\'')
			prevHyphen = false
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(r)
			prevHyphen = false
			continue
		}
		// ignore any punctuation or special character
	}
	slug := strings.Trim(sb.String(), "-")
	return slug
}
