package config

type Pave struct {
	// ProjectName represents the project that tinker was initialized under.
	// Ex. If tinker was initialized in a directory by the name of
	// `ProjectFoo` the project name will reflect this name.
	// ProjectName will be used as a namespace for downstream configurations.
	ProjectName string `yaml:"project_name"`

	// ProjectLang is based on a list of supported languages. (Ex. go, ruby, node, etc...)
	// Nice to have would be the ability to infer project language from language managers.
	// TODO: Detect language and version from lang managers auto-load files .<lang>-version
	ProjectLang string `yaml:"project_lang"` // Format LANG:VERSION

	// DockerSupport dictates whether tinker will manage the Dockerfile.
	// Pave will generate a Dockerfile policy which will enforce supported base images, artifact, and
	// container generation.
	DockerEnabled bool `yaml:"docker_enabled"`

	// TerraformEnabled dictates whether tinker will manage infra as code via terraform.
	// The idea is that tinker will assist the developer, via scaffolding, by providing the
	// operator a menu of aws resources via supported modules. Since tinker will mainly support containers
	// the default manifest will be that of the ECS resources.
	TerraformEnabled bool `yaml:"terraform_enabled"`
}
