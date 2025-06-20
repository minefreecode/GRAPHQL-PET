В этом руководстве вы узнаете, как создать сервер GraphQL с помощью gqlgen, который может:

*   Возвращает список текущих задач
*   Создание новых задач
*   Отмечайте задачи по мере их выполнения

Готовый код для этого урока вы можете найти [здесь](https://translated.turbopages.org/proxy_u/en-ru.ru.2248daa8-6854fb86-ff99f625-74722d776562/https/github.com/vektah/gqlgen-tutorials/tree/master/gettingstarted)

## Настройка проекта[](https://translated.turbopages.org/proxy_u/en-ru.ru.2248daa8-6854fb86-ff99f625-74722d776562/https/gqlgen.com/getting-started/#set-up-project)

Создайте каталог для своего проекта и [инициализируйте его как модуль Go](https://translated.turbopages.org/proxy_u/en-ru.ru.2248daa8-6854fb86-ff99f625-74722d776562/https/golang.org/doc/tutorial/create-module):

    mkdir gqlgen-todos
    cd gqlgen-todos
    go mod init github.com/[username]/gqlgen-todos


Затем создайте файл `tools.go` и добавьте gqlgen в качестве [инструментальной зависимости для вашего модуля](https://translated.turbopages.org/proxy_u/en-ru.ru.2248daa8-6854fb86-ff99f625-74722d776562/https/go.dev/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module).

    //go:build tools
    
    package tools
    
    import (
    	_ "github.com/99designs/gqlgen"
    )


Чтобы автоматически добавить зависимость в ваш `go.mod` запуск

    go mod tidy


По умолчанию вы будете использовать последнюю версию gqlgen, но если вы хотите указать конкретную версию, вы можете использовать `go get` (заменив `VERSION` на нужную версию)

    go get -d github.com/99designs/gqlgen@VERSION


## Создание сервера[](#building-the-server)

### Создайте скелет проекта[](#create-the-project-skeleton)

    go run github.com/99designs/gqlgen init


Это создаст предложенную нами структуру пакета. При необходимости вы можете изменить эти пути в gqlgen.yml.

    ├── go.mod
    ├── go.sum
    ├── gqlgen.yml               - The gqlgen config file, knobs for controlling the generated code.
    ├── graph
    │   ├── generated            - A package that only contains the generated runtime
    │   │   └── generated.go
    │   ├── model                - A package for all your graph models, generated or otherwise
    │   │   └── models_gen.go
    │   ├── resolver.go          - The root graph resolver type. This file wont get regenerated
    │   ├── schema.graphqls      - Some schema. You can split the schema into as many graphql files as you like
    │   └── schema.resolvers.go  - the resolver implementation for schema.graphql
    └── server.go                - The entry point to your app. Customize it however you see fit


### Определите свою схему[](#define-your-schema)

gqlgen — это библиотека, ориентированная на схему: перед написанием кода вы описываете свой API с помощью [языка определения схемы](https://translated.turbopages.org/proxy_u/en-ru.ru.2248daa8-6854fb86-ff99f625-74722d776562/graphql.org/learn/schema/) GraphQL. По умолчанию это делается в файле с именем `schema.graphqls`, но вы можете разбить его на столько файлов, сколько захотите.

Сгенерированная для нас схема выглядела следующим образом:

    type Todo {
      id: ID!
      text: String!
      done: Boolean!
      user: User!
    }
    
    type User {
      id: ID!
      name: String!
    }
    
    type Query {
      todos: [Todo!]!
    }
    
    input NewTodo {
      text: String!
      userId: String!
    }
    
    type Mutation {
      createTodo(input: NewTodo!): Todo!
    }


### Реализуйте преобразователи[](#implement-the-resolvers)

При выполнении команда `generate` в gqlgen сравнивает файл схемы (`graph/schema.graphqls`) с моделями `graph/model/*` и, где это возможно, напрямую привязывается к модели. Это уже было сделано при запуске `init` . Мы отредактируем схему позже в этом руководстве, а пока давайте посмотрим, что уже было сгенерировано.

Если мы заглянем в `graph/schema.resolvers.go` , то увидим, что gqlgen не смог их сопоставить. У нас это произошло дважды:

    func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
    	panic(fmt.Errorf("not implemented"))
    }
    
    func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
    	panic(fmt.Errorf("not implemented"))
    }


Нам просто нужно реализовать эти два метода, чтобы наш сервер заработал:

Сначала нам нужно где-то хранить наше состояние, давайте поместим его в `graph/resolver.go`. В файле `graph/resolver.go` мы объявляем зависимости нашего приложения, например, нашу базу данных. Он инициализируется один раз в `server.go` при создании графа.

    type Resolver struct{
    	todos []*model.Todo
    }


Возвращаясь к `graph/schema.resolvers.go`, давайте реализуем тела этих автоматически сгенерированных функций-решателей. Для `CreateTodo` мы будем использовать [`crypto.rand`пакет](https://translated.turbopages.org/proxy_u/en-ru.ru.2248daa8-6854fb86-ff99f625-74722d776562/https/pkg.go.dev/crypto/rand#Int), чтобы просто возвращать задачу со случайным идентификатором и сохранять её в списке задач в памяти — в реальном приложении вы, скорее всего, будете использовать базу данных или какой-либо другой серверный сервис.

    func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
    	randNumber, _ := rand.Int(rand.Reader, big.NewInt(100))
    	todo := &model.Todo{
    		Text: input.Text,
    		ID:   fmt.Sprintf("T%d", randNumber),
    		User: &model.User{ID: input.UserID, Name: "user " + input.UserID},
    	}
    	r.todos = append(r.todos, todo)
    	return todo, nil
    }
    
    func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
    	return r.todos, nil
    }


### Запустите сервер[](#run-the-server)

Теперь у нас есть работающий сервер, чтобы запустить его:

    go run server.go


Откройте http://localhost:8080 в браузере. Вот несколько запросов, которые можно попробовать, начиная с создания задачи:

    mutation createTodo {
      createTodo(input: { text: "todo", userId: "1" }) {
        user {
          id
        }
        text
        done
      }
    }


А затем запрашиваю его:

    query findTodos {
      todos {
        text
        done
        user {
          name
        }
      }
    }


### Не торопись с выводом о пользователе[](#dont-eagerly-fetch-the-user)

Этот пример великолепен, но в реальном мире получение большинства объектов обходится дорого. Мы не хотим загружать данные о пользователе в todo-приложение, если пользователь этого не запросил. Поэтому давайте заменим сгенерированную `Todo` модель на что-то более реалистичное.

Сначала давайте включим `autobind`, чтобы gqlgen использовал ваши пользовательские модели, если сможет их найти, а не генерировал их. Для этого нужно раскомментировать строку `autobind` в `gqlgen.yml`:

    # gqlgen will search for any type names in the schema in these go packages
    # if they match it will use them, otherwise it will generate them.
    autobind:
     - "github.com/[username]/gqlgen-todos/graph/model"


И добавьте конфигурацию преобразователя полей `Todo` в `gqlgen.yml` для создания преобразователя для поля `user`

    # This section declares type mapping between the GraphQL and go type systems
    #
    # The first line in each type will be used as defaults for resolver arguments and
    # modelgen, the others will be allowed when binding to fields. Configure them to
    # your liking
    models:
      ID:
        model:
          - github.com/99designs/gqlgen/graphql.ID
          - github.com/99designs/gqlgen/graphql.Int
          - github.com/99designs/gqlgen/graphql.Int64
          - github.com/99designs/gqlgen/graphql.Int32
      Int:
        model:
          - github.com/99designs/gqlgen/graphql.Int32
      Todo:
        fields:
          user:
            resolver: true


Затем создайте новый файл под названием `graph/model/todo.go`

    package model
    
    type Todo struct {
    	ID     string `json:"id"`
    	Text   string `json:"text"`
    	Done   bool   `json:"done"`
    	UserID string `json:"userId"`
    	User   *User  `json:"user"`
    }


И беги `go run github.com/99designs/gqlgen generate`.

> Если вы столкнулись с этой ошибкой `package github.com/99designs/gqlgen: no Go files` при выполнении приведенной выше команды `generate`, следуйте инструкциям в [этом](https://translated.turbopages.org/proxy_u/en-ru.ru.2248daa8-6854fb86-ff99f625-74722d776562/https/github.com/99designs/gqlgen/issues/800#issuecomment-888908950) комментарии, чтобы найти возможное решение.

Теперь, если мы заглянем в `graph/schema.resolvers.go` и увидим новый преобразователь, давайте реализуем его и исправим `CreateTodo`.

    func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
    	randNumber, _ := rand.Int(rand.Reader, big.NewInt(100))
    	todo := &model.Todo{
    		Text:   input.Text,
    		ID:     fmt.Sprintf("T%d", randNumber),
    		UserID: input.UserID,
    	}
    	r.todos = append(r.todos, todo)
    	return todo, nil
    }
    
    func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
    	return &model.User{ID: obj.UserID, Name: "user " + obj.UserID}, nil
    }


## Последние штрихи[](#finishing-touches)

Вверху нашего `resolver.go`, между `package` и `import`, добавьте следующую строку:

    //go:generate go run github.com/99designs/gqlgen generate


Этот магический комментарий сообщает `go generate` о том, какую команду нужно запустить, когда мы хотим обновить наш код. Чтобы запустить рекурсивную генерацию всего проекта, используйте эту команду:

    go generate ./...