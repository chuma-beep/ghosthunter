package engine



type EntityType int 

const (
	EntityGhost EntityType = iota 
	EntityWizard
	EntityDemon
  EntityWraith
  EntityReaper 
)


type EntityState int

const (
    StateChase EntityState = iota
    StateAttack
    StatePain
    StateDeath
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


type Entity struct {
    X, Y        float64
    VX, VY      float64
    Frame        int
    FrameTimer   int
    FadeTimer    int
    Dead         bool
    Type         EntityType
    Health       int
    Speed        float64
    Damage       int
    State        EntityState
    StateTimer   int
    FacingAngle  float64
    LOSTimer int
    HasLOS   bool
}


type HealthPickup struct {
    X, Y   float64
    Active bool
}
