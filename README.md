# Проект таск-трекера на Go

Описание эндпоинтов:

- **GET /all** - возвращает список созданных задач
- **POST /task** - создание задачи, принимает на вход:
  {
    "Title": "string",
    "Description": "string",
    "Status": bool,
    "Priority": number
  }
- **PUT /task/:id** - редактирование задачи по id
- **DELETE /delete/:id** - удаление задачи по id
- **GET /filter** - возвращает задачи с фильтром по status и priority
- **GET /list?page=** - возращает список задач постранично, с указанием общего количества задач
