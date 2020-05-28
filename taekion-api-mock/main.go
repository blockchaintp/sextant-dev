package main

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

type Opts struct {
	Bind string `short:"b" long:"bind" description:"Bind port for API" default:":8000"`
}

type returnVal struct {
	Action  string      `json:"action"`
	Object  string      `json:"object"`
	Payload interface{} `json:"payload"`
}

func main() {
	var opts Opts

	parser := flags.NewParser(&opts, flags.Default)
	remaining, err := parser.Parse()

	if err != nil {

		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Errorf("Failed to parse args: %v\n", err)
			os.Exit(-1)
		}

	}

	if len(remaining) > 0 {
		log.Errorf("Error: Unrecognized arguments passed: %v\n", remaining)
		os.Exit(-1)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Path("/volume").
		Methods(http.MethodGet).
		HandlerFunc(volume)
	router.Path("/snapshot").
		Methods(http.MethodGet).
		HandlerFunc(snapshot)
	s := router.
		PathPrefix("/debug").
		Methods(http.MethodGet).
		Subrouter()
	s.HandleFunc("/address", address)
	s.HandleFunc("/bundle", bundle)
	s.HandleFunc("/datablock", datablock)
	s.HandleFunc("/directory", directory)
	s.HandleFunc("/inode", inode)
	s.HandleFunc("/volume", debugVolume)
	s.HandleFunc("/wrapper", wrapper)

	log.Info("Binding to " + opts.Bind)
	srv := new(http.Server)
	srv.Addr = opts.Bind
	corsHandler := cors.Default().Handler(router)
	loggingHandler := handlers.LoggingHandler(os.Stdout, corsHandler)
	srv.Handler = loggingHandler
	connsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		// interrupt - shut down.  Stop accepting new requests and wait
		// for ones in progress to complete.
		log.Info("Shutting down...")

		if err := srv.Shutdown(context.Background()); err != nil {
			// Error closing listeners, or timeout
			log.Print(err)
		}

		close(connsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener
		log.Print(err)
		close(connsClosed)
	}

	// wait for shutdown to complete before exiting
	<-connsClosed
}

// address simulates querying the debug API for various data structure addresses.
// Query Arguments:
// "object" - string, accepts the following values: "inode", "datablock", "volume", "directory"
// "id" - UUID in xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx format, example: 11111111-1111-1111-1111-111111111111
// --or-- string, example: volume - name, datablock hash - 2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea

func address(w http.ResponseWriter, req *http.Request) {
	var response returnVal

	response.Action = "debug/address"
	object := strings.ToLower(req.FormValue("object"))
	id := strings.ToLower(req.FormValue("id"))
	hexAddress := hexDigest(512, "demo")[:6] + "00" + hexDigest(512, id)[:62]

	switch object {
	case "volume", "datablock":
		response.Object = object

		switch object {
		case "volume":
			volMatched, _ := regexp.MatchString(`^[\w]{1,128}$`, req.FormValue("id"))

			if volMatched {
				response.Payload = map[string]string{
					"name":    id,
					"id":      uuid.New().String(),
					"address": hexAddress,
				}
			} else {
				errorHandler(w, 400, "Invalid volume name submitted.")
				return
			}
		case "datablock":
			hashVal := strings.ToLower(req.FormValue("hash"))
			_, err := hex.DecodeString(hashVal)
			hashLen := len(hashVal)

			if hashLen != 64 || err != nil {
				if hashLen != 64 {
					errorHandler(w, 400, "Invalid hash length submitted")
					return
				}

				if err != nil {
					errorHandler(w, 400, err.Error())
					return
				}

			} else {
				response.Payload = map[string]string{
					"name":    "datalock",
					"id":      id,
					"address": hexAddress,
				}
			}
		}
	case "inode", "directory":

		if _, err := uuid.Parse(id); err != nil {
			errorHandler(w, 400, "Invalid UUID supplied")
			return
		}

		response.Object = object

		switch object {
		case "inode":
			response.Payload = map[string]string{
				"name":    "inode",
				"id":      id,
				"address": hexAddress,
			}
		case "directory":
			response.Payload = map[string]string{
				"name":    "directory",
				"id":      id,
				"address": hexAddress,
			}
		}

	default:
		errorHandler(w, 400, "Invalid object type selected")
		return
	}

	if err := jsonReponse(w, response, 200); err != nil {
		errorHandler(w, 500, "Internal Server Error")
		return
	}

	return
}

