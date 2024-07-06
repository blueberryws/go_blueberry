package go_blueberry 

import (
    "embed"
    "html/template"
    "path"
)

func InitTemplates(templateFS embed.FS, templateBasePath string) *template.Template {
    var pageTemplates *template.Template = template.New("pages")
    if len(templateBasePath) == 0 {
        templateBasePath = "."
    }
    templates, err := templateFS.ReadDir(templateBasePath)
    if err != nil {
        panic(err)
    }
    for _, templateFile := range templates {
        templateFullPath := path.Join(templateBasePath, templateFile.Name())
        template.Must(pageTemplates.ParseFS(templateFS, templateFullPath))
    }
    return pageTemplates
}
