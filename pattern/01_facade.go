package pattern

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern

The facade pattern is a software design pattern that provides a simple interface to a complex subsystem or set of interfaces. Here are some of the pros and cons of using the facade pattern in software development:
Pros:
1. Simplifies the interface: The facade pattern simplifies the interface between the client and the subsystem by providing a single, simplified interface to the client. This makes it easier for clients to use the subsystem.
2. Encapsulates complexity: The facade pattern encapsulates the complexity of the subsystem behind a simple interface. This makes it easier to change the subsystem without affecting the clients that use it.
3. Provides abstraction: The facade pattern provides a level of abstraction between the client and the subsystem. This means that the client does not need to know the details of how the subsystem works.
4. Improves maintainability: The facade pattern improves the maintainability of the code by reducing the coupling between the client and the subsystem. This makes it easier to modify the subsystem without affecting the client.
5. Improves testability: The facade pattern makes it easier to test the subsystem because the client does not need to know the details of how the subsystem works.
Cons:
1. Can limit flexibility: The facade pattern can limit the flexibility of the subsystem by providing a fixed interface. This can make it difficult to add new functionality to the subsystem without changing the facade.
2. Can add complexity: The facade pattern can add complexity to the code by introducing an additional layer of abstraction.
3. Can reduce performance: The facade pattern can reduce performance by introducing an additional layer of abstraction.
4. Can violate the single responsibility principle: The facade pattern can violate the single responsibility principle by providing a single interface to a complex subsystem. This can make it difficult to maintain and modify the code.
5. Can create unnecessary dependencies: The facade pattern can create unnecessary dependencies between the client and the subsystem by providing a fixed interface. This can make it difficult to modify the subsystem without affecting the client.

*/

import (
	"fmt"
)

// CPU .
type CPU struct{}

// Freeze .
func (*CPU) Freeze() {
	fmt.Println("CPU Freeze")
}

// Jump .
func (*CPU) Jump(position int) {
	fmt.Printf("CPU Jump to %d\n", position)
}

// Execute .
func (*CPU) Execute() {
	fmt.Println("CPU Execute")
}

// Memory .
type Memory struct{}

// Load .
func (*Memory) Load(position int, data string) {
	fmt.Printf("Memory Load data '%s' to position %d\n", data, position)
}

// HardDrive .
type HardDrive struct{}

// Read .
func (*HardDrive) Read(position int, size int) string {
	data := fmt.Sprintf("HardDrive Read data from position %d with size %d", position, size)
	fmt.Println(data)
	return data
}

// ComputerFacade .
type ComputerFacade struct {
	cpu       *CPU
	memory    *Memory
	hardDrive *HardDrive
}

// NewComputerFacade .
func NewComputerFacade() *ComputerFacade {
	return &ComputerFacade{
		cpu:       &CPU{},
		memory:    &Memory{},
		hardDrive: &HardDrive{},
	}
}

// Start .
func (c *ComputerFacade) Start() {
	c.cpu.Freeze()
	c.memory.Load(0, "boot_loader")
	c.cpu.Jump(0)
	c.cpu.Execute()
}

func main() {
	computer := NewComputerFacade()
	computer.Start()
}
