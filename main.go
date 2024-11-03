package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	certmanagermetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
)

var GroupName = os.Getenv("GROUP_NAME")

type lexiconCommandOptions struct {
	Provider    string
	authUser    string
	authToken   string
	command     []string
	usePassword bool
}
type DNSLexiconDnsRecord struct {
	RecordType string `json:"type"`
	Name       string `json:"name"`
	TTL        int16  `json:"ttl"`
	Content    string `json:"content"`
	Id         string `json:"id"`
}

type DNSLexiconCredentials struct {
	key   string
	token string
}

// func lexiconList(cmd lexiconCommandOptions) ([]DNSLexiconDnsRecord, error) {
// 	var cmdArgs = []string{
// 		cmd.Provider,
// 		"--auth-username", cmd.authUser,
// 		"--auth-token", cmd.authToken,
// 	}
// 	cmdArgs = append(cmdArgs, cmd.command...)
// 	cmdArgs = append(cmdArgs,
// 		"--output", "JSON",
// 	)
// 	lexCmd := exec.Command(
// 		"lexicon",
// 		cmdArgs...,
// 	)

// 	output, err := lexCmd.CombinedOutput()
// 	if err != nil {
// 		return nil, fmt.Errorf("error running lexicon command. %v", err)
// 		// fmt.Printf("Error running command. %v\n", err)
// 	}
// 	cfg := []DNSLexiconDnsRecord{}
// 	if err := json.Unmarshal(output, &cfg); err != nil {
// 		return cfg, fmt.Errorf("error decoding solver config: %v", err)
// 	}

// 	return cfg, nil
// }

func lexiconCmd(cmd lexiconCommandOptions) (bool, error) {

	pwParam := "--auth-token"
	if cmd.usePassword {
		pwParam = "--auth-password"
	}

	var cmdArgs = []string{
		cmd.Provider,
		"--auth-username", cmd.authUser,
		pwParam, cmd.authToken,
	}
	var err error
	cmdArgs = append(cmdArgs, cmd.command...)
	cmdArgs = append(cmdArgs,
		"--output", "JSON",
	)
	fmt.Println("Lexicon command:", cmdArgs)
	fmt.Println()
	lexCmd := exec.Command(
		"lexicon",
		cmdArgs...,
	)
	var output, stderr bytes.Buffer
	lexCmd.Stdout = &output
	lexCmd.Stderr = &stderr

	err = lexCmd.Run()

	if len(stderr.String()) > 0 {
		fmt.Println("STDERR output: ", stderr.String())
	}

	if err != nil {
		printError(err)
		return false, fmt.Errorf("error running lexicon command. %v", err)
	}

	lines := strings.Split(strings.TrimSpace(output.String()), "\n")

	resultStr := strings.TrimSpace(lines[len(lines)-1])

	fmt.Println("lexicon output:", output.String())

	lexResult, err := strconv.ParseBool(resultStr)

	if err != nil {
		printError(err)
	}
	if !lexResult {
		return false, nil
	}
	return true, nil
}

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	// This will register our custom DNS provider with the webhook serving
	// library, making it available as an API under the provided GroupName.
	// You can register multiple DNS provider implementations with a single
	// webhook, where the Name() method will be used to disambiguate between
	// the different implementations.
	cmd.RunWebhookServer(GroupName,
		&DNSLexiconDNSProviderSolver{},
	)
}

// DNSLexiconDNSProviderSolver implements the provider-specific logic needed to
// 'present' an ACME challenge TXT record for your own DNS provider.
// To do so, it must implement the `github.com/jetstack/cert-manager/pkg/acme/webhook.Solver`
// interface.
type DNSLexiconDNSProviderSolver struct {
	// If a Kubernetes 'clientset' is needed, you must:
	// 1. uncomment the additional `client` field in this structure below
	// 2. uncomment the "k8s.io/client-go/kubernetes" import at the top of the file
	// 3. uncomment the relevant code in the Initialize method below
	// 4. ensure your webhook's service account has the required RBAC role
	//    assigned to it for interacting with the Kubernetes APIs you need.
	client *kubernetes.Clientset
}

