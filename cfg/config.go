package cfg

import (
	"flag"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/BurntSushi/toml"
)

var (
	configPath        = flag.String("config", "config/config.toml", "path to streaming configuration")
	staticsPath       = flag.String("statics", "", "Path to directory with statics")
	videosPath        = flag.String("videos", "", "Path to directory with videos")
	postersPath       = flag.String("poster", "", "Path to directory with posters")
	host              = flag.String("host", "", "host address")
	httpDrainInterval = flag.String("http-drain-interval", "", "Http drain interval")
)

// Config holds information about the streaming configuration.
type Config struct {
	Host              string
	HttpDrainInterval string
	StaticsPath       string
	VideosPath        string
	PostersPath       string
	Build 			  Build
}

// Build holds information about the build.
type Build struct {
	Version    string `json:"version"`
	BuildTime  string `json:"buildTime"`
	LastCommit Commit `json:"lastCommit"`
}

// Commit holds information about last git commit.
type Commit struct {
	Id     string `json:"id"`
	Time   string `json:"time"`
	Author string `json:"author"`
}


func LoadCfg() (*Config, error) {
	flag.Parse()

	bytes, err := ioutil.ReadFile(*configPath)
	if err != nil {
		return nil, err
	}

	conf := Config{}
	if err := toml.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}

	if *staticsPath != "" {
		conf.StaticsPath = *staticsPath
	}

	if *videosPath != "" {
		conf.VideosPath = *videosPath
	}

	if *postersPath != "" {
		conf.PostersPath = *postersPath
	}

	if *host != "" {
		conf.Host = *host
	}

	if *httpDrainInterval != "" {
		conf.HttpDrainInterval = *httpDrainInterval
	}

	return &conf, nil
}

func (c Config) Print() {

	s := reflect.ValueOf(&c).Elem()
	typeOfT := s.Type()

	log.Println("[info] configuration:")
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		log.Printf("[info] %s %s = %v \n", typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}
