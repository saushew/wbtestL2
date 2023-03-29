package pattern

import (
	"fmt"
)

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern

The State pattern is a design pattern used in software engineering to model the state of an object, and to allow it to change its behavior based on its internal state. Here are some pros and cons of using the State pattern:
Pros:
1. Flexibility: The State pattern allows you to change the behavior of an object at runtime by changing its state, which provides a high degree of flexibility and adaptability to changing requirements.
2. Encapsulation: The State pattern encapsulates the behavior of an object in a separate class, which makes it easier to modify and maintain the code, and reduces the likelihood of errors.
3. Reusability: The State pattern promotes code reuse, since the same state objects can be used across multiple objects.
4. Readability: The State pattern makes code more readable and understandable by separating out different states into their own classes.
Cons:
1. Complexity: The State pattern can add complexity to the code, especially if there are many states and state transitions.
2. Overhead: The State pattern can introduce overhead, since it involves creating separate classes for each state.
3. Coupling: The State pattern can increase coupling between classes, since the state objects need to be aware of the context in which they are used.
4. Performance: The State pattern can affect performance, since it involves more objects and method calls than a simpler implementation.
In summary, the State pattern can be a powerful tool for modeling the behavior of an object that changes over time, but it comes with some trade-offs in terms of complexity, overhead, coupling, and performance.
*/

type State interface {
	Handle()
}

type Context struct {
	state State
}

func (c *Context) SetState(state State) {
	c.state = state
}

func (c *Context) Request() {
	c.state.Handle()
}

type ConcreteStateA struct{}

func (s *ConcreteStateA) Handle() {
	fmt.Println("Handling request in state A.")
}

type ConcreteStateB struct{}

func (s *ConcreteStateB) Handle() {
	fmt.Println("Handling request in state B.")
}

func main() {
	context := &Context{}

	stateA := &ConcreteStateA{}
	stateB := &ConcreteStateB{}

	context.SetState(stateA)
	context.Request()

	context.SetState(stateB)
	context.Request()
}
