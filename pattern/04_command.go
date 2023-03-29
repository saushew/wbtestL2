package pattern

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern

The Command Pattern is a design pattern that is used to encapsulate a request as an object, thereby letting you parameterize clients with different requests, queue or log requests, and support undoable operations.
Pros:
1. Decouples sender and receiver: The Command Pattern separates the sender and receiver of a command, which promotes loose coupling between objects. This makes it easier to modify, extend or replace the components of a system.
2. Supports undo and redo: Since the Command Pattern encapsulates a request as an object, it makes it easier to implement undo and redo operations. Each command object can store the state before and after the execution of the command, which enables undo and redo operations.
3. Encapsulates complexity: The Command Pattern encapsulates the details of how a request is processed into a separate object, which simplifies the client code. This also enables the use of composite commands, which are made up of multiple commands.
4. Supports logging and auditing: Since the Command Pattern represents a request as an object, it enables logging and auditing of the commands. This allows for better tracking of system activities and easier debugging of issues.
Cons:
1. Increases code complexity: The Command Pattern can increase the complexity of the code, as it requires the creation of multiple objects for each command. This can lead to increased memory usage and slower performance.
2. Requires careful design: The Command Pattern requires careful design to ensure that the command objects are properly encapsulated and that the system is not overly complex. This can make it harder to understand and modify the system.
3. Can be overused: The Command Pattern can be overused, leading to unnecessary complexity and reduced performance. It is important to carefully consider whether the pattern is necessary for a particular use case.

*/

import (
	"fmt"
)

// Command interface
type Command interface {
	Execute()
}

// Light - Receiver
type Light struct {
	isOn bool
}

// TurnOn .
func (l *Light) TurnOn() {
	l.isOn = true
	fmt.Println("Light turned on")
}

// TurnOff .
func (l *Light) TurnOff() {
	l.isOn = false
	fmt.Println("Light turned off")
}

// TurnOnLightCommand - Concrete Command
type TurnOnLightCommand struct {
	light *Light
}

// Execute .
func (c *TurnOnLightCommand) Execute() {
	c.light.TurnOn()
}

// TurnOffLightCommand - Concrete Command
type TurnOffLightCommand struct {
	light *Light
}

// Execute .
func (c *TurnOffLightCommand) Execute() {
	c.light.TurnOff()
}

// Invoker
type Switch struct {
	onCommand  Command
	offCommand Command
}

func (s *Switch) SetOnCommand(c Command) {
	s.onCommand = c
}

func (s *Switch) SetOffCommand(c Command) {
	s.offCommand = c
}

func (s *Switch) On() {
	s.onCommand.Execute()
}

func (s *Switch) Off() {
	s.offCommand.Execute()
}

// Client code
func main() {
	light := &Light{}
	turnOnCommand := &TurnOnLightCommand{light: light}
	turnOffCommand := &TurnOffLightCommand{light: light}
	switcher := &Switch{}
	switcher.SetOnCommand(turnOnCommand)
	switcher.SetOffCommand(turnOffCommand)

	// Turn on the light
	switcher.On()

	// Turn off the light
	switcher.Off()
}
