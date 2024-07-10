DROP TABLE IF EXISTS forecasts;
DROP TABLE IF EXISTS cities;

CREATE TABLE cities
(
    id SERIAL PRIMARY KEY,
    city CHARACTER VARYING,
    country CHARACTER VARYING,
    lat DECIMAL(9,6),
    long DECIMAL(9,6)
);

CREATE TABLE forecasts
(
    id SERIAL PRIMARY KEY,
    temp INT,
    date DATE,
    additional_info JSONB,
    city_id INT,
    CONSTRAINT forecasts_city_id_fkey FOREIGN KEY (city_id)
        REFERENCES cities(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_city_date UNIQUE (city_id, date)
);
