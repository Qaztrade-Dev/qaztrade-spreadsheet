# backend

## Создание таблицы для заполнения

Создание таблицы для заполнения требует использования OAuth 2.0 токена.
Это нужно, чтобы таблица была создана от имени пользователя. Если автором таблицы является
сервисный аккаунт, а не пользователь, тогда Google Apps Script не будет работать.

Модуль создания таблицы заполнения должен иметь действительный access_token и refresh_token.
Данные токены имеют свойство истечения, поэтому нужно будет чтобы главный пользователь 
периодично авторизовывался для обновления токена в системе.

Также с этим токеном нужно добавить сервисный аккаунт как `editor` на protected range.
 
### Обновление Access Token

Access Token работает несколько часов. Для его обновления нужно использовать Refresh Token,
чтобы получить свежий Access Token. В качестве решения этой проблемы есть два способа:
- Периодично обновлять access token.
- Пытаться делать запрос, если access token истек - запросить новый. Если refresh token истек, то ничего не поделаешь.

## Поделиться публичной ссылкой для заполнения

Для получения публичной ссылки достаточно использовать сервисный аккаунт. Через него можно 
управлять правами публичной ссылки - давать права на изменение или чтение.

## Роуты

### Создание новой заявки

Создает новый spreadsheet для пользователя и возвращает ссылку для заполнения.

```
> POST /sheets
> 
> Authorization: Bearer jwt({uid: "user_id"})

< 200 OK
< {spreadsheet_url: "url"}
```

### Добавление новой записи в заявку

Добавляет новую запись в таблицу заявки. Действует конкурентное добавление записей в рамках одного spreadsheet.

```
> POST /sheets/records
> 
> Authorization: Bearer jwt({sid: "spreadsheet_id"})
> 
> {PAYLOAD}

< 200 OK
```

### Добавление информации о заявке из tally

```
> POST /sheets/tally
> 
> {PAYLOAD WITH JWT}

< 200 OK
```


```sql
CREATE OR REPLACE VIEW applications_agg AS
select
    app.no as "№",
    app.sign_at as "Дата подписания",
    appst.value as "Статус",
    sheet.value->>'title' as "Название",
    sheet.value->>'expenses' as "Заявленная сумма",
    sheet.value->>'rows' as "Строки",
    app.attrs->'application'->>'from' as "От кого",
    app.attrs->'application'->>'gov_reg' as "Гос. регистрация",
    app.attrs->'application'->>'fact_addr' as "Фактический адрес",
    app.attrs->'application'->>'bin' as "БИН",
    app.attrs->'application'->>'industry' as "Отрасль",
    app.attrs->'application'->>'industry_other' as "Отрасль (другое)",
    app.attrs->'application'->>'activity' as "Вид деятельности",
    app.attrs->'application'->>'emp_count' as "Кол-во сотрудников",
    app.attrs->'application'->>'tax_sum' as "Сумма налогов",
    app.attrs->'application'->>'product_capacity' as "Производственная мощность",
    app.attrs->'application'->>'manufacturer' as "Производитель",
    app.attrs->'application'->>'item' as "Наименование",
    app.attrs->'application'->>'item_volume' as "Объем",
    app.attrs->'application'->>'fact_volume_earnings' as "Фактический объем выручки",
    app.attrs->'application'->>'fact_workload' as "Фактическая загрузка",
    app.attrs->'application'->>'chief_lastname' as "Фамилия руководителя",
    app.attrs->'application'->>'chief_firstname' as "Имя руководителя",
    app.attrs->'application'->>'chief_middlename' as "Отчество руководителя",
    app.attrs->'application'->>'chief_position' as "Должность руководителя",
    app.attrs->'application'->>'chief_phone' as "Телефон руководителя",
    app.attrs->'application'->>'cont_lastname' as "Фамилия контактного лица",
    app.attrs->'application'->>'cont_firstname' as "Имя контактного лица",
    app.attrs->'application'->>'cont_middlename' as "Отчество контактного лица",
    app.attrs->'application'->>'cont_position' as "Должность контактного лица",
    app.attrs->'application'->>'cont_phone' as "Телефон контактного лица",
    app.attrs->'application'->>'cont_email' as "Email контактного лица",
    app.attrs->'application'->>'info_manufactured_goods' as "Информация о производимых товарах",
    app.attrs->'application'->>'name_of_goods' as "Наименование товара",
    app.attrs->'application'->>'has_agreement' as "Наличие договора",
    app.attrs->'application'->>'spend_plan' as "План расходов",
    app.attrs->'application'->>'spend_plan_other' as "План расходов (другое)",
    app.attrs->'application'->>'metrics_2022' as "Показатели 2022",
    app.attrs->'application'->>'metrics_2023' as "Показатели 2023",
    app.attrs->'application'->>'metrics_2024' as "Показатели 2024",
    app.attrs->'application'->>'metrics_2025' as "Показатели 2025",
    app.attrs->'application'->>'agreement_file' as "Файл договора",
    app.attrs->'application'->>'expenses_sum' as "Сумма расходов",
    app.attrs->'application'->>'expenses_list' as "Список расходов",
    app.attrs->'application'->>'application_date ' as "Дата подачи заявки"
from applications app
cross join jsonb_array_elements(app.attrs -> 'sheets') as sheet
join application_statuses appst on appst.id = app.status_id
where app.no > 0
order by app.no asc
```