package darknet

//Object represents an object detected in an image
type Object struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Prob  int    `json:"prob"`
}

//Result is a darknet detection result
type Result struct {
	Result  string   `json:"res"`
	Objects []Object `json:"objects"`
}
