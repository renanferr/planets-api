package listing

// Planet defines the storage form of a planet
type Planet struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Climate     string `json:"climate"`
	Terrain     string `json:"terrain"`
	Appearances int    `json:"appearances"`
}
