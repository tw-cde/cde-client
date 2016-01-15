package cmd
import (
	"fmt"
	"github.com/cde/apisdk/api"
	"github.com/cde/apisdk/net"
	"github.com/cde/client/config"
)

// AppCreate creates an app.
func AppCreate(name string, stack string, memory int, disk int, instances int) error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	appParams := api.AppParams{
		Name: name,
		Stack: stack,
		Mem: memory,
		Disk:disk,
		Instances:instances}
	fmt.Println(appParams)
	createdApp, err := appRepository.Create(appParams)
	fmt.Println(createdApp)
	return err
}

func AppsList() error {
	configRepository := config.NewConfigRepository(func(error) {})
	appRepository := api.NewAppRepository(configRepository,
		net.NewCloudControllerGateway(configRepository))
	apps, err := appRepository.GetApps()
	if err != nil {
		return err
	}
	fmt.Printf("=== Apps%s", len(apps.Items()))

	for _, app := range apps.Items() {
		fmt.Println(app)
	}
	return nil
}