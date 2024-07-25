# Golang Site Cookie
 
## Репозиторий создан как пример работы с куки и запросами

### В основе лежит две структуры:
```
type User struct {
	ID       uint8
	Username string
	Password string
}

type MyHandler struct {
	users    map[string]*User
	sessions []string
}
```
### К типу `MyHandler` подвязаны основные методы взаимодействия с сервером:
- ```func (api *MyHandler) Login(w http.ResponseWriter, r *http.Request)```
- `func (api *MyHandler) Logout(w http.ResponseWriter, r *http.Request)`
- `func (api *MyHandler) MainPage(w http.ResponseWriter, r *http.Request)`

### В качестве роутера используется кастомный `gorilla/mux`
