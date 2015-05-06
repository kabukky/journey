package structure

// Navigation: an entry in the navigation menu
type Navigation struct {
	Label string `json:"label"`
	Url   string `json:"url"`
	Slug  string `json:"-"`
}
