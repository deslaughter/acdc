package acdc

type Conditions struct {
	WindSpeed            float64 // Wind speed (m/s)
	BladePitch           float64 // Blade pitch (deg)
	RotorSpeed           float64 // Rotor speed in (rpm)
	TowerTopDispForeAft  float64 // Tower Top Displacement Fore-Aft (m)
	TowerTopDispSideSide float64 // Tower Top Displacement Side-Side (m)
}
