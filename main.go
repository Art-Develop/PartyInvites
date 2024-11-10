package main

import (
	"fmt"
	"html/template"
	"net/http"
)

/*
создаем структуру Rsvp которая имеет пользовательские значения
Name, Email, Phone это имена, а string и bool это тип данных
*/
type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

/*
создали переменную которая будет хранить срез указателей на структуру Rsvp
1. аргумент это тип данных, в нашем случае это указатель на структуру Rsvp
2. это длина среза - показывает есть ли уже элементы в срезе. Если он равен 0, то срез не содержит элементов
3. это емкость - он определяет сколько элементов можно хранить
*/
var responses = make([]*Rsvp, 0, 10)

/*
создаем переменную которая хранит тип данных карту(map)
string это тип ключа
*template.Template - это тип значения
3 - это емкость(размер) карты, т.е. карта вмещает до 3 элементов ключ-значение
*/
var templates = make(map[string]*template.Template, 3)

// создаем функцию для загрузки HTML-шаблонов и сохранять их в карте templates
func loadTemplates() {
	//создаем массив из 5 элементов типа string
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	//создан цикл for range для того чтобы загрузить шаблоны и пройтись по ним по очереди
	for index, name := range templateNames {
		/*
			функция template.ParseFiles загружает и анализирует HTML-шаблон
			"layout.html" - это основной файл макета, общий для всех страниц
			name+".html" - это файл конкретного шаблона
			Результат сохраняется в переменную t которая будет указателем на *tempate.Template.
			Почему так - это происходит из за того что функция template.ParseFiles возращает указатель на *tempate.Template.
			Если будет ошибка, то она запишется в переменную err
		*/
		t, err := template.ParseFiles("layout.html", name+".html")
		// проверяем есть ли ошибка
		if err == nil {
			//сохраняется загруженный шаблон t в карту templates, использую name в качестве ключа
			templates[name] = t
			//выводим сообщение, а также индекс и название
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

/*
создаем обработчик HTTP-запросов, где:
writer http.ResponseWriter - это объект, используемый для отправки ответа клиенту.
request *http.Request - это объект, представляющий запрос, сделанный клиентом
*/
func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	/*
		templates["welcome"] - обращаемся к карте templates и извлекаем шаблон "welcome"
		метод Execute рендерит(выполняет) шаблон и записывает его в writer
		writer - это поток данных для записи ответа, через него мы отправляем результат шаблона пользователю
		nil - это данные которые передает шаблон
	*/
	templates["welcome"].Execute(writer, nil)
}

func listHandler(writer http.ResponseWriter, request *http.Request) {
	templates["list"].Execute(writer, responses)
}

// создаем еще одну структуру ссылаясь через указатель на RSVP, а также если будут ошибки то будем записывать их в строковой срез
type formData struct {
	*Rsvp
	Errors []string
}

func formHandler(writer http.ResponseWriter, request *http.Request) {
	/*
		if request.Method == http.MethodGet проверяет, яв-ся ли запрос методом GET(обычно когда пользователь окрывает страницу)
		Если это GET запрос то тогда подготавливаем начальные данные(пустую форму)
	*/
	if request.Method == http.MethodGet {
		templates["form"].Execute(writer,
			formData{
				Rsvp:   &Rsvp{},    // создаем новый объект Rsvp
				Errors: []string{}, // Пустой срез для ошибок
			})
		// if request.Method == http.MethodPost  проверяет, яв-ся ли запрос методом POST(обычно когда пользователь отправляет какие нибудь данные)
	} else if request.Method == http.MethodPost {
		// обрабатывает данные формы, после чего к ним можно обратиться через request.Form. К примеру request.Form["name"]
		request.ParseForm()
		/*
			Создается объект Rsvp(его имя responseData) который заполняется значениями из полей
			Name: берется значение из поля формы "name" и другие значения беруться точно также
		*/
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}

		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Please enter your name")
		}
		if responseData.Email == "" {
			errors = append(errors, "Please enter email adress")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Please enter phone number")
		}
		if len(errors) > 0 {
			templates["form"].Execute(writer,
				formData{
					Rsvp:   &responseData,
					Errors: errors,
				})
		} else {
			/*
				в responses мы будем хранить ответы из responseData
				функция append используется для добавления данных в срезу
				амперсанд используется для создания указателя созданное значение Rsvp, а также что бы данные не дублировалиьсь
			*/
			responses = append(responses, &responseData)
			// сдесь проверяется будем ли пользователь на вечеринке
			if responseData.WillAttend {
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				templates["sorry"].Execute(writer, responseData.Name)
			}
		}
	}
}

func main() {
	loadTemplates()
	/*
		назначем обработчики для маршрутов
		http.HandleFunc - это функция которая связывает обработчики(welcomeHandler, listHadnler) с маршрутом(URL-путями) веб сервера
	*/
	http.HandleFunc("/", welcomeHandler) // сдесь маршрут - "/"(корневой путь сайта, например http://localhost:8080/) связавается в функцией welcomeHandler
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	//создали переменную которая запускает веб сервер по порту :5000
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
