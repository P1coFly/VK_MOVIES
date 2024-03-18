DROP TABLE IF EXISTS public."ACTORS" CASCADE;

CREATE TABLE public."ACTORS"
(
    id bigserial NOT NULL,
	name text,
	sex Char(1),
	birthday DATE,
    PRIMARY KEY (id)
);

ALTER TABLE public."ACTORS"
    OWNER to postgres;		

---------------------------------------------------------
DROP TABLE IF EXISTS public."MOVIES" CASCADE;

CREATE TABLE public."MOVIES"
(
	id bigserial NOT NULL,
	title VARCHAR(150) CHECK (LENGTH(title) >0),
	description VARCHAR(1000),
	date_of_issue DATE,
	rating DECIMAL(3,1) CHECK (rating >= 0 AND rating <= 10),
    PRIMARY KEY (id)
);

ALTER TABLE public."MOVIES"
    OWNER to postgres;		

---------------------------------------------------------
DROP TABLE IF EXISTS public."ACTORS_MOVIES" CASCADE;

CREATE TABLE public."ACTORS_MOVIES"
(
    actor_id BIGINT,
    movie_id BIGINT,
    PRIMARY KEY (actor_id, movie_id),
    FOREIGN KEY (actor_id) REFERENCES public."ACTORS" (id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES public."MOVIES" (id) ON DELETE CASCADE
);

ALTER TABLE public."ACTORS_MOVIES"
    OWNER to postgres;
			

---------------------------------------------------------

	
-- Добавление данных для фильмов
INSERT INTO public."MOVIES" (title, description, date_of_issue, rating)
VALUES
    ('Однажды в Голливуде', 'Комедия о золотой эпохе киноиндустрии и уходящей эпохе Голливуда', '2019-07-26', 7.7),
    ('Бойцовский клуб', 'Фильм о мире подпольных боевых поединков, основанный на одноименном романе', '1999-10-15', 8.7),
	('Быстрее пули','Наёмник Божья Коровка отправляется на новую миссию: в он должен сесть в скоростной поезд и выкрасть чемоданчик','2022-07-18',7.7),
	('Волк с Уолл-стрит','Джордан прожигает жизнь: лавирует от одной вечеринки к другой. Однажды наступает момент, когда быстрым обогащением Белфорта начинает интересоваться агент ФБР...','2013-12-09',8.0);

-- Добавление данных для актеров
INSERT INTO public."ACTORS" (name, sex, birthday)
VALUES
    ('Брэд Питт', 'M', '1963-12-18'),
    ('Леонардо Ди Каприо', 'M', '1974-11-11'),
	('Марго Робби', 'Ж','1990-07-02'),
	('Джона Хилл', 'М','1983-12-20');


-- Добавление данных в таблицу-связь ACTORS_MOVIES
-- Создаем связи между актерами и фильмами
INSERT INTO public."ACTORS_MOVIES" (actor_id, movie_id)
VALUES
    (1, 1), -- Брэд Питт в фильме Однажды в Голливуде
    (2, 1), -- Леонардо Ди Каприо в фильме Однажды в Голливуде
	(3, 1), -- Марго Робби в фильме Однажды в Голливуде
    (1, 2), -- Брэд Питт в фильме Бойцовский клуб
    (2, 2), -- Леонардо Ди Каприо в фильме Бойцовский клуб
	(1, 3), -- Бред Питт в фильме Быстрее пули
	(2, 4), -- Леонардо Ди Каприо в фильме Волк с Уолл-стрит
	(3, 4), -- Брэд Питт в фильме Волк с Уолл-стрит
	(4, 4); -- Джона Хилла в фильме Волк с Уолл-стрит


CREATE ROLE api_service WITH LOGIN PASSWORD "12345678";
GRANT CONNECT ON DATABASE "VK_MOVIES" TO api_service;
GRANT UPDATE, SELECT, INSERT, DELETE ON ALL TABLES IN SCHEMA public TO api_service;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO api_service;


-- REVOKE ALL PRIVILEGES ON DATABASE "VK_MOVIES" FROM api_service;
-- REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM api_service;
-- REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM api_service;
-- DROP ROLE api_service