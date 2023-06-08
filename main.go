package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	auth "github.com/abbot/go-http-auth"
	"github.com/gin-gonic/gin"
	flags "github.com/jessevdk/go-flags"
)

// Options
type Options struct {
	Verbose  bool   `short:"v" long:"verbose" description:"display verbose debug output"`
	All      bool   `short:"a" long:"all" description:"allow all (read/write/delete) operations"`
	Read     bool   `short:"r" long:"read" description:"allow read operations"`
	Write    bool   `short:"w" long:"write" description:"allow write operations"`
	Delete   bool   `short:"d" long:"delete" description:"allow delete operations"`
	Listen   string `short:"l" long:"listen" description:"[host]:port to listen on" default:":3137"`
	Passwd   string `short:"p" long:"passwd" description:"htpasswd file for authentication (optional)"`
	Cert     string `short:"c" long:"cert" description:"path to tls certificate file (optional)"`
	Key      string `short:"k" long:"key" description:"path to tls key file (optional)"`
	PostHook string `long:"post-hook" description:"path to (executable) script to run after each operation"`
	Args     struct {
		Directory string `required:"yes"`
	} `positional-args:"yes"`
}

type Env struct {
	Verbose   bool
	Read      bool
	Write     bool
	Delete    bool
	Passwd    string
	PostHook  string
	Directory string
}

// BasicAuth middleware
// https://github.com/gin-gonic/gin/issues/2326
func basicAuth(a *auth.BasicAuth) gin.HandlerFunc {
	if a == nil {
		return func(c *gin.Context) {}
	}

	realmHeader := "Basic realm=" + strconv.Quote(a.Realm)
	return func(c *gin.Context) {
		user := a.CheckAuth(c.Request)
		if user == "" {
			c.Header("WWW-Authenticate", realmHeader)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func (env *Env) getHandler(c *gin.Context) {
	path := env.Directory + "/" + c.Param("filename")
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, "File Not Found")
			return
		}
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	c.File(path)
	c.Status(http.StatusOK)

	if env.PostHook != "" {
		postHook(env.PostHook, "GET", path, env.Verbose)
	}
}

func (env *Env) putHandler(c *gin.Context) {
	defer c.Request.Body.Close()
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	path := env.Directory + "/" + c.Param("filename")
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Status(http.StatusNoContent)

	if env.PostHook != "" {
		postHook(env.PostHook, "PUT", path, env.Verbose)
	}
}

func (env *Env) deleteHandler(c *gin.Context) {
	path := env.Directory + "/" + c.Param("filename")
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, "File Not Found")
			return
		}
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	err = os.Remove(path)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
	c.Status(http.StatusNoContent)

	if env.PostHook != "" {
		postHook(env.PostHook, "DELETE", path, env.Verbose)
	}
}

func checkOptions(opts *Options) {
	if !opts.All && !opts.Read && !opts.Write && !opts.Delete {
		log.Fatal("Error: must specify at least one of [-a|-r|-w|-d]")
	}
	if opts.All {
		opts.Read = true
		opts.Write = true
		opts.Delete = true
	}
	if opts.Cert != "" && opts.Key == "" {
		log.Fatal("Error: must specify --key with --cert")
	}
	if opts.Key != "" && opts.Cert == "" {
		log.Fatal("Error: must specify --cert with --key")
	}
	if opts.PostHook != "" {
		err := checkPostHook(opts.PostHook)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func (env *Env) setupRouter(r *gin.Engine) {
	// Authentication
	var authenticator *auth.BasicAuth
	if env.Passwd != "" {
		htpasswd := auth.HtpasswdFileProvider(env.Passwd)
		authenticator = auth.NewBasicAuthenticator("Protected", htpasswd)
	}
	auth := r.Group("/", basicAuth(authenticator))

	// Routes
	if env.Read {
		auth.GET("/:filename", env.getHandler)
	}
	if env.Write {
		auth.PUT("/:filename", env.putHandler)
	}
	if env.Delete {
		auth.DELETE("/:filename", env.deleteHandler)
	}
}

func main() {
	log.SetFlags(0)
	// Parse default options are HelpFlag | PrintErrors | PassDoubleDash
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}

		// Does PrintErrors work? Is it not set?
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		parser.WriteHelp(os.Stderr)
		os.Exit(2)
	}
	checkOptions(&opts)

	// Check we're not running as root
	if os.Geteuid() == 0 {
		log.Fatal("Error: not to be run as root - use a non-privileged user account instead")
	}

	// Setup environment
	env := Env{
		Verbose:   opts.Verbose,
		Read:      opts.Read,
		Write:     opts.Write,
		Delete:    opts.Delete,
		Passwd:    opts.Passwd,
		PostHook:  opts.PostHook,
		Directory: opts.Args.Directory,
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	env.setupRouter(r)

	// Event loop
	fmt.Println("listening on", opts.Listen)
	if opts.Cert == "" {
		err = r.Run(opts.Listen)
	} else {
		err = r.RunTLS(opts.Listen, opts.Cert, opts.Key)
	}
	if err != nil {
		log.Fatal(err)
	}
}
