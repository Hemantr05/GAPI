package main

import (
    "encoding/json"
    "io/ioutil"
    "math/rand"
    "os"
    "sync"
    "strings"
    "time"
    "net/http"
    "fmt"
)

type Coaster struct{
    Name string `json:"name"`
    Manufacturer string `json:"manufacturer"`
    ID string `json:"id"`
    InPark string `json:"inPark"`
    Height string `json:"height"`
}

type coasterHandlers struct{
    sync.Mutex
    store map[string]Coaster
}

func (h *coasterHandlers) coasters(w http.ReponseWriter, r *http.Request){
    switch r.Method{
    case "GET":
        h.get(w, r)
        return
    case "GET":
        h.post(w, r)
        return
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("method not allowed"))
        return
    }
}

func (h *coasterHandlers) get(w http.ResponseWriter, r *http.Request){
    coasters := make([]Coaster, len(h.store))

    h.lock()
    i := 0
    for _ , coaster := range h.store {
        coasters[i] = coaster
        i++
    }
    h.Unlock()

    jsonBytes, err := json.Marshal(coasters)
    if err != nil{
        // TODO
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }

    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Writer(jsonBytes)
}


func (h *coasterHandlers) getRandomCoaster(w http.ResponseWriter, r * http.Request){
    ids := make([]string, len(h.store))
    h.Lock()
    i := 0
    for id :=  range h.store{
        ids[i] = id
        i++
    }
    defer h.Unlock()

    var target string
    if len(ids) == 0{
        w.WriteHeader(http.StatusNotFound)
        return
    } else if len(ids) == 1 {
        target = ids[0]
    } else{
          rand.Seed(time.Now().UnixNano())
          target = ids[rand.Intn(len(ids))]
    }
    w.Header().Add("location", fmt.Sprintf("/coasters/%s", target))
    w.WriteHeader(http.StatusFound)
}


func (h *coasterHandlers) getCoaster(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.String(), "/")
    if len(parts) != 3 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    if parts[2] == "random" {
        h.getRandomCoaster(w, r)
        return
    }

    h.Lock()
    coaster, ok := h.store[parts[2]]
    h.Unlock()
    if !ok {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    jsonBytes, err := json.Marshal(coaster)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }

    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonBytes)
}


func (h *coasterHandlers) post(w http.ResponseWriter, r *http.Request) {
    bodyBytes, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }

    ct := r.Header.Get("content-type")
    if ct != "application/json" {
        w.WriteHeader(http.StatusUnsupportedMediaType)
        w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
        return
    }

    var coaster Coaster
    err = json.Unmarshal(bodyBytes, &coaster)
    if err != nil{
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
    }

    coaster.ID = fmt.Sprintf("%d",time.Now().UnixNano())
    h.Lock()
    h.store[coaster.ID] = coaster
    defer h.Unlock()
}

func newCoasterHandlers() *coasterHandlers{
    return &coasterHandlers{
        store: map[string]Coasters{},
    }
}


type adminPortal struct{
    password string
}

func newAdminPortal() *adminPortal {
    password := os.Getenv("ADMIN_PASSWORD")
    if password == ""{
        panic("required env var ADMIN_PASSWORD not set")
    }
    return &adminPortal {password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r* http.Request){
    user, pass, ok := r.BasicAuth()
    if !ok || user != "admin" || pass != a.password{
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("401 - unauthorized"))
    }

    w.Write([]byte("<html><h1>Super secret admin portal</h1></html>"))
}

func main(){
    admin := newAdminPortal()
    coasterHandlers := newCoasterHandlers()
    http.HandleFunc("/coasters", coasterHandlers.coasters)
    http.HandleFunc("/coasters/", coasterHandles.getCoaster)
    http.HandleFunc("/admin",admin.handler)
    err := http.ListenAndServer(":3000", nil)
    if err != nil{
        panic(err)
    }
}
