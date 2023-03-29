package pattern

import "fmt"

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern

The strategy pattern is a behavioral design pattern that allows you to define a family of algorithms, encapsulate each one as an object, and make them interchangeable at runtime. Here are some of the pros and cons of using the strategy pattern:
Pros:
1. Open/Closed Principle: The strategy pattern is an implementation of the Open/Closed Principle, which means that you can add new strategies to the system without modifying existing code. This makes the code more maintainable and extensible.
2. Flexibility: The strategy pattern provides a flexible way to change the behavior of an object at runtime. This allows you to vary the algorithm being used based on the context.
3. Encapsulation: The strategy pattern encapsulates the algorithm and makes it easy to swap it out for a different one without affecting the rest of the code.
4. Testability: The strategy pattern makes it easy to test individual strategies in isolation, which can help improve the quality and reliability of the code.
Cons:
1. Increased Complexity: The strategy pattern can add some complexity to the codebase, especially if the number of strategies is large. It can also increase the number of classes and interfaces in the system.
2. Runtime Overhead: The strategy pattern involves creating and managing objects at runtime, which can result in some runtime overhead.
3. Increased Memory Usage: Using the strategy pattern can lead to increased memory usage because it involves creating multiple objects.
4. Requires Careful Design: The strategy pattern requires careful design to ensure that the strategies are interchangeable and that the system can handle changes to the strategy interface without breaking existing code.
*/

// Define the strategy interface
type SortingStrategy interface {
	Sort([]int) []int
}

// Define the bubble sort strategy
type BubbleSort struct{}

func (bs *BubbleSort) Sort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	return arr
}

// Define the quick sort strategy
type QuickSort struct{}

func (qs *QuickSort) Sort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	pivot := arr[0]
	var left, right []int
	for _, v := range arr[1:] {
		if v < pivot {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}
	left = qs.Sort(left)
	right = qs.Sort(right)
	return append(append(left, pivot), right...)
}

// Define the context struct that uses the strategy
type Sorter struct {
	strategy SortingStrategy
}

func (s *Sorter) SetStrategy(strategy SortingStrategy) {
	s.strategy = strategy
}

func (s *Sorter) Sort(arr []int) []int {
	if s.strategy == nil {
		panic("sorting strategy not set")
	}
	return s.strategy.Sort(arr)
}

// Example usage
func main() {
	arr := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	sorter := &Sorter{}

	// Use bubble sort
	sorter.SetStrategy(&BubbleSort{})
	fmt.Println(sorter.Sort(arr)) // Output: [1 1 2 3 3 4 5 5 5 6 9]

	// Use quick sort
	sorter.SetStrategy(&QuickSort{})
	fmt.Println(sorter.Sort(arr)) // Output: [1 1 2 3 3 4 5 5 5 6 9]
}
