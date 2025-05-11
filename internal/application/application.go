package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/enigmasterr/final_project/pkg/calculation"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

type Config struct {
	Addr string
}

var PORT string

func ConfigFromEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	PORT = config.Addr
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

// Функция запуска приложения
// тут будем чиать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение

// func (a *Application) Run() error {
// 	for {
// 		// читаем выражение для вычисления из командной строки
// 		log.Println("input expression")
// 		reader := bufio.NewReader(os.Stdin)
// 		text, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Println("failed to read expression from console")
// 		}
// 		// убираем пробелы, чтобы оставить только вычислемое выражение
// 		text = strings.TrimSpace(text)
// 		// выходим, если ввели команду "exit"
// 		if text == "exit" {
// 			log.Println("aplication was successfully closed")
// 			return nil
// 		}
// 		//вычисляем выражение
// 		result, err := calculation.Calc(text)
// 		if err != nil {
// 			log.Println(text, " calculation failed wit error: ", err)
// 		} else {
// 			log.Println(text, "=", result)
// 		}
// 	}
// }

type TaskF struct {
	ID             int     `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time int     `json:"operation_time"`
}
type Task struct {
	Task TaskF `json:"task"`
}

var allTasks = map[int]TaskF{}
var allresults = map[int]float64{}

func Calc(expression string, id int) (float64, error) {

	ans, err := calculation.Get_expression(expression)
	if err != nil {
		return 0, err
	}

	fmt.Printf("TEST! %v\n", ans)
	var stk []float64
	for _, v := range ans {
		if v == "+" || v == "-" || v == "*" || v == "/" {
			if len(stk) < 2 {
				return 0, calculation.ErrInvalidExpression
			}
			a := stk[len(stk)-1]
			stk = stk[:len(stk)-1]
			b := stk[len(stk)-1]
			stk = stk[:len(stk)-1]
			if v == "+" {
				task := TaskF{ID: id, Arg1: b, Arg2: a, Operation: "+", Operation_time: 1}
				allTasks[id] = task
				//stk = append(stk, b+a) // нужно отправить таск на "+"
			} else if v == "-" {
				task := TaskF{ID: id, Arg1: b, Arg2: a, Operation: "-", Operation_time: 1}
				allTasks[id] = task
				//stk = append(stk, b-a) // нужно отправить таск на "-"
			} else if v == "*" {
				task := TaskF{ID: id, Arg1: b, Arg2: a, Operation: "*", Operation_time: 1}
				allTasks[id] = task
				//stk = append(stk, b*a) // нужно отправить таск на "*"
			} else if v == "/" {
				if a == 0 {
					return 0, calculation.ErrDivisionByZero
				}
				task := TaskF{ID: id, Arg1: b, Arg2: a, Operation: "/", Operation_time: 1}
				allTasks[id] = task
				//stk = append(stk, b/a) // нужно отправить таск на "/"
			}
			for {
				addr := fmt.Sprintf("http://localhost:%v/internal/getresult/%d", PORT, id)
				resp, err := http.Get(addr)
				//fmt.Println(resp)
				if err != nil {
					fmt.Errorf("Some trouble with getting answer")
				}
				if resp.StatusCode == http.StatusOK {
					type resJSON struct {
						ID     int     `json:"id"`
						Result float64 `json:"result"`
					}
					var res resJSON
					err = json.NewDecoder(resp.Body).Decode(&res)
					//fmt.Println(res)
					if err != nil {
						return 0, err
					}
					stk = append(stk, res.Result)
					delete(allresults, res.ID)
					break
				}
				time.Sleep(2 * time.Second)
			}
			// надо получить ответы и закинуть в стек
			// stk = append(stk, res)
		} else {
			num, _ := strconv.ParseFloat(v, 64)
			stk = append(stk, num)
		}
	}
	if len(stk) != 1 {
		return 0, calculation.ErrInvalidExpression
	}
	return stk[0], nil
}

type Request struct {
	Expression string `json:"expression"`
}
type expressionJSON struct {
	ID     int     `json:"id"`
	Status int     `json:"status"`
	Result float64 `json:"result"`
}
type Expressions struct {
	Expressions []expressionJSON `json:"expressions"`
}

var allExpressions Expressions
var curID int

func changeStatus(expr expressionJSON) {
	for i := 0; i < len(allExpressions.Expressions); i++ {
		if allExpressions.Expressions[i].ID == expr.ID {
			allExpressions.Expressions[i].Status = expr.Status
		}
	}
}

func addAnswer(expr expressionJSON) {
	for i := 0; i < len(allExpressions.Expressions); i++ {
		if allExpressions.Expressions[i].ID == expr.ID {
			allExpressions.Expressions[i].Status = expr.Status
			allExpressions.Expressions[i].Result = expr.Result
		}
	}
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	var mu sync.Mutex
	mu.Lock()
	curID++
	mu.Unlock()
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	log.Println("get request - ", request)

	type AnsJSON struct {
		ID int `json:"id"`
	}

	if err != nil {
		newExpres := expressionJSON{ID: curID, Status: http.StatusBadRequest, Result: 0}
		allExpressions.Expressions = append(allExpressions.Expressions, newExpres)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AnsJSON{ID: curID})
		return
	} else { // если само выражение получено не важно какое, то добавим в map со всеми выражениями AllExpressions
		newExpres := expressionJSON{ID: curID, Status: http.StatusCreated, Result: 0}
		allExpressions.Expressions = append(allExpressions.Expressions, newExpres)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(AnsJSON{ID: curID})
	}

	result, err := Calc(request.Expression, curID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, calculation.ErrInvalidExpression) {
			newExpres := expressionJSON{ID: curID, Status: http.StatusBadRequest, Result: 0}
			changeStatus(newExpres)
			log.Printf("err: %s", err.Error())
		} else if errors.Is(err, calculation.ErrStrangeSymbols) {
			newExpres := expressionJSON{ID: curID, Status: http.StatusUnprocessableEntity, Result: 0}
			changeStatus(newExpres)
			log.Printf("err: %s", err.Error())
		} else {
			newExpres := expressionJSON{ID: curID, Status: http.StatusInternalServerError, Result: 0}
			changeStatus(newExpres)
			log.Printf("err: %s", err.Error())
		}

	} else {
		newExpres := expressionJSON{ID: curID, Status: http.StatusOK, Result: result}
		addAnswer(newExpres)
		log.Printf("send json {\"result\": \"%s\"}", string(fmt.Sprintf("%d", curID)))
	}
}

func ExprHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allExpressions)
	log.Printf("send JSON {\"expressions\": [{},{},{}...]}")
}

func ExprIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	type AnsJSON struct {
		Expression expressionJSON `json:"expression"`
	}
	found := false
	for _, expresn := range allExpressions.Expressions {
		if expresn.ID == id {
			found = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(AnsJSON{Expression: expresn})
			break
		}
	}
	// Не найден ID
	if !found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AnsJSON{Expression: expressionJSON{ID: id, Status: 404, Result: 0}})
	}
}

func TaskHandlerGET(w http.ResponseWriter, r *http.Request) {
	var task TaskF
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	if len(allTasks) > 0 { // в этом блоке у нас есть задача
		for _, value := range allTasks {
			task = value
			break // Выходим из цикла после первого элемента
		}
		// удаляем задачу
		delete(allTasks, task.ID)
		// передаем задачу агенту
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Данные отправлены агенту: %+v\n", task)
	}
}

func TaskHandlerPOST(w http.ResponseWriter, r *http.Request) {
	type taskAns struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}

	var data taskAns
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Данные от агента получены: %+v\n", data)
	w.WriteHeader(http.StatusOK)
	allresults[data.ID] = data.Result
}

func GetResultOperation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	type resJSON struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}
	if res, found := allresults[id]; found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resJSON{ID: id, Result: res})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resJSON{})
	}
}

func test() {
	ctx := context.TODO()

	db, err := sql.Open("sqlite", "file:/internal/db/store.db?cache=shared")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
}

func (a *Application) RunServer() error {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/calculate", CalcHandler).Methods("GET", "POST")
	router.HandleFunc("/api/v1/expressions", ExprHandler).Methods("GET")
	router.HandleFunc("/api/v1/expressions/:{id}", ExprIDHandler).Methods("GET")
	router.HandleFunc("/internal/task", TaskHandlerGET).Methods("GET")
	router.HandleFunc("/internal/task", TaskHandlerPOST).Methods("POST")
	router.HandleFunc("/internal/getresult/{id}", GetResultOperation).Methods("GET")
	fmt.Println("Calculator is ready!")
	test()
	return http.ListenAndServe(":"+a.config.Addr, router)
}
