package vo

type Point struct {
	X int
	Y int
}

type ClickInfo struct {
	TotalImageHeight int
	TotalImageWidth  int
	Points           []Point
}
