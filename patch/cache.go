package patch

import (
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
	// Save saves the res.Files generated content to a temporary directory
	Save(res *pluginpb.CodeGeneratorResponse) error
	// Load loads the cached files in the res.Files
	Load(res *pluginpb.CodeGeneratorResponse) error
	// CleanResFiles cleans the files added from cache into the response files
	CleanResFiles(res *pluginpb.CodeGeneratorResponse)
}

type cache struct {
	tmp string
	originals []string
}

func NewCache() Cache {
	return &cache{tmp: os.TempDir()}
}

func (c *cache) Save(res *pluginpb.CodeGeneratorResponse) error {
	for _, v := range res.File {
		fpath := path.Join(c.tmp, *v.Name)
		dir := path.Dir(fpath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Printf("ERROR failed to create dir %s: %v\n", dir, err)
			return err
		}
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
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		res.File = append(res.File, &pluginpb.CodeGeneratorResponse_File{
			Name:           proto.String(strings.Replace(path, c.tmp, "", 1)),
			Content:        proto.String(string(b)),
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
		for _, vv := range c.originals {
			if *v.Name == vv {
				out = append(out, v)
			}
		}
	}
	res.File = out
	return
}
