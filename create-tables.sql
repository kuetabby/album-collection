DROP TABLE IF EXISTS Album;
DROP TABLE IF EXISTS User;

CREATE TABLE "User" (
  id           SERIAL PRIMARY KEY,
  address      VARCHAR(128) NOT NULL UNIQUE,
  username     VARCHAR(255) NOT NULL,
  createdAt    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  role         VARCHAR(128) NOT NULL DEFAULT 'user'
);

CREATE TABLE "Album" (
  id         SERIAL PRIMARY KEY,
  user_id  INT NOT NULL,
  title      VARCHAR(128) NOT NULL UNIQUE,
  artist     VARCHAR(255) NOT NULL,
  price      DECIMAL(5,2) NOT NULL,
  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "User"(id)
);

-- INSERT INTO album 
--   (title, artist, price)
-- VALUES
--   ('Blue Train', 'John Coltrane', 56.99),
--   ('Giant Steps', 'John Coltrane', 63.99),
--   ('Jeru', 'Gerry Mulligan', 17.99),
--   ('Sarah Vaughan', 'Sarah Vaughan', 34.98);

-- INSERT INTO User (address, username) VALUES ('0xCD75375183e40543aee583dc2a8DBdbE3Ca15CDA', 'unnamed')

-- ALTER TABLE Album 
-- DROP CONSTRAINT fk_user_id, 
-- ADD CONSTRAINT user_id_fk 
-- FOREIGN KEY (user_id) 
-- REFERENCES User (id) ON DELETE SET NULL ON UPDATE CASCADE;

-- SELECT * FROM "User" JOIN "Album" ON "User".id = "Album".user_id;

-- // getting album within user based on relation
-- SELECT
--   "User".id,
--   "User".address,
--   "User".username,
--   "User".createdAt,
--   json_agg(json_build_object(
--     'id', "Album".id,
--     'title', "Album".title,
--     'artist', "Album".artist,
--     'price', "Album".price,
--     'user_id', "Album".user_id
--   )) AS album
-- FROM
--   "User"
-- LEFT JOIN
--   "Album" ON "User".id = "Album".user_id
-- WHERE
--   "User".id = <user_id>
-- GROUP BY
--   "User".id;

    

