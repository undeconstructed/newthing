package machine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type command0 struct {
	command Command
}

type contextKey int

const (
	argsKey contextKey = iota
)

// GetArgs pulls args from a request.
func GetArgs(r *http.Request) Args {
	return r.Context().Value(argsKey).(Args)
}

// Machine is a sort of event sourcing thing with HTTP on the front.
type Machine struct {
	operations Operations
	root       Resource
	commandCh  chan command0
	server     *http.Server
}

// New makes a Machine.
func New(operations Operations, root Resource) (*Machine, error) {
	commandCh := make(chan command0, 100)
	root1 := compileResource("/", root)
	return &Machine{
		operations: operations,
		root:       root1,
		commandCh:  commandCh,
	}, nil
}

func compileResource(path string, r Resource) Resource {
	r.path = parsePathex(path)
	children := Resources{}
	for p, c := range r.Children {
		c1 := compileResource(p, c)
		children[p] = c1
	}
	r.Children = children
	return r
}

// Run runs.
func (m *Machine) Run() {
	m.server = &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(m.serveHTTP),
	}
	go m.server.ListenAndServe()

	go m.commandLoop()

	fmt.Printf("ready\n")
}

func (m *Machine) commandLoop() {
	for {
		c := <-m.commandCh
		fmt.Printf("command: %v\n", c.command)
		if op, exists := m.operations[c.command.Operation]; exists {
			err := op(c.command)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		}
	}
}

func (m *Machine) serveHTTP(w http.ResponseWriter, r *http.Request) {
	path := splitPath(r.URL.Path)
	if len(path) == 1 && path[0] == "docs" {
		m.serveDocs(w, r)
		return
	}
	if rmatch, exists := m.findResource(m.root, path); exists {
		if action, exists := rmatch.resource.Actions[r.Method]; exists {
			m.serveAction(w, r, rmatch, action)
		} else {
			m.serveNothing(w, r, http.StatusMethodNotAllowed)
		}
	} else {
		m.serveNothing(w, r, http.StatusBadRequest)
	}
}

func (m *Machine) serveDocs(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	buf.WriteString("<html><head><title>docs</title></head><body><h1>docs</h1>")
	writeDoc(&buf, "/", m.root)
	buf.WriteString("</body></html>")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func writeDoc(buf *bytes.Buffer, path string, r Resource) {
	buf.WriteString("<dl><dt>")
	buf.WriteString(path)
	buf.WriteString("</dt><dd>")
	for p, c := range r.Children {
		writeDoc(buf, p, c)
	}
	buf.WriteString("</dd></dl>")
}

type resourceMatch struct {
	resource Resource
	vars     map[string]string
}

func (m *Machine) findResource(r Resource, path []string) (resourceMatch, bool) {
	if vars, rest, matches := r.path.match(path); matches {
		if len(rest) == 0 {
			return resourceMatch{r, vars}, true
		}
		for _, child := range r.Children {
			if r0, e0 := m.findResource(child, rest); e0 {
				for k, v := range r0.vars {
					vars[k] = v
				}
				return resourceMatch{r0.resource, vars}, true
			}
		}
	}
	return resourceMatch{}, false
}

func (m *Machine) serveAction(w http.ResponseWriter, r *http.Request, rmatch resourceMatch, action Action) {
	if action.Handler != nil {
		action.Handler(w, r)
	} else if action.Getter != nil {
		response, err := action.Getter(Args{rmatch.vars}, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.MarshalIndent(response.Body, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if response.Status == 0 {
			response.Status = http.StatusOK
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(j)
	} else if action.Commander != nil {
		command, error := action.Commander(Args{rmatch.vars})
		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		command.id = RandomString(8)
		m.commandCh <- command0{command}
		j, _ := json.Marshal(map[string]interface{}{
			"id": command.id,
		})
		w.WriteHeader(http.StatusAccepted)
		w.Write(j)
	} else {
		m.serveNothing(w, r, http.StatusBadRequest)
	}
}

func (m *Machine) serveNothing(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}
