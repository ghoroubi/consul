package consul

import (
	"errors"
	"github.com/go-kit/kit/sd"
	consulSd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"strconv"
)

// Register , registers a service in consul discovery service
func Register(consulAddress string,
	serviceAddress string,
	httpPort string,
	grpcPort string,
	serviceID, serviceName string,
	tags []string,
	interval, timeout string) sd.Registrar {

	config := api.DefaultConfig()           // Assign default config to a config variable.
	config.Address = consulAddress          // assign address to the config var.
	apiClient, err := api.NewClient(config) // Create new consul api client.
	if err != nil {
		panic(err)
	}
	// Create new consul
	client := consulSd.NewClient(apiClient)

	check := api.AgentServiceCheck{
		HTTP:     "http://" + serviceAddress + ":" + httpPort + "/health",
		Interval: interval,
		Timeout:  timeout,
		Notes:    "Basic health check",
	}

	// gRPC port casting to integer.
	port, _ := strconv.Atoi(grpcPort)
	asr := api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceAddress,
		Port:    port,
		Tags:    tags,
		Check:   &check,
	}

	// create and return a new registrar.
	return NewRegistrar(client, &asr)
}

// GetServerAddress , returns the address of gRPC service that already has been registered in consul
// Using the service name and one of the service tags
func GetServerAddress(consulAddr, serviceName, tag string,
	onlyHealthy bool, query *api.QueryOptions) (addr string, port string, err error) {
	// define a config.
	config := api.DefaultConfig()

	config.Address = consulAddr          // Assign address
	client, err := api.NewClient(config) // Creating consul api client
	if err != nil {
		return "", "", err
	}

	// Getting service info from consul agent
	services, _, err := client.Health().
		Service(serviceName, tag, onlyHealthy, query)
	if err != nil {
		return "", "", err
	}

	// Check fo empty service slice.
	if len(services) == 0 {
		return "", "", errors.New("no service available.")
	}

	// Extracting address and port from info
	addr = services[0].Service.Address
	p := services[0].Service.Port

	// Convert integer port to string
	port = strconv.Itoa(p)

	return addr, port, nil
}
