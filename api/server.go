package api

import(
    "fmt"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
)

type(
  Server struct {
    Id uint16 `json:"id"`
    Name string `json:"name" pg:",notnull,unique"`
    Type string `json:"type" pg:",notnull"`
    Signature string `json:"_" pg:",notnull,unique"`
  }
)

func CreateServer(w http.ResponseWriter, r *http.Request) {
    data := DecodeJsonRequest(r)
    server := createServer(
        data["name"].(string),
        data["type"].(string),
        data["signature"].(string),
    )
    factions := server.createFactions(data["factions"].([]interface{}))
    server.generateMap(factions, uint16(data["map_size"].(float64)))
    SendJsonResponse(w, 201, server)
}

func RemoveServer(w http.ResponseWriter, r *http.Request) {
	serverId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    data := DecodeJsonRequest(r)
	server := getServerBySignature(data["signature"].(string))

	if server.Id != uint16(serverId) {
        panic(NewHttpException(400, "servers.mismatch", nil))
	}
	server.remove()
	
	w.WriteHeader(204)
	w.Write([]byte(""))
}

func getServerBySignature(signature string) *Server {
    server := &Server{}
    if err := Database.Model(server).Where("signature = ?", signature).Select(); err != nil {
        panic(NewHttpException(404, "servers.not_found", err))
    }
    return server
}

func createServer(name string, serverType string, signature string) *Server {
    server := &Server{
        Name: name,
        Type: serverType,
        Signature: signature,
    }
    if err := Database.Insert(server); err != nil {
        panic(NewHttpException(500, "Server could not be created", err))
    }
    return server
}

func (s *Server) remove() {
    // All data dependencies are removed by cascade operation and triggers
    if err := Database.Delete(s); err != nil {
        panic(NewHttpException(500, "Server could not be deleted", err))
    }
}

func (s Server) String() string {
  return fmt.Sprintf(
    "Server<Id=%d Name=%q Type=%q Signature=%q>",
    s.Id,
    s.Name,
    s.Type,
    s.Signature,
  )
}
