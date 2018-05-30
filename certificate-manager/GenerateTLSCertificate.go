package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

type Subject struct {
	Organization  []string
	Country       []string
	Province      []string
	Locality      []string
	StreetAddress []string
	PostalCode    []string
	Domain        string
}

func (s *Subject) GenerateCertificate(interfaceName string) {
	var err error

	template := x509.Certificate{
		Subject: pkix.Name{
			Organization: s.Organization,
		},
		NotBefore: time.Now(),

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		// if no ca is given, we have to set IsCA to self sign
		IsCA: true,
	}

	template.Subject.CommonName = GetIPv4ForInterfaceName(interfaceName)
	ip := net.ParseIP(GetIPv4ForInterfaceName(interfaceName))
	//Setting this to resolve "cannot validate certificate for <ip> because it doesn't contain any IP SANs"
	template.IPAddresses = append(template.IPAddresses, ip)

	template.NotAfter = template.NotBefore.Add(time.Duration(365) * time.Hour * 24)
	//anyKey()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		os.Exit(1)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	template.SerialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Println("Failed to generate serial number:", err)
		os.Exit(1)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		fmt.Println("Failed to create certificate:", err)
		os.Exit(1)
	}
	if _, err := os.Stat("certs"); os.IsNotExist(err) {
		if err := os.MkdirAll("certs", os.ModePerm); err != nil {
			fmt.Println("Unable to create certs directory: ", err)
			os.Exit(1)
		}
	}

	certOut, err := os.OpenFile("./certs/server.crt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("Failed to open selfsigned server.crt for writing:", err)
		os.Exit(1)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()


	keyOut, err := os.OpenFile("./certs/server.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("failed to open selfsigned server.key for writing:", err)
		os.Exit(1)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
}

func GetIPv4ForInterfaceName(ifname string) (ip string) {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.Name == ifname {
			if addrs, err := inter.Addrs(); err == nil {
				for _, addr := range addrs {
					switch ip := addr.(type) {
					case *net.IPNet:
						if ip.IP.DefaultMask() != nil {
							return (ip.IP.String())
						}
					}
				}
			}
		}
	}
	return ""
}

func main() {
	ifName := os.Args[1]
	subject := Subject{
		Organization:	[]string{"RVBD"},
		Country:		[]string{"IN"},
		Province:		[]string{"KA"},
		Locality:		[]string{"BLR"},
		StreetAddress:	[]string{},
		PostalCode:		[]string{},
		Domain:			"riverbed.com",
	}
	subject.GenerateCertificate(string(ifName))
}
