package ec

type Controller struct {
	hvac        HVAC
	gauge       Gauge
	blowerTimer int
	coolerTimer int
}

type Gauge interface {
	CurrentTemperature() int // Current ambient temperature rounded to the nearest degree (Fahrenheit).
}

type HVAC interface {
	SetBlower(state bool) // Turns the blower on or off.
	SetCooler(state bool) // Turns the cooler on or off.
	SetHeater(state bool) // Turns the heater on or off.

	IsBlowing() bool // Is the blower currently on or off?
	IsCooling() bool // Is the cooler currently on or off?
	IsHeating() bool // Is the heater currently on or off?
}

func newController(hvac HVAC, thermo Gauge) *Controller {
	hvac.SetBlower(off)
	hvac.SetCooler(off)
	hvac.SetHeater(off)
	return &Controller{hvac: hvac, gauge: thermo}
}

func (c *Controller) Regulate() {
	c.DecrementTimers()
	if c.TemperatureTooHot() {
		c.CoolTheRoom()
	} else if c.TemperatureTooCold() {
		c.HeatTheRoom()
	} else {
		c.NoHeatingOrCoolingTheRoom()
	}
}

func (c *Controller) DecrementTimers() {
	if c.blowerTimer > 0 {
		c.blowerTimer--
	}
	if c.coolerTimer > 0 {
		c.coolerTimer--
	}
}

func (c *Controller) CoolTheRoom() {
	c.hvac.SetBlower(on)
	c.TurnOnCooler()
	c.hvac.SetHeater(off)
}
func (c *Controller) HeatTheRoom() {
	c.hvac.SetBlower(on)
	c.TurnOnHeater()
	c.TurnOffCooler()
}
func (c *Controller) NoHeatingOrCoolingTheRoom() {
	c.TurnOffBlower()
	c.hvac.SetHeater(off)
	c.TurnOffCooler()
}

func (c *Controller) TemperatureTooHot() bool {
	return c.gauge.CurrentTemperature() > 75
}
func (c *Controller) TemperatureTooCold() bool {
	return c.gauge.CurrentTemperature() < 65
}

func (c *Controller) TurnOffCooler() {
	if c.hvac.IsCooling() {
		c.coolerTimer = 3
	}
	c.hvac.SetCooler(off)
}
func (c *Controller) TurnOnCooler() {
	if c.coolerTimer == 0 {
		c.hvac.SetCooler(on)
	}
}
func (c *Controller) TurnOnHeater() {
	c.hvac.SetHeater(on)
	c.blowerTimer = 5
}
func (c *Controller) TurnOffBlower() {
	if c.blowerTimer == 0 {
		c.hvac.SetBlower(false)
	}
}

const on = true
const off = false
