package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"sort"

	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"

	"encoding/json"
	"io"

	"path"

	"github.com/gin-gonic/gin"
)

type GitHubResponse struct {
	Sha string `json:"sha"`
}

func downloadCsv() {
	folderPath := "./downloads"
	githubRepo := "sapics/ip-location-db"

	urls := []string{
		fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s", githubRepo, "geo-whois-asn-country/geo-whois-asn-country-ipv4-num.csv"),
		fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s", githubRepo, "geo-asn-country/geo-asn-country-ipv6-num.csv"),
	}

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	for idx := range urls {
		file := path.Base(urls[idx])

		filePath := folderPath + "/" + file
		if _, err := os.Stat(filePath); err == nil {
			localSha, err := getSHA256(filePath)
			if err != nil {
				panic(err)
			}

			remoteSha, err := getRemoteSHA256(githubRepo, file)
			if err != nil {
				panic(err)
			}

			if localSha != remoteSha {
				downloadFile(filePath, urls[idx])
			}
		} else {
			downloadFile(filePath, urls[idx])
		}
	}
}

func getRemoteSHA256(repo, file string) (string, error) {
	resp, err := http.Get("https://api.github.com/repos/" + repo + "/contents/" + file)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var gitHubResponse GitHubResponse
	err = json.NewDecoder(resp.Body).Decode(&gitHubResponse)
	if err != nil {
		return "", err
	}

	return gitHubResponse.Sha, nil
}

func downloadFile(filePath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getSHA256(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(bytes)
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash, nil
}

func Ip2Int(ip net.IP) *big.Int {
	i := big.NewInt(0)
	i.SetBytes(ip)
	return i
}

func ReadAndGet(fn string) []IpItem {
	var ip_items []IpItem = []IpItem{}
	f, _ := os.Open("./downloads" + "/" + fn)
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		start := new(big.Int)
		start, _ = start.SetString(record[0], 10)
		end := new(big.Int)
		end, _ = end.SetString(record[1], 10)
		ip_items = append(ip_items, IpItem{start, end, record[2]})
	}
	f.Close()
	return ip_items
}

func main() {
	downloadCsv()

	var ip_items []IpItem = []IpItem{}
	ip_items = append(ip_items, ReadAndGet("geo-whois-asn-country-ipv4-num.csv")...)
	ip_items = append(ip_items, ReadAndGet("geo-asn-country-ipv6-num.csv")...)

	sort.Slice(ip_items, func(i, j int) bool {
		return ip_items[i].start.Cmp(ip_items[j].start) == -1
	})

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.GET("/getIpInfo", func(c *gin.Context) {
		addr := net.ParseIP(c.Query("addr"))
		if addr != nil {
			ip_num := big.NewInt(0)

			if addr.To4() != nil {
				ip_num = new(big.Int).SetUint64(uint64(binary.BigEndian.Uint32(addr.To4())))
			} else {
				ip_num.SetBytes(addr)
			}
			idx, _ := Binary(ip_items, ip_num, 0, len(ip_items))
			if idx != -1 && ip_num.Cmp(big.NewInt(0)) != 0 {
				c.JSON(http.StatusOK, gin.H{
					"ok":      true,
					"country": ip_items[idx].country,
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"ok": false,
		})
	})
	router.NoMethod(catchAll)
	router.NoRoute(catchAll)
	router.Run(":8080")
}

func catchAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok": false,
	})
}
