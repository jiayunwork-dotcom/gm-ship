package hull

import (
	"math"
	"testing"
)

func TestDisplacementTotalMass(t *testing.T) {
	d := Displacement{LightShip: 1000, Deadweight: 500, Consumables: 100}
	if math.Abs(d.TotalMass()-1600) > 1e-9 {
		t.Fatalf("total mass should be 1600, got %v", d.TotalMass())
	}
}

func TestDisplacementVolume(t *testing.T) {
	d := Displacement{LightShip: 1025, Deadweight: 0, Consumables: 0}
	v := d.DisplacementVolume( 1025)
	if math.Abs(v-1.0) > 1e-9 {
		t.Fatalf("volume should be 1.0, got %v", v)
	}
}

func TestCargoMassToDraft(t *testing.T) {
	d := Displacement{}
	delta := d.CargoMassToDraft(1025, 100, 1025)
	if math.Abs(delta-0.01) > 1e-9 {
		t.Fatalf("1 tonne over 100 m² should sink 1 cm, got %v", delta)
	}
}

func TestTonsPerCentimetreImmersion(t *testing.T) {
	tpc := TonsPerCentimetreImmersion(100, 1025)
	expected := 1025 * 100 * 0.01
	if math.Abs(tpc-expected) > 1e-9 {
		t.Fatalf("TPC should be %v, got %v", expected, tpc)
	}
}

func TestDeadweightFraction(t *testing.T) {
	d := Displacement{LightShip: 1000, Deadweight: 1000, Consumables: 0}
	if math.Abs(d.DeadweightFraction()-0.5) > 1e-9 {
		t.Fatalf("deadweight fraction should be 0.5, got %v", d.DeadweightFraction())
	}
}

func TestVolumetricUtilisation(t *testing.T) {
	d := Displacement{LightShip: 100, Deadweight: 500, Consumables: 0}
	u := d.VolumetricUtilisation(2, 100) // 200 kg in hold
	if math.Abs(u-0.4) > 1e-9 {
		t.Fatalf("utilisation should be 0.4, got %v", u)
	}
}

func TestFreeboard(t *testing.T) {
	if math.Abs(Freeboard(10, 6)-4.0) > 1e-9 {
		t.Fatalf("freeboard should be 4, got %v", Freeboard(10, 6))
	}
	if Freeboard(5, 6) != 0 {
		t.Fatalf("submerged deck freeboard must be 0")
	}
}

func TestBlockCoefficientFromDisplacement(t *testing.T) {
	cb := BlockCoefficientFromDisplacement(1025, 1025, 10, 5, 2)
	if math.Abs(cb-0.01) > 1e-9 {
		t.Fatalf("Cb should be 0.01, got %v", cb)
	}
}

func TestReserveBuoyancyFraction(t *testing.T) {
	rb := ReserveBuoyancyFraction(10, 6, 5, 10)
	if rb <= 0 || rb >= 1 {
		t.Fatalf("reserve buoyancy out of range: %v", rb)
	}
}

func TestListFromShift(t *testing.T) {
	list := ListFromShift(100, 5, 10000, 1.0)
	if list <= 0 || list > 90 {
		t.Fatalf("list angle out of range: %v", list)
	}
	if ListFromShift(100, 5, 0, 1) != 0 {
		t.Fatalf("zero displacement must yield 0")
	}
}
