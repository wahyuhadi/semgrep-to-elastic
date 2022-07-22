package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projectdiscovery/gologger"
)

var (
// elasticURI   = os.Getenv("elastic_uri")
// indexElastic = "semgrep"
// Username     = os.Getenv("elastic_user")
// Password     = os.Getenv("elastic_pass")
)

type ElasticModels struct {
	Timestamp time.Time `json:"timestamp"`
	RepoURI   string    `json:"repo_url"`
	CheckID   string    `json:"check_id"`
	End       struct {
		Col    int `json:"col"`
		Line   int `json:"line"`
		Offset int `json:"offset"`
	} `json:"end"`
	Extra struct {
		Fingerprint string `json:"fingerprint"`
		IsIgnored   bool   `json:"is_ignored"`
		Lines       string `json:"lines"`
		Message     string `json:"message"`
		Metadata    struct {
		} `json:"metadata"`
		Metavars struct {
			VAR struct {
				AbstractContent string `json:"abstract_content"`
				End             struct {
					Col    int `json:"col"`
					Line   int `json:"line"`
					Offset int `json:"offset"`
				} `json:"end"`
				Start struct {
					Col    int `json:"col"`
					Line   int `json:"line"`
					Offset int `json:"offset"`
				} `json:"start"`
				UniqueID struct {
					Md5Sum string `json:"md5sum"`
					Type   string `json:"type"`
				} `json:"unique_id"`
			} `json:"$VAR"`
		} `json:"metavars"`
		Severity string `json:"severity"`
	} `json:"extra"`
	Path  string `json:"path"`
	Start struct {
		Col    int `json:"col"`
		Line   int `json:"line"`
		Offset int `json:"offset"`
	} `json:"start"`
}

func main() {
	r := gin.Default()
	r.POST("/elastic", elastic)
	r.Run(":8888") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func elastic(c *gin.Context) {
	var obj ElasticModels
	obj.Timestamp = time.Now()
	err := c.ShouldBind(&obj)
	if err != nil {
		c.String(http.StatusForbidden, "error binding body")
		return
	}
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	// resp, err := Elastic(obj)

	f, err := os.OpenFile("logs/semgrep.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		fmt.Println("Error push to API ", err.Error())
		c.String(http.StatusForbidden, "Error push into elastic")
		return
	}
	defer f.Close()
	if _, err := f.WriteString(string(b) + "\n"); err != nil {
		log.Println(err)
		fmt.Println("Error push to API ", err.Error())
		c.String(http.StatusForbidden, "Error push into elastic")
		return
	}

	if err != nil {
		fmt.Println("Error push to API ", err.Error())
		c.String(http.StatusForbidden, "Error push into elastic")
		return
	}
	gologger.Info().Str("Info", "Elastic").Msg("Error")

	c.String(http.StatusOK, "success")

}

// func Elastic(data ElasticModels) (esapi.Response, error) {
// 	cfg := elasticsearch.Config{
// 		Addresses: []string{elasticURI},
// 		Username:  Username,
// 		Password:  Password,
// 		Transport: &http.Transport{
// 			MaxIdleConnsPerHost:   10,
// 			ResponseHeaderTimeout: time.Second,
// 			TLSClientConfig: &tls.Config{
// 				MinVersion: tls.VersionTLS11,
// 				// ...
// 			},
// 		},
// 	}
// 	c, err := elasticsearch.NewClient(cfg)
// 	// PushexamplePushData(c)
// 	if err != nil {
// 		gologger.Info().Str("Info", "Elastic auth").Msg("Problem with auth to elastic")
// 	}

// 	data.Timestamp = time.Now()
// 	datas := esutil.NewJSONReader(data)
// 	// Push data to elastic
// 	fmt.Println(datas)

// 	resp, err := es.PushData(c, indexElastic, datas)

// 	// Error handling when input data to elastic serach
// 	return esapi.Response(*resp), err

// }
