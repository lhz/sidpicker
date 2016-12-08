package hvsc

import "strings"

var FilteredTunes = make([]SidTune, 0)
var NumFilteredTunes = 0

func FilterAll() {
	FilteredTunes = Tunes
	NumFilteredTunes = NumTunes
}

func Filter(terms string) bool {
	tunes := make([]SidTune, 0)
	for _, tune := range Tunes {
		if filterTune(&tune, terms) {
			tunes = append(tunes, tune)
		}
	}
	if len(tunes) == 0 {
		return false
	}
	FilteredTunes = tunes
	NumFilteredTunes = len(FilteredTunes)
	return true
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
			exclude = exclude || !strings.Contains(strings.ToUpper(allText(tune)), strings.ToUpper(term))
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
	case 's':
		return strings.Join(tune.Info, " ")
	case 't':
		return tune.Header.Name
	}
	return ""
}

func allText(tune *SidTune) string {
	text := []string{
		tune.Header.Author,
		tune.Header.Name,
		tune.Path,
		tune.Header.Released}
	text = append(text, tune.Info...)
	return strings.Join(text, " ")
}
