package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	Message string `json:"message"`
}

type World struct {
	Id           uint16 `json:"id"`
	RandomNumber uint16 `json:"randomNumber"`
}

type Fortune struct {
	Id      uint16 `json:"id"`
	Message string `json:"message"`
}

const (
	// Content
	helloWorldString = "Hello, World!"
	fortuneHTML      = `<!DOCTYPE html>
<html>
<head>
<title>Fortunes</title>
</head>
<body>
<table>
<tr>
<th>id</th>
<th>message</th>
</tr>
{{range .}}
<tr>
<td>{{.Id}}</td>
<td>{{.Message}}</td>
</tr>
{{end}}
</table>
</body>
</html>`

	// Databases
	//
	// `interpolateParams=true` enables client side parameter interpolation.
	// It reduces round trips without prepared statement.
	//
	// We can see difference between prepared statement and interpolation by comparing go-raw and go-raw-interpolate
	connectionString = "benchmarkdbuser:benchmarkdbpass@tcp(%s:3306)/hello_world?interpolateParams=true"
	worldSelect      = "SELECT id, randomNumber FROM World WHERE id = ?"
	worldUpdate      = "UPDATE World SET randomNumber = ? WHERE id = ?"
	fortuneSelect    = "SELECT id, message FROM Fortune"
	worldRowCount    = 10000
	maxConnections   = 256
)

var (
	helloWorldBytes = []byte(helloWorldString)

	// Templates
	tmpl = template.Must(template.New("fortune.html").Parse(fortuneHTML))

	// Database
	db                    *sql.DB
	worldSelectPrepared   *sql.Stmt
	worldUpdatePrepared   *sql.Stmt
	fortuneSelectPrepared *sql.Stmt
)

var prefork = flag.Bool("prefork", false, "use prefork")
var child = flag.Bool("child", false, "is child proc")
var nodb = flag.Bool("nodb", false, "do not connect to the db")

