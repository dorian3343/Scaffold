# External Documentation

---
This folder contains all the doc's for the api and how to configure your Scaffold 
application. 

## Content's:
1. [Components](component's.md)
2. [Good practices](good-practices.md)


## Flow Chart for Scaffold's process
Scaffold has a simple yet effective process.
It starts with an API consumer making a request to a route with an attached controller. 
The controller then makes its first decision: should it serve static content (HTML, files, images)? 
If not, it checks if it has an attached model. If it doesn't, it returns the fallback; 
otherwise, it tries to retrieve data from the database 
and returns the fallback if it encounters an error.
---
![process-flowchart.png](process-flowchart.png)