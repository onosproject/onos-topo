package reader

import (
    "io/ioutil"
    "gopkg.in/yaml.v3"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

type Underlay struct {
	Networks []Network `yaml:"underlay_networks,flow"`
}

type Network struct {
    Entity_Id       string      `yaml:"entity_id"`
    Display_Name    string      `yaml:"display_name"`
    Switches        []Switches  `yaml:"switches,flow"`
	Links			[]Link		`yaml:"links,flow"`
}

type Switches struct {
    Entity_Id               string          `yaml:"entity_id"`
    Display_Name 			string	        `yaml:"display_name"`
	Model_Id				string	        `yaml:"model_id"`
	Role					string	        `yaml:"role"`
	Management_Endpoint		string	        `yaml:"management_endpoint"`
	P4RT_Server_Endpoint	string	        `yaml:"p4rt_server_endpoint"`
	TLS_insecure			int	        	`yaml:"tls_insecure"`
    Ports                 	[]Ports   		`yaml:"ports,flow"`
}

type Ports struct {
	Entity_Id		string	`yaml:"entity_id"`
	Display_Name	string	`yaml:"display_name"`
	Speed			string	`yaml:"speed"`
	Port_Number		int		`yaml:"port_number"`
	Channel_Number	int		`yaml:"channel_number"`
}

type Link struct {
	Source				string	`yaml:"src"`
	Destination			string	`yaml:"dst"`
	Link_Type			string	`yaml:"link_type"`
}

var log = logging.GetLogger()

func ReadFile(filename string) Underlay {
	// reading in the human-readable schema
	var result Underlay
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(err.Error())
		return result
    }
	err = yaml.Unmarshal(content, &result)
    if err != nil {
        log.Fatal("Failed to parse file ", err)
    }
	return result
}
