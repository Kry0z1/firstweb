## Auth middleware

Модуль аутентификации пользователя с токеном доступа `Bearer token`

### Usage
Endpoint `LoginForToken` возвращает токен и тип токена в формате json в ответ на запрос содержащий логин и пароль:
- Запрос:
    ```bash
    curl --request POST --data '{"username": "fill", "password": "fill2"}' http://localhost:21000/login
    ```
- Ответ:
    ```bash
    {"AccessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDM2MjI0OTgsInN1YiI6ImZpbGwifQ.m1AZYsO2c38kHHdsOvnGP0u4vnxWLNOrPD-CRi7S9mk","TokenType":"Bearer"}
    ```
Авторизация пользователя при последющих запросах осуществляется при помощи middleware `CheckAuth`. Для успешной авторизации клиент должен выставить следующий header в следующем формате:

`Authorization: Bearer {token}`

Далее методы, требующие авторизацию могут вызвать метод `ContextUser` для получения текущего пользователя. 