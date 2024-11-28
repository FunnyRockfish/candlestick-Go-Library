package internal

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

// compareFiles сравнивает содержимое двух файлов по их путям.
func compareFiles(path1, path2 string) bool {
	file1, err := ioutil.ReadFile(path1)

	if err != nil {
		log.Fatal(err) // Завершает выполнение программы, если произошла ошибка чтения первого файла.
	}

	file2, err := ioutil.ReadFile(path2)

	if err != nil {
		log.Fatal(err) // Завершает выполнение программы, если произошла ошибка чтения второго файла.
	}

	return bytes.Equal(file1, file2) // Возвращает true, если содержимое файлов совпадает.
}

// TestImage сравнивает тестовый файл <NAME>.png с эталонным файлом <NAME>_golden.png
// и вызывает ошибку теста, если их содержимое не совпадает.
func TestImage(t *testing.T, testFile string) {
	if !compareFiles(testFile, strings.TrimSuffix(testFile, filepath.Ext(testFile))+"_golden.png") {
		t.Errorf("несоответствие изображения для %s\n", testFile)
	}
}
