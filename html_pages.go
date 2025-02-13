package go_blueberry

import (
	"bytes"
	"embed"
	"github.com/gorilla/schema"
	"html/template"
	"net/http"
	"path"
	"strings"
)

var decoder = schema.NewDecoder()

type Validatable interface {
    Validate() bool
}

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
		if strings.HasSuffix(templateFile.Name(), "swp") {
			continue
		}
		templateFullPath := path.Join(templateBasePath, templateFile.Name())
		template.Must(pageTemplates.ParseFS(templateFS, templateFullPath))
	}
	return pageTemplates
}

func HtmlTemplateHandler[T any](pageTemplates *template.Template, templatePath string, dataGetter func() T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateData := dataGetter()
		pageTemplates.ExecuteTemplate(w, templatePath, templateData)
	}
}

func HandleContactFormRequest[FormType Validatable](mail Mail, templateCollection *template.Template, subjectTemplate string, emailTemplate string, thankYouTemplate string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ParseForm() != nil {
            w.WriteHeader(http.StatusBadRequest)
            templateCollection.ExecuteTemplate(w, "badRequest", "") 
			return
		}
		var form FormType
		err := decoder.Decode(&form, r.PostForm)
		if err != nil || !form.Validate() {
            w.WriteHeader(http.StatusBadRequest)
            templateCollection.ExecuteTemplate(w, "badRequest", "") 
            return
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
