INSERT INTO ebook.publishers (id, name) VALUES (1, 'OReilly');
INSERT INTO ebook.languages (id, name) VALUES (1, 'English');
INSERT INTO ebook.books (id, title, subtitle, description, isbn10, isbn13, asin, pages, edition,
                         language_id, publisher_id, publisher_url, pub_date, book_file_name, book_file_size, cover_file_name)
VALUES (1, 'CockroachDB', 'The Definitive Guide', 'Get the lowdown on CockroachDB', '1234567890',
        9781234567890, 'BH34567890', 256, 2, 1, 1, 'https://amazon.com/dp/1234567890.html', '2022-07-19',
        'OReilly.CockroachDB.2nd.Edition.1234567890.zip', 5192, '1234567890.jpg');
