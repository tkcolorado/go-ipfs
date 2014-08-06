package http

import (
	"fmt"
	"github.com/gorilla/mux"
	core "github.com/jbenet/go-ipfs/core"
	"github.com/jbenet/go-ipfs/importer"
	mh "github.com/jbenet/go-multihash"
	"net/http"
)

type ipfsHandler struct {
	node *core.IpfsNode
}

func Serve(address string, node *core.IpfsNode) error {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { ipfsPostHandler(w, r, node) }).Methods("POST")
	r.PathPrefix("/").Handler(&ipfsHandler{node}).Methods("GET")
	http.Handle("/", r)

	return http.ListenAndServe(address, nil)
}

func (i *ipfsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	nd, err := i.node.Resolver.ResolvePath(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", nd.Data)
}

func ipfsPostHandler(w http.ResponseWriter, r *http.Request, node *core.IpfsNode) {

	root, err := importer.NewDagFromReader(r.Body, 1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	k, err := node.DAG.Put(root)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO: return json representation of list instead
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", mh.Multihash(k).B58String())
}
