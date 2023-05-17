package game

type Player struct {
	Name       string `json:"name"`
	Hp         int    `json:"hp"`
	Gender     string `json:"gender"`
	Preference string `json:"preference"`
}
