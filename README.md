# Scaffold
<img src="https://img.shields.io/badge/Sqlite-003B57?style=for-the-badge&logo=sqlite&logoColor=white" />   <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" /> ![Ruby](https://img.shields.io/badge/ruby-%23CC342D.svg?style=for-the-badge&logo=ruby&logoColor=white)   ![YAML](https://img.shields.io/badge/yaml-%23ffffff.svg?style=for-the-badge&logo=yaml&logoColor=151515)

---
A repetitive backend generator. Generate endpoints, quickly and with as few LOC's as possible.

## Documentation
1. [External Documentation](./docs/external/README.md)
2. [Internal Documentation](./docs/internal/README.md)


## How to setup Scaffold
1. Download the Latest Scaffold Release from the github repository.
2. Decompress and open the Scaffold folder.
3. Edit the main.yml file to work on your project
4. Start your project by running Scaffold.exe (prebuilt binary)
5. Test your backend

## Features

* Easy / Fast to learn Syntax (YAML & SQL)
* Golang's Speed
* Easy to use database (Embedded sqlite)
* Fast development time

## Example
This example define's two endpoint's, one to create user's and one to retrieve all the user's.
```yaml
database:
  init-query: |
    CREATE TABLE IF NOT EXISTS user (
      id INTEGER PRIMARY KEY,
      name TEXT NOT NULL,
      age INTEGER
    );
  path: ./main.db
model(s):
  - name: add_user_model
    query-template: INSERT INTO user (name, age) VALUES ('%s', %s)
    json-template:
      - Name: name
        Type: string
      - Name: age
        Type: integer
  - name: main_model
    query-template: SELECT * FROM user;
controller(s):
  - name: main_controller
    fallback: Something went wrong
    model: main_model
  - name: second_controller
    fallback: Something went wrong
    model: add_user_model
server:
  port: 8080
  service(s):
    - controller: main_controller
      route: /get_user
    - controller: second_controller
      route: /post_user

```


