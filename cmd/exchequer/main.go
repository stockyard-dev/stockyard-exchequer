package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-exchequer/internal/server";"github.com/stockyard-dev/stockyard-exchequer/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./exchequer-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("exchequer: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Exchequer — Self-hosted budget tracker\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("exchequer: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
