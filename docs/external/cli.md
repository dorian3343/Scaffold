# Scaffold CLI

The cli is the basic interface for running your Scaffold app. This page explains all the commands you can use.

## version
Displays the version of Scaffold you're using.
```batch
Scaffold version
```

## run
Start's your scaffold project, it requires a main.yml file to start.
```batch
Scaffold run [project name]
```
You can use this to run the current directory.
```
Scaffold run .
```

## init
Create a new Scaffold project.
```batch
Scaffold init [project name]
```

## auto-doc (experimental feature)
Automatically generates documentation in a auto-doc.md file,this gives a very rough documenation of the endpoint's.
```batch
Scaffold auto-doc [project name]
```
