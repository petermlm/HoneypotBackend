package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	mathRand "math/rand"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	dockerApiTypes "github.com/docker/docker/api/types"
	dockerContainer "github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	dockerNAT "github.com/docker/go-connections/nat"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

const itemsNum = 1439 // Looks "random" enough
const dockerImage = "docker.elastic.co/elasticsearch/elasticsearch:7.10.1"
const containerName = "temp-es"
const containerHost = "http://127.0.0.1:9250/"
const containerHTTPPortMapping = "127.0.0.1:9250:9200"
const containerCommPortMapping = "127.0.0.1:9350:9300"
const containerReadyRetries = 20
const saltSize = 16
const namesList = "names.txt"

var scripts = []string{"./get_indexes.sh", "./get_docs.sh", "./get_mapping.sh"}
var names []string

const usersMapping = `{
	"properties": {
		"email": { "type": "text" },
		"first_name": { "type": "text" },
		"last_name": { "type": "text" },
		"passhash": { "type": "text" },
		"passsalt": { "type": "text" }
	}
}`

type esItem struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Passhash  string `json:"passhash"`
	Passsalt  string `json:"passsalt"`
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() (ferr error) {
	var err error

	ctx := context.Background()

	names, err = getNamesList()
	if err != nil {
		return err
	}

	// Start docker
	cli, err := createDockerCLI()
	if err != nil {
		return err
	}
	containerID, err := startESDocker(ctx, cli)
	if err != nil {
		return err
	}
	defer func() {
		if err := stopESDocker(ctx, cli, containerID); err != nil {
			ferr = err
		}
	}()

	// Wait until Elasticsearch is ready to be used
	if err := waitForES(ctx, cli); err != nil {
		return err
	}

	// Create Elasticsearch client
	es, err := createClient()
	if err != nil {
		return err
	}
	if err := clear(es); err != nil {
		return err
	}
	bi, err := createBulkIndexer(es)
	if err != nil {
		return err
	}
	bulkIndexerClosed := false
	defer func() {
		if !bulkIndexerClosed {
			if err := closeBulkIndexer(bi); err != nil {
				ferr = err
			}
		}
	}()

	// Write entries to Elasticsearch and get desired output
	if err := addItems(bi); err != nil {
		return err
	}
	if err := closeBulkIndexer(bi); err != nil {
		return err
	}
	bulkIndexerClosed = true
	es.Indices.Refresh()

	for _, script := range scripts {
		if err := getCURLOutput(script); err != nil {
			return err
		}
	}

	return nil
}

func createClient() (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:     []string{containerHost},
		RetryOnStatus: []int{502, 503, 504, 429},
		MaxRetries:    5,
	})
	if err != nil {
		return nil, err
	}
	return es, nil
}

func clear(es *elasticsearch.Client) error {
	res, err := es.Indices.Delete([]string{"users"}, es.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil || res.IsError() {
		return fmt.Errorf("Cannot delete index: %s", err)
	}
	res.Body.Close()

	res, err = es.Indices.Create("users")
	if err != nil {
		return fmt.Errorf("Cannot create index: %s", err)
	}
	if res.IsError() {
		return fmt.Errorf("Cannot create index: %s", res)
	}
	res.Body.Close()

	reader := strings.NewReader(usersMapping)
	res, err = es.Indices.PutMapping([]string{"users"}, reader)
	if err != nil {
		return fmt.Errorf("Cannot put mapping: %s", err)
	}
	res.Body.Close()
	return nil
}

func createBulkIndexer(es *elasticsearch.Client) (esutil.BulkIndexer, error) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         "users",
		Client:        es,
		NumWorkers:    runtime.NumCPU(),
		FlushBytes:    5e+6,
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating the indexer: %s", err)
	}
	return bi, nil
}

func closeBulkIndexer(bi esutil.BulkIndexer) error {
	if err := bi.Close(context.Background()); err != nil {
		return fmt.Errorf("Unexpected error: %s", err)
	}
	return nil
}

