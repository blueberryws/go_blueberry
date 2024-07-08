package go_blueberry 

import (
    "bytes"
    "embed"
    "path"
    "html/template"
    "net/http"
    "github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

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

func HandleContactFormRequest[FormType any](mail Mail, templateCollection *template.Template, subjectTemplate string, emailTemplate string, thankYouTemplate string) http.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) {
        if r.ParseForm() != nil {
            http.Error(w, "Could not parse form.", http.StatusBadRequest)
            return
        }
        var form FormType
        err := decoder.Decode(&form, r.PostForm)
        if err != nil {
            panic(err)
        }

        var subject bytes.Buffer
        templateCollection.ExecuteTemplate(&subject, subjectTemplate, form)

        var body bytes.Buffer
        templateCollection.ExecuteTemplate(&body, emailTemplate, form)
        mail.SendMail(
            subject.String(),
            body.String(),
        )

        templateCollection.ExecuteTemplate(w, thankYouTemplate, "")
    }
}
