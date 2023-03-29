package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern

The Visitor pattern is a design pattern used in object-oriented programming that allows adding new operations to a class hierarchy without modifying the classes themselves. The pattern involves defining a visitor interface with methods that correspond to the classes in the hierarchy, and then implementing the visitor interface to perform the desired operations on each class.
Here are some pros and cons of using the Visitor pattern:
Pros:
1. Adds new operations without modifying existing classes: With the Visitor pattern, you can add new operations to a class hierarchy without changing the existing classes. This makes the code more flexible and easier to maintain.
2. Separates behavior from objects: The Visitor pattern separates the behavior of an object from its data structure, allowing you to change or add new behaviors without affecting the structure of the object.
3. Improves extensibility: The Visitor pattern improves the extensibility of the code because new operations can be added by implementing new visitor classes.
4. Enables double dispatch: The Visitor pattern enables double dispatch, which allows you to call a method based on both the runtime type of an object and the type of a parameter.
Cons:
1. Increases complexity: The Visitor pattern can make the code more complex because it requires additional classes and interfaces.
2. Violates encapsulation: The Visitor pattern can violate encapsulation because it requires exposing the internal structure of the visited classes to the visitor.
3. Can be less efficient: The Visitor pattern can be less efficient than other approaches because it involves additional method calls and object creations.
4. Requires careful design: The Visitor pattern requires careful design to ensure that it is implemented correctly and to avoid creating a tight coupling between the visitor and the visited classes.
Overall, the Visitor pattern can be useful in situations where you need to add new operations to a class hierarchy without modifying the existing classes. However, it should be used with caution and only after careful consideration of its pros and cons.
*/

import "fmt"

// Element .
type Element interface {
	Accept(visitor Visitor)
}

// Visitor .
type Visitor interface {
	VisitConcreteElementA(element ConcreteElementA)
	VisitConcreteElementB(element ConcreteElementB)
}

// ConcreteElementA .
type ConcreteElementA struct{}

// Accept .
func (e ConcreteElementA) Accept(visitor Visitor) {
	visitor.VisitConcreteElementA(e)
}

// ConcreteElementB .
type ConcreteElementB struct{}

// Accept .
func (e ConcreteElementB) Accept(visitor Visitor) {
	visitor.VisitConcreteElementB(e)
}

// ConcreteVisitor .
type ConcreteVisitor struct{}

// VisitConcreteElementA .
func (v ConcreteVisitor) VisitConcreteElementA(element ConcreteElementA) {
	fmt.Println("Visited ConcreteElementA")
}

// VisitConcreteElementB .
func (v ConcreteVisitor) VisitConcreteElementB(element ConcreteElementB) {
	fmt.Println("Visited ConcreteElementB")
}

func main() {
	elementA := ConcreteElementA{}
	elementB := ConcreteElementB{}
	visitor := ConcreteVisitor{}

	elementA.Accept(visitor)
	elementB.Accept(visitor)
}
