INSERT INTO ebook.publishers (id, name) VALUES (1, 'OReilly'), (2, 'Manning');
INSERT INTO ebook.languages (id, name) VALUES (1, 'English'), (2, 'German');
INSERT INTO ebook.authors (id, name) VALUES (1, 'John Doe'), (2, 'Amanda Lee');
INSERT INTO ebook.categories (id, name, parent_id)
    VALUES (1, 'Computer Science', null), (2, 'Computers', 1), (3, 'Programming', 2);
INSERT INTO ebook.file_types (id, name) VALUES (1, 'pdf'), (2, 'epub');
INSERT INTO ebook.tags (id, name) VALUES (1, 'programming'), (2, 'database');

INSERT INTO ebook.books (id, title, subtitle, description, isbn10, isbn13, asin, pages, edition, language_id,
                         publisher_id, publisher_url, pub_date, book_file_name, book_file_size, cover_file_name)
VALUES (1, 'Book 01', 'Book 01 Subtitle', 'Book 01 Description', '1111111111',
        9781111111111, 'BH11111111', 256, 1, 1, 1, 'https://amazon.com/dp/1111111111.html',
        '2022-07-19','OReilly.Book.01.1st.Edition.1111111111.zip', 5192, '1111111111.jpg');
INSERT INTO ebook.book_author (book_id, author_id) VALUES (1, 1), (1, 2);
INSERT INTO ebook.book_category (book_id, category_id) VALUES (1, 1), (1, 2), (1, 3);
INSERT INTO ebook.book_file_type (book_id, file_type_id) VALUES (1, 1), (1, 2);
INSERT INTO ebook.book_tag (book_id, tag_id) VALUES (1, 1), (1, 2);

INSERT INTO ebook.books (id, title, subtitle, description, isbn10, isbn13, asin, pages, edition, language_id,
                         publisher_id, publisher_url, pub_date, book_file_name, book_file_size, cover_file_name)
VALUES (2, 'Book 02', 'Book 02 Subtitle', 'Book 02 Description', '2222222222',
        9782222222222, 'BH22222222', 256, 1, 1, 1, 'https://amazon.com/dp/2222222222.html',
        '2021-02-11','OReilly.Book.02.1st.Edition.1111111111.zip', 5192, '2222222222.jpg');
INSERT INTO ebook.book_author (book_id, author_id) VALUES (2, 1);
INSERT INTO ebook.book_category (book_id, category_id) VALUES (2, 1);
INSERT INTO ebook.book_file_type (book_id, file_type_id) VALUES (2, 1);
INSERT INTO ebook.book_tag (book_id, tag_id) VALUES (2, 1);

INSERT INTO ebook.books (id, title, subtitle, description, isbn10, isbn13, asin, pages, edition, language_id,
                         publisher_id, publisher_url, pub_date, book_file_name, book_file_size, cover_file_name)
VALUES (3, 'Book 03', 'Book 03 Subtitle', 'Book 03 Description', '3333333333',
        9783333333333, 'BH33333333', 356, 2, 2, 2, 'https://amazon.com/dp/3333333333.html',
        '2022-05-21','OReilly.Book.03.2nd.Edition.1111111111.zip', 5193, '3333333333.jpg');
INSERT INTO ebook.book_author (book_id, author_id) VALUES (3, 2);
INSERT INTO ebook.book_category (book_id, category_id) VALUES (3, 2);
INSERT INTO ebook.book_file_type (book_id, file_type_id) VALUES (3, 2);
INSERT INTO ebook.book_tag (book_id, tag_id) VALUES (3, 2);
