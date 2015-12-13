package bdgate

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

func HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s to %s...\n", indexPage, r.RemoteAddr)

	tmplIndex.ExecuteTemplate(w, "index.html", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
