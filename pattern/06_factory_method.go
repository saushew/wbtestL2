package pattern

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern

The Factory Method pattern is a design pattern that allows an object to create instances of other classes, without the need for the object to know the exact class of the object it is creating. This pattern is used to abstract the process of creating objects, which can make code more flexible and easier to maintain. Here are some pros and cons of the Factory Method pattern:
Pros:
1. Encapsulation: The Factory Method pattern encapsulates the creation of objects within a separate object. This allows the object to focus on its primary responsibility, rather than the details of object creation.
2. Flexibility: The Factory Method pattern allows for flexible object creation, as it enables the creation of new objects without changing the code of the object that uses them.
3. Code reuse: The Factory Method pattern can be used to reuse existing code by providing a standard interface for object creation.
4. Extensibility: The Factory Method pattern can be extended to create new objects or modify the way existing objects are created. This makes it easy to add new functionality without breaking existing code.
Cons:
1. Increased complexity: The Factory Method pattern can add complexity to a system by introducing new objects and classes.
2. Performance overhead: The Factory Method pattern can introduce performance overhead due to the additional objects and classes involved in the process of object creation.
3. Dependency injection: The Factory Method pattern can be seen as an alternative to dependency injection. However, dependency injection is generally considered to be a more flexible and scalable approach to object creation.
4. Tight coupling: If the Factory Method pattern is not implemented correctly, it can lead to tight coupling between objects, which can make the code more difficult to maintain.
Overall, the Factory Method pattern is a useful tool for managing object creation, but it should be used judiciously to avoid introducing unnecessary complexity or performance overhead.
*/

import (
	"fmt"
)

// Define an interface for a product
type Product interface {
	GetName() string
}

// Define a concrete product type
type ConcreteProduct struct{}

func (c *ConcreteProduct) GetName() string {
	return "Concrete Product"
}

// Define a factory interface
type Factory interface {
	CreateProduct() Product
}

// Define a concrete factory type
type ConcreteFactory struct{}

func (c *ConcreteFactory) CreateProduct() Product {
	return &ConcreteProduct{}
}

func main() {
	// Create a factory object
	factory := &ConcreteFactory{}

	// Use the factory to create a product object
	product := factory.CreateProduct()

	// Call a method on the product object
	fmt.Println(product.GetName()) // Output: Concrete Product
}
