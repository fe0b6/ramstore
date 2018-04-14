package ramstore

import (
	"log"
	"sync"
	"time"
)

const (
	mapCount       = 29
	deletedTimeout = 7 * 24
)

var (
	data   []stor
	exited bool
	wg     sync.WaitGroup
)

// InitStore - Инициализация хранилища
func InitStore() (exitChan chan bool) {
	data = make([]stor, mapCount)
	for i := 0; i < mapCount; i++ {
		data[i] = stor{Data: make(map[string]Obj)}
	}

	// Канал для оповещения о выходе
	exitChan = make(chan bool)

	go waitExit(exitChan)

	readData()

	return
}

// Ждем сигнал о выходе
func waitExit(exitChan chan bool) {
	_ = <-exitChan

	exited = true

	log.Println("[info]", "Завершаем работу store ramkv")

	// Ждем пока все запросы завершатся
	wg.Wait()

	saveData()

	log.Println("[info]", "Работа store ramkv завершена корректно")
	exitChan <- true
}

// Set - Добавление хэша в хранилище
func Set(key string, obj Obj) (err string) {
	if exited {
		return "store down"
	}

	if key == "" {
		return "specify key"
	}

	if obj.Deleted {
		if time.Now().After(time.Unix(0, obj.Time).Add(deletedTimeout * time.Hour)) {
			return "delete timeout"
		}
	}

	wg.Add(1)

	num := getArrNum(key)

	data[num].Lock()

	v, ok := data[num].Data[key]
	if ok {
		if v.Time >= obj.Time {
			data[num].Unlock()
			wg.Done()
			return "wrong time"
		}
	}

	data[num].Data[key] = obj
	data[num].Unlock()
	wg.Done()
	return
}

// Get - Получение хэша из хранилища
func Get(key string) (Obj, string) {
	num := getArrNum(key)

	data[num].RLock()

	obj, ok := data[num].Data[key]
	data[num].RUnlock()

	if !ok || obj.Deleted {
		return obj, "key not found"
	}

	return obj, ""
}

// Foreach - Перебираем все эелементы
func Foreach(f func(string, Obj)) {

	for i := range data {
		data[i].RLock()

		for k, v := range data[i].Data {
			f(k, v)
		}

		data[i].RUnlock()

		f("", Obj{})
	}
}

// getArrNum - получаем номер массива с объектом
func getArrNum(key string) int {
	var sum int
	for _, v := range key {
		sum += int(v)
	}

	return sum % mapCount
}
