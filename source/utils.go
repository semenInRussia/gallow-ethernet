package utils

// template ...
type template struct {
	templateName string
}

// handle ...
func (s *template) handle (w http.ResponseWriter, r *http.Request){
	// create ...
	t, err := template.ParseFiles("templates/"+s.templateName+".html")
	if err != nil {
		log.Fatal(err) 
		fmt.Fprintf(err.Error())
	}

	t.ExecuteTemplate(w, s.templateName, nil)

}