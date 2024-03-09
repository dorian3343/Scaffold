# Good practice's

---
Scaffold is a opinionated framework, naming decisions have already been made when it comes to naming.
## Before Deploying
Before deploying your Scaffold app, it's a good idea to run these two commands to generate an auto-doc.md and find potential bug's.
```batch
Scaffold auto-doc [project name]
Scaffold audit [project name]
```
## Naming Conventions

---
When naming reusable component's (Model's,Controller's etc) standard practices are:

* All lowercase with snake_case
* Ending with the component's name 'example_controller'
* Make name's clear and concise. A good idea is to match controller name's to route's
* Capitalize field name's in JSON template's


## Project Structure

---
When it come's to project structure, the standard practice is a hierarchical structure.
```

# Good:

    COMPONENT A:
    - do stuff
    COMPONENT B:
    - do thing
    - also do thing
    COMPONENT C:           Component C is the entry point and 
    - attach:              is the first thing you'll read upon opening a file
        - COMPONENT A      so it goes at the bottom where it inherits both A and C.
        - COMPONENT B       
        
         
# Bad:                     
    COMPONENT C:
    - attach:
        - COMPONENT A
        - COMPONENT B
    COMPONENT A:
    - do stuff
    COMPONENT B:
    - do thing
    - also do thing
   
```