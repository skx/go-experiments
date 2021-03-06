/**
 * Implementation for 'publishr serve'.
 *
 * This starts a HTTP-server which accepts uploads,
 * and serves downloads.
 */

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/dgryski/dgoogauth"
	"github.com/gorilla/mux"
	"github.com/rakyll/magicmime"
	"github.com/speps/go-hashids"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
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
func (r cmd_serve) help(extended bool) string {
	short := "Launch our HTTP-server."
	if extended {
		fmt.Printf("%s\n", short)
		fmt.Printf("Extra Options:\n")
		fmt.Printf("  --port=N Specify the port to listen upon. (8081)\n")
		fmt.Printf("  --host=N Specify the IP to listen upon.  (127.0.0.1)\n")
		fmt.Printf("\n")
	}

	return short
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
 * Get the remote IP-address, taking account of X-Forwarded-For.
 */
func getRemoteIP(r *http.Request) string {

	hdrForwardedFor := r.Header.Get("X-Forwarded-For")
	if hdrForwardedFor != "" {
		parts := strings.Split(hdrForwardedFor, ",")
		// TODO: should return first non-local address
		return parts[0]
	}

	// Fall-back
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return (ip)
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

	//  Remove any suffix that might be present.
	extension := filepath.Ext(fname)
	fname = fname[0 : len(fname)-len(extension)]

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
			if _, err = io.Copy(outfile, infile); nil != err {
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
			md := &UploadMetaData{MIME: mimetype, IP: getRemoteIP(req), AT: time.Now().Format(time.RFC850)}
			data_json, _ := json.Marshal(md)

			var meta *os.File
			defer meta.Close()
			if meta, err = os.Create("./public/" + sn + ".meta"); nil != err {
				status = http.StatusInternalServerError
				return
			}
			meta.WriteString(string(data_json)) //mimetype)

			//
			// Write out the redirection - using the host
			// scheme, and the end-point of the new upload.
			//
			hostname := req.Host
			scheme := "http"

			if strings.HasPrefix(req.Proto, "HTTPS") {
				scheme = "https"
			}
			if req.Header.Get("X-Forwarded-Proto") == "https" {
				scheme = "https"
			}

			res.Write([]byte(scheme + "://" + hostname + "/get/" + sn + "\n"))
		}
	}
}

/**
 * Called via GET /robots.txt
 */
func RobotsHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "User-agent: *\nDisallow: /")
}

/**
 * Fallback handler, returns 404 for all requests.
 */
func MissingHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(res, "publishr - 404 - content is not hosted here.")
}

/**
 * Launch our HTTP-server.
 */
func (r cmd_serve) execute(args ...string) int {

	f1 := flag.NewFlagSet("f1", flag.ContinueOnError)
	port := f1.String("port", "8081", "The port to bind to.")
	host := f1.String("host", "127.0.0.1", "The host to listen upon.")
	f1.Parse(args)

	/* Create a router */
	router := mux.NewRouter()

	/* Get a previous upload */
	router.HandleFunc("/get/{id}", GetHandler).Methods("GET")

	/* Post a new one */
	router.HandleFunc("/upload", UploadHandler).Methods("POST")

	/* Robots.txt handler */
	router.HandleFunc("/robots.txt", RobotsHandler).Methods("GET")

	/* Error-Handler - Return a 404 on all requests */
	router.PathPrefix("/").HandlerFunc(MissingHandler)

	/* Load the routers beneath the server root */
	http.Handle("/", router)

	/* Build up the bind-string */
	bind := *host + ":" + *port

	/* Launch the server */
	fmt.Printf("Launching the server on http://%s\n", bind)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", *host, *port), nil)
	if err != nil {
		panic(err)
	}
	return 0
}

func init() {
	CMDS = append(CMDS, cmd_serve{})
}
