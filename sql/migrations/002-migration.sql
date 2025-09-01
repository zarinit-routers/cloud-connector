-- title is for article, name is for things
ALTER TABLE nodes
RENAME COLUMN title TO name;

ALTER TABLE tags
RENAME COLUMN title TO tag;