// bundle simulates requesting bundle status from the debug API.
// Query Arguments:
// "id" - UUID in xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx format, example: 11111111-1111-1111-1111-111111111111

func bundle(w http.ResponseWriter, req *http.Request) {
	var response returnVal

	response.Action = "debug/bundle"
	response.Object = ""
	id := strings.ToLower(req.FormValue("id"))

	if _, err := uuid.Parse(id); err != nil {
		errorHandler(w, 400, "Invalid UUID supplied")
		return
	} else {

		if rand.Intn(2) == 0 {
			response.Payload = map[string]string{
				"bundle_status": "committed",
			}
		} else {
			response.Payload = map[string]string{
				"bundle_status": "not found",
			}
		}

	}

	if err := jsonReponse(w, response, 200); err != nil {
		errorHandler(w, 500, "Internal Server Error")
		return
	}

	return
}

// datablock simulates requesting the information from a datablock from the debug API.
// Query Arguments:
// "hash" - Hex string, example: 2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea
// "raw" - string, example: ""

func datablock(w http.ResponseWriter, req *http.Request) {
	var response returnVal

	response.Action = "debug/datablock"
	response.Object = "datablock"
	hashVal := strings.ToLower(req.FormValue("hash"))
	_, err := hex.DecodeString(hashVal)
	hashLen := len(hashVal)

	if hashLen != 64 || err != nil {
		if hashLen != 64 {
			errorHandler(w, 400, "Invalid hash length submitted")
			return
		}

		if err != nil {
			errorHandler(w, 400, err.Error())
			return
		}

	} else {
		empty, val := req.URL.Query()["raw"]
		token := make([]byte, 64)
		rand.Read(token)

		if val && empty[0] == "" {
			finalVal := hex.Dump(token)
			response.Payload = map[string]string{"Data": finalVal}
		} else {
			finalVal := token
			response.Payload = map[string][]byte{"Data": finalVal}
		}
	}

	if err := jsonReponse(w, response, 200); err != nil {
		errorHandler(w, 500, "Internal Server Error")
		return
	}

	return
}

// debugVolume simulates requesting volume meta data from the debug API.
// Query Arguments:
// "name" - string, example: "test"

func debugVolume(w http.ResponseWriter, req *http.Request) {
	var response returnVal

	response.Action = "debug/volume"
	response.Object = "volume"
	volMatched, _ := regexp.MatchString(`^[\w]{1,128}$`, req.FormValue("name"))

	if volMatched {
		lastHash := hexDigest(512, "demo")[:6] + "00" + hexDigest(512, req.FormValue("name"))[:62]
		fingerprint := hexDigest(256, req.FormValue("name"))[:32]
		response.Payload = map[string]string{
			"volume name":      req.FormValue("name"),
			"root inode uuid":  uuid.New().String(),
			"last hash":        lastHash,
			"compression type": "LZ4",
			"encryption type":  "AES_GCM",
			"key fingerprint":  fingerprint,
		}

		if err := jsonReponse(w, response, 200); err != nil {
			errorHandler(w, 500, "Internal Server Error")
			return
		}

		return

	} else {
		errorHandler(w, 400, "Invalid volume name submitted.")
		return
	}
}

// directory simulates requesting directory meta data from the debug API.
// Query Arguments:
// "id" - UUID in xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx format, example: 11111111-1111-1111-1111-111111111111

func directory(w http.ResponseWriter, req *http.Request) {
	var response returnVal

	response.Action = "debug/directory"
	response.Object = "directory"
	id := strings.ToLower(req.FormValue("id"))

	if _, err := uuid.Parse(id); err != nil {
		errorHandler(w, 400, "Invalid UUID supplied.")
		return
	} else {
		response.Payload = map[string]map[string]uuid.UUID{"files": {"example.txt": uuid.New(), "demo.txt": uuid.New()}}
	}

	if err := jsonReponse(w, response, 200); err != nil {
		errorHandler(w, 500, "Internal Server Error")
		return
	}

	return
}

// directory simulates requesting directory meta data from the debug API.
// Query Arguments:
// "id" - UUID in xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx format, example: 11111111-1111-1111-1111-111111111111

