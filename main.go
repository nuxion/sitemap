package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
)

type route struct {
	Loc        string
	ChangeFreq string `toml:"changefreq"`
	Mobile     bool
	Priority   float32
	Sitemap    bool
	Robot      bool
}

type tomlConfig struct {
	Name        string
	FullURL     string `toml:"fullurl"`
	SitemapPath string `toml:"sitemap_path"`
	Compress    bool
	Output      string `toml:"output_dir"`
	Routes      []route
}

type Robots struct {
	Rules   []string
	Sitemap string
}

func readConfig(fp string) (*tomlConfig, error) {
	var conf tomlConfig
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(content))
	if _, err := toml.Decode(string(content), &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func makeRobot(fp, dst string, rbt *Robots) {
	t, err := template.ParseFiles(fp)
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create(dst)
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	defer f.Close()

	err = t.Execute(f, rbt)
	if err != nil {
		log.Print("execute: ", err)
		return
	}

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	confFile := flag.String("config", "config.toml", "Config file to use")
	flag.Parse()

	conf, err := readConfig(*confFile)
	if err != nil {
		log.Fatal(err)
	}

	sm := stm.NewSitemap(1)

	sm.SetDefaultHost(conf.FullURL)
	// The directory to write sitemaps to locally
	sm.SetPublicPath(conf.Output)
	// The remote host where your sitemaps will be hosted
	// sm.SetSitemapsHost(conf.SitemapPath))
	sm.SetSitemapsPath(conf.SitemapPath)
	sm.SetCompress(conf.Compress)

	// Create method must be called first before adding entries to
	// the sitemap.
	sm.Create()

	for _, r := range conf.Routes {

		if r.Sitemap {

			sm.Add(stm.URL{
				{"loc", r.Loc},
				{"changefreq", r.ChangeFreq},
				{"mobile", r.Mobile},
				{"priority", r.Priority},
			})
		}
	}

	sm.Finalize()

	robot := &Robots{
		Sitemap: fmt.Sprintf("%s%ssitemap.xml", conf.FullURL, conf.SitemapPath),
		Rules:   []string{"Allow: /"},
	}

	for _, r := range conf.Routes {
		if r.Loc != "/" {
			var rule string
			if r.Robot {
				rule = fmt.Sprintf("Allow: /%s", r.Loc)

			} else {
				rule = fmt.Sprintf("Disallow: /%s", r.Loc)
			}

			robot.Rules = append(robot.Rules, rule)
		}
	}
	makeRobot("robots.txt.tpl",
		fmt.Sprintf("%s/robots.txt", conf.Output),
		robot,
	)

}
