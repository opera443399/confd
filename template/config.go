package template

import (
	"os"

	"github.com/kelseyhightower/memkv"
	"github.com/opera443399/confd/backends"
)

// Config for template
type Config struct {
	ConfDir       string
	ConfigDir     string
	KeepStageFile bool
	SyncOnly      bool
	Noop          bool
	Prefix        string
	StoreClient   backends.StoreClient
	TemplateDir   string
}

// TemplateResourceConfig holds the parsed template resource.
type TemplateResourceConfig struct {
	TemplateResource TemplateResource `toml:"template"`
}

// TemplateResource is the representation of a parsed template resource.
type TemplateResource struct {
	Src           string
	Dest          string
	FileMode      os.FileMode
	Mode          string
	Uid           int
	Gid           int
	StageFile     *os.File
	keepStageFile bool
	syncOnly      bool
	noop          bool
	Prefix        string
	Keys          []string
	storeClient   backends.StoreClient
	store         memkv.Store
	funcMap       map[string]interface{}
	lastIndex     uint64
	CheckCmd      string `toml:"check_cmd"`
	ReloadCmd     string `toml:"reload_cmd"`
}
