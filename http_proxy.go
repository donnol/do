package do

import (
	"bufio"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

type HTTPProxyOption struct {
	Director       func(req *http.Request)
	ModifyResponse func(r *http.Response) error
	ErrorHandler   func(w http.ResponseWriter, r *http.Request, err error)
}

func (opt *HTTPProxyOption) complete(rp *httputil.ReverseProxy) {
	if rp == nil {
		return
	}
	if opt != nil {
		if opt.Director != nil {
			rp.Director = opt.Director
		}
		if opt.ModifyResponse != nil {
			rp.ModifyResponse = opt.ModifyResponse
		}
		if opt.ErrorHandler != nil {
			rp.ErrorHandler = opt.ErrorHandler
		}
	}
}

// HTTPProxy listen localAddr and transfer any request to remoteAddr. We can use handlers to specify one custom func to transfer data.
func HTTPProxy(localAddr, remoteAddr string, opt *HTTPProxyOption) (err error) {
	url, err := url.Parse(remoteAddr)
	if err != nil {
		return err
	}

	rp := httputil.NewSingleHostReverseProxy(url)
	opt.complete(rp)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rp.ServeHTTP(w, r)
	})

	s := &http.Server{
		Addr:    localAddr,
		Handler: mux,
	}
	return s.ListenAndServe()
}

// createCert creates a new certificate/private key pair for the given domains,
// signed by the parent/parentKey certificate. hoursValid is the duration of
// the new certificate's validity.
func createCert(dnsNames []string, parent *x509.Certificate, parentKey crypto.PrivateKey, hoursValid int) (cert []byte, priv []byte) {
	privateKey := Must1(ecdsa.GenerateKey(elliptic.P256(), rand.Reader))

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber := Must1(rand.Int(rand.Reader, serialNumberLimit))

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Sample MITM proxy"},
		},
		DNSNames:  dnsNames,
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Duration(hoursValid) * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	derBytes := Must1(x509.CreateCertificate(rand.Reader, &template, parent, &privateKey.PublicKey, parentKey))
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCert == nil {
		panic(fmt.Errorf("failed to encode certificate to PEM"))
	}

	privBytes := Must1(x509.MarshalPKCS8PrivateKey(privateKey))
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		panic(fmt.Errorf("failed to encode key to PEM"))
	}

	return pemCert, pemKey
}

// loadX509KeyPair loads a certificate/key pair from files, and unmarshals them
// into data structures from the x509 package. Note that private key types in Go
// don't have a shared named interface and use `any` (for backwards
// compatibility reasons).
func loadX509KeyPair(certFile, keyFile string) (cert *x509.Certificate, key any, err error) {
	cf, err := os.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}

	kf, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	certBlock, _ := pem.Decode(cf)
	cert, err = x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	keyBlock, _ := pem.Decode(kf)
	key, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}

// mitmProxy is a type implementing http.Handler that serves as a MITM proxy
// for CONNECT tunnels. Create new instances of mitmProxy using createMitmProxy.
type mitmProxy struct {
	caCert *x509.Certificate
	caKey  any
}

// createMitmProxy creates a new MITM proxy. It should be passed the filenames
// for the certificate and private key of a certificate authority trusted by the
// client's machine.
func createMitmProxy(caCertFile, caKeyFile string) *mitmProxy {
	caCert, caKey := Must2(loadX509KeyPair(caCertFile, caKeyFile))
	log.Printf("loaded CA certificate and key; IsCA=%v\n", caCert.IsCA)

	return &mitmProxy{
		caCert: caCert,
		caKey:  caKey,
	}
}

func (p *mitmProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		p.proxyConnect(w, req)
	} else {
		http.Error(w, "this proxy only supports CONNECT", http.StatusMethodNotAllowed)
	}
}

