package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

The Chain-of-responsibility pattern is a behavioral design pattern that allows multiple objects to handle a request in a chain-like fashion, where each object in the chain can either handle the request or pass it on to the next object in the chain.
Here are some pros and cons of using the Chain-of-responsibility pattern:
Pros:
1. Decouples sender and receiver: The pattern decouples the sender of the request from its receiver, which means that the sender does not need to know who will handle the request.
2. Flexible and extensible: The pattern allows you to add or remove handlers dynamically at runtime without affecting the other parts of the system.
3. Promotes loose coupling: Since each handler only knows about its immediate successor, there is a low level of coupling between the handlers. This makes the system more flexible and easier to maintain.
4. Simplifies object creation: The pattern simplifies object creation by eliminating the need to create complex chains of objects and also allows for the creation of different chains for different use cases.
Cons:
1. Can cause performance issues: The pattern can cause performance issues if the chain becomes too long or if each handler takes a lot of time to process the request.
2. Can lead to spaghetti code: If not implemented carefully, the pattern can lead to a tangled mess of code, making it difficult to debug and maintain.
3. Can lead to overuse: The pattern should be used judiciously, as it can be overused, leading to unnecessary complexity in the code.
4. Can be difficult to test: The pattern can be difficult to test since it involves a chain of objects, and it can be hard to isolate and test individual handlers.
Overall, the Chain-of-responsibility pattern can be a powerful tool in designing flexible and extensible systems, but it should be used with care and attention to avoid potential issues.
*/

import (
	"fmt"
)

type Handler interface {
	SetNext(handler Handler)
	Handle(request string) string
}

type BaseHandler struct {
	nextHandler Handler
}

func (h *BaseHandler) SetNext(handler Handler) {
	h.nextHandler = handler
}

func (h *BaseHandler) Handle(request string) string {
	if h.nextHandler != nil {
		return h.nextHandler.Handle(request)
	}
	return "No handler found"
}

type ConcreteHandler1 struct {
	BaseHandler
}

func (h *ConcreteHandler1) Handle(request string) string {
	if request == "handler1" {
		return "Handled by Handler 1"
	}
	return h.BaseHandler.Handle(request)
}

type ConcreteHandler2 struct {
	BaseHandler
}

func (h *ConcreteHandler2) Handle(request string) string {
	if request == "handler2" {
		return "Handled by Handler 2"
	}
	return h.BaseHandler.Handle(request)
}

func main() {
	handler1 := &ConcreteHandler1{}
	handler2 := &ConcreteHandler2{}

	handler1.SetNext(handler2)

	fmt.Println(handler1.Handle("handler1")) // Output: Handled by Handler 1
	fmt.Println(handler1.Handle("handler2")) // Output: Handled by Handler 2
	fmt.Println(handler1.Handle("handler3")) // Output: No handler found
}
