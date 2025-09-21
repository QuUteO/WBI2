package main

import (
	"fmt"
	"os"
)

/*
Функция Foo возвращает интерфейс error.
Интерфейсный тип в Go хранит внутри себя как значение конкретного типа, так и конкретное значение этого типа.
err в функции Foo имеет значение значение конкретного типа *os.ParhError, а значение этого типа равно nil
*/
func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)        // выведится nil, тк значение конкретного типа равно nil
	fmt.Println(err == nil) // выведится false, потому что конкретный типа не равен nil
	// true выведится только в том случае, если и значение и тип будут равны nil
}
