package packs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cloud66/starter/common"
)

const templateHeader = "# Generated by Cloud66 Starter"

type TemplateWriterBase struct {
	TemplateDir  string
	OutputDir    string
	ShouldPrompt bool
}

func (w *TemplateWriterBase) WriteTemplate(templateName string, filename string, context interface{}) error {
	tmpl, err := template.ParseFiles(filepath.Join(w.TemplateDir, templateName))
	if err != nil {
		return err
	}

	destFullPath := filepath.Join(w.OutputDir, filename)
	if !w.shouldOverwriteExistingFile(destFullPath) {
		newName := filename + ".old"
		err = os.Rename(destFullPath, filepath.Join(w.OutputDir, newName))
		if err != nil {
			return err
		}
		fmt.Println(common.MsgL2, fmt.Sprintf("----> Renaming %s to %s...", filename, newName), common.MsgReset)
	}

	destFile, err := os.Create(destFullPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			fmt.Printf("%s Cannot close file %s due to: %s\n", common.MsgError, filename, err.Error())
		}
	}()

	fmt.Println(common.MsgL2, fmt.Sprintf("----> Writing %s...", filename), common.MsgReset)
	err = tmpl.Execute(destFile, context)
	if err != nil {
		return err
	}

	return w.prependToFile(destFullPath, templateHeader)
}

func (w *TemplateWriterBase) shouldOverwriteExistingFile(filename string) bool {
	if !common.FileExists(filename) {
		return true
	}
	isStarterTemplate := w.isStarterTemplate(filename)
	if isStarterTemplate {
		return true
	}
	if !w.ShouldPrompt {
		return isStarterTemplate
	}

	message := fmt.Sprintf(" %s cannot be written as it already exists. What to do? [o: overwrite, R: rename] ", filepath.Base(filename))
	answer := "none"
	for answer != "o" && answer != "r" && answer != "" {
		fmt.Print(common.MsgL1, message, common.MsgReset)
		if _, err := fmt.Scanln(&answer); err != nil {
			return false
		}
		answer = strings.TrimSpace(strings.ToLower(answer))
	}
	return answer == "o"
}

func (w *TemplateWriterBase) isStarterTemplate(filename string) bool {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}

	lines := strings.Split(string(content), "\n")
	return len(lines) > 0 && lines[0] == templateHeader
}

func (w *TemplateWriterBase) prependToFile(filename string, text string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	lines = append([]string{text}, lines...)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(strings.Join(lines, "\n"))
	return err
}
