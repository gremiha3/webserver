package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"
)

type srv struct { //структура хранения данных
	mu    *sync.RWMutex   //защищаем работу с данными мьютексом из=за многопоточного режима работы
	stats map[uint64]uint //храним данные вида "ID кандидата"/"количество голосов за этого кандидата"
}

func (s *srv) Vote(w http.ResponseWriter, r *http.Request) { //метод голосования
	if r.Method != http.MethodPost { //если метод был не POST, то
		w.WriteHeader(http.StatusMethodNotAllowed) //отдаем ошибку 405
		return
	}
	//пример краткого объявления переменной:
	req := struct { //структура для размаршиливания входных данных json
		CandidateID uint64 `json:"candidate_id"` //поле json будет корректно мапиться в поле структуры
		Passport    string `json:"passport"`
	}{}

	raw, err := io.ReadAll(r.Body) //вычитываем тело запроса
	if err != nil {                //если при вычитывании будет ошибка, то
		w.WriteHeader(http.StatusInternalServerError) //отправляем ошибку 500
		return
	}
	//краткая запись условия:
	if err := json.Unmarshal(raw, &req); err != nil { //если при размаршаливании произойдет ошибка, то
		w.WriteHeader(http.StatusInternalServerError) //отправляем ошибку 500
		return
	}
	if len(req.Passport) == 0 { //если нет данных паспорта, то
		w.WriteHeader(http.StatusBadRequest) //отправляем ошибку 400
		return
	}
	if len(req.Passport) == 0 || req.CandidateID == 0 { //если нет данных паспорта или кандидата, то
		w.WriteHeader(http.StatusBadRequest) //отправляем ошибку 400
		return
	}
	s.mu.Lock()                //лочим мапу для работы с ней
	s.stats[req.CandidateID]++ //добавляем один голос за кандидата
	s.mu.Unlock()              //разлочиваем мапу для работы других горутин

	w.WriteHeader(http.StatusOK) //посылаем статус успешного завершения операции
}
func (s *srv) Stats(w http.ResponseWriter, r *http.Request) { //метод статистика
	const candID = "candidate_id"
	if r.Method != http.MethodGet { //если метод был не GET, то
		w.WriteHeader(http.StatusMethodNotAllowed) //отдаем ошибку 405
		return
	}

	vals := r.URL.Query() //вычитываем параметры из запроса
	if vals.Has(candID) { //если в запросе есть выбор по кандидату, то
		id, err := strconv.ParseUint(vals.Get(candID), 10, 64) //вычитываем ID кандидата и конвертим его в INT
		if err != nil {                                        //если при получении ID кандидата получаем ошибку, то
			w.WriteHeader(http.StatusInternalServerError) //отправляем ошибку 500
			return
		}
		s.mu.RLock()             //лочим переменную на чтение
		candStats := s.stats[id] //вычитываем данные во внутреннюю переменную для уменьшения времени работы мьютекса
		s.mu.RUnlock()           //анлочим переменную для работы с другими горутинами

		raw, err := json.Marshal(candStats) //замаршаливаем данные
		if err != nil {                     //если при замаршаливании произойдет ошибка, то
			w.WriteHeader(http.StatusInternalServerError) //отправляем ошибку 500
			return
		}
		w.Write(raw) //отправляем данные в ответ. Заголовок уже включен в эту функцию
		return
	}

	s.mu.RLock()     //лочим переменную на чтение
	stats := s.stats //вычитываем данные во внутреннюю переменную для уменьшения времени работы мьютекса
	s.mu.RUnlock()   //анлочим переменную для работы с другими горутинами

	raw, err := json.Marshal(stats) //замаршаливаем данные
	if err != nil {                 //если при замаршаливании произойдет ошибка, то
		w.WriteHeader(http.StatusInternalServerError) //отправляем ошибку 500
		return
	}

	w.Write(raw) //отправляем данные в ответ. Заголовок уже включен в эту функцию
}

func New() srv { //конструктор
	return srv{ //инициализируем экземпляр структуры
		mu:    &sync.RWMutex{},
		stats: make(map[uint64]uint),
	}
}
