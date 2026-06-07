

INSERT INTO genres (name, slug) VALUES
('Инди-рок', 'indi-rock'),
('Пост-панк', 'post-punk'),
('Электронная', 'electronic'),
('Техно', 'techno'),
('Хип-хоп', 'hip-hop'),
('Метал', 'metal'),
('Альтернатива', 'alternative'),
('Шугейз', 'shoegaze'),
('Nintendo-core', 'nintendo-core'),
('Чиптюн', 'chiptune'),
('Экспериментальная', 'experimental'),
('Инди-поп', 'indi-pop'),
('Дарквейв', 'darkwave'),
('Панк', 'punk'),
('Джаз', 'jazz')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO cities (name, slug) VALUES
('Москва', 'moskva'),
('Санкт-Петербург', 'sankt-peterburg'),
('Екатеринбург', 'ekaterinburg'),
('Новосибирск', 'novosibirsk'),
('Казань', 'kazan'),
('Нижний Новгород', 'nizhniy-novgorod'),
('Красноярск', 'krasnoyarsk'),
('Владивосток', 'vladivostok'),
('Ростов-на-Дону', 'rostov-na-donu'),
('Самара', 'samara')
ON CONFLICT (slug) DO NOTHING;

SELECT 'Жанров добавлено:' || COUNT(*) FROM genres;
SELECT 'Городов добавлено:' || COUNT(*) FROM cities;