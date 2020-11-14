package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func main() {
	c := chi.NewRouter()
	worker := NewWorker()

	go worker.Work()

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
		j := NewJob(script)
		worker.Push(j)
		output, err := j.GetOutput()
		if err != nil {
			render.PlainText(w, r, err.Error())
			return
		}
		render.PlainText(w, r, output)
	})

	log.Fatalf("Error occured: %v", http.ListenAndServe(":3000", c))
}