// proxyConnect implements the MITM proxy for CONNECT tunnels.
func (p *mitmProxy) proxyConnect(w http.ResponseWriter, proxyReq *http.Request) {
	log.Printf("CONNECT requested to %v (from %v)", proxyReq.Host, proxyReq.RemoteAddr)

	// "Hijack" the client connection to get a TCP (or TLS) socket we can read
	// and write arbitrary data to/from.
	hj, ok := w.(http.Hijacker)
	if !ok {
		panic(fmt.Errorf("http server doesn't support hijacking connection"))
	}

	clientConn, _ := Must2(hj.Hijack())

	// proxyReq.Host will hold the CONNECT target host, which will typically have
	// a port - e.g. example.org:443
	// To generate a fake certificate for example.org, we have to first split off
	// the host from the port.
	host, _ := Must2(net.SplitHostPort(proxyReq.Host))

	// Create a fake TLS certificate for the target host, signed by our CA. The
	// certificate will be valid for 10 days - this number can be changed.
	pemCert, pemKey := createCert([]string{host}, p.caCert, p.caKey, 240)
	tlsCert := Must1(tls.X509KeyPair(pemCert, pemKey))

	// Send an HTTP OK response back to the client; this initiates the CONNECT
	// tunnel. From this point on the client will assume it's connected directly
	// to the target.
	Must1(clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n")))
	log.Printf("connect ok")

	// Configure a new TLS server, pointing it at the client connection, using
	// our certificate. This server will now pretend being the target.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:               tls.VersionTLS13,
		Certificates:             []tls.Certificate{tlsCert},
	}

	tlsConn := tls.Server(clientConn, tlsConfig)
	defer tlsConn.Close()

	// Create a buffered reader for the client connection; this is required to
	// use http package functions with this connection.
	connReader := bufio.NewReader(tlsConn)

	// Run the proxy in a loop until the client closes the connection.
	for {
		log.Printf("into read write")

		// Read an HTTP request from the client; the request is sent over TLS that
		// connReader is configured to serve. The read will run a TLS handshake in
		// the first invocation (we could also call tlsConn.Handshake explicitly
		// before the loop, but this isn't necessary).
		// Note that while the client believes it's talking across an encrypted
		// channel with the target, the proxy gets these requests in "plain text"
		// because of the MITM setup.
		r, err := http.ReadRequest(connReader)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		// We can dump the request; log it, modify it...
		if b, err := httputil.DumpRequest(r, false); err == nil {
			log.Printf("incoming request:\n%s\n", string(b))
		}

		// Take the original request and changes its destination to be forwarded
		// to the target server.
		changeRequestToTarget(r, proxyReq.Host)

		log.Printf("will reqeust %s", r.URL)
		// Send the request to the target server and log the response.
		resp := Must1(http.DefaultClient.Do(r))
		if b, err := httputil.DumpResponse(resp, false); err == nil {
			log.Printf("target response:\n%s\n", string(b))
		}
		defer resp.Body.Close()

		// Send the target server's response back to the client.
		if err := resp.Write(tlsConn); err != nil {
			log.Println("error writing response back:", err)
		}
	}
}

// changeRequestToTarget modifies req to be re-routed to the given target;
// the target should be taken from the Host of the original tunnel (CONNECT)
// request.
func changeRequestToTarget(req *http.Request, targetHost string) {
	targetUrl := addrToUrl(targetHost)
	targetUrl.Path = req.URL.Path
	targetUrl.RawQuery = req.URL.RawQuery
	req.URL = targetUrl
	// Make sure this is unset for sending the request through a client
	req.RequestURI = ""
}

func addrToUrl(addr string) *url.URL {
	if !strings.HasPrefix(addr, "https") {
		addr = "https://" + addr
	}
	u := Must1(url.Parse(addr))
	return u
}

type (
	HTTPSProxyOption struct {
		CertFile, KeyFile     string
		CaCertFile, CaKeyFile string
	}
)

func HTTPSProxy(localAddr string, opt *HTTPSProxyOption) (err error) {
	if opt == nil {
		return fmt.Errorf("should pass cert and key file")
	}
	if opt.CertFile == "" || opt.KeyFile == "" || opt.CaCertFile == "" || opt.CaKeyFile == "" {
		return fmt.Errorf("should pass cert and key file")
	}

	proxy := createMitmProxy(opt.CaCertFile, opt.CaKeyFile)

	log.Println("Starting https proxy server on", localAddr)
	return http.ListenAndServeTLS(localAddr, opt.CertFile, opt.KeyFile, proxy)
}
