package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// Вариант 4, Науменко К. А.

// Task 1
// Алгоритмы и структуры данных: Реализуйте алгоритм проверки, является ли
// одна строка анаграммой другой без использования встроенных функций
// сортировки. Опишите ваш алгоритм и приведите пример его работы для строк
// “listen” и “silent”.

func isAnagram(target, comp string) bool {
	// Если длины слов разные - они не являются анаграммами, дальше можно не смотреть
	if len(target) != len(comp) {
		return false
	}

	// Проходимся по слову, берем каждый отдельный символ (руну)
	for _, r := range target {
		// Если хотя бы одна буква из target не содержится в comp, то дальше можно не проверять, они не анаграммы
		if !strings.ContainsRune(comp, r) {
			return false
		}
	}
	// в противном случае - анаграммы
	return true
}

// Task 2
// Go Routines и каналы: Создайте программу, которая имитирует работу
// очереди задач. Горутина “производитель” генерирует случайные задачи
// (просто строки) и отправляет их в канал. Горутина “потребитель” получает
// задачи из канала и выполняет их (просто выводит на экран с пометкой
// “Выполнено”). Создайте несколько “потребителей” и убедитесь, что задачи
// распределяются между ними.

const (
	// Кол-во генерируемых тасок
	TASKS_COUNT = 30
	// Количество консьюмеров
	CONSUMERS_COUNT = 3
	// Длина генерируемого id таски
	TASK_ID_LEN = 5
	// Размер буфера канала для тасок
	BUFFER_SIZE = 10
)

// Рандомная строка длинной n
func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Рандомное число в диапазоне от n до m
func randomInt(n, m int) int {
	if n > m {
		n, m = m, n
	}

	return rand.Intn(m-n+1) + n
}

// Функция продюсер, которая генерирует задачи
func producer(tasks chan<- string, wg *sync.WaitGroup) {
	// Вызываем wg.Done(), чтобы декрементить счетчик wait группы (после выполнения цикла)
	defer wg.Done()
	for range TASKS_COUNT {
		// Генерируем рандомный таск id
		t := randomString(TASK_ID_LEN)
		fmt.Printf("Produced task with id %s\n", t)
		// Передаем его в канал
		tasks <- t
		// Имитируем деятельность
		time.Sleep(time.Millisecond * 100) // Задержка для имитации времени генерации задачи
	}
}

// Функция потребитель, которая выполняет задачи
func consumer(id int, tasks <-chan string, wg *sync.WaitGroup) {
	// Так же декрементим счетчик
	defer wg.Done()
	for task := range tasks {
		// "Выполняем" таску (просто печатаем её айдишник и ждем от 500 до 1500 мс)
		fmt.Printf("Consumer %d: Execute task %s\n", id, task)
		time.Sleep(time.Millisecond * time.Duration(randomInt(500, 1500))) // Задержка для имитации времени выполнения задачи
	}
}

func main() {
	// Task 1
	// слова для сравнения
	t, c := "listen", "silent"

	if isAnagram(t, c) {
		fmt.Printf("target word %s is anagram for %s\n", t, c)
	} else {
		fmt.Printf("target word %s is NOT anagram for %s\n", t, c)
	}

	// Task 2
	// wait группы для консьюмеров и продюсера
	var consumerWg sync.WaitGroup
	var producerWg sync.WaitGroup

	// канал для передачи тасок между горутинами с заданным буффером
	// буфер нужен, чтобы не блочить продюсер, без буфера продюсер будет
	// ждать, пока один из консьюмеров не получит переданное в канал значение
	tasks := make(chan string, BUFFER_SIZE)

	// инкрементим счетчик и запускаем продюсер
	producerWg.Add(1)
	go producer(tasks, &producerWg)

	// Запуск нескольких потребителей
	for i := range CONSUMERS_COUNT {
		consumerWg.Add(1)
		go consumer(i, tasks, &consumerWg)
	}

	// Ждем завершение работы продюсера, чтобы закрыть канал
	go func() {
		producerWg.Wait()
		close(tasks)
	}()

	// Ждем завершение работы консьюмеров, чтобы не завершить программу, пока таски ещё обрабатываются
	consumerWg.Wait()
}
