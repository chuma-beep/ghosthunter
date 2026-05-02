package engine

// WeaponStateID identifies each state in the weapon state table
type WeaponStateID int

const (
	S_NULL WeaponStateID = iota

	// Pistol states
	S_PISTOL_UP
	S_PISTOL_DOWN
	S_PISTOL_READY
	S_PISTOL_ATK1
	S_PISTOL_ATK2
	S_PISTOL_ATK3

	// Shotgun states
	S_SHOTGUN_UP
	S_SHOTGUN_DOWN
	S_SHOTGUN_READY
	S_SHOTGUN_ATK1
	S_SHOTGUN_ATK2
	S_SHOTGUN_ATK3
	S_SHOTGUN_ATK4

	// Machinegun states
	S_MACHINEGUN_UP
	S_MACHINEGUN_DOWN
	S_MACHINEGUN_READY
	S_MACHINEGUN_ATK1
	S_MACHINEGUN_ATK2
	S_MACHINEGUN_ATK3
	S_MACHINEGUN_ATK4
	S_MACHINEGUN_ATK5
	S_MACHINEGUN_ATK6
	S_MACHINEGUN_ATK7
	S_MACHINEGUN_ATK8
)

// WeaponAction is a function called when entering a state
type WeaponAction func(g *Game)

// WeaponState is one entry in the state table
type WeaponState struct {
	Weapon    int           // which weapon animation folder (0=pistol, 1=shotgun, 2=machinegun)
	Frame     int           // frame index within that weapon's animation
	Tics      int           // how many game ticks to stay in this state (-1 = forever)
	Action    WeaponAction  // function to call on entry (nil = no action)
	NextState WeaponStateID // state to go to after Tics expire
}

// WeaponDef defines a weapon's key states
type WeaponDef struct {
	ReadyState WeaponStateID
	AtkState   WeaponStateID
	UpState    WeaponStateID
	DownState  WeaponStateID
}

// Weapons defines the 3 weapons
var Weapons = []WeaponDef{
	{ // 0 - Pistol
		ReadyState: S_PISTOL_READY,
		AtkState:   S_PISTOL_ATK1,
		UpState:    S_PISTOL_UP,
		DownState:  S_PISTOL_DOWN,
	},
	{ // 1 - Shotgun
		ReadyState: S_SHOTGUN_READY,
		AtkState:   S_SHOTGUN_ATK1,
		UpState:    S_SHOTGUN_UP,
		DownState:  S_SHOTGUN_DOWN,
	},
	{ // 2 - Machinegun
		ReadyState: S_MACHINEGUN_READY,
		AtkState:   S_MACHINEGUN_ATK1,
		UpState:    S_MACHINEGUN_UP,
		DownState:  S_MACHINEGUN_DOWN,
	},
}

