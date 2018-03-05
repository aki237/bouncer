package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Bouncer struct {
	sites         map[string]string
	listenAddress string
}

func NewBouncer(la string) *Bouncer {
	return &Bouncer{sites: make(map[string]string), listenAddress: la}
}

func (b *Bouncer) ReadConfig(filename string) error {

	// Read all the contents of the file and handle errors.
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Split the content into lines
	lines := strings.Split(string(bs), "\n")

	// iterate through lines
	for i, line := range lines {
		// If the line starts with a #, it is a comment so ignore it.
		if strings.HasPrefix(strings.TrimSpace(line), "#") ||
			strings.TrimSpace(line) == "" {
			continue
		}

		// get the parts separated by numerous whitespace charaters
		parts := strings.Fields(line)
		// only 2 parts, error... :)
		if len(parts) != 2 {
			return fmt.Errorf("In file %s, at line %d :: parse error", filename, i)
		}

		// First part is error.
		host := parts[0]
		// Second part is error.
		remote := parts[1]

		// if the host already configured, throw error...
		if _, ok := b.sites[host]; ok {
			return fmt.Errorf("In file %s, at line %d :: host configuration already loaded", filename, i)
		}
		// else store it.
		b.sites[host] = remote
	}

	return nil
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+":443"+r.RequestURI, http.StatusMovedPermanently)
}

func (b *Bouncer) Serve(cert string, priv string) {
	if FileExists(cert) && FileExists(priv) {
		log("provided valid certificates and private key files. Enabling SSL support.")
		go func() {
			if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
				fatal(err.Error())
			}
		}()
		err := http.ListenAndServeTLS(":443", cert, priv, http.HandlerFunc(b.bounce))
		if err != nil {
			fatal(err.Error())
		}
	} else {
		http.ListenAndServe(b.listenAddress, http.HandlerFunc(b.bounce))
	}
}

func (b *Bouncer) bounce(w http.ResponseWriter, r *http.Request) {
	local, ok := b.sites[r.Host]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if strings.HasPrefix(local, ":") {
		local = "localhost" + local
	}

	fmt.Println("http://" + local + r.URL.EscapedPath())

	c := &http.Client{}
	req, err := http.NewRequest(r.Method, "http://"+local+r.URL.EscapedPath(), r.Body)
	if err != nil {
		warn(err.Error())
		w.WriteHeader(502)
		w.Write([]byte("Some error occured"))
		return
	}

	for k, s := range r.Header {
		req.Header.Set(k, strings.Join(s, ""))
	}

	resp, err := c.Do(req)
	if err != nil {
		warn(err.Error())
		w.WriteHeader(502)
		w.Write([]byte("Some error occured"))
		return
	}

	for k, s := range resp.Header {
		w.Header().Set(k, strings.Join(s, ""))
	}

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)

	resp.Body.Close()
}
