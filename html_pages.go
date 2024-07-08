package go_blueberry 

import (
    "embed"
    "path"
    "html/template"
    "net/http"
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


func HtmlTemplateHandler(pageTemplates *template.Template, template_path string, template_data any) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pageTemplates.ExecuteTemplate(w, template_path, template_data)
    }
}
