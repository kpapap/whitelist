package whitelist

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Define the structure of the YAML data
// ConfigMapMapping represents the mapping of a ConfigMap to a namespace.
// It contains the name of the ConfigMap and the namespace it belongs to.
type ConfigMapMapping struct {
	// Name is the name of the ConfigMap.
	Name string `yaml:"name"`
	// Namespace is the namespace in which the ConfigMap is located.
	Namespace string `yaml:"namespace"`
}

type nbcmrReceiver struct {
	config *Config
	logger *zap.Logger
	nextConsumer consumer.Logs
}

func (c *nbcmrReceiver) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (c *nbcmrReceiver) ConsumeLogs(ctx context.Context, ld consumer.Logs) error {
	return nil
}

func (r *nbcmrReceiver) Start(ctx context.Context, host component.Host) error {
	// Start the receiver
		// Load Kubernetes cluster configuration
		log.Println("Loading Kubernetes cluster configuration")
		clusterconfig, err := clientcmd.BuildConfigFromFlags("", "")
		if err != nil {
			log.Fatalf("Error building kubeconfig: %s", err.Error())
		}
		log.Println("Kubernetes cluster configuration loaded")
	
		// Create Kubernetes client
		log.Println("Creating Kubernetes client")
		clientset, err := kubernetes.NewForConfig(clusterconfig)
		if err != nil {
			log.Fatalf("Error creating Kubernetes client: %s", err.Error())
		}
		log.Println("Kubernetes client created")
	
		// Read the YAML data from the environment variable
		log.Println("Reading YAML data from environment variable")
		yamlData := os.Getenv("CONFIGMAP_LIST")
		if yamlData == "" {
			log.Fatal("CONFIGMAP_LIST environment variable is not set")
		}
		log.Println("YAML data read from environment variable")
	
		// Decode the YAML data into a struct
		log.Println("Decoding YAML data")
		var configYAML []ConfigMapMapping
		err = yaml.Unmarshal([]byte(yamlData), &configYAML)
		if err != nil {
			log.Fatalf("Error decoding YAML data: %s", err.Error())
		}
		log.Println("YAML data decoded")
	
		// Create a map of ConfigMap names and their corresponding namespaces
		log.Println("Creating map of ConfigMap names and namespaces")
		configMapMap := make(map[string]string)
	
		// Populate the configMapMap from the YAML data
		for _, mapping := range configYAML {
			configMapMap[mapping.Name] = mapping.Namespace
		}
		log.Println("Map of ConfigMap names and namespaces created")
	
		// Create a ticker to repeat the code
		log.Println("Creating ticker")
		// Get the repeat time from environment variable
		repeatTimeStr := os.Getenv("INTERVAL")
		if repeatTimeStr == "" {
			log.Fatal("INTERVAL environment variable is not set")
		}
		repeatTime, err := time.ParseDuration(repeatTimeStr)
		if err != nil {
			log.Fatalf("Error parsing INTERVAL environment variable: %s", err.Error())
		}
		ticker := time.NewTicker(repeatTime)
		defer ticker.Stop()
		log.Println("Ticker created")
	
		// Run the code in a loop with a ticker
		log.Println("Starting loop")
		for range ticker.C {
			log.Println("Listing selected ConfigMaps:")
			for name, namespace := range configMapMap {
				log.Printf("Getting ConfigMap %s in namespace %s", name, namespace)
				configmap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
				if err != nil {
					log.Printf("Error getting ConfigMap %s in namespace %s: %s", name, namespace, err.Error())
					continue
				}
				if configmap != nil {
					log.Printf("Namespace: %s, Name: %s, Data: %v", configmap.Namespace, configmap.Name, configmap.Data)
				}
			}
		}
		return nil
	}

// Shutdown shuts down the receiver.
func (r *nbcmrReceiver) Shutdown(ctx context.Context) error {
	// Shutdown the receiver
	// Log a message indicating that the receiver is shutting down.
	log.Println("Shutting down receiver")
	// Return nil to indicate that the receiver shut down successfully.
	return nil
}

