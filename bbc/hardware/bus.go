package hardware

import (
	"bbc/logical"
	"bbc/utils"
	"fmt"
)

type Bus struct {
	Clock        *Clock
	watchers     map[string]Component
	addressables map[string]AddressableComponent
}

type Component interface {
	GetName() string
	Start() error
	Reset() error
	Stop() error
	PlugToBus(*Bus)
}

type AddressableComponent interface {
	Component
	IsWritable() bool
	IsReadable() bool
	GetSegment() *utils.Segment
}

type ReadableComponent interface {
	AddressableComponent
	DirectRead(uint16) (byte, error)
	OffsetRead(uint16, uint8) (byte, uint16, error)
}

type WritableComponent interface {
	AddressableComponent
	DirectWrite(byte, uint16) error
	OffsetWrite(byte, uint16, uint8) (uint16, error)
}

func (bus *Bus) componentReadAt(addr uint16) ReadableComponent {
	for _, component := range bus.addressables {
		if !component.IsReadable() {
			continue
		}

		readComponent := component.(ReadableComponent)
		if readComponent.GetSegment().IsIn(addr) {
			return readComponent
		}
	}
	return nil
}

func (bus *Bus) componentWriteAt(addr uint16) WritableComponent {
	for _, component := range bus.addressables {
		if !component.IsWritable() {
			continue
		}
		writeComponent := component.(WritableComponent)
		if writeComponent.GetSegment().IsIn(addr) {
			return writeComponent
		}
	}
	return nil
}

func (bus *Bus) Reset() {
	bus.Clock.Reset()
	for _, component := range bus.watchers {
		component.Reset()
	}
	for _, component := range bus.addressables {
		component.Reset()
	}
}

func (bus *Bus) Tick() error {
	return bus.Clock.Tick()
}

// 1 cycle
func (bus *Bus) DirectRead(addr uint16) (byte, error) {
	readComponent := bus.componentReadAt(addr)
	if readComponent == nil {
		return 0, fmt.Errorf("reading garbage as no component answer for this address %x", addr)
	}
	if err := bus.Tick(); err != nil {
		return 0, err
	}
	return readComponent.DirectRead(addr)
}

// 1 cycle, +1 if page crossed or forced
func (bus *Bus) OffsetRead(addr uint16, offset uint8, forceTick bool) (byte, uint16, error) {
	readComponent := bus.componentReadAt(addr)
	if readComponent == nil {
		return 0, 0, fmt.Errorf("reading garbage as no component answer for this address %x", addr)
	}
	if err := bus.Tick(); err != nil {
		return 0, 0, err
	}
	if forceTick || utils.IsPageCrossed(addr, offset) {
		if err := bus.Tick(); err != nil {
			return 0, 0, err
		}
	}
	return readComponent.OffsetRead(addr, offset)
}

// 1 cycle
func (bus *Bus) DirectWrite(value byte, addr uint16) error {
	writeComponent := bus.componentWriteAt(addr)
	if writeComponent == nil {
		return fmt.Errorf("writing in void as no component answer for this address %x", addr)
	}
	if err := bus.Tick(); err != nil {
		return err
	}
	return writeComponent.DirectWrite(value, addr)
}

// 1 cycle, +1 if page crossed or forced
func (bus Bus) OffsetWrite(value byte, addr uint16, offset uint8, forceTick bool) (uint16, error) {
	writeComponent := bus.componentWriteAt(addr)
	if writeComponent == nil {
		return 0, fmt.Errorf("writing in void as no component answer for this address %x", addr)
	}
	if err := bus.Tick(); err != nil {
		return 0, err
	}
	if forceTick || utils.IsPageCrossed(addr, offset) {
		if err := bus.Tick(); err != nil {
			return 0, err
		}
	}
	return writeComponent.OffsetWrite(value, addr, offset)
}

func (bus *Bus) AddComponent(component Component) error {
	addrComponent, ok := component.(AddressableComponent)
	if !ok {
		fmt.Printf("Adding new watcher component %s\n", component.GetName())
		_, okW := bus.watchers[component.GetName()]
		if okW {
			return fmt.Errorf("component already registered with name %s", component.GetName())
		}
		bus.watchers[component.GetName()] = component
	} else {
		fmt.Printf("Adding new addressable component %s\n", component.GetName())
		segment := addrComponent.GetSegment()
		for _, registered := range bus.addressables {
			if segment.Intersect(registered.GetSegment()) {
				return fmt.Errorf("cannot register new component on bus, segment intersects with %s one", registered.GetName())
			}
			_, okA := bus.addressables[component.GetName()]
			_, okW := bus.watchers[component.GetName()]
			if okA || okW {
				return fmt.Errorf("component already registered with name %s", component.GetName())
			}
		}
		bus.addressables[component.GetName()] = addrComponent
	}
	component.PlugToBus(bus)
	return nil
}

func (bus *Bus) GetComponent(name string) Component {
	if component, ok := bus.addressables[name]; ok {
		return component
	}
	if component, ok := bus.watchers[name]; ok {
		return component
	}
	return nil
}

func (bus *Bus) WriteMultiple(values []byte, start uint16) error {
	if len(values)+int(start) > int(logical.AdressableSegment.End) {
		return fmt.Errorf("cannot write outside memory bound")
	}
	addr := start
	for _, v := range values {
		if err := bus.DirectWrite(v, addr); err != nil {
			return err
		}
		addr++
	}
	return nil
}

func NewBus(clock *Clock, components ...Component) (*Bus, error) {
	bus := Bus{
		Clock:        clock,
		watchers:     map[string]Component{},
		addressables: map[string]AddressableComponent{},
	}
	for _, component := range components {
		if err := bus.AddComponent(component); err != nil {
			return nil, err
		}
	}
	return &bus, nil
}
