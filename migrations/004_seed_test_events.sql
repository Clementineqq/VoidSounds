






-- DELETE FROM EVENTS
-- WHERE TITLE = 'Шум и Выходки в баре «Подвал»'

-- DELETE FROM EVENTS
-- WHERE TITLE = 'Mitski';

-- DELETE FROM EVENTS
-- WHERE TITLE = 'name';


-- INSERT INTO events (
--     title, 
--     description, 
--     date, 
--     address, 
--     price, 
--     available, 
--     status,
--     organizer_id
-- ) VALUES 
-- (
--     'Шум и Выходки в баре «Подвал»',
--     'Сольный концерт Nintendocore группы Шум и Выходки. RIPRIPRIPRIPRIPRIPRIPRIPRIPRIPRIPRIPRIPRIP.',
--     '2026-05-15 20:00:00+03',
--     'Бар «Подвал», Мытищи',
--     800,
--     87,
--     'published',
--     1
-- ),
-- (
--     'Mitski',
--     'Лютый Арт перфоманс Митски в нашем доме!',
--     '2026-05-09 22:00:00+03',
--     'квартира жинки любимой, где-то в Самаре',
--     9999,
--     0,
--     'published',
--     1
-- ),
-- (
--     'name',
--     'description',
--     '2026-06-03 19:30:00+03',
--     'address',
--     1200,
--     120,
--     'published',
--     1
-- )
-- ON CONFLICT DO NOTHING;

-- SELECT 'оп Добавлено мероприятий: ' || COUNT(*) FROM events;