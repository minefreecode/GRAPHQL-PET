# Использование библиотеки GraphQL [gqlgen](https://gqlgen.com/getting-started/)

# Использованные технологии
- `GO` основной язык реализации
- `gqlgen` для GraphQL
- `viper` для конфигураций


# Генерация новых моделей
- `go run github.com/99designs/gqlgen generate`

# Порядок добавления новых действий следующий:
1. Добавление функции в `task-controller.go`, либо нужно создать файл в папке `controller` и туда добавить функцию-обработчик
2. Описание схемы в файле `schema.graphqls` согласно документации `graphql.org/learn/schema/`. При этом обязательно нужно указать 
куда относится запрос, к `Query` или `Mutation`
3. Выполнение генерации из схемы по команде `go run github.com/99designs/gqlgen generate`
4. Уточнение функции контроллера для нового сгенерированного типа данных
5. Написание запроса к серверу согласно документации `https://graphql.org/learn/queries/` или
`https://graphql.org/learn/mutations/`


# Список всех запросов:
- Получение всех тасков: `query GetAllTasks{
    tasks{
      _id
      title
      description
      company
      url
    }
  }`
- Создание тасков: `CreateTaskListingInput!){
  createTaskListing(input:$input){
    _id
    title
    description
    company
    url
  }`

`
{
  "input": {
    "title": "Task1",
    "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt",
    "company": "Alkona",
    "url": "alkona.ru"
  }
}
`
- Получение таска
`query GetTask($id: ID!){
  task(id:$id){
    _id,
    title,
    description,
    url,
    company
  }
}
`
`
{
  "id": "68554937e6bf6f0a85a2ff68"
}
`
- Обновление таска: `
  mutation UpdateTask($id: ID!,$input: UpdateTaskListingInput!) {
  updateTaskListing(id:$id,input:$input){
    title
    description
    _id
    company
    url
  }
  }
`
`
{
  "id": "68554937e6bf6f0a85a2ff68",
  "input": {
    "title": "wewe"
  }
}
`
- Удаление таска:
`
  mutation DeleteQuery($id: ID!) {
  deleteTaskListing(id:$id){
    deleteTaskId
  }
  }
`
`
{
  "id": "68554937e6bf6f0a85a2ff68"
}
`