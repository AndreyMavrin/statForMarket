CREATE TABLE events (
    eventID Int64,
    eventType String,
    userID Int64,
    eventTime DateTime,
    payload String
) ENGINE = MergeTree
ORDER BY (eventID, eventTime);


-- Выборки всех уникальных eventType у которых более 1000 событий.
    SELECT eventType FROM events GROUP BY eventType HAVING COUNT(eventType) > 1000;

-- Выборки событий которые произошли в первый день каждого месяца.
    SELECT * FROM events WHERE day(eventTime) = 1;

-- Выборки пользователей которые совершили более 3 различных eventType.
    SELECT userID FROM events GROUP BY userID HAVING COUNT(DISTINCT eventType) > 3;
