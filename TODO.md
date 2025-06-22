# api

1. Добавить ручки макисльно быстро в одном файле без подлючения к базе, данные храни в map
Используем fiber

    1. `POST` `/user` req {name, age, ...}; resp id(uuid), 201, 400, 500
    2. `GET` `/user/:id`; resp {id, name, age, ...}, 200, 400, 500
    3. `DELETE` `/user/:id`; 200, 400, 500
    4. `PUT` `/user` req {id, name, age, ...}; resp id(uuid), 200, 400, 500

