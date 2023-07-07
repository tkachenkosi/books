package ini

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var sec map[string]string

// парсим только одну секцию
func Parser(section string, keys ...string) bool {
	var (
		ok      bool
		secName = "[" + section + "]"
		key     string
	)

	if len(keys) > 0 {
		// keys []string - не обязательный параметр когда нужно прочитать один ключ
		key = strings.Join(keys, "")
	}

	// удалить старые записи
	DelVals()
	sec = make(map[string]string)

	f, err := os.Open("app.ini")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		l := strings.TrimSpace(s.Text())

		if len(l) == 0 || strings.HasPrefix(l, "#") {
			// удалим коментарии
			continue
		}

		if !ok && strings.Contains(l, secName) {
			// нашли нужную секцию
			ok = true
			continue
		}

		if ok && strings.HasPrefix(l, "[") && strings.HasSuffix(l, "]") {
			// нашли другую секцию
			ok = false
			continue
		}

		if ok && strings.Contains(l, "=") {
			fl := strings.Split(l, "=")
			fl[0] = strings.TrimSpace(fl[0])
			fl[1] = strings.TrimSpace(fl[1])

			// считываем один параметр или всю секцию
			if len(fl[0]) > 0 && len(fl[1]) > 0 && (key == "" || key == fl[0]) {
				sec[fl[0]] = fl[1]
			}
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	f.Close()

	return len(sec) > 0
}

// получить значение когда в секции одна запись
func Get() (val string) {
	if len(sec) == 1 {
		for k := range sec {
			val = sec[k]
			delete(sec, k)
		}
	}
	return
}

// возвращаем значение строкой
func Val(key string) (val string) {
	val = sec[key]
	delete(sec, key)
	return
}

func Int(key string) (val int) {
	val, err := strconv.Atoi(sec[key])
	delete(sec, key)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Bool(key string) bool {
	k := strings.ToLower(sec[key])
	delete(sec, key)
	if k == "true" || k == "1" {
		return true
	} else {
		return false
	}
}

// удалить все значения
func DelVals() {
	for _, k := range sec {
		delete(sec, k)
	}
}
