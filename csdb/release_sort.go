package csdb

type ByDate []Release

func (s ByDate) Len() int {
	return len(s)
}

func (s ByDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByDate) Less(i, j int) bool {
	if s[i].Date == s[j].Date {
		return s[i].Name < s[j].Name
	} else {
		return s[i].Date < s[j].Date
	}
}
