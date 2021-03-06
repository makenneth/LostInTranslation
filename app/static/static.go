package static

import (
  "../api"
  "../csrf"
  "html/template"
  "net/http"
)

type TemplateData struct {
  CSRFToken string
}

func TemplateHandler(w http.ResponseWriter, r *http.Request){
  _, err := api.GetUserFromCookie(w, r)
  if err != nil {
    http.Redirect(w, r, "/login", http.StatusSeeOther)
    return
  }

  t, _ := template.ParseFiles("index.html")
  token := csrf.SetCSRF(w, r)
  d := TemplateData{CSRFToken: token}
  t.Execute(w, d)
}

func AuthPageHandler(w http.ResponseWriter, r *http.Request){
  _, err := api.GetUserFromCookie(w, r)
  if err == nil {
    http.Redirect(w, r, "/", http.StatusSeeOther)
  }

  t, _ := template.ParseFiles("templates/authPage.html")
  token := csrf.SetCSRF(w, r)
  d := TemplateData{CSRFToken: token}
  t.Execute(w, d)
}
