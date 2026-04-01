package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-exchequer/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListBudgets();if list==nil{list=[]store.Budget{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var b store.Budget;json.NewDecoder(r.Body).Decode(&b);if b.Name==""||b.Period==""{writeError(w,400,"name and period required");return};s.db.CreateBudget(&b);writeJSON(w,201,b)}
func(s *Server)handleAddActual(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var a store.Actual;json.NewDecoder(r.Body).Decode(&a);a.BudgetID=id;if a.Description==""{writeError(w,400,"description required");return};s.db.AddActual(&a);writeJSON(w,201,a)}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteBudget(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
