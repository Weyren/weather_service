Функционал апи

http://localhost:8080/api/cities - список городов (айди, город, страна, широта, долгота)

http://localhost:8080/api/cities/:id/forecasts/shortforecast/ - краткий прогноз на 5 дней (страна, город, средняя температура на 5 дней((средняя по дневной)), список доступных дат)

ПРИМЕР: http://localhost:8080/api/cities/1/forecasts/shortforecast/ 


http://localhost:8080/api/cities/:id/forecasts/fullforecast/:date/ - прогноз погоды для города по дате на весь день

ПРИМЕР:  http://localhost:8080/api/cities/1/forecasts/fullforecast/2024-07-11/ 


http://localhost:8080/api/cities/:id/forecasts/fullforecast/:date/:time/ -прогноз погоды для города на конкретное время

ПРИМЕР: http://localhost:8080/api/cities/1/forecasts/fullforecast/2024-07-11/12:00:00/


При запуске через докер композ установить в файле конфигурации database: host: “db”
При локальном запуске в файле конфигурации установить database: host: “localhost”
