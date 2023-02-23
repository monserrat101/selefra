package module

import (
	"fmt"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/utils"
)

// ------------------------------------------------- --------------------------------------------------------------------

// One of the root-level blocks
type SelefraBlock struct {

	// Name of project
	Name string `yaml:"name" mapstructure:"name"`

	// selefra CloudBlock-related configuration
	CloudBlock *CloudBlock `yaml:"cloud" mapstructure:"cloud"`

	// The version of the cli used by the project
	CliVersion string `yaml:"cli_version" mapstructure:"cli_version"`

	// Global log level. This level is used when the provider does not specify a log level
	LogLevel string `yaml:"log_level" mapstructure:"log_level"`

	//What are the providers required for operation
	RequireProvidersBlock RequireProvidersBlock `yaml:"providers" mapstructure:"providers"`

	// The configuration required to connect to the database
	ConnectionBlock *ConnectionBlock `yaml:"connection" mapstructure:"connection"`

	*LocatableImpl `yaml:"-"`
}

var _ Block = &SelefraBlock{}
var _ MergableBlock[*SelefraBlock] = &SelefraBlock{}

func NewSelefraBlock() *SelefraBlock {
	return &SelefraBlock{
		LocatableImpl: NewLocatableImpl(),
	}
}

func (x *SelefraBlock) Merge(other *SelefraBlock) (*SelefraBlock, *schema.Diagnostics) {

	diagnostics := schema.NewDiagnostics()
	mergedSelefraBlock := &SelefraBlock{}

	// CloudBlock
	if x.CloudBlock != nil && other.CloudBlock != nil {
		// TODO block location
		diagnostics.AddErrorMsg("selefra.cloud block duplicated")
	} else if x.CloudBlock != nil {
		mergedSelefraBlock.CloudBlock = x.CloudBlock
	} else {
		mergedSelefraBlock.CloudBlock = other.CloudBlock
	}

	// Name
	if x.Name != "" && other.Name != "" {
		// TODO block location
		diagnostics.AddErrorMsg("selefra.name duplicated")
	} else if x.Name != "" {
		mergedSelefraBlock.Name = x.Name
	} else {
		mergedSelefraBlock.Name = other.Name
	}

	// CliVersion
	if x.CliVersion != "" && other.CliVersion != "" {
		// TODO block location
		diagnostics.AddErrorMsg("selefra.cli_version duplicated")
	} else if x.CliVersion != "" {
		mergedSelefraBlock.CliVersion = x.CliVersion
	} else {
		mergedSelefraBlock.CliVersion = other.CliVersion
	}

	// LogLevel
	if x.LogLevel != "" && other.LogLevel != "" {
		// TODO block location
		diagnostics.AddErrorMsg("selefra.log_level duplicated")
	} else if x.LogLevel != "" {
		mergedSelefraBlock.LogLevel = x.LogLevel
	} else {
		mergedSelefraBlock.LogLevel = other.LogLevel
	}

	// only RequireProvidersBlock can merge
	if x.RequireProvidersBlock != nil && other.RequireProvidersBlock != nil {
		merge, d := x.RequireProvidersBlock.Merge(other.RequireProvidersBlock)
		diagnostics.AddDiagnostics(d)
		if d == nil || !d.HasError() {
			mergedSelefraBlock.RequireProvidersBlock = merge
		}
	} else if x.RequireProvidersBlock != nil {
		mergedSelefraBlock.RequireProvidersBlock = x.RequireProvidersBlock
	} else {
		mergedSelefraBlock.RequireProvidersBlock = other.RequireProvidersBlock
	}

	// ConnectionBlock
	if x.ConnectionBlock != nil && other.ConnectionBlock != nil {
		// TODO block location
		diagnostics.AddErrorMsg("selefra.connection duplicated")
	} else if x.ConnectionBlock != nil {
		mergedSelefraBlock.ConnectionBlock = x.ConnectionBlock
	} else {
		mergedSelefraBlock.ConnectionBlock = other.ConnectionBlock
	}

	return mergedSelefraBlock, diagnostics
}

func (x *SelefraBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {

	diagnostics := schema.NewDiagnostics()

	// The local name of the project
	if x.Name == "" {
		diagnostics.AddErrorMsg("selefra.cloud.name can not be empty")
	}

	// cloud block is optional, but if it is configured, it must be legal
	if x.CloudBlock != nil {
		diagnostics.AddDiagnostics(x.CloudBlock.Check(module, validatorContext))
	}

	if x.ConnectionBlock != nil {
		x.ConnectionBlock.Check(module, validatorContext)
	}

	if len(x.RequireProvidersBlock) == 0 {
		diagnostics.AddErrorMsg("selefra.providers can not be empty")
	} else {
		diagnostics.AddDiagnostics(x.RequireProvidersBlock.Check(module, validatorContext))
	}

	return diagnostics
}

func (x *SelefraBlock) IsEmpty() bool {
	return x.Name == "" &&
		(x.CloudBlock == nil || x.CloudBlock.IsEmpty()) &&
		x.CliVersion == "" &&
		x.LogLevel == "" &&
		len(x.RequireProvidersBlock) == 0 &&
		x.ConnectionBlock == nil
}

// ------------------------------------------------- --------------------------------------------------------------------

// CloudBlock CloudBlock-related configuration
type CloudBlock struct {

	// Which project in the cloud is associated with
	Project string `yaml:"project" mapstructure:"project"`

	//
	Organization string `yaml:"organization" mapstructure:"organization"`

	// Debug parameters, temporarily masked
	HostName string `yaml:"hostname" mapstructure:"hostname"`

	*LocatableImpl `yaml:"-"`
}

var _ Block = &CloudBlock{}

func NewCloudBlock() *CloudBlock {
	return &CloudBlock{
		LocatableImpl: NewLocatableImpl(),
	}
}

func (x *CloudBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {
	diagnostics := schema.NewDiagnostics()

	// check project name
	if x.Project == "" {
		diagnostics.AddErrorMsg("selefra.cloud.project can not be empty")
	}

	return diagnostics
}

func (x *CloudBlock) IsEmpty() bool {
	return x.Project == "" && x.Organization == "" && x.HostName == ""
}

// ------------------------------------------------- --------------------------------------------------------------------

// ConnectionBlock for db connection
// Example:
//
//	connection:
//	   type: postgres
//	   username: postgres
//	   password: pass
//	   host: localhost
//	   port: 5432
//	   database: postgres
//	   sslmode: disable
type ConnectionBlock struct {
	// These params are mutually exclusive with DSN
	Type     string   `yaml:"type,omitempty" json:"type,omitempty"`
	Username string   `yaml:"username,omitempty" json:"username,omitempty"`
	Password string   `yaml:"password,omitempty" json:"password,omitempty"`
	Host     string   `yaml:"host,omitempty" json:"host,omitempty"`
	Port     *uint64  `yaml:"port,omitempty" json:"port,omitempty"`
	Database string   `yaml:"database,omitempty" json:"database,omitempty"`
	SSLMode  string   `yaml:"sslmode,omitempty" json:"sslmode,omitempty"`
	Extras   []string `yaml:"extras,omitempty" json:"extras,omitempty"`

	*LocatableImpl `yaml:"-"`
}

var _ Block = &ConnectionBlock{}

func NewConnectionBlock() *ConnectionBlock {
	return &ConnectionBlock{
		LocatableImpl: NewLocatableImpl(),
	}
}

// ParseConnectionBlockFromDSN convert dsn to connection block
func ParseConnectionBlockFromDSN(dsn string) *ConnectionBlock {
	// TODO
	return nil
}

func (x *ConnectionBlock) BuildDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=%s", x.Host, x.Username, x.Password, x.Port, x.Database, x.SSLMode)
}

func (x *ConnectionBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {
	diagnostics := schema.NewDiagnostics()

	if x.Type == "" {
		// TODO block location
		diagnostics.AddErrorMsg("connection.type can not be empty")
	}

	if x.Host == "" {
		// TODO block location
		diagnostics.AddErrorMsg("connection.host can not be empty")
	}

	// Add safety Tips
	if x.Username != "" && x.Password == "" {
		// TODO block location
		diagnostics.AddWarn("For security reasons, it is not recommended that you use an empty password when connecting to the database")
	}

	if x.Database == "" {
		// TODO block location
		diagnostics.AddErrorMsg("connection.database can not be empty")
	}

	return diagnostics
}

func (x *ConnectionBlock) IsEmpty() bool {
	return x.Type == "" &&
		x.Username == "" &&
		x.Password == "" &&
		x.Host == "" &&
		x.Port == nil &&
		x.Database == "" &&
		x.SSLMode == "" &&
		len(x.Extras) == 0
}

// ------------------------------------------------- --------------------------------------------------------------------

type RequireProvidersBlock []*RequireProviderBlock

var _ MergableBlock[RequireProvidersBlock] = &RequireProvidersBlock{}
var _ Block = &RequireProvidersBlock{}

func (x RequireProvidersBlock) ToNameMap() map[string]*RequireProviderBlock {
	m := make(map[string]*RequireProviderBlock)
	for _, r := range x {
		m[r.Name] = r
	}
	return m
}

func (x RequireProvidersBlock) Merge(other RequireProvidersBlock) (RequireProvidersBlock, *schema.Diagnostics) {

	diagnostics := schema.NewDiagnostics()

	providerNameSet := make(map[string]struct{})
	mergedRequireProvidersBlock := make(RequireProvidersBlock, 0)

	// merge self
	for _, requireProviderBlock := range x {
		if _, exists := providerNameSet[requireProviderBlock.Name]; exists {
			// TODO block location
			diagnostics.AddErrorMsg("require provider name duplicated")
			continue
		}
		providerNameSet[requireProviderBlock.Name] = struct{}{}
		mergedRequireProvidersBlock = append(mergedRequireProvidersBlock, requireProviderBlock)
	}

	// merge other
	for _, requireProviderBlock := range other {
		if _, exists := providerNameSet[requireProviderBlock.Name]; exists {
			// TODO block location
			diagnostics.AddErrorMsg("require provider name duplicated")
			continue
		}
		providerNameSet[requireProviderBlock.Name] = struct{}{}
		mergedRequireProvidersBlock = append(mergedRequireProvidersBlock, requireProviderBlock)
	}

	return mergedRequireProvidersBlock, diagnostics
}

func (x RequireProvidersBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {

	diagnostics := schema.NewDiagnostics()

	providerNameSet := make(map[string]struct{})
	providerSourceSet := make(map[string]struct{})

	for _, requireProviderBlock := range x {

		if _, exists := providerNameSet[requireProviderBlock.Name]; exists {
			// TODO block location
			diagnostics.AddErrorMsg("require provider name duplicated")
			continue
		}
		providerNameSet[requireProviderBlock.Name] = struct{}{}

		if _, exists := providerSourceSet[requireProviderBlock.Source]; exists {
			// TODO block location
			diagnostics.AddErrorMsg("require provider name duplicated")
			continue
		}
		providerSourceSet[requireProviderBlock.Name] = struct{}{}

		diagnostics.AddDiagnostics(requireProviderBlock.Check(module, validatorContext))
	}

	return diagnostics
}

func (x RequireProvidersBlock) IsEmpty() bool {
	return len(x) == 0
}

func (x RequireProvidersBlock) GetNodeLocation(selector string) *NodeLocation {
	panic("not supported")
}

func (x RequireProvidersBlock) SetNodeLocation(selector string, nodeLocation *NodeLocation) error {
	panic("not supported")
}

// ------------------------------------------------- --------------------------------------------------------------------

// RequireProviderBlock Specifies the version of the Provider to be installed
type RequireProviderBlock struct {

	// The name of this constraint
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// Where does the Provider load from
	Source string `yaml:"source,omitempty" json:"source,omitempty"`

	// Version requirements for this provider
	Version string `yaml:"version,omitempty" json:"version,omitempty"`

	// The debug parameter, if configured, uses the given path instead of downloading
	Path string `yaml:"path" json:"path"`

	//runtime *RequireProviderBlockRuntime
	*LocatableImpl `yaml:"-"`
}

var _ Block = &RequireProviderBlock{}

//var _ HaveRuntime[*RequireProviderBlockRuntime] = &RequireProviderBlock{}

func NewRequireProviderBlock() *RequireProviderBlock {
	x := &RequireProviderBlock{
		LocatableImpl: NewLocatableImpl(),
	}
	//x.runtime = NewRequireProviderBlockRuntime(x)
	return x
}

func (x *RequireProviderBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {

	diagnostics := schema.NewDiagnostics()

	if x.Name == "" {
		// TODO block location
		diagnostics.AddErrorMsg("selefra.providers.name can not be empty")
	}

	if x.Source == "" {
		// TODO block location
		diagnostics.AddErrorMsg("selefra.providers.source can not be empty")
	}

	if x.Version == "" {
		// TODO block location
		diagnostics.AddErrorMsg("selefra.providers.version can not be empty")
	}

	// check file is exists
	if x.Path != "" {
		if !utils.ExistsFile(x.Path) {
			// TODO block location
			diagnostics.AddErrorMsg("selefra.providers.path file not found")
		}
	}

	//diagnostics.AddDiagnostics(x.runtime.check())

	return diagnostics
}

func (x *RequireProviderBlock) IsEmpty() bool {
	return x.Name == "" && x.Source == "" && x.Version == "" && x.Path == ""
}

//func (x *RequireProviderBlock) Runtime() *RequireProviderBlockRuntime {
//	return x.runtime
//}

// ------------------------------------------------- --------------------------------------------------------------------
//
//type RequireProviderBlockRuntime struct {
//	block *RequireProviderBlock
//
//	// Parsed version constraint
//	Constraints version.Constraints
//}
//
//func NewRequireProviderBlockRuntime(block *RequireProviderBlock) *RequireProviderBlockRuntime {
//	return &RequireProviderBlockRuntime{
//		block:       block,
//		Constraints: nil,
//	}
//}
//
//func (x *RequireProviderBlockRuntime) check() *schema.Diagnostics {
//	return x.ensureConstraints()
//}
//
//// IsConstraintsAllow Determines whether the given version conforms to the version constraint
//func (x *RequireProviderBlockRuntime) IsConstraintsAllow(version *version.Version) (bool, *schema.Diagnostics) {
//	d := x.ensureConstraints()
//	if utils.HasError(d) {
//		return false, d
//	}
//
//	// Any version can meet the constraints
//	for _, c := range x.Constraints {
//		if c.Check(version) {
//			return true, nil
//		}
//	}
//	return false, nil
//}
//
//func (x *RequireProviderBlockRuntime) ensureConstraints() *schema.Diagnostics {
//	if x.Constraints != nil {
//		return nil
//	}
//	// Parse the version into structured information
//	constraint, err := version.NewConstraint(x.block.Version)
//	if err != nil {
//		// TODO block location
//		return schema.NewDiagnostics().AddErrorMsg("parse version constraints error")
//	}
//	x.Constraints = constraint
//	return nil
//}

// ------------------------------------------------- --------------------------------------------------------------------

// TODO wait discussion
//type RegistryBlock struct {
//	Type        string
//	Private     bool
//	RegistryUrl string
//	Source      string
//	Token       string
//	TokenEnv    string
//}

// ------------------------------------------------ ---------------------------------------------------------------------
