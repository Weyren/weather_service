Функционал апи

http://localhost:8080/api/cities - список городов (json формат содержит название страну и координаты)

http://localhost:8080/api/cities/:id/forecasts/shortforecast/  -краткий прогноз на 5 дней(Страна, город, средняя температура на 5 дней((средняя по дневной)), список доступных дат)
ПРИМЕР: http://localhost:8080/api/cities/1/forecasts/shortforecast/ 

http://localhost:8080/api/cities/:id/forecasts/fullforecast/:date/ - прогноз погоды для города по дате на весь день
ПРИМЕР:  http://localhost:8080/api/cities/1/forecasts/fullforecast/2024-07-11/ 

http://localhost:8080/api/cities/:id/forecasts/fullforecast/:date/:time/ -прогноз погоды для города на конкретное время
ПРИМЕР: http://localhost:8080/api/cities/1/forecasts/fullforecast/2024-07-11/12:00:00/
