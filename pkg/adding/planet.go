package adding

// Planet defines the properties of a planet to be added
type Planet struct {
	Name        string `json:"name" valid:"length(2|128)"`
	Climate     string `json:"climate" valid:"length(2|128)"`
	Terrain     string `json:"terrain" valid:"length(2|128)"`
	Appearances int    `json:"appearances" valid:"optional"`
}
