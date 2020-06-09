package sender_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rhdedgar/pleg-watcher/api"
	"github.com/rhdedgar/pleg-watcher/config"
	"github.com/rhdedgar/pleg-watcher/docker"
	"github.com/rhdedgar/pleg-watcher/models"
	. "github.com/rhdedgar/pleg-watcher/sender"
)

// postClamScanResult handles received clamAV scan result data in json format. It's accessed with:
// POST /api/clam/scanresult
func postClamScanResult(c echo.Context) error {
	var scanResult api.ScanResult

	loadData(scanResult, "clam_scan_result_example.json")

	if err := c.Bind(&scanResult); err != nil {
		log.Println("Error binding received scan result data:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process scan result"}
	}

	//fmt.Println("Scan result bound to:")
	log.Printf("%+v", scanResult)

	return c.NoContent(http.StatusOK)
}

// postCrioPodLog handles received crictl inspect data in json format. It's accessed with:
// POST /api/crio/log
func postCrioPodLog(c echo.Context) error {
	var container models.Container

	loadData(container, "crictl_inspect_example.json")

	if err := c.Bind(&container); err != nil {
		log.Println("Error binding received crio data:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process crio content"}
	}

	return c.NoContent(http.StatusOK)
}

// postDockerPodLog handles received docker inspect data in json format. It's accessed with:
// POST /api/docker/log
func postDockerPodLog(c echo.Context) error {
	var container docker.DockerContainer

	loadData(container, "docker_inspect_example.json")

	if err := c.Bind(&container); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process docker content"}
	}

	return c.NoContent(http.StatusOK)
}

// loadData reads the specified file and Unmarshals its JSON contents into the provided data structure
func loadData(ds interface{}, filePath string) error {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
		return err
	}

	err = json.Unmarshal(fileBytes, &ds)
	if err != nil {
		fmt.Println("Error Unmarshaling secrets json: ", err)
		return err
	}
	return nil
}

var _ = Describe("Sender", func() {
	var (
		e = echo.New()
	)

	config.ClamURL = "http://localhost:8080/api/clam/scanresult"
	config.CrioURL = "http://localhost:8080/api/crio/log"
	config.DockerURL = "http://localhost:8080/api/docker/log"

	BeforeEach(func() {
		go func() {
			e = echo.New()

			e.Use(middleware.Logger())
			e.Use(middleware.Recover())

			e.POST("/api/clam/scanresult", postClamScanResult)
			e.POST("/api/crio/log", postCrioPodLog)
			e.POST("/api/docker/log", postDockerPodLog)

			e.Use(middleware.Logger())

			e.Logger.Info(e.Start(":8080"))
		}()
	})

	AfterEach(func() {
		//e.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	})

	Describe("SendClamData", func() {
		Context("Validate clam scan results can be formatted and POSTed as JSON", func() {
			It("Should correctly marshal and POST the clam data as JSON", func() {
				var sRes api.ScanResult

				result, err := SendClamData(sRes)

				Expect(err).To(BeNil())
				Expect(result).To(Equal(200))
			})
		})
	})

	Describe("SendCrioData", func() {
		Context("Validate crio container definitions can be formatted and POSTed as JSON", func() {
			It("Should correctly marshal and POST the crio data as JSON", func() {
				var mStat models.Status

				result, err := SendCrioData(mStat)

				Expect(err).To(BeNil())
				Expect(result).To(Equal(200))
			})
		})
	})

	Describe("SendDockerData", func() {
		Context("Validate docker container definitions can be formatted and POSTed as JSON", func() {
			It("Should correctly marshal and POST the docker data as JSON", func() {
				var dCon docker.DockerContainer

				result, err := SendDockerData(dCon)

				Expect(err).To(BeNil())
				Expect(result).To(Equal(200))
			})
		})
	})
})
