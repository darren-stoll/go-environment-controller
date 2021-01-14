package ec

import "testing"

func TestAtStartupEverythingTurnedOff(t *testing.T) {
	f := NewFixture(t)
	f.AssertAllOff()
}

func TestWhenTooHotEngageCoolerAndBlower(t *testing.T) {
	f := NewFixture(t)
	f.MakeTooHot()
	f.AssertCooling()
}

func TestWhenTooColdEngageHeaterAndBlower(t *testing.T) {
	f := NewFixture(t)
	f.MakeTooCold()
	f.AssertHeating()
}

func TestWhenTooColdThenComfortable(t *testing.T) {
	f := NewFixture(t)
	f.MakeTooCold()
	f.MakeComfortable()
	f.AssertBlowing()
}

func TestWhenTooColdThenTooHot(t *testing.T) {
	f := NewFixture(t)
	f.MakeTooCold()
	f.MakeTooHot()
	f.AssertCooling()
}

func TestWhenTooHotThenTooCold(t *testing.T) {
	f := NewFixture(t)
	f.MakeTooHot()
	f.MakeTooCold()
	f.AssertHeating()
}

func TestWhenTooHotThenComfortable(t *testing.T) {
	f := NewFixture(t)
	f.MakeTooHot()
	f.MakeComfortable()
	f.AssertAllOff()
}

func TestWhenTooColdThenBlowerStaysOnForFiveMinutesAfterHeating(t *testing.T) {
	f := NewFixture(t)
	f.MakeTooCold()
	for i := 0; i < 4; i++ {
		f.MakeComfortable()
		f.AssertBlowing()
	}
	f.MakeComfortable()
	f.AssertAllOff()
}

func TestWhenTooHotThenComfortableThenTooHotThenCoolerStaysOffForThreeMinutes(t *testing.T) {
	f := NewFixture(t)
	f.MakeTooHot()
	f.MakeComfortable()
	for i := 0; i < 2; i++ {
		f.MakeTooHot()
		f.AssertBlowing()
	}
	f.MakeTooHot()
	f.AssertCooling()
}

type Fixture struct {
	hvac       *FakeHVAC
	thermo     *FakeThermometer
	controller *Controller
}

func NewFixture(t *testing.T) *Fixture {
	hvac := NewHVAC(t)
	thermo := NewThermometer()
	controller := newController(hvac, thermo)
	return &Fixture{hvac: hvac, thermo: thermo, controller: controller}
}

func (f *Fixture) MakeTooHot() {
	f.thermo.temp = 76
	f.controller.Regulate()
}

func (f *Fixture) MakeTooCold() {
	f.thermo.temp = 64
	f.controller.Regulate()
}

func (f *Fixture) MakeComfortable() {
	f.thermo.temp = 70
	f.controller.Regulate()
}

func (f *Fixture) AssertAllOff()  { f.hvac.AssertState("bch") }
func (f *Fixture) AssertBlowing() { f.hvac.AssertState("Bch") }
func (f *Fixture) AssertCooling() { f.hvac.AssertState("BCh") }
func (f *Fixture) AssertHeating() { f.hvac.AssertState("BcH") }

type FakeThermometer struct{ temp int }

func NewThermometer() *FakeThermometer             { return &FakeThermometer{temp: 70} }
func (t *FakeThermometer) CurrentTemperature() int { return t.temp }

type FakeHVAC struct {
	blowing, cooling, heating bool
	t                         *testing.T
}

func NewHVAC(t *testing.T) *FakeHVAC {
	return &FakeHVAC{
		blowing: true,
		cooling: true,
		heating: true,
		t:       t,
	}
}

func (h *FakeHVAC) SetBlower(state bool) { h.blowing = state }
func (h *FakeHVAC) SetCooler(state bool) { h.cooling = state }
func (h *FakeHVAC) SetHeater(state bool) { h.heating = state }
func (h *FakeHVAC) IsBlowing() bool      { return h.blowing }
func (h *FakeHVAC) IsCooling() bool      { return h.cooling }
func (h *FakeHVAC) IsHeating() bool      { return h.heating }
func (h *FakeHVAC) AssertBlower(expectedState bool) {
	if h.IsBlowing() != expectedState {
		h.t.Errorf("Expected %t, actual %t for blowing", expectedState, h.IsBlowing())
	}
}
func (h *FakeHVAC) AssertCooler(expectedState bool) {
	if h.IsCooling() != expectedState {
		h.t.Errorf("Expected %t, actual %t for cooling", expectedState, h.IsCooling())
	}
}
func (h *FakeHVAC) AssertHeater(expectedState bool) {
	if h.IsHeating() != expectedState {
		h.t.Errorf("Expected %t, actual %t for heating", expectedState, h.IsHeating())
	}
}

// AssertState expected string should contain three characters: one for blowing, one for cooling, and one for heating. Uppercase is true, while lowercase is false
func (h *FakeHVAC) AssertState(expected string) {
	h.AssertBlower(expected[0] == 'B')
	h.AssertCooler(expected[1] == 'C')
	h.AssertHeater(expected[2] == 'H')
}
