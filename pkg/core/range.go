package core

type LineRange struct {
	File  string `json:"file"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}
