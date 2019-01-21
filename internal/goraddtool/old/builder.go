package old

import (
	"net/http"
	"bytes"
	"path/filepath"
)


func serveBuilder() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		buf := new (bytes.Buffer)
		cmd2 := r.FormValue("cmd")
		if cmd2 != "" {
			cmd = cmd2
			results = ""
			stop = false
		}
		switch cmd2 {
		case "codegen":
			go startCodegen()
		case "run":
			go buildAndRunApp();
		}
		drawBuilder(buf)
		w.Write(buf.Bytes())
	}
	return http.HandlerFunc(fn)

}

func startCodegen() {
	results += "Generating codegen.go...\n"
	codegenLoc := filepath.Join(goraddPath(), "buildtools", "codegen.go")
	cmdResult, errStr, err := executeCmd("go", "generate", codegenLoc)

	if err != nil {
		results += errStr
		stop = true
		return
	}
	results += cmdResult

	results += "Success!"
	stop = true
}

func buildAndRunApp() {
	buildApp()
	runApp()
}

func buildApp() {

}

func runApp() {

}