// DNSLexiconDNSProviderConfig is a structure that is used to decode into when
// solving a DNS01 challenge.
// This information is provided by cert-manager, and may be a reference to
// additional configuration that's needed to solve the challenge for this
// particular certificate or issuer.
// This typically includes references to Secret resources containing DNS
// provider credentials, in cases where a 'multi-tenant' DNS solver is being
// created.
// If you do *not* require per-issuer or per-certificate configuration to be
// provided to your webhook, you can skip decoding altogether in favour of
// using CLI flags or similar to provide configuration.
// You should not include sensitive information here. If credentials need to
// be used by your provider here, you should reference a Kubernetes Secret
// resource and fetch these credentials using a Kubernetes clientset.
type DNSLexiconDNSProviderConfig struct {
	// Change the two fields below according to the format of the configuration
	// to be decoded.
	// These fields will be set by users in the
	// `issuer.spec.acme.dns01.providers.webhook.config` field.

	// Provider to use with DNS lexicon
	Provider string `json:"provider"`

	APIKeyRef    certmanagermetav1.SecretKeySelector `json:"apiKeyRef"`
	APISecretRef certmanagermetav1.SecretKeySelector `json:"apiSecretRef"`
	TTL          *int                                `json:"ttl"`
	Sandbox      bool                                `json:"sandbox"`
	UsePassword  bool                                `json:"usePassword"`
	//Secrets directly in config - not recomended -> use secrets!
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource.
// This should be unique **within the group name**, i.e. you can have two
// solvers configured with the same Name() **so long as they do not co-exist
// within a single webhook deployment**.
// For example, `cloudflare` may be used as the name of a solver.
func (c *DNSLexiconDNSProviderSolver) Name() string {
	return "lexicon"
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (c *DNSLexiconDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	fmt.Printf("\n>>>Present: fqdn:[%s] zone:[%s]\n", ch.ResolvedFQDN, ch.ResolvedZone)
	cfg, err := c.loadConfig(ch)
	if err != nil {
		printError(err)
		return err
	}

	// TODO: do something more useful with the decoded configuration
	fmt.Printf("\n\nDecoded configuration %v\n", cfg)

	// set a record in the DNS provider

	success, err := lexiconCmd(lexiconCommandOptions{
		authUser:    cfg.APIKey,
		authToken:   cfg.APISecret,
		Provider:    cfg.Provider,
		usePassword: cfg.UsePassword,
		command: []string{
			"create",
			ch.ResolvedZone,
			"TXT",
			"--name",
			ch.ResolvedFQDN,
			"--ttl",
			strconv.Itoa(*cfg.TTL),
			"--content=" + ch.Key,
		},
	})

	if !success {
		return fmt.Errorf("lexicon result returned false")
	}
	return err
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (c *DNSLexiconDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	fmt.Printf("\n>>>Cleanup: fqdn:[%s] zone:[%s]\n", ch.ResolvedFQDN, ch.ResolvedZone)
	cfg, err := c.loadConfig(ch)
	if err != nil {
		printError(err)
		return err
	}

	// clears a record in the DNS provider

	lexiconCmd(lexiconCommandOptions{
		authUser:    cfg.APIKey,
		authToken:   cfg.APISecret,
		Provider:    cfg.Provider,
		usePassword: cfg.UsePassword,
		command: []string{
			"delete",
			ch.ResolvedZone,
			"TXT",
			"--name",
			ch.ResolvedFQDN,
			"--content",
			ch.Key,
		},
	})

	return nil
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (c *DNSLexiconDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	///// UNCOMMENT THE BELOW CODE TO MAKE A KUBERNETES CLIENTSET AVAILABLE TO
	///// YOUR CUSTOM DNS PROVIDER

	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		printError(err)
		return err
	}

	c.client = cl

	///// END OF CODE TO MAKE KUBERNETES CLIENTSET AVAILABLE
	return nil
}

// loadConfig is a small helper function that decodes JSON configuration into
// the typed config struct.
func (c *DNSLexiconDNSProviderSolver) loadConfig(ch *v1alpha1.ChallengeRequest) (DNSLexiconDNSProviderConfig, error) {
	cfg := DNSLexiconDNSProviderConfig{}
	cfgJSON := ch.Config
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	//API Key
	apiKey := cfg.APIKey
	if apiKey == "" {
		ref := cfg.APIKeyRef
		if ref.Key == "" || ref.Name == "" {
			return cfg, fmt.Errorf("no apiKeyRef for %q in secret '%s/%s'", ref.Name, ref.Key, ch.ResourceNamespace)
		}
		secret, err := c.client.CoreV1().Secrets(ch.ResourceNamespace).Get(context.TODO(), ref.Name, metav1.GetOptions{})
		if err != nil {
			return cfg, err
		}
		apiKeyRef, ok := secret.Data[ref.Key]
		if !ok {
			return cfg, fmt.Errorf("no apiKeyRef for %q in secret '%s/%s'", ref.Name, ref.Key, ch.ResourceNamespace)
		}
		apiKey = string(apiKeyRef)
		cfg.APIKey = apiKey
	}

	//API Secret
	apiSecret := cfg.APISecret
	if apiSecret == "" {
		ref := cfg.APISecretRef
		if ref.Key == "" || ref.Name == "" {
			return cfg, fmt.Errorf("no apiSecretRef for %q in secret '%s/%s'", ref.Name, ref.Key, ch.ResourceNamespace)
		}
		secret, err := c.client.CoreV1().Secrets(ch.ResourceNamespace).Get(context.TODO(), ref.Name, metav1.GetOptions{})
		if err != nil {
			return cfg, err
		}
		apiSecretRef, ok := secret.Data[ref.Key]
		if !ok {
			return cfg, fmt.Errorf("no accessKeySecret for %q in secret '%s/%s'", ref.Name, ref.Key, ch.ResourceNamespace)
		}
		apiSecret = string(apiSecretRef)
		cfg.APISecret = apiSecret
	}

	return cfg, nil
}

func printError(err error) {
	fmt.Printf("\n\nERROR\n %v \n\n", err)
}
