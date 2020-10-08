package adding

// Planet defines the properties of a planet to be added
type Planet struct {
	Name        string `json:"name"`
	Climate     string `json:"climate"`
	Terrain     string `json:"terrain"`
	Appearances int64  `json:"Appearances"`
}
