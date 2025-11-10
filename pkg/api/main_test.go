package api

import (
	"context"
	"fmt"
	"github.com/dalet-oss/opensearch-cli/pkg/api/types/replication"
	"github.com/dalet-oss/opensearch-cli/pkg/appconfig"
	"github.com/dalet-oss/opensearch-cli/pkg/consts"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/creds"
	"github.com/dalet-oss/opensearch-cli/pkg/utils/fp"
	printutils "github.com/dalet-oss/opensearch-cli/pkg/utils/print"
	"github.com/google/uuid"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	tcopensearch "github.com/testcontainers/testcontainers-go/modules/opensearch"
	"github.com/testcontainers/testcontainers-go/network"
	"strings"
	"sync"
	"testing"
	"time"
)

// this file contains constants and helpers for api testing

const (
	MainContainer   = "main"
	LeaderContainer = "leader"
	// cName used everywhere in the config
	cName = "testcontainers::opensearch"
	// osUsername username in the opensearch server
	osUsername = "admin"
	// osPassword password in the opensearch server
	osPassword = "admin"
	// vaultUserKey username key in the vault
	vaultUserKey = "username"
	// vaultPassword password key in the vault
	vaultPasswordKey = "password"
	// vaultPassword password of the vault
	vaultPassword     = "admin"
	defaultHTTPPort   = 9200
	transportPort     = 9300
	leaderClusterName = "leader-cluster"
	afName            = "test-autofollow"
	testIndex         = "test-index"
)

var (
	testIndexAutofollowRule = replication.CreateAutofollowReq{
		Header: nil,
		Body: replication.CreateAutofollowBody{
			Name:         afName,
			LeaderAlias:  leaderClusterName,
			IndexPattern: testIndex,
		},
	}

	// OSImage image of the opensearch server
	OSImage = "opensearchproject/opensearch"
	// OSVersion version of the opensearch server
	OSVersion = "2.19.3"
	// vaultData data of the vault
	vaultData = map[string]string{
		vaultUserKey:     osUsername,
		vaultPasswordKey: vaultPassword,
	}
	osContainers = []string{MainContainer, LeaderContainer}
	osCtrx       = map[string]*tcopensearch.OpenSearchContainer{}
	// opensearchContainer container of the opensearch server
	opensearchContainer *tcopensearch.OpenSearchContainer
	// config configuration of the opensearch-cli used for tests
	config *appconfig.AppConfig
	// contextWithPassword context with the password of the vault [prevents interactive password ask]
	contextWithPassword = context.WithValue(context.Background(), consts.VaultPasswordFlag, vaultPassword)
)

func getCCR() CCRCreateOpts {
	return CCRCreateOpts{
		Type:       "",
		Mode:       "",
		RemoteName: leaderClusterName,
		RemoteAddr: fmt.Sprintf("%s:%d", InternalIP(osCtrx[LeaderContainer]), transportPort),
	}
}

func spinOpenSearch(netAliases []string, net *testcontainers.DockerNetwork, name string) *tcopensearch.OpenSearchContainer {
	ctx := context.Background()
	ctr, err := tcopensearch.Run(
		ctx,
		fmt.Sprintf("%s:%s", OSImage, OSVersion),
		tcopensearch.WithUsername(osUsername),
		tcopensearch.WithPassword(osPassword),
		network.WithNetwork(netAliases, net),
		testcontainers.WithLabels(map[string]string{
			"org.testcontainers.service":        "opensearch",
			"org.testcontainers.container-name": name,
		}),
	)
	if err != nil {
		log.Fatal().Err(err)
	}
	return ctr
}

func ConfigTContainer(c *tcopensearch.OpenSearchContainer) *appconfig.AppConfig {
	addr, err := c.Address(context.Background())
	if err != nil {
		log.Fatal().Err(err)
	}
	return &appconfig.AppConfig{
		ApiVersion: "v1",
		CliParams: &appconfig.CliParams{
			ServerTimeoutSeconds: fp.AsPointer(120),
		},
		Clusters: []appconfig.ClusterConfig{
			{
				Name: cName,
				Params: appconfig.ClusterParams{
					Server: addr,
					Tls:    false,
				},
			},
		},
		Users: []appconfig.UserConfig{
			{
				Name: cName,
				User: appconfig.User{
					Vault: &appconfig.VaultConfig{
						VaultString: creds.CreateVault(vaultData, vaultPassword),
						Username:    vaultUserKey,
						Password:    vaultPasswordKey,
					},
				},
			},
		},
		Contexts: []appconfig.ContextConfig{
			{
				Name:    cName,
				Cluster: cName,
				User:    cName,
			},
		},
		Current: cName,
	}
}

// testWrapper - generates a wrapper for the opensearch-cli with the given config and contextWithPassword
func testWrapper() *OpensearchWrapper {
	W, err := New(*config, contextWithPassword)
	if err != nil {
		log.Fatal().Err(err)
	}
	return W
}

// wrapperForContainer - generates a wrapper for the opensearch-cli with the given config and contextWithPassword for the given container
func wrapperForContainer(containerName string) *OpensearchWrapper {
	W, err := New(*ConfigTContainer(osCtrx[containerName]), contextWithPassword)
	if err != nil {
		log.Fatal().Err(err)
	}
	return W
}

// InternalIP - returns the internal IP of the container
func InternalIP(ctr *tcopensearch.OpenSearchContainer) string {
	addr, err := ctr.ContainerIP(context.Background())
	if err != nil {
		log.Fatal().Err(err)
	}
	return addr
}

