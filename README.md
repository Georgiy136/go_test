# WbSchool_DB



# Работа с магазином

### Добавление нового товара
```sql
SELECT  petshop.good_upd('{
"name": "Кот Сфинкс",
"good_type_id":  1,
"description": "Сфинкс – бесшерстная порода кошек небольшого размера с бархатной кожей. Вес - 5 кг, рост - 50 см.",
"selling_price":  35000
    }'::jsonb, _staff_id := 2);
```

Пример ответа при правильном выполнении:
```jsonb
{"data" : null}
```
### Редактирование данных по существующему товару (по nm_id)
```sql
SELECT  petshop.good_upd('{
"nm_id": 2,
"name": "Кот Сфинкс",
"good_type_id":  1,
"description": "Сфинкс – бесшерстная порода кошек небольшого размера с бархатной кожей. Вес - 5 кг, рост - 50 см.",
"selling_price":  55000
    }'::jsonb, _staff_id := 2);
```

Пример ответа при правильном выполнении:
```jsonb
{"data" : null}
```
