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
    app.id as "ID",
    app.no as "№",
    app.sign_at as "Дата подписания",
    appst.value as "Статус",
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
    app.attrs->'application'->>'application_date ' as "Дата подачи заявки",
    ass."Вид затрат",
    ass."Всего строк",
    ass."Заявленная сумма",
    ass."Оцифровка",
    ass."Юр.часть",
    ass."Фин.часть"
from applications app
join application_statuses appst on appst.id = app.status_id
join assignments_agg ass on ass."Идентификатор заявления" = app.id
where app.no > 0
order by app.no asc
```

```sql
CREATE OR REPLACE VIEW assignments_agg AS
select
	ap.id as "Идентификатор заявления",
	ap.no as "Номер заявления",
	ap.attrs->'application'->>'from' as "От кого",
	ass.sheet_title as "Вид затрат",
	ass.total_rows as "Всего строк",
	ass.total_sum as "Заявленная сумма",
	MAX(CASE WHEN ass.type = 'digital' THEN u.email ELSE NULL END) as "Оцифровка",
	MAX(CASE WHEN ass.type = 'legal' THEN u.email ELSE NULL END) as "Юр.часть",
	MAX(CASE WHEN ass.type = 'finance' THEN u.email ELSE NULL END) as "Фин.часть"
from assignments ass
join applications ap on ap.id = ass.application_id
join users u on u.id = ass.user_id
group by 
	ap.id, 
	ap.no, 
	ap.attrs->'application'->>'from',
	ass.sheet_title,
	ass.total_rows,
	ass.total_sum
order by ap.no asc
```

Миграции чтобы объеденить задания по заявлениям, а не по выгрузкам:

```sql
with application_sheets as (
    select
        app.id,
        string_agg(distinct ass.sheet_title, ', ') as sheet_title
    from assignments ass
    join applications app on app.id = ass.application_id
    group by app.id
)
update assignments as ass set
    sheet_title = appsh.sheet_title
from application_sheets appsh
where 
    appsh.id = ass.application_id
;

with assignment_totals as (
    select
        app.id,
        ass.type,
        sum(ass.total_rows) as total_rows,
        sum(ass.total_sum) as total_sum
    from assignments ass
    join applications app on app.id = ass.application_id
    group by 
        app.id, ass.type
)
update assignments as ass set
    total_rows = asstot.total_rows,
    total_sum = asstot.total_sum
from assignment_totals asstot
where 
    asstot.id = ass.application_id
    and asstot.type = ass.type
;

WITH assignments_rn AS (
    SELECT
        ass.id,
        ROW_NUMBER() OVER(PARTITION BY ass.sheet_title, ass.application_id, ass.type ORDER BY (SELECT NULL)) AS rn
    FROM assignments ass
)
DELETE FROM assignments
WHERE id IN (
    SELECT id
    FROM assignments_rn
    WHERE rn > 1
);
```
