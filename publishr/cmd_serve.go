package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgryski/dgoogauth"
	"github.com/gorilla/mux"
	"github.com/rakyll/magicmime"
	"github.com/speps/go-hashids"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * types for each sub-command
 */
type cmd_serve struct{}

/**
 * Implementation for "serve"
 */
func (r cmd_serve) name() string {
	return "serve"
}

/**
 * Server.
 */
func (r cmd_serve) help() string {
	return "Launch our HTTP-server."
}

/**
 * Meta-data structure for every uploaded-file.
 *
 * MIME -> The MIME-type of the uploaded file.
 *
 * IP -> The remote host that performed the upload.
 *
 * AT -> The date/time of the upload.
 *
 */
type UploadMetaData struct {
	MIME string `json:"MIME"`
	IP   string `json:"IP"`
	AT   string `json:"AT"`
}

/**
 * Get the next short-ID.
 *
 * Do that by loading the state, increasing the count, and saving it.
 */
func NextShortID() string {

	state, _ := LoadState()

	/**
	 * Increase the count, and hash it.
	 */
	state.Count += 1

	hd := hashids.NewData()
	hd.Salt = "I hope this is secure"
	hd.MinLength = 1
	h := hashids.NewWithData(hd)

	numbers := []int{99}
	numbers[0] = state.Count
	hash, _ := h.Encode(numbers)

	/**
	 * Write out the body
	 */
	SaveState(state)
	return hash
}

/**
 * Called via GET /get/XXXXXX
 */
func GetHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
		}
	}()

	vars := mux.Vars(req)
	fname := vars["id"]
	fname = "./public/" + fname

	if !Exists(fname) || !Exists(fname+".meta") {
		http.NotFound(res, req)
		return
	}

	file, _ := ioutil.ReadFile(fname + ".meta")

	var md UploadMetaData

	if err := json.Unmarshal(file, &md); err != nil {
		status = 500
		err = errors.New("Loading JSON failed")
		return
	}

	/**
	 * Serve the file, with the correct MIME-type
	 */
	res.Header().Set("Content-Type", md.MIME)
	http.ServeFile(res, req, fname)
}

/**
 * Upload a file to ./public/ - with a short-name.
 *
 * Each file will also have a ".meta" file created, to contain
 * some content.
 */
func UploadHandler(res http.ResponseWriter, req *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
		}
	}()

	/**
	 * Get the authentication-header.  If missing we abort.
	 */
	auth := string(req.Header.Get("X-Auth"))
	if len(auth) < 1 {
		status = 401
		err = errors.New("Missing X-Auth header")
		return
	}
	auth = strings.TrimSpace(auth)

	/**
	 * Load the secret.
	 */
	state, err := LoadState()
	if err != nil {
		status = 500
		err = errors.New("Loading state failed")
		return
	}

	/**
	 * Test the token.
	 */
	otpc := &dgoogauth.OTPConfig{
		Secret:      state.Secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	val, err := otpc.Authenticate(auth)
	if err != nil {
		status = 401
		err = errors.New("Failed to use X-Auth header")
		return
	}

	/**
	 * If it failed then we're done.
	 */
	if !val {
		status = 401
		err = errors.New("Invalid X-Auth header")
		return
	}

	/**
	 ** At ths point we know we have an authorized submitter.
	 **/

	/**
	 * Parse the incoming request
	 */
	const _24K = (1 << 20) * 24
	if err = req.ParseMultipartForm(_24K); nil != err {
		status = http.StatusInternalServerError
		return
	}

	/**
	 * Get the short-ID
	 */
	sn := NextShortID()

	for _, fheaders := range req.MultipartForm.File {
		for _, hdr := range fheaders {
			// open uploaded
			var infile multipart.File
			if infile, err = hdr.Open(); nil != err {
				status = http.StatusInternalServerError
				return
			}
			// open destination
			var outfile *os.File
			if outfile, err = os.Create("./public/" + sn); nil != err {
				status = http.StatusInternalServerError
				return
			}
			// 32K buffer copy
			var written int64
			if written, err = io.Copy(outfile, infile); nil != err {
				status = http.StatusInternalServerError
				return
			}

			// Get the MIME-type of the uploaded file.
			err = magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR)
			if err != nil {
				status = http.StatusInternalServerError
				return
			}
			defer magicmime.Close()
			mimetype, _ := magicmime.TypeByFile("./public/" + sn)

			//
			// Write out the meta-data - which is a structure
			// containing the following members.
			//
			md := &UploadMetaData{MIME: mimetype, IP: req.RemoteAddr, AT: time.Now().Format(time.RFC850)}
			data_json, _ := json.Marshal(md)

			var meta *os.File
			defer meta.Close()
			if meta, err = os.Create("./public/" + sn + ".meta"); nil != err {
				status = http.StatusInternalServerError
				return
			}
			meta.WriteString(string(data_json)) //mimetype)

			res.Write([]byte("uploaded file:" + hdr.Filename + ";link " + sn + " ;length:" + strconv.Itoa(int(written)) + " ;mime-type:" + mimetype))
		}
	}
}

/**
 * Launch our HTTP-server.
 */
func (r cmd_serve) execute(args ...string) int {

	/* Create a router */
	router := mux.NewRouter()

	/* Get a previous upload */
	router.HandleFunc("/get/{id}", GetHandler).Methods("GET")

	/* Post a new one */
	router.HandleFunc("/upload", UploadHandler).Methods("POST")

	/* Load the routers beneath the server root */
	http.Handle("/", router)

	/* Launch the server */
	fmt.Printf("Launching the server on http://0.0.0.0:8081\n")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
	return 0
}