func addItems(bi esutil.BulkIndexer) error {
	var ferr error
	for i := 0; i < itemsNum; i++ {
		docID := fmt.Sprintf("%d", i)

		data, err := makeRandomItem()
		if err != nil {
			return fmt.Errorf("Can't generate item: %s", err)
		}

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: docID,
				Body:       bytes.NewReader(data),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					commonMsg := "Error adding item:"
					if err != nil {
						ferr = fmt.Errorf("%s %s", commonMsg, err)
					} else {
						ferr = fmt.Errorf("%s: %s: %s", commonMsg, res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			return fmt.Errorf("Unexpected error: %s", err)
		}
		if ferr != nil {
			break
		}
	}
	return ferr
}

func makeRandomItem() ([]byte, error) {
	firstName := names[mathRand.Intn(len(names))]
	lastName := names[mathRand.Intn(len(names))]
	password := names[mathRand.Intn(len(names))]
	email := fmt.Sprintf("%s.%s@gmail.com", firstName, lastName)

	salt, err := generateRandomSalt()
	if err != nil {
		return nil, err
	}
	hash := hashPassword(password, salt)

	item := esItem{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Passhash:  string(hash),
		Passsalt:  string(salt),
	}

	js, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func getCURLOutput(script string) error {
	// Run command
	cmd := exec.Command(script)
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return err
	}

	// Get output
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)

	// Write it to file
	outputFile := fmt.Sprintf("%s.output", script)
	if err := ioutil.WriteFile(outputFile, buf.Bytes(), 0666); err != nil {
		return fmt.Errorf("Can't write output file of %s: %s", script, err)
	}
	return nil
}

func createDockerCLI() (*docker.Client, error) {
	cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("Can't create docker client: %s", err)
	}
	return cli, nil
}

func startESDocker(ctx context.Context, cli *docker.Client) (string, error) {
	exposedPorts, portBindings, _ := dockerNAT.ParsePortSpecs([]string{
		containerHTTPPortMapping,
		containerCommPortMapping,
	})
	containerConfig := &dockerContainer.Config{
		ExposedPorts: exposedPorts,
		Env:          []string{"discovery.type=single-node"},
		Image:        dockerImage,
	}
	hostConfig := &dockerContainer.HostConfig{
		PortBindings: portBindings,
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("Can't create docker container: %s", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, dockerApiTypes.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("Can't start docker container: %s", err)
	}

	return resp.ID, nil
}

func stopESDocker(ctx context.Context, cli *docker.Client, containerID string) error {
	if err := cli.ContainerKill(ctx, containerID, "SIGKILL"); err != nil {
		return fmt.Errorf("Error killing ES: %s", err)
	}
	if err := cli.ContainerRemove(ctx, containerID, dockerApiTypes.ContainerRemoveOptions{}); err != nil {
		return fmt.Errorf("Error removing ES: %s", err)
	}
	return nil
}

func waitForES(ctx context.Context, cli *docker.Client) error {
	var err error
	for i := containerReadyRetries; i > 0; i-- {
		_, err = http.Get(containerHost)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return fmt.Errorf("Can't connect to Elasticsearch container: %s", err)
	}
	return nil
}

func generateRandomSalt() ([]byte, error) {
	var salt = make([]byte, saltSize)
	_, err := rand.Read(salt[:])
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func hashPassword(password string, salt []byte) string {
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()

	passwordBytes = append(passwordBytes, salt...)
	sha512Hasher.Write(passwordBytes)

	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash
}

func getNamesList() ([]string, error) {
	namesBuf, err := ioutil.ReadFile(namesList)
	if err != nil {
		return nil, fmt.Errorf("Can't read names list: %s", err)
	}
	namesStr := string(namesBuf)
	names := strings.Split(namesStr, "\n")
	return names, nil
}
