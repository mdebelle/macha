package apiv1
 
import (
    "net/http"
 
    "github.com/zenazn/goji/web"
 
    "github.com/mdebelle/macha/controller"
)
 
// GetAuthor shows author data
func GetAuthor(c web.C, w http.ResponseWriter, r *http.Request) {
    data := controller.NewResponse()
    data.Add("object", "author")
    data.Add("name", c.URLParams["name"])
    controller.RenderJSON(w, data)
}