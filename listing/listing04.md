Что выведет программа? Объяснить вывод программы.

```go
package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		println(n)
	}
}
```

Ответ:
```
Программа выведет числа от 0 до 9 и упадет с паникой из-за deadlock`а
Это случится потому что range будет читать из канала в который после определенного момента 100% больше не придут данные и поэтому мы зависнем в бесконечном ожидании 
```
