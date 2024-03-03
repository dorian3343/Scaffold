# Components

Components are the basic building blocks of Scaffold applications. They configure every element from your database to your server, making them extremely important to understand.

## Specification

### Server
There can only be one server component, because it is the entry point for your application's construction. 
Here you configure the web server and define route's.
```YAML
server: 
  port: 80 
  static: ./static
  target-log: ./main.json 
  $service: 
    - controller: example 
      route: /Example 
  
  # port -> Set's the server's port to the int value.
  # static -> display's the static content of the input server @ path '/'
  # target-log -> Set's the target file for logging, if left empty it only prints to stdout
  # $service -> Connects an endpoint to a Scaffold Controller
  # controller -> Set the controller for the specific service. These can be reused. Use the controller's name.
  # route -> Exposes an endpoint to handle a service.
  ```

### Controller
Controllers are the point of entry for your application's users. They attach basic logic to a route (which can be extended with Models).

```yaml
$controller:
  - name: name1 
    fallback: example 
    model: model1 
    cors: "*" 

    # name -> This is the name of the controller. YOou use this to attach it to other component's
    
    # fallback -> This is the value that the endpoint will return if the model isn't set OR fails.
    # Whatever you set there is returned as JSON allowing for sending Objects.
    
    # model -> Attaches data handling to a controller, read up on them at the 'Model' section.
    
    # cors -> Sets a cors value to input string, without setting it, nothing gets set
```
### Database
```yaml
database: 
  init-query: |
    CREATE TABLE IF NOT EXISTS table1 (
      id INTEGER PRIMARY KEY,
      name TEXT NOT NULL,
      age INTEGER
    );
  path: ./main.db 

  # init-query -> The query it runs upon starting up, used to initialize databases.
  # path -> The path Scaffold should look for a database, if it doesn't exist. It get's Created
  # Currently only sqlite is supported as a database due to ease of embedding
```
### Model
Model's handle data operations, they communicate through the Controller's.
```yaml
$model:
  - name: add_user_model
    query-template: INSERT INTO user (name, age) VALUES ('%s', %s)
    json-template:
      - Name: name
        Type: string
      - Name: age
        Type: integer

  # name -> This is the model's name. You use this to attach it to other component's
        
  # query-template -> As the name suggest's its a template to fill out 
  #using data recieved via JSON during the request. It uses '%s' as value placeholders.
  
  # json-template -> This is the template for JSON request's. 
  # It takes in an array of 'object's' with two values:
  # Name -> the name of the data field. The Name's MUST be capitalized or will fail
  # Type -> the data type of the field. (Currently only supports : string, integer)
  # If the request doesn't match the spec it throw's a status 400.
  # If left empty it doesn't fill out the template.
```

## Multiple Component's
Using YAML's Array syntax you can define multiple Component's
if multiple component's are allowed the name will end with (s) 'controller(s) , service(s) etc'
```yaml
# Example with Controller's
$controller:
  - fallback: Hello world
    name: main_controller
  - fallback: Hello from scaffold
    name: second_controller
  - fallback: 79
    name: int_controller
  - fallback: {"key":"value"}
    name: obj_controller

#Example with Service's
$service:
  - controller: main_controller
    route: /Greeting
  - controller: second_controller
    route: /new_example
  - controller: int_controller
    route: /get_int
  - controller: obj_controller
    route: /get_obj
```