func inode(w http.ResponseWriter, req *http.Request) {
	var response returnVal

	response.Action = "debug/inode"
	response.Object = "inode"
	id := strings.ToLower(req.FormValue("id"))

	if _, err := uuid.Parse(id); err != nil {
		errorHandler(w, 400, "Invalid UUID supplied.")
		return
	}

	response.Payload = map[string]string{
		"uuid":           uuid.New().String(),
		"file mode":      "0777",
		"user id":        "1000",
		"group id":       "1000",
		"size":           "1024",
		"created time":   time.Now().String(),
		"modified time":  time.Now().String(),
		"accessed time":  time.Now().String(),
		"directory uuid": uuid.New().String(),
	}

	if err := jsonReponse(w, response, 200); err != nil {
		errorHandler(w, 500, "Internal Server Error")
		return
	}

	return
}

// Wrapper simulates querying the debug API for TFS wrapper data about various data structures.
// Query Arguments:
// "object" - string, accepts the following values: "inode", "datablock", "volume", "directory"
// "id" - UUID in xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx format, example: 11111111-1111-1111-1111-111111111111
// --or-- string, example: datablock hash - 2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea
// "dump" - string, example: ""

func wrapper(w http.ResponseWriter, req *http.Request) {
	var response returnVal

	response.Action = "debug/wrapper"
	response.Object = "wrapper"
	object := strings.ToLower(req.FormValue("object"))
	id := strings.ToLower(req.FormValue("id"))

	switch object {
	case "datablock":
		break
	case "inode", "directory":

		if _, err := uuid.Parse(id); err != nil {
			errorHandler(w, 400, "Invalid UUID supplied.")
			return
		}

	default:
		errorHandler(w, 400, "Invalid object type selected.")
		return
	}

	data := ""
	empty, val := req.URL.Query()["dump"]

	if val && empty[0] == "" {
		token := make([]byte, 64)
		rand.Read(token)
		data = hex.Dump(token)
	}

	response.Payload = map[string]string{
		"crc":              "0xFFFFFFFF",
		"compression type": "LZ4",
		"encryption type":  "AES_GCM",
		"key fingerprint":  "2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea",
		"data size":        "64",
		"data":             data,
	}

	if err := jsonReponse(w, response, 200); err != nil {
		errorHandler(w, 500, "Internal Server Error")
		return
	}

	return
}

// snapShot simulates creating a volume snapshot or listing a volume snapshot from the TFS API.
// Query Arguments:
// "create" - string - name of new snapshot to create
// "volume" - string - name of volume that snapshot represents
// "list" - string, example: ""

func snapshot(w http.ResponseWriter, req *http.Request) {
	var response returnVal

	response.Action = "snapshot"
	response.Object = "snapshot"
	snapMatched, _ := regexp.MatchString(`^[\w]{1,128}$`, req.FormValue("create"))
	snapshotNameLength := len(req.FormValue("create"))
	volMatched, _ := regexp.MatchString(`^[\w]{1,128}$`, req.FormValue("volume"))
	volumeNameLength := len(req.FormValue("volume"))

	if !volMatched && volumeNameLength == 0 {
		errorHandler(w, 400, "Volume name must be supplied.")
		return
	} else if !volMatched && volumeNameLength > 0 {
		errorHandler(w, 400, "Invalid volume name submitted.")
		return
	}

	if snapMatched {
		createSnapshot := make(map[string]string)
		createSnapshot["name"] = req.FormValue("create")

		if volMatched {
			createSnapshot["volume"] = req.FormValue("volume")
		} else if volumeNameLength == 0 {
			errorHandler(w, 400, "Volume name must be supplied.")
			return
		} else if volumeNameLength > 0 {
			errorHandler(w, 400, "Invalid volume name submitted.")
			return
		}

		response.Payload = createSnapshot

		if err := jsonReponse(w, response, 200); err != nil {
			errorHandler(w, 500, "Internal Server Error")
			return
		}

		return
	} else if snapshotNameLength > 0 {
		errorHandler(w, 400, "Invalid snapshot name submitted.")
		return
	}

	empty, val := req.URL.Query()["list"]

	if val && empty[0] == "" {

		if volMatched {
			response.Payload = map[string]map[string]string{
				"Data": {
					"volume":       req.FormValue("volume"),
					"testSnapshot": time.Now().String(),
					"demoSnapshot": time.RFC822,
				},
			}

		}
		if err := jsonReponse(w, response, 200); err != nil {
			errorHandler(w, 500, "Internal Server Error")
			return
		}
		return
	} else {
		errorHandler(w, 400, "\"list\" parameter does not accept a value")
		return
	}

}

