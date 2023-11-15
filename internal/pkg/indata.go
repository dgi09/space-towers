package pkg

type InDataDetails struct {
	Name string `json:"name"`
}

type InDataRoom struct {
	Room string `json:"room"`
}

type InDataRoundFinished struct {
	Combos []int `json:"combos"`
}
