package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//Терминальная утилита RBS-EX2.3 используется для анализа размера содержимого для указанной директории.

func main() {
	start := time.Now()
	var root = flag.String("root", "", "path")
	flag.Parse()
	duration := time.Since(start)
	var filesOfDir []string
	files, err := os.ReadDir(*root)
	if err != nil {
		fmt.Printf("Ошибка чтения директории %e", err)
	}
	path, err := filepath.Abs(*root)
	if err != nil {
		fmt.Printf("Ошибка назначенного пути %e", err)
	}
	filepath.Abs(*root)
	for _, file := range files {
		filesOfDir = append(filesOfDir, filepath.Join(path, file.Name()))
	}
	var wg sync.WaitGroup
	wg.Add(len(filesOfDir))
	for _, dirEntered := range filesOfDir {

		go func(dirEntered string) {
			defer wg.Done()
			dirSize(dirEntered)
		}(dirEntered)
	}
	wg.Wait()
	fmt.Println(duration)
}

/*
dirSize() функция принимает путь к директории, определяет тип содержимого (файл или папка)
и возвращает размер содержимого для файла либо сумму размеров содержимого для папки
*/
func dirSize(path string) float64 {
	sizes := make(chan int64)
	// booler := make(chan bool)
	size := int64(0)
	var resultSize float64
	readSize := func(path string, file os.FileInfo, err error) error {
		if err != nil || file == nil {
			return err
		}
		if !file.IsDir() {
			sizes <- file.Size()
		}
		// if file.IsDir() {
		// 	size = 0
		// }
		return err
	}

	go func() {
		filepath.Walk(path, readSize)
		close(sizes)
	}()

	for s := range sizes {
		size += s
	}

	resultSize = float64(size)
	kb := math.Pow(1024, 1)
	mb := math.Pow(1024, 2)
	gb := math.Pow(1024, 3)
	tb := math.Pow(1024, 4)

	switch {
	case resultSize > (kb) && resultSize <= (mb):
		resultSize = float64(size) / kb
		fmt.Printf("%s\tРазмер %.2f Кб\n", path, resultSize)
	case resultSize > (mb) && resultSize <= (gb):
		resultSize = float64(size) / mb
		fmt.Printf("%s\tРазмер %.2f Мб\n", path, resultSize)
	case resultSize > (gb) && resultSize <= (tb):
		resultSize = float64(size) / gb
		fmt.Printf("%s\tРазмер %.2f Гб\n", path, resultSize)
	case resultSize > (tb):
		resultSize = float64(size) / tb
		fmt.Printf("%s\tРазмер %.2f Тб\n", path, resultSize)
	case resultSize == 0:
		fmt.Printf("%s\t Папка пуста %.f байт\n", path, resultSize)
	}

	return resultSize

}