func initDB() {
	var err error
	var dbhost = os.Getenv("DBHOST")
	if dbhost == "" {
		dbhost = "localhost"
	}
	ds := fmt.Sprintf(connectionString, dbhost)
	log.Println("Connecting to", ds)
	db, err = sql.Open("mysql", ds)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	db.SetMaxIdleConns(maxConnections)
	db.SetMaxOpenConns(maxConnections)

	worldSelectPrepared, err = db.Prepare(worldSelect)
	if err != nil {
		log.Fatal(err)
	}
	worldUpdatePrepared, err = db.Prepare(worldUpdate)
	if err != nil {
		log.Fatal(err)
	}
	fortuneSelectPrepared, err = db.Prepare(fortuneSelect)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Println("runtime version:", runtime.Version())

	var listener net.Listener
	flag.Parse()
	if !*prefork {
		newProcs := runtime.NumCPU()
		oldProcs := runtime.GOMAXPROCS(newProcs)
		log.Println("Changed GOMAXPROCS from", oldProcs, "to", newProcs)
	} else {
		listener = doPrefork()
	}

	if !*nodb {
		initDB()
	}

	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/db", dbHandler)
	http.HandleFunc("/dbInterpolate", dbInterpolateHandler)
	http.HandleFunc("/queries", queriesHandler)
	http.HandleFunc("/queriesInterpolate", queriesInterpolateHandler)
	http.HandleFunc("/fortune", fortuneHandler)
	http.HandleFunc("/fortuneInterpolate", fortuneInterpolateHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/updateInterpolate", updateInterpolateHandler)
	http.HandleFunc("/plaintext", plaintextHandler)
	// profiling & tracing
	http.HandleFunc("/profile/start", startProfileHandler)
	http.HandleFunc("/profile/stop", stopProfileHandler)
	http.HandleFunc("/trace/start", startTraceHandler)
	http.HandleFunc("/trace/stop", stopTraceHandler)
	if !*prefork {
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		log.Fatal(http.Serve(listener, nil))
	}
}

func doPrefork() (listener net.Listener) {
	var err error
	var fl *os.File
	var tcplistener *net.TCPListener
	if !*child {
		var addr *net.TCPAddr
		addr, err = net.ResolveTCPAddr("tcp", ":8080")
		if err != nil {
			log.Fatal(err)
		}
		tcplistener, err = net.ListenTCP("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		fl, err = tcplistener.File()
		if err != nil {
			log.Fatal(err)
		}
		children := make([]*exec.Cmd, runtime.NumCPU()/2)
		for i := range children {
			children[i] = exec.Command(os.Args[0], "-prefork", "-child")
			children[i].Stdout = os.Stdout
			children[i].Stderr = os.Stderr
			children[i].ExtraFiles = []*os.File{fl}
			err = children[i].Start()
			if err != nil {
				log.Fatal(err)
			}
		}
		for _, ch := range children {
			var err error = ch.Wait()
			if err != nil {
				log.Print(err)
			}
		}
		os.Exit(0)
	} else {
		fl = os.NewFile(3, "")
		listener, err = net.FileListener(fl)
		if err != nil {
			log.Fatal(err)
		}
		runtime.GOMAXPROCS(2)
	}
	return listener
}

var profileFile *os.File
var profileMu sync.Mutex

func startProfileHandler(w http.ResponseWriter, r *http.Request) {
	profileMu.Lock()
	defer profileMu.Unlock()

	if profileFile != nil {
		http.Error(w, "Profiling in progress", http.StatusBadRequest)
		return
	}
	f := r.URL.Query().Get("f")
	if f == "" {
		http.Error(w, "f query parameter is required", http.StatusBadRequest)
		return
	}
	if !strings.HasSuffix(f, ".prof") {
		f += ".prof"
	}

	var err error
	profileFile, err = os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = pprof.StartCPUProfile(profileFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Started CPU profiling:", f)
}

func stopProfileHandler(w http.ResponseWriter, r *http.Request) {
	profileMu.Lock()
	defer profileMu.Unlock()

	if profileFile == nil {
		http.Error(w, "Profiling not in progress", http.StatusBadRequest)
		return
	}
	pprof.StopCPUProfile()
	if err := profileFile.Close(); err != nil {
		log.Fatal(err)
	}
	profileFile = nil
	log.Println("Stopped CPU profiling")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("profile stopped"))
}

var traceFile *os.File
var traceMu sync.Mutex

func startTraceHandler(w http.ResponseWriter, r *http.Request) {
	traceMu.Lock()
	defer traceMu.Unlock()

	if traceFile != nil {
		http.Error(w, "Tracing in progress", http.StatusBadRequest)
		return
	}
	f := r.URL.Query().Get("f")
	if f == "" {
		http.Error(w, "f query parameter is required", http.StatusBadRequest)
		return
	}
	if !strings.HasSuffix(f, ".trace") {
		f += ".trace"
	}

	var err error
	traceFile, err = os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = trace.Start(traceFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Started tracing:", f)
}

func stopTraceHandler(w http.ResponseWriter, r *http.Request) {
	traceMu.Lock()
	defer traceMu.Unlock()

	if traceFile == nil {
		http.Error(w, "Trace not in progress", http.StatusBadRequest)
		return
	}
	trace.Stop()
	if err := traceFile.Close(); err != nil {
		log.Fatal(err)
	}
	traceFile = nil

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Stopped tracing"))
}

func getQueriesParam(r *http.Request) int {
	n := 1
	if nStr := r.URL.Query().Get("queries"); len(nStr) > 0 {
		n, _ = strconv.Atoi(nStr)
	}

	if n < 1 {
		n = 1
	} else if n > 500 {
		n = 500
	}
	return n
}

// Test 1: JSON serialization
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Message{helloWorldString})
}

// Test 2: Single database query
func dbHandler(w http.ResponseWriter, r *http.Request) {
	var world World
	err := worldSelectPrepared.QueryRow(rand.Intn(worldRowCount)+1).Scan(&world.Id, &world.RandomNumber)
	if err != nil {
		log.Fatalf("Error scanning world row: %s", err.Error())
	}

	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&world)
}

