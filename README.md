# Scaffold
<img src="https://img.shields.io/badge/Sqlite-003B57?style=for-the-badge&logo=sqlite&logoColor=white" />   <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" /> ![Ruby](https://img.shields.io/badge/ruby-%23CC342D.svg?style=for-the-badge&logo=ruby&logoColor=white)   ![YAML](https://img.shields.io/badge/yaml-%23ffffff.svg?style=for-the-badge&logo=yaml&logoColor=151515)

---

## What is Scaffold?
Scaffold is an easy to use HTTP API generator that abstract's away writing any code and replacing it with an elegant configuration system. It allow's to quickly setup simple endpoint's be it for prototyping or getting boring / repetetive work done quickly.

## Documentation
1. [External Documentation (For User's)](./docs/external/README.md)
2. [Internal Documentation (For Developer's)](./docs/internal/README.md)

## Features

* Easy / Fast to learn Syntax (YAML & SQL)
* Golang's Performance
* Easy to use database (Embedded sqlite)
* Fast development time

## How to setup Scaffold
1. Download the Latest Scaffold Release from the github repository.
2. Decompress and open the Scaffold folder.
3. Edit the main.yml file to work on your project
4. Start your project by running Scaffold.exe (prebuilt binary)
5. Test your backend

## Example
This example define's four endpoint's, one to create user's,one to retrieve a user, one to delete a user and one to get all the Users.
```yaml
database:
  init-query: |
    CREATE TABLE IF NOT EXISTS user (
      id INTEGER PRIMARY KEY,
      name TEXT NOT NULL,
      age INTEGER
    );
  path: ./main.db
$model:
  - name: add_user_model
    query-template: INSERT INTO user (name, age) VALUES ('%s', %s)
    json-template:
      - Name: Name
        Type: string
      - Name: Age
        Type: integer
  - name: main_model
    query-template: SELECT * FROM user;
  - name: delete_user_model
    query-template: DELETE FROM user WHERE name='%s'
    json-template:
      - Name: Name
        Type: string
  - name: select_user_model
    query-template: SELECT * FROM user WHERE name='%s'
    json-template:
      - Name: Name
        Type: string
$controller:
  - name: main_controller
    fallback: Something went wrong
    model: main_model
  - name: second_controller
    fallback: Something went wrong
    model: add_user_model
  - name: third_controller
    fallback: Something went wrong
    model: delete_user_model
  - name: fourth_controller
    fallback: Something went wrong
    model: select_user_model
server:
  port: 8080
  $service:
    - controller: main_controller
      route: /show_user
    - controller: second_controller
      route: /add_user
    - controller: third_controller
      route: /delete_user
    - controller: fourth_controller
      route: /select
```


