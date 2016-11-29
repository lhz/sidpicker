package hvsc

import "strings"

var FilteredTunes = make([]SidTune, 0)
var NumFilteredTunes = 0

func FilterAll() {
	FilteredTunes = Tunes
	NumFilteredTunes = NumTunes
}

func Filter(terms string) {
	FilteredTunes = make([]SidTune, 0)
	for _, tune := range Tunes {
		if filterTune(&tune, terms) {
			FilteredTunes = append(FilteredTunes, tune)
		}
	}
	NumFilteredTunes = len(FilteredTunes)
}

func filterTune(tune *SidTune, terms string) bool {
	exclude := false
	for _, term := range strings.Split(terms, " ") {
		if len(term) > 1 && term[1] == ':' {
			prefix := term[0]
			term = term[2:]
			if prefix == 'y' {
				exclude = exclude || !filterYear(tune.YearMin, tune.YearMax, term)
			} else {
				value := valueByFilterPrefix(tune, prefix)
				exclude = exclude || !strings.Contains(strings.ToUpper(value), strings.ToUpper(term))
			}
		} else {
			exclude = exclude || !strings.Contains(tune.Header.Author, term)
		}
	}

	return !exclude
}

func filterYear(min, max int, term string) bool {
	strict := false
	if term[len(term)-1] == '!' {
		term = term[:len(term)-1]
		strict = true
	}
	parts := strings.Split(term, "-")
	yearFrom := parseYear(parts[0], 1900)
	if len(parts) == 1 {
		if strict {
			if yearFrom != min || yearFrom != max {
				return false
			}
		} else {
			if yearFrom < min || yearFrom > max {
				return false
			}
		}
	} else {
		yearTo := parseYear(parts[1], 9999)
		if strict {
			if yearFrom > max || yearTo < min || min == 1900 || max == 9999 {
				return false
			}
		} else {
			if yearFrom > max || yearTo < min {
				return false
			}
		}
	}
	return true
}

func valueByFilterPrefix(tune *SidTune, prefix byte) string {
	switch prefix {
	case 'a':
		return tune.Header.Author
	case 'n':
		return tune.Header.Name
	case 'p':
		return tune.Path
	case 'r':
		return tune.Header.Released
	case 't':
		return tune.Header.Name
	}
	return ""
}
