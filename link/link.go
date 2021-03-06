package api

import (
        "context"
        "fmt"
        "log"
        "os"
        "net/http"

        "cloud.google.com/go/firestore"
        "github.com/gorilla/mux"
      )

var router = newRouter()

func newRouter() *mux.Router {
  router := mux.NewRouter()
  router.HandleFunc("/{id}", LinkController)
  router.HandleFunc("/api/{id}", LinkController)
	return router
}

//F represents cloud function entry point
func Link(w http.ResponseWriter, r *http.Request) {
  router.ServeHTTP(w, r)
}

// GCLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("PROJECT_ID")

func LinkController(w http.ResponseWriter, r *http.Request) {

  // https://github.com/golang/go/issues/15867#issuecomment-223748637
  cs := w.Header().Get("Set-Cookie")
  cs += "; SameSite=None; Secure"
  w.Header().Set("Set-Cookie", cs)

  url := r.URL.Path
  id := mux.Vars(r)["id"]
  ctx := context.Background()

  log.Println("ingress url:", url, "id:", id)

  client, err := firestore.NewClient(ctx, projectID)
  if err != nil {
    log.Fatalln(err)
    http.Error(w, "500 - Internal Server Error " + url, http.StatusInternalServerError)
    return
  }

  dsnap, err := client.Collection("links").Doc(id).Get(ctx)
  if err != nil {
    log.Fatalln(err)
    http.Error(w, "404 - Not Found url", http.StatusNotFound)
    fmt.Fprintf(w, "Hello, World!")
    return
  }

  m := dsnap.Data()
  linkTo := m["link"].(string)

  log.Printf("Redirect to: %s -> %s", url, linkTo)
  http.Redirect(w, r, linkTo, http.StatusFound)
}
