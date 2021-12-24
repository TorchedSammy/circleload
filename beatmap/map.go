package beatmap

type Mapset struct {
	SetID   int
	Artist string
	Title string
	Mapper string `json:"Creator"`
	Status int `json:"RankedStatus"`
}

