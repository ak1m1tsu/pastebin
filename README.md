# Pastebin

## High-level design

![hl design](assets/hl-design.jpg)

## Database schema

![db schema](assets/db-schema.jpg)

### Users table

| column     | type                        | attrs                     |
| ---------- | --------------------------- | ------------------------- |
| id         | uuid                        | default = uuidv4, pk      |
| username   | text                        | not null, unique          |
| email      | citext                      | not null, unique          |
| avatar     | path                        | not null                  |
| deleted    | bool                        | not null, default = false |
| created_at | timestamp(0) with time zone | not null, current         |
| updated_at | timestamp(0) with time zone | not null, current         |

### Pastes table

| column        | type                        | attrs                                 |
| ------------- | --------------------------- | ------------------------------------- |
| hash          | varchar(8)                  | not null, pk                          |
| user_id       | uuid                        | ref to users(id)                      |
| title         | varchar(255)                |                                       |
| format        | varchar(255)                |                                       |
| password_hash | bytea                       |                                       |
| expires_at    | timestamp(0) with time zone | not null, default = current + 2 years |
| created_at    | timestamp(0) with time zone | not null, default = current           |
| updated_at    | timestamp(0) with time zone | not null, default = current           | 

## How to run?
