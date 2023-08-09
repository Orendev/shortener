package pkg1

import (
	"fmt"
	"os"
)

func osExitCheckFunc() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	fmt.Println("Hello world!")
	os.Exit(1) // want "direct use of os.Exit found"
}
