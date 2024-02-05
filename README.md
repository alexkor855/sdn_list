# Приложение для поиска людей в SDN листе

Для запуска сервера выполните make run

Перед использованием необходимо применить миграции.
Для запуска миграций БД вызовите <http://localhost:8080/migrations>

# Задание С

Описать алгоритм для более эффективного обновления данных при
повторном вызове метода localhost:8080/update

Можно использовать несколько подходов одновременно:

1) При публикации обновлений в xml документе указывается дата последнего обновления. В приложении реализовано сохранение истории всех попыток загрузки вместе с указанной датой. При очередной загрузке можно сравнить даты публикации, и если дата не изменилась, то ничего не загружать.

2) Если предположить, что список ФИО в AkaList может только расширяться, то можно предварительно загрузить количество записей по каждому uid из базы данных.
В случае если количество по uid не изменилось, данные по этому человеку можно не обновлять.
