# Scaffold
<img src="https://img.shields.io/badge/Sqlite-003B57?style=for-the-badge&logo=sqlite&logoColor=white" />   <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />   

---
A repetitive backend generator. Generate endpoints, quickly and with as few LOC's as possible.

## Usage

---
Scaffold comes with a couple of configurable and reusable pieces.
### Controller
```YAML
controller(s): # This is the entry point of the controller definition
  - fallback: example #This is the fallback json the endpoint will return if the model is empty or returns an error.
    name: name1 #This the name of the controller, think of it as the id that is used to call it
    model: model1 # Put the name of the model you want to attach here
    CORS: true # Bool value that enables CORS, defaults to false
  #*Do not name multiple controller's the same way.
```
### Server
```YAML
server: # This is the entry point of the server definition
  port: 80 # The port that the endpoint server should use.
  target-log: ./main.json #This sets the output for logging, if its empty it just logs to stdout
  service(s): # This is the entry point of the service(s) definition, where you attach endpoints to logic.
    - controller: example # This calls the controller with the name 'example'
      route: /Example # This exposes the server endpoint at 'http://yourIpHere:port/Example'
```
### Database
```yaml
database: # This is the entry point of the database definition
  init-query: |  #This is the first SQL Query it runs upon starting, use it to setup the database
    CREATE TABLE IF NOT EXISTS table1 (
      id INTEGER PRIMARY KEY,
      name TEXT NOT NULL,
      age INTEGER
    );
  path: ./main.db #This is the database it operates on

#*Currently only sqlite is supported as a database
```
### Model
```Yaml
Currently model's can't be generated. WIP
```

## Defining multiple values (Controllers etc )
Using YAML's Array syntax you can define multiple values, 
if multiple values are allowed the name will end with (s) 'controller(s) , service(s) etc'
```yaml
# Example with Controller's
controller(s):
  - fallback: Hello world
    name: main_controller
  - fallback: Hello from scaffold
    name: second_controller
  - fallback: 79
    name: int_controller
  - fallback: {"key":"value"}
    name: obj_controller

#Example with Service's
service(s):
  - controller: main_controller
    route: /Greeting
  - controller: second_controller
    route: /new_example
  - controller: int_controller
    route: /get_int
  - controller: obj_controller
    route: /get_obj
```
## Examples

---
### Hello world
```YAML
controller(s):
  - fallback: Hello world
    name: main_controller
server:
  port: 8080
  target-log: ./main.json
  service(s):
    - controller: main_controller
      route: /Greeting

```
### Empty
```YAML
database:
  init-query: 
  path:
model(s):
  - query-template:
      json-template:
      name:
controller(s):
  - fallback:
    name: 
    model:
    CORS:
server:
  target-log:
  port: 
  service(s):
    - controller:
      route:
 
```
