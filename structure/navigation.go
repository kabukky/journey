package structure

// Navigation is an entry in the navigation menu
type Navigation struct {
	Label string `json:"label"`
	URL   string `json:"url"`
	Slug  string `json:"-"`
}
