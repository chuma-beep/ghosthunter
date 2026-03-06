package engine



type EntityType int 

const (
	EntityGhost EntityType = iota 
	EntityWizard
	EntityDemon 
)


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

type Entity struct{
	X  float64
	Y  float64
	VX float64
	VY float64
	FadeTimer int 
	Dead bool 
	Type EntityType 
	Health int 
	Speed  float64
	Damage int 
}



