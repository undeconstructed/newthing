package main

import (
	"fmt"
	"net/http"

	m "github.com/undeconstructed/newthing/machine"
)

func main() {
	operations := m.Operations{
		"updateThing": updateThing,
	}
	resources := m.Resource{
		Actions: m.Actions{
			"GET": {
				Handler: serveRoot,
			},
		},
		Children: m.Resources{
			"test": {
				Actions: m.Actions{
					"GET": {
						Handler: serveTest,
					},
				},
				Children: m.Resources{
					"{key}": {
						Actions: m.Actions{
							"PUT": {
								Commander: thingPutter,
							},
							"GET": {
								Getter: thingGetter,
							},
						},
					},
				},
			},
			"flob": {},
			"norb": {},
		},
	}
	machine, _ := m.New(operations, resources)
	machine.Run()

	ch := make(chan struct{})
	<-ch
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}

func serveTest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test\n"))
}

func thingGetter(args m.Args, store m.Store) (m.Response, error) {
	r := map[string]string{args.Vars["key"]: "???"}
	return m.Response{Status: http.StatusOK, Body: r}, nil
}

func thingPutter(args m.Args) (m.Command, error) {
	key := args.Vars["key"]
	r := fmt.Sprintf("putting %s", key)
	return m.Command{Category: "write", Key: key, Operation: "updateThing", Message: r}, nil
	// return m.Command{}, error.New("not implemented")
}

func updateThing(command m.Command) error {
	fmt.Printf("update: %s\n", command.Message)
	return nil
}
