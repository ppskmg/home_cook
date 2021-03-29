// Немного накину
// Предположим у нас несколько роутеров быстрый и функциональный
// Сигнатура у роутеров отличается

// Кейс 1: Имеем больше чем 2 эндпоинта типичный случай 15+ с некоторой долей вариативности внутри домена
  // Имеется ввиду что один роутер может обрабатывать users/:id
  // другой что то типа user/{monky}/create и пихать кучу всего в заголовки, иногда это нужно
  // третий статику, вариатвность ещё и в этом плане

// Кейс 2: 
  // чем больше эндпоинтов, тем больше ручек и кишок внутри одного файла
  // предположим нужно выделить часть эндпоинтов в одельный сервис
    // взять и положить в другое место уже просто не выйдет, застряли на одном толстом конфиге и кучей хендлеров в придачу


// донор для наброса file internal/app/apiserver/server.go
type server struct {
	router *mux.Router // Хардкодим роутер
	logger *logrus.Logger 
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter() // кишки всего роута из api
	return s
}

// торчат хендлеры на уровне конфигурации всего api
func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")
}

func (s *server) handleUsersCreate() http.HandlerFunc {} //... куча хендлеров 

// + структура вида 

// cmd/apiserver
// configs/
// internal/app/apiserver/кишки
// internal/app/model/снова кишки (+ на уровень каждой точки по 3-4+ файла + само название вводит в ступор меньше всего ожидаешь тут middleware)
// internal/app/store
// ... тут у многих ещё больше кишок в профиль


// собственно объективно почему люди не используют простое выделение на смысловые куски хотя бы
// cmd
// configs
// вложенность я бы тоже пофиксил, но чтобы не терять контекст
// internal/app/apiserver/{domain}/model <= обычно действительно модель описывающая данные, интерфейс/структура/ORM важное слово «описывающая»
// internal/app/apiserver/{domain}/handler <= тут те самые handlerFunc
// internal/app/apiserver/{domain}/router <= не плохое место чтобы применять частные middleware, но не хранить их здесь
// internal/app/apiserver/{domain}/other... <= дальше от бизнес требований
// internal/app/store
// а внутри домена обмазываться internal так как вздумается или обмазать сам домен или несколько доменов ( паблик сревис/внутренний сервис )
// domain в контексте api может быть понятным эндпоинтом для бизнеса с чёткими границами ответсвенноости
// всё что в него входит извне это взаимодействие с базой через прокладки и параметры из урла или контекста
// возможно ещё бы выделил отдельную папку для sql/nsql запросов
// router. в контексте можно воспринимать как config на уровне пакета, короче тут уже можно городить внутри зоны ответсвенности

// теперь чтобы юзать разные роутеры по хорошему бы на уровне сервера собирать роуты из пакетов, а не пытаться конфигурировать весь api внутри конфига сервера

// app/apiserver/server.go
mux := http.NewServeMux()
mux.Handle("/api/", routerHandler) // где сам routerHandler формируется уже внутри домена с нужным роутером и персональными кишками 

// внутри самого домена можно делать такого плана 

// ...{domain}/router
hs := httprouter.New()
hs.POST("/api/auth/", handlerFuncAuth)
hs.GET("/api/rout2/create", handlerFuncRout2)
hs.GET("/api/rout2/:param", handlerFuncParam)

// ...{domain2}/router
hs := httprouter.New()
hs.POST("/api/routd2", handlerFuncAuth)
hs.GET("/api/routd2/create", handlerFuncRout2)
hs.GET("/api/routd2/:param", handlerFuncParam)

// а в сервере из пакетов уровня doamin дёргать только сам роутер, серверу не нужно знать как работает пакет внутри, но надо уметь его втянуть в конфиг
// app/apiserver/server.go
mux := http.NewServeMux()
mux.Handle("/api/auth/", routerDomainHandler)
mux.Handle("/api/rout1/", routerDomainHandler2)
mux.Handle("/api/rout2/", routerDomainHandler3) // слушаем точку конкретным роутером, а остальное разгребается на уровне домена

// с логером в данном контексте примерно такие же проблемы, всё здорово пока пару точек и пару ручек. 

// как итог получим понятный менеджер пакетов на уровне api 
// не тратим времени на поиск и попытку понять зону ответственности
// при желании нужную логику можно перенести простым копипастом и подключением к новому конфигу в удну строчку
// проще дебажить и разрабатывать

// ps то чувтсво когда чистую архитектуру каждый понял по разному
