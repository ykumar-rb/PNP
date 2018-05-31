package helper

import (
	"github.com/ZTP/pnp/util/color"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net"
	"strings"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		color.Fatalf("Error: ", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		color.Fatalf("Error: ", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		color.Fatalf("Error: ", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		color.Fatalf("Error: ", err)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		color.Fatalf("Error: ", err)
	}
	return plaintext
}

func GetIPv4ForInterfaceName(ifname string) (ifaceip *net.IPNet) {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.Name == ifname {
			if addrs, err := inter.Addrs(); err == nil {
				for _, addr := range addrs {
					switch ip := addr.(type) {
					case *net.IPNet:
						if ip.IP.DefaultMask() != nil {
							return (ip)
						}
					}
				}
			}
		}
	}
	color.Fatal("Check the interface name provided")
	return (nil)
}

func GetIPFromIPwithCIDR(ipCidr string) string {
	ipCidrArr := strings.Split(ipCidr, "/")
	ip := ipCidrArr[0]
	return ip
}