// WeaponStates
var WeaponStates = map[WeaponStateID]WeaponState{
	S_NULL: {Weapon: 0, Frame: 0, Tics: -1, Action: nil, NextState: S_NULL},

	// Pistol
	S_PISTOL_UP:    {Weapon: 0, Frame: 0, Tics: 1, Action: nil, NextState: S_PISTOL_READY},
	S_PISTOL_DOWN:  {Weapon: 0, Frame: 0, Tics: 1, Action: nil, NextState: S_PISTOL_READY},
	S_PISTOL_READY: {Weapon: 0, Frame: 0, Tics: -1, Action: nil, NextState: S_PISTOL_READY},
	S_PISTOL_ATK1:  {Weapon: 0, Frame: 1, Tics: 4, Action: A_FirePistol, NextState: S_PISTOL_ATK2},
	S_PISTOL_ATK2:  {Weapon: 0, Frame: 1, Tics: 4, Action: nil, NextState: S_PISTOL_ATK3},
	S_PISTOL_ATK3:  {Weapon: 0, Frame: 0, Tics: 4, Action: nil, NextState: S_PISTOL_READY},

	// shotgun
	S_SHOTGUN_UP:    {Weapon: 0, Frame: 0, Tics: 1, Action: nil, NextState: S_SHOTGUN_READY},
	S_SHOTGUN_DOWN:  {Weapon: 0, Frame: 0, Tics: 1, Action: nil, NextState: S_SHOTGUN_READY},
	S_SHOTGUN_READY: {Weapon: 0, Frame: 0, Tics: -1, Action: nil, NextState: S_SHOTGUN_READY},
	S_SHOTGUN_ATK1:  {Weapon: 0, Frame: 1, Tics: 4, Action: A_FireShotgun, NextState: S_SHOTGUN_ATK2},
	S_SHOTGUN_ATK2:  {Weapon: 0, Frame: 1, Tics: 6, Action: nil, NextState: S_SHOTGUN_ATK3},
	S_SHOTGUN_ATK3:  {Weapon: 0, Frame: 1, Tics: 6, Action: nil, NextState: S_SHOTGUN_ATK4},
	S_SHOTGUN_ATK4:  {Weapon: 0, Frame: 0, Tics: 8, Action: nil, NextState: S_SHOTGUN_READY},

	// Machinegun
	S_MACHINEGUN_UP:    {Weapon: 2, Frame: 0, Tics: 1, Action: nil, NextState: S_MACHINEGUN_READY},
	S_MACHINEGUN_DOWN:  {Weapon: 2, Frame: 0, Tics: 1, Action: nil, NextState: S_MACHINEGUN_READY},
	S_MACHINEGUN_READY: {Weapon: 2, Frame: 0, Tics: -1, Action: nil, NextState: S_MACHINEGUN_READY},
	S_MACHINEGUN_ATK1:  {Weapon: 2, Frame: 1, Tics: 3, Action: A_FireMachinegun, NextState: S_MACHINEGUN_ATK2},
	S_MACHINEGUN_ATK2:  {Weapon: 2, Frame: 2, Tics: 3, Action: A_FireMachinegun, NextState: S_MACHINEGUN_ATK3},
	S_MACHINEGUN_ATK3:  {Weapon: 2, Frame: 3, Tics: 3, Action: A_FireMachinegun, NextState: S_MACHINEGUN_ATK4},
	S_MACHINEGUN_ATK4:  {Weapon: 2, Frame: 4, Tics: 3, Action: A_FireMachinegun, NextState: S_MACHINEGUN_ATK5},
	S_MACHINEGUN_ATK5:  {Weapon: 2, Frame: 5, Tics: 3, Action: A_FireMachinegun, NextState: S_MACHINEGUN_ATK6},
	S_MACHINEGUN_ATK6:  {Weapon: 2, Frame: 6, Tics: 3, Action: A_FireMachinegun, NextState: S_MACHINEGUN_ATK7},
	S_MACHINEGUN_ATK7:  {Weapon: 2, Frame: 7, Tics: 3, Action: A_FireMachinegun, NextState: S_MACHINEGUN_ATK8},
	S_MACHINEGUN_ATK8:  {Weapon: 2, Frame: 8, Tics: 3, Action: nil, NextState: S_MACHINEGUN_READY},
}

//  Action functions

func A_FirePistol(g *Game) {
	if g.Ammo < 1 {
		return
	}
	g.Ammo--
	g.GunKick = 8
	g.ScreenShake = 8
	PlaySound("assets/shoot.wav")
	g.shootRay(g.Angle, 1)
}

func A_FireShotgun(g *Game) {
	if g.Ammo < 3 {
		return
	}
	g.Ammo -= 3
	g.GunKick = 8
	g.ScreenShake = 8
	PlaySound("assets/shoot.wav")
	for s := -2; s <= 2; s++ {
		g.shootRay(g.Angle+float64(s)*0.05, 2)
	}
}

func A_FireMachinegun(g *Game) {
	if g.Ammo < 1 {
		return
	}
	g.Ammo--
	g.GunKick = 4
	g.ScreenShake = 4
	PlaySound("assets/shoot.wav")
	g.shootRay(g.Angle, 1)
}

// SetWeaponState transitions to a new state and calls its action
func (g *Game) SetWeaponState(id WeaponStateID) {
	g.WeaponStateID = id
	g.WeaponStateTics = WeaponStates[id].Tics
	if WeaponStates[id].Action != nil {
		WeaponStates[id].Action(g)
	}
}

// TickWeapon advances the weapon state machine each game tick
func (g *Game) TickWeapon() {
	state := WeaponStates[g.WeaponStateID]
	if state.Tics == -1 {
		// waiting for input — check for fire
		canFire := false
		aiFiring := (g.AI != nil && g.AI.Enabled && g.AI.ShootRequested) || (g.NeuralAI != nil && g.NeuralAI.Enabled && g.NeuralAI.ShootRequested)
		if g.WeaponType == 2 {
			canFire = isSpacePressed() || aiFiring
		} else {
			canFire = isSpaceJustPressed() || aiFiring
		}
		if canFire && g.Ammo > 0 {
			g.SetWeaponState(Weapons[g.WeaponType].AtkState)
		}
		return
	}
	g.WeaponStateTics--
	if g.WeaponStateTics <= 0 {
		g.SetWeaponState(state.NextState)
	}
}
