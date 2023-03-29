package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern

The Builder pattern is a design pattern used to simplify the creation of complex objects. Here are some of the pros and cons of using the Builder pattern:
Pros:
1. Encapsulates the creation of objects: The Builder pattern encapsulates the process of object creation and separates it from the object's representation. This makes the code more maintainable and easier to read.
2. Improves code readability: The use of a Builder class can make the code more readable and self-explanatory, especially when dealing with complex object creation.
3. Supports fluent interface: The Builder pattern can be designed to support a fluent interface, which allows for a more natural and readable code style.
4. Provides flexibility: The Builder pattern provides flexibility in object creation by allowing the user to customize the object creation process. The user can decide which properties to set and which ones to leave with their default values.
5. Reduces dependency on constructors: Using the Builder pattern reduces the dependency on constructors, which can make it easier to change the object creation process in the future.
Cons:
1. Increased complexity: The use of the Builder pattern can increase the complexity of the code, especially when dealing with a large number of properties to set.
2. Extra code: Implementing the Builder pattern requires writing additional code, which can increase the code size.
3. Can be overkill for simple objects: Using the Builder pattern for simple objects can be overkill and can add unnecessary complexity to the code.
4. Increases memory usage: The use of the Builder pattern can increase memory usage as it involves creating additional objects.
5. Can be time-consuming: Implementing the Builder pattern can be time-consuming, especially for large and complex objects.
*/

// Pizza .
type Pizza struct {
	dough    string
	sauce    string
	cheese   string
	toppings []string
}

// PizzaBuilder .
type PizzaBuilder interface {
	SetDough(string) PizzaBuilder
	SetSauce(string) PizzaBuilder
	SetCheese(string) PizzaBuilder
	SetToppings([]string) PizzaBuilder
	Build() *Pizza
}

// ConcretePizzaBuilder .
type ConcretePizzaBuilder struct {
	pizza *Pizza
}

// NewConcretePizzaBuilder .
func NewConcretePizzaBuilder() *ConcretePizzaBuilder {
	return &ConcretePizzaBuilder{pizza: &Pizza{}}
}

// SetDough .
func (cpb *ConcretePizzaBuilder) SetDough(dough string) PizzaBuilder {
	cpb.pizza.dough = dough
	return cpb
}

// SetSauce .
func (cpb *ConcretePizzaBuilder) SetSauce(sauce string) PizzaBuilder {
	cpb.pizza.sauce = sauce
	return cpb
}

// SetCheese .
func (cpb *ConcretePizzaBuilder) SetCheese(cheese string) PizzaBuilder {
	cpb.pizza.cheese = cheese
	return cpb
}

// SetToppings .
func (cpb *ConcretePizzaBuilder) SetToppings(toppings []string) PizzaBuilder {
	cpb.pizza.toppings = toppings
	return cpb
}

// Build .
func (cpb *ConcretePizzaBuilder) Build() *Pizza {
	return cpb.pizza
}

// Director .
type Director struct {
	builder PizzaBuilder
}

// NewDirector .
func NewDirector(builder PizzaBuilder) *Director {
	return &Director{builder: builder}
}

// Construct .
func (d *Director) Construct() *Pizza {
	return d.builder.SetDough("Thin Crust").SetSauce("Tomato").SetCheese("Mozzarella").SetToppings([]string{"Mushrooms", "Olives", "Onions"}).Build()
}

func main() {
	builder := NewConcretePizzaBuilder()
	director := NewDirector(builder)
	_ = director.Construct()

}
