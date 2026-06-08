package utils

import (
	"TravelSphere/models"
	"fmt"
	"strings"
	"unicode"
)

// FormatPopulation population কে readable format এ convert করে
// যেমন: 170000000 → "170M", 2400000 → "2.4M", 47400 → "47.4K"
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

// FormatCurrencies currency map কে readable string এ convert করে
// যেমন: {"BDT": {Name: "Bangladeshi taka", Symbol: "৳"}} → "BDT (Bangladeshi taka)"
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

// FormatCurrenciesWithCode currency code সহ format করে
// যেমন: "BDT (Bangladeshi taka)"
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

// FormatLanguages language map কে readable string এ convert করে
// যেমন: {"ben": "Bengali", "eng": "English"} → "Bengali, English"
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

// FormatCapital capital slice এর প্রথম item নেয়
func FormatCapital(capitals []string) string {
	if len(capitals) == 0 {
		return "N/A"
	}
	return capitals[0]
}

// FormatRegion region ও subregion combine করে
// যেমন: "Asia", "Southern Asia" → "Asia - Southern Asia"
func FormatRegion(region, subregion string) string {
	if subregion == "" {
		return region
	}
	return fmt.Sprintf("%s - %s", region, subregion)
}

// FormatKinds OpenTripMap kinds string কে clean list এ convert করে
// যেমন: "museums,historic_architecture" → ["museums", "historic architecture"]
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

// SlugToName slug কে display name এ convert করে
// যেমন: "united-states" → "United States"
func SlugToName(slug string) string {
	words := strings.Split(slug, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

// NameToSlug country name কে URL slug এ convert করে
// যেমন: "United States" → "united-states"
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