func dbInterpolateHandler(w http.ResponseWriter, r *http.Request) {
	var world World
	err := db.QueryRow(worldSelect, rand.Intn(worldRowCount)+1).Scan(&world.Id, &world.RandomNumber)
	if err != nil {
		log.Fatalf("Error scanning world row: %s", err.Error())
	}

	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&world)
}

// Test 3: Multiple database queries
func queriesHandler(w http.ResponseWriter, r *http.Request) {
	n := getQueriesParam(r)

	world := make([]World, n)
	for i := 0; i < n; i++ {
		err := worldSelectPrepared.QueryRow(rand.Intn(worldRowCount)+1).Scan(&world[i].Id, &world[i].RandomNumber)
		if err != nil {
			log.Fatalf("Error scanning world row: %v", err)
		}
	}

	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(world)
}

func queriesInterpolateHandler(w http.ResponseWriter, r *http.Request) {
	n := getQueriesParam(r)

	world := make([]World, n)
	for i := 0; i < n; i++ {
		err := db.QueryRow(worldSelect, rand.Intn(worldRowCount)+1).Scan(&world[i].Id, &world[i].RandomNumber)
		if err != nil {
			log.Fatalf("Error scanning world row: %v", err)
		}
	}

	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(world)
}

// Test 4: Fortunes
func fortuneHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := fortuneSelectPrepared.Query()
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}

	fortunes := fetchFortunes(rows)
	fortunes = append(fortunes, &Fortune{Message: "Additional fortune added at request time."})

	sort.Sort(ByMessage{fortunes})
	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, fortunes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func fortuneInterpolateHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(fortuneSelect)
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}

	fortunes := fetchFortunes(rows)
	fortunes = append(fortunes, &Fortune{Message: "Additional fortune added at request time."})

	sort.Sort(ByMessage{fortunes})
	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, fortunes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func fetchFortunes(rows *sql.Rows) Fortunes {
	defer rows.Close()

	fortunes := make(Fortunes, 0, 16)
	for rows.Next() { //Fetch rows
		fortune := Fortune{}
		if err := rows.Scan(&fortune.Id, &fortune.Message); err != nil {
			log.Fatalf("Error scanning fortune row: %s", err.Error())
		}
		fortunes = append(fortunes, &fortune)
	}
	return fortunes
}

// Test 5: Database updates
func updateHandler(w http.ResponseWriter, r *http.Request) {
	n := getQueriesParam(r)

	world := make([]World, n)
	for i := 0; i < n; i++ {
		if err := worldSelectPrepared.QueryRow(rand.Intn(worldRowCount)+1).Scan(&world[i].Id, &world[i].RandomNumber); err != nil {
			log.Fatalf("Error scanning world row: %v", err)
		}
		world[i].RandomNumber = uint16(rand.Intn(worldRowCount) + 1)
		if _, err := worldUpdatePrepared.Exec(world[i].RandomNumber, world[i].Id); err != nil {
			log.Fatalf("Error updating world row: %v", err)
		}
	}

	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(world)
}

func updateInterpolateHandler(w http.ResponseWriter, r *http.Request) {
	n := getQueriesParam(r)

	world := make([]World, n)
	for i := 0; i < n; i++ {
		if err := db.QueryRow(worldSelect, rand.Intn(worldRowCount)+1).Scan(&world[i].Id, &world[i].RandomNumber); err != nil {
			log.Fatalf("Error scanning world row: %v", err)
		}
		world[i].RandomNumber = uint16(rand.Intn(worldRowCount) + 1)
		if _, err := db.Exec(worldUpdate, world[i].RandomNumber, world[i].Id); err != nil {
			log.Fatalf("Error updating world row: %v", err)
		}
	}

	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(world)
}

// Test 6: Plaintext
func plaintextHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go")
	w.Header().Set("Content-Type", "text/plain")
	w.Write(helloWorldBytes)
}

type Fortunes []*Fortune

func (s Fortunes) Len() int      { return len(s) }
func (s Fortunes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByMessage struct{ Fortunes }

func (s ByMessage) Less(i, j int) bool { return s.Fortunes[i].Message < s.Fortunes[j].Message }
