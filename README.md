# Hello world
```Yaml
controller(s):
  - fallback: Hello world
    name: main_controller
server:
  port: 8080
  service(s):
    - controller: main_controller
      route: /Greeting

```
# Empty
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
  - fallback: 
    name: 
    model:
server:
  port: 8080
  service(s):
    - controller:
      route:
 
```