// volume simulates creating or listing a volume from the TFS API.
// Query Arguments:
// "create" - string - name of new snapshot to create
// "compression" - string - example "LZ4", "none"
// "encryption" - string - example "aes_gcm", "none"
// "fingerprint" - string - example "2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea"
// "list" - string, example: ""

func volume(w http.ResponseWriter, req *http.Request) {
	var response returnVal

	response.Action = "volume"
	response.Object = "volume"
	matched, _ := regexp.MatchString(`^[\w]{1,128}$`, req.FormValue("create"))
	volumeNameLength := len(req.FormValue("create"))

	if matched {
		createVolume := make(map[string]string)
		createVolume["name"] = req.FormValue("create")
		validCompression := map[string]int{"lz4": 2, "none": 1}

		if validCompression[strings.ToLower(req.FormValue("compression"))] > 0 {
			createVolume["compression"] = req.FormValue("compression")
		} else {
			errorHandler(w, 400, "Invalid compression type supplied")
			return
		}

		validEncryption := map[string]int{"aes_gcm": 2, "none": 1}

		if validEncryption[strings.ToLower(req.FormValue("encryption"))] > 0 {
			createVolume["encryption"] = req.FormValue("encryption")
			hashVal := req.FormValue("fingerprint")
			if strings.ToLower(req.FormValue("encryption")) != "none" {
				_, err := hex.DecodeString(hashVal)
				hashLen := len(hashVal)

				if hashLen != 64 || err != nil {
					if hashLen != 64 {
						errorHandler(w, 400, "Invalid hash length submitted")
						return
					}

					if err != nil {
						errorHandler(w, 400, err.Error())
						return
					}
				}
			}

			createVolume["fingerprint"] = req.FormValue("fingerprint")
		} else {
			errorHandler(w, 400, "Invalid encryption type supplied.")
			return
		}

		response.Payload = createVolume

		if err := jsonReponse(w, response, 200); err != nil {
			errorHandler(w, 500, "Internal Server Error")
			return
		}
		return

	} else if volumeNameLength > 0 {
		errorHandler(w, 400, "Invalid volume name submitted.")
		return
	}

	empty, val := req.URL.Query()["list"]

	if val && empty[0] == "" {
		response.Payload = map[string]map[string]map[string]string{
			"volumes": {
				"exampleVol0": {
					"compression": "LZ4",
					"encryption":  "AES-GCM",
					"fingerprint": "2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea",
				},
				"exampleVol1": {
					"compression": "none",
					"encryption":  "none",
					"fingerprint": "none",
				},
			},
		}

		if err := jsonReponse(w, response, 200); err != nil {
			errorHandler(w, 500, "Internal Server Error")
			return
		}

		return
	} else if len(empty) > 0 && len(empty[0]) > 0 {
		errorHandler(w, 400, "\"list\" parameter does not accept a value")
		return
	} else {
		errorHandler(w, 400, "no valid parameters supplied.")
		return
	}

}

// jsonResponse is a helper function to return JSON to client.

func jsonReponse(w http.ResponseWriter, response interface{}, status int) error {

	if resp, err := json.Marshal(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("charset", "UTF-8")
		w.WriteHeader(status)
		_, err = w.Write(resp)
		if err != nil {
			return err
		}
	}

	return nil
}

// hexDigest is a helper function to compute hex strings.

func hexDigest(bit int, str string) string {
	var hashVal hash.Hash

	if bit == 512 {
		hashVal = sha512.New()
	} else {
		hashVal = sha256.New()
	}

	hashVal.Write([]byte(str))
	hashBytes := hashVal.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

// errorHandler is a helper function to construct error messages returned via the API.

func errorHandler(w http.ResponseWriter, code int, message string) {
	output := fmt.Sprintf("Error %d: %s", code, message)
	http.Error(w, output, code)
	log.Errorf(output)
	return
}
