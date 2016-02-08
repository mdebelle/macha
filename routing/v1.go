package routing
 
import (
    "fmt"
    "net/http"
 
    "github.com/zenazn/goji/web"
)
 
// SetV1 sets api routing ver1
func SetV1(r *web.Mux) {
    r.Get("/foo", func(c web.C, w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "custom route!")
    })
}