package engine


type Sprite struct{
	X float64
	Y float64
	VX float64
	VY float64
    FadeTimer int 
	Dead bool
}

type AmmoPickup struct {
    X float64
    Y float64
    Active bool
}




