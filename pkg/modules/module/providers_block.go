package module

import (
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"gopkg.in/yaml.v3"
	"strconv"
)

// ------------------------------------------------- --------------------------------------------------------------------

// ProvidersBlock The root level providers block
type ProvidersBlock []*ProviderBlock

var _ MergableBlock[ProvidersBlock] = (*ProvidersBlock)(nil)
var _ Block = (*ProvidersBlock)(nil)

func (x ProvidersBlock) ToProviderNameMap() map[string]*ProviderBlock {
	m := make(map[string]*ProviderBlock)
	for _, p := range x {
		m[p.Provider] = p
	}
	return m
}

func (x ProvidersBlock) Merge(other ProvidersBlock) (ProvidersBlock, *schema.Diagnostics) {

	diagnostics := schema.NewDiagnostics()

	nameSet := make(map[string]struct{}, 0)
	mergedProviders := make([]*ProviderBlock, 0)

	// merge self
	for _, providerBlock := range x {
		if _, exists := nameSet[providerBlock.Name]; exists {
			// TODO block location
			diagnostics.AddErrorMsg("merge providers block error, find same provider %s", providerBlock.Name)
			continue
		}
		mergedProviders = append(mergedProviders, providerBlock)
	}

	// merge other
	for _, providerBlock := range other {
		if _, exists := nameSet[providerBlock.Name]; exists {
			// TODO block location
			diagnostics.AddErrorMsg("merge providers block error, find same provider %s", providerBlock.Name)
			continue
		}
		mergedProviders = append(mergedProviders, providerBlock)
	}

	return mergedProviders, diagnostics
}

func (x ProvidersBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {
	diagnostics := schema.NewDiagnostics()

	providerNameSet := make(map[string]struct{}, 0)
	for _, providerBlock := range x {
		if _, exists := providerNameSet[providerBlock.Name]; exists {
			// TODO block location
			diagnostics.AddErrorMsg("provider name conflict error, find same provider %s", providerBlock.Name)
			continue
		}
		diagnostics.AddDiagnostics(providerBlock.Check(module, validatorContext))
		providerNameSet[providerBlock.Name] = struct{}{}
	}
	return diagnostics
}

func (x ProvidersBlock) IsEmpty() bool {
	return len(x) == 0
}

func (x ProvidersBlock) GetNodeLocation(selector string) *NodeLocation {
	//TODO implement me
	panic("implement me")
}

func (x ProvidersBlock) SetNodeLocation(selector string, nodeLocation *NodeLocation) error {
	//TODO implement me
	panic("implement me")
}

// ------------------------------------------------- --------------------------------------------------------------------

// ProviderBlock An element in the providers block array at the root level
type ProviderBlock struct {

	// Name of the current block
	Name string

	// How long is the cache
	Cache string

	// Which of the selefra.providers is associated with
	Provider string

	// What is the maximum concurrency when pulling data
	MaxGoroutines *uint64

	// What resources need to be pulled? If you do not write, the default is to pull all resources
	Resources []string

	// What are the self-defined configurations of the provider? These should be passed to the provider through
	ProvidersConfigYamlString string

	*LocatableImpl `yaml:"-"`
}

var _ yaml.Marshaler = (*ProviderBlock)(nil)
var _ Block = (*ProviderBlock)(nil)

func NewProviderBlock() *ProviderBlock {
	return &ProviderBlock{
		LocatableImpl: NewLocatableImpl(),
	}
}

func (x *ProviderBlock) MarshalYAML() (interface{}, error) {
	configurationMappingNode := &yaml.Node{
		Kind: yaml.MappingNode,
	}

	// name
	configurationMappingNode.Content = append(configurationMappingNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "name"})
	configurationMappingNode.Content = append(configurationMappingNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: x.Name})

	// cache
	configurationMappingNode.Content = append(configurationMappingNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "cache"})
	configurationMappingNode.Content = append(configurationMappingNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: x.Cache})

	// provider
	configurationMappingNode.Content = append(configurationMappingNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "provider"})
	configurationMappingNode.Content = append(configurationMappingNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: x.Provider})

	// max_goroutines
	if x.MaxGoroutines != nil {
		configurationMappingNode.Content = append(configurationMappingNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "max_goroutines"})
		configurationMappingNode.Content = append(configurationMappingNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: strconv.Itoa(int(*x.MaxGoroutines))})
	}

	// resources
	if len(x.Resources) != 0 {
		configurationMappingNode.Content = append(configurationMappingNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: "resources"})
		resourcesNode := &yaml.Node{Kind: yaml.SequenceNode}
		for _, resourceName := range x.Resources {
			resourcesNode.Content = append(resourcesNode.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: resourceName})
		}
		configurationMappingNode.Content = append(configurationMappingNode.Content, resourcesNode)
	}

	if x.ProvidersConfigYamlString != "" {
		var customProviderConfiguration yaml.Node
		err := yaml.Unmarshal([]byte(x.ProvidersConfigYamlString), &customProviderConfiguration)
		if err != nil {
			return nil, err
		}
		for _, node := range customProviderConfiguration.Content[0].Content {
			configurationMappingNode.Content = append(configurationMappingNode.Content, node)
		}
	}

	return configurationMappingNode, nil
	//out, err := yaml.Marshal(configurationMappingNode)
	//if err != nil {
	//	return nil, err
	//}
	//s := string(out)
	//return strings.TrimLeft(s, "|"), nil
}

// If the provider is not configured, this is the default configuration
func GetDefaultProviderConfigYamlConfiguration(providerName, providerVersion string) string {
	block := ProviderBlock{
		Name:          "default-" + providerName,
		Cache:         "1d",
		Provider:      providerName,
		MaxGoroutines: pointer.ToUInt64Pointer(50),
	}
	out, _ := yaml.Marshal(block)
	return string(out)
}

func (x *ProviderBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {

	diagnostics := schema.NewDiagnostics()

	if x.Name == "" {
		diagnostics.AddErrorMsg("provider block name must not be empty")
	}

	if x.Cache != "" {
		// TODO check cache style
	}

	if x.Provider == "" && !module.HasRequiredProviderName(x.Provider) {
		diagnostics.AddErrorMsg("")
	}

	if x.MaxGoroutines != nil {
		if *x.MaxGoroutines > 3000 {

		} else if *x.MaxGoroutines < 0 {

		}
	}

	// TODO check other field

	return diagnostics
}

func (x *ProviderBlock) IsEmpty() bool {
	return x.Name == "" &&
		x.Cache == "" &&
		x.Provider == "" &&
		x.MaxGoroutines == nil &&
		len(x.Resources) == 0 &&
		x.ProvidersConfigYamlString == ""
}

// ------------------------------------------------- --------------------------------------------------------------------
