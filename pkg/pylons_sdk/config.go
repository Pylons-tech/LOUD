package pylonssdk

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
)

const DefaultHome = "~/.loud"

type Config struct {
	From          string
	OutputFormat  string
	Height        int64
	TrustNode     bool
	BroadcastMode string
}

func InitConfig() {

}

func loadConfigFile(cfgFile string) (*toml.Tree, error) {
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s does not exist\n", cfgFile)
		return toml.Load(``)
	}

	bz, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	tree, err := toml.LoadBytes(bz)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func saveConfigFile(cfgFile string, tree *toml.Tree) error {
	fp, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = tree.WriteTo(fp)
	return err
}