// InternalNetAddr - returns the internal network address of the container
func InternalNetAddr(ctr *tcopensearch.OpenSearchContainer) string {
	addr := InternalIP(ctr)
	return fmt.Sprintf("%s:%d", addr, defaultHTTPPort)
}

type TestCase struct {
	Name          string
	Wrapper       *OpensearchWrapper
	ConfigureFunc func(t *testing.T, c *OpensearchWrapper)
	PostFunc      func(t *testing.T, c *OpensearchWrapper)
	CaseInput     interface{}
	WantErr       bool
}
type ReplicationTesting struct {
	Name                  string
	Wrapper               *OpensearchWrapper
	Shotgun               *shotgun
	DocumentCount         *int
	ConfigureLeaderFunc   func(t *testing.T, c *OpensearchWrapper)
	ConfigureFollowerFunc func(t *testing.T, c *OpensearchWrapper)
	PostLeaderFunc        func(t *testing.T, c *OpensearchWrapper)
	PostFollowerFunc      func(t *testing.T, c *OpensearchWrapper)
	CaseInput             interface{}
	ExtraValidationFunc   func(t *testing.T, execResult any)
	WantErr               bool
}

// OSDocGenerator type alias around a function which generates a document(JSON string)
type OSDocGenerator func() string

// shotgun
// type used exclusively in the testing infrastructure to apply some data to the opensearch
type shotgun struct {
	Wrapper     *OpensearchWrapper
	CreateIndex bool
	IndexName   string
	// function which generates the index data
	GeneratorFunc  OSDocGenerator
	DocIdGenerator func() string
	IgnoreErr      bool
	Delay          time.Duration
	Documents      map[string][]string
	mutex          sync.Mutex
}

// uuidGen - wraps basic UUID generator
func uuidGen() string {
	uuidGenerator, _ := uuid.NewV7()
	return uuidGenerator.String()
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewShotgun(wrapper *OpensearchWrapper, createIndex bool, indexName string, generatorFunc OSDocGenerator, ignoreErr bool, delay time.Duration) *shotgun {
	return &shotgun{
		Wrapper:       wrapper,
		CreateIndex:   createIndex,
		IndexName:     indexName,
		GeneratorFunc: generatorFunc,
		IgnoreErr:     ignoreErr,
		Delay:         delay,
		Documents:     map[string][]string{},
		mutex:         sync.Mutex{},
	}
}

func (s *shotgun) WithIndex(indexName string) *shotgun {
	s.IndexName = indexName
	return s
}

func (s *shotgun) WithDelay(delay time.Duration) *shotgun {
	s.Delay = delay
	return s
}

func (s *shotgun) WithGeneratorFunc(generatorFunc OSDocGenerator) *shotgun {
	s.GeneratorFunc = generatorFunc
	return s
}

// shoot - uploads the given number of documents to the index.
// The documents are generated by the given function.
// to override the target index, pass the index name as the second argument(non empty).
func (s *shotgun) shoot(documentCount int, index *string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.CreateIndex {
		if err := s.Wrapper.CreateIndex(s.IndexName); err != nil {
			log.Warn().Err(err)
		}
	}

	targetIndex := s.IndexName
	if index != nil && len(*index) > 0 {
		targetIndex = *index
	}
	for i := 1; i <= documentCount; i++ {
		ctx, cancelFunc := context.WithTimeout(context.TODO(), 2*time.Second)
		defer cancelFunc()
		documentId := uuidGen()
		documentRequest := opensearchapi.DocumentCreateReq{
			Index:      targetIndex,
			DocumentID: documentId,
			Body:       strings.NewReader(s.GeneratorFunc()),
			Params:     opensearchapi.DocumentCreateParams{},
		}
		var writeResult opensearchapi.DocumentCreateResp
		rsp, createErr := s.Wrapper.Client.Do(ctx, documentRequest, &writeResult)
		if createErr != nil {
			if !s.IgnoreErr {
				return createErr
			}
			log.Warn().Err(createErr)
		}
		if rsp.IsError() {
			if !s.IgnoreErr {
				return fmt.Errorf("fail to push document:%s", rsp.String())
			}
			log.Warn().Msgf("fail to write document:\n%s\nerror:\n%s", printutils.MarshalJSONOrDie(documentRequest), rsp.String())
		} else {
			s.Documents[targetIndex] = append(s.Documents[targetIndex], documentId)
			log.Info().Msgf("[%d/%d]document written successfully", i, documentCount)
		}
		time.Sleep(s.Delay)
	}
	return nil
}

func (s *shotgun) Shoot(t *testing.T, count int, index *string) {
	if s.IgnoreErr {
		err := s.shoot(count, index)
		if err != nil {
			t.Logf("ignoring errors:%v", err)
		}
	} else {
		assert.NoError(t, s.shoot(count, index), "no errors expected during the upload of test data")
	}
}
func shotgunBasicDocument() string {
	return fmt.Sprintf(`{"name":"test-index","type":"test-type","unique_field":"%s"}`, uuid.New().String())
}

// TestShotgun test shotgun behaviour
func TestShotgun(t *testing.T) {
	gun := NewShotgun(
		wrapperForContainer(MainContainer),
		true,
		"shotgun-populated-index",
		shotgunBasicDocument,
		true,
		10*time.Millisecond,
	)
	start := time.Now()
	gun.Shoot(t, 100, nil)
	idx, err := gun.Wrapper.GetIndexList()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("index list:\n%v", idx)
	t.Logf("done in %s", time.Since(start))
}
