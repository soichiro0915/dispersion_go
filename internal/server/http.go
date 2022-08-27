package server 

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")
	r.HandleFunc("/", httpsrv.handleDelete).Methods("DELETE")
	return &http.Server{
		Addr: addr,
		Handler: r,
	}
}

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer {
		Log: NewLog(),
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

type DeleteRequest struct {
	Offset uint64 `json:"offset"`
}

type DeleteResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	//リクエストのJSONボディをアンマーシャルして構造体に変換
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//ロジックをリクエストに対して実行
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//結果をマーシャルしてレスポンスに書き込む
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//リクエストのJSONボディをアンマーシャルして構造体に変換
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//ロジックをリクエストに対して実行
	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//結果をマーシャルしてレスポンスに書き込む
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//リクエストのJSONボディをアンマーシャルして構造体に変換
	var req DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//ロジックをリクエストに対して実行
	record, err := s.Log.Delete(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//結果をマーシャルしてレスポンスに書き込む
	res := DeleteResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}