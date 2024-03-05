package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"service/configuration"
	"slices"
	"strings"
)

func PrintGuide() {
	fmt.Println(`Scaffold is a tool for building APIs fast and easy

Usage:
	Scaffold <command> [argument]

List of commands:
	version   print out your scaffold version
	run       run the scaffold from a config in a specified directory
	init	  creates a new project from a template  
	auto-doc  generates api documentation for your app`)
	fmt.Println("\nTool by Dorian Kalaczyński")
	os.Exit(0)

}

func PrintVersion() {
	body, err := os.ReadFile("VERSION")
	if err != nil {
		fmt.Println("Something went wrong reading version: " + err.Error())
	} else {
		fmt.Println(string(body))
	}
}

func ProjectInit(x string) {
	// Create directory x relative to the current working directory
	err := os.Mkdir(x, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create main.yml file within the newly created directory
	filename := filepath.Join(x, "main.yml")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	yamlString := `
$controller:
  - name: main_controller
    fallback: hello world
server:
  port: 8080
  $service:
    - controller: main_controller
      route: /api
`
	_, err = file.WriteString(yamlString)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("Project Created Successfully")
}

func GenerateDoc(path string) {
	var docString strings.Builder

	// Generate the config struct first
	conf, _ := configuration.Setup(path + "/main.yml")
	filename := filepath.Join(path, "auto-doc.md")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error: could not find project")
		return
	}
	docString.WriteString("# " + path + " auto docs\n")
	docString.WriteString("This is auto generated documentation using Scaffold's auto-doc command\n")

	if conf.Server.Static != "" {
		docString.WriteString("## " + " Static content\n")
		docString.WriteString("Static content is being served at the server's address on the route '/' \n")
	}

	for key, value := range conf.Server.Services {
		docString.WriteString("## " + key + "\n")
		if value.Model == nil && slices.Equal(value.Fallback, []byte("null")) {
			docString.WriteString("This route does nothing\n")
		} else if value.Model == nil {
			docString.WriteString("This route returns:\n ```JSON\n" + string(value.Fallback) + "```\n")
		} else {
			docString.WriteString("This route runs the query:\n ```SQL\n" + value.Model.GetQuery() + "\n```\n")
			docString.WriteString("and fallsback to:\n ```JSON\n" + string(value.Fallback) + "\n```\n")
		}
	}

	_, err = file.WriteString(docString.String())
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
}
