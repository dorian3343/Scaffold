package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"service/configuration"
	"service/misc"
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
	auto-doc  generates api documentation for your app
	audit	  checks your project for potential error's'`)
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
	if conf.Database != nil {
		if conf.Database.InitQuery != "" {
			docString.WriteString("## Database\n The database initializes using this query :\n```SQL\n" + conf.Database.InitQuery + "\n```\n")
		}
	}

	if conf.Server.Static != "" {
		docString.WriteString("## " + " Static content\n")
		docString.WriteString("Static content is being served at the server's address on the route '/' \n")
	}

	for key, value := range conf.Server.Services {
		docString.WriteString("## " + key + "\n")
		if value.Model == nil && slices.Equal(value.Fallback, []byte("null")) {
			docString.WriteString("This route does nothing\n")
		} else if value.Model == nil {
			fmt.Println(string(value.Fallback))
			docString.WriteString("This route returns:\n```JSON\n" + string(value.Fallback) + "\n```\n")
		} else {
			docString.WriteString("This route runs the query:\n ```SQL\n" + value.Model.GetQuery() + "\n```\n")
			if value.Model.GetJsonTemplate().Len() != 0 {
				jsonT := value.Model.GetJsonTemplate()
				// Generate JSON example
				docString.WriteString("JSON Specification:\n```json\n{\n")
				for _, Name := range jsonT.Keys() {
					T, _ := jsonT.Get(Name)
					docString.WriteString(fmt.Sprintf("  %s : %s\n", Name, T))
				}
				docString.WriteString("}\n```\n")
			}
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
func Audit(path string) {
	// Generate the config struct first
	conf, _ := configuration.Setup(path + "/main.yml")

	// Track seen names
	seenNames := make(map[string]bool)

	// Track if any warnings were found
	foundWarning := false

	// Check controllers
	for _, controller := range conf.Controllers {
		if !strings.HasSuffix(controller.Name, "_controller") {
			fmt.Printf("Naming warning -> Controller %s should end with '_controller'\n", controller.Name)
			foundWarning = true
		}
		if controller.Name != strings.ToLower(controller.Name) {
			fmt.Printf("Naming warning -> Controller %s should be all lowercase\n", controller.Name)
			foundWarning = true
		}
		if seenNames[controller.Name] {
			fmt.Printf("Duplicate warning -> Controller %s is duplicated\n", controller.Name)
			foundWarning = true
		} else {
			seenNames[controller.Name] = true
		}
		if slices.Equal(controller.Fallback, []byte("null")) {
			fmt.Printf("General warning -> Controller %s has an empty fallback\n", controller.Name)
			foundWarning = true
		}
	}

	// Check models
	for _, model := range conf.Models {
		if !strings.HasSuffix(model.Name, "_model") {
			fmt.Printf("Naming warning -> Model %s should end with '_model'\n", model.Name)
			foundWarning = true
		}
		if model.Name != strings.ToLower(model.Name) {
			fmt.Printf("Naming warning -> Model %s should be all lowercase\n", model.Name)
			foundWarning = true
		}

		jsonT := model.GetJsonTemplate()
		for _, Name := range jsonT.Keys() {
			if Name != misc.Capitalize(Name) {
				fmt.Printf("Naming warning -> Model %s has non-capitalized JSON field '%s'\n", model.Name, Name)
				foundWarning = true
			}
		}

		if seenNames[model.Name] {
			fmt.Printf("Duplicate warning -> Model %s is duplicated\n", model.Name)
			foundWarning = true
		} else {
			seenNames[model.Name] = true
		}
	}

	// If no warnings were found, print a message
	if !foundWarning {
		fmt.Println("Success -> No warnings found during audit.")
	}
}
