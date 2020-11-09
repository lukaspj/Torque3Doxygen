package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func main() {
	c := chi.NewRouter()

	c.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, r,`
<html>
<form action="/" method="post">
<textarea rows="20" cols="100" name="script">
%asd = 2+2;
echo(%asd @ "qwe");
return 3+3;
</textarea>
         <input type = "submit" name = "submit" value = "Submit" />
</form>
</html>
`)
	})

	c.Post("/", func(w http.ResponseWriter, r *http.Request) {
		script := r.PostFormValue("script")
		log.Println("Script is: ", script)
		render.PlainText(w, r, EvaluateScript(script))
	})

	log.Fatalf("Error occured: %v", http.ListenAndServe(":3000", c))
}
