package patch

import (
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type Cache interface {
	// Exists return true if cache has content for the CodeGeneratorRequest
	Exists() (bool, error)
	// Ensure run protoc-gen-go if no cache already exists
	Ensure(plugin string) error
	// Save saves the res.Files generated content to a temporary directory
	Save(res *pluginpb.CodeGeneratorResponse) error
	// Load loads the cached files in the res.Files
	Load(res *pluginpb.CodeGeneratorResponse) error
	// CleanResFiles cleans the files added from cache into the response files
	CleanResFiles(res *pluginpb.CodeGeneratorResponse)
}

type cache struct {
	tmp       string
	req       *pluginpb.CodeGeneratorRequest
	originals []string
}

func NewCache(req *pluginpb.CodeGeneratorRequest) Cache {
	var content string
	for _, v := range req.ProtoFile {
		content += v.String()
	}
	hash := hash(content)
	tmp := path.Join(os.TempDir(), hash)
	log.Printf("Cache - dir: \t%s\n", tmp)
	return &cache{tmp: tmp, req: proto.Clone(req).(*pluginpb.CodeGeneratorRequest)}
}

func (c *cache) Exists() (bool, error) {
	if _, err := os.Stat(c.tmp); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *cache) Ensure(plugin string) error {
	if plugin == "go" {
		log.Println("Cache -\tPlugin is protoc-gen-go: skipping")
		return nil
	}
	ok, err := c.Exists()
	if err != nil {
		return err
	}
	if ok {
		log.Println("Cache -\tExits")
		return nil
	}
	log.Println("Cache -\tNo go definitions found: generating definitions")
	c.req.Parameter = proto.String("")
	res, err := RunPlugin("go", c.req, nil)
	if err != nil {
		return err
	}
	log.Println("Cache -\tDefinitions generated")
	if err := c.Save(res); err != nil {
		return err
	}
	return nil
}

func (c *cache) Save(res *pluginpb.CodeGeneratorResponse) error {
	for _, v := range res.File {
		fpath := path.Join(c.tmp, *v.Name)
		dir := path.Dir(fpath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Printf("ERROR failed to create dir %s: %v\n", dir, err)
			return err
		}
		log.Printf("Cache - save: \t%s: %s\n", *v.Name, fpath)
		f, err := os.Create(fpath)
		if err != nil {
			log.Printf("ERROR failed to create file %s: %v\n", *v.Name, err)
			return err
		}
		if _, err := f.Write([]byte(*v.Content)); err != nil {
			log.Printf("ERROR failed to write file %s: %v\n", *v.Name, err)
			return err
		}
	}
	return nil
}

func (c *cache) Load(res *pluginpb.CodeGeneratorResponse) error {
	if _, err := os.Stat(c.tmp); err != nil {
		return err
	}
	if len(res.File) == 0 {
		return nil
	}
	for _, v := range res.File {
		c.originals = append(c.originals, *v.Name)
	}
	if err := filepath.Walk(filepath.Join(c.tmp, path.Dir(*res.File[0].Name)), func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		log.Printf("Cache - load:\t%s\n", path)
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(strings.Replace(path, c.tmp, "", 1)),
			Content: proto.String(string(b)),
		})
		return nil
	}); err != nil {
		log.Printf("ERROR failed to add cached files: %v\n", err)
		return err
	}
	return nil
}

func (c *cache) CleanResFiles(res *pluginpb.CodeGeneratorResponse) {
	var out []*pluginpb.CodeGeneratorResponse_File
	for _, v := range res.File {
		if contains(c.originals, *v.Name) {
			out = append(out, v)
		} else {
			log.Printf("Cache - removing:\t%s\n", *v.Name)
		}
	}
	res.File = out
	return
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func hash(s string) string {
	sha := sha1.New()
	sha.Write([]byte(s))
	b := sha.Sum(nil)
	h := base64.RawStdEncoding.EncodeToString(b)
	h = strings.Replace(h, "/", "_", -1)
	h = strings.Replace(h, "+", "-", -1)
	return h
}
