package manifest

type Remote struct {
	Name string `yaml:"name"`
	Alias string `yaml:"alias,omitempty"`
	Fetch string `yaml:"fetch"`
	PushUrl string `yaml:"push-url,omitempty"`
	Review string `yaml:"review,omitempty"`
	Revision string `yaml:"revision,omitempty"`
}

type Default struct {
	Remote string `yaml:"remote,omitempty"`
	Revision string `yaml:"revision,omitempty"`
	DestBranch string `yaml:"dest-branch,omitempty"`
	Upstream string `yaml:"upstream,omitempty"`
	SyncJ bool `yaml:"sync-j,omitempty"`
	SyncC bool `yaml:"sync-c,omitempty"`
	SyncS bool `yaml:"sync-s,omitempty"`
	SyncTags bool `yaml:"sync-tags,omitempty"`
}

type Server struct {
	Url string `yaml:"url"`
}

type Project struct {
	Name string `yaml:"name"`
	Path string `yaml:"path,omitempty"`
	Remote string `yaml:"remote,omitempty"`
	Revision string `yaml:"revision,omitempty"`
	DestBranch string `yaml:"dest-branch,omitempty"`
	Groups []string `yaml:"groups,omitempty"`
	SyncC bool `yaml:"sync-c,omitempty"`
	SyncS bool `yaml:"sync-s,omitempty"`
	SyncTags bool `yaml:"sync-tags,omitempty"`
	Upstream string `yaml:"upstream,omitempty"`
	CloneDepth int `yaml:"clone-depth,omitempty"`
	ForcePath bool `yaml:"force-path,omitempty"`
}

type Manifest struct {
	Manifest struct{
		Notice string `yaml:"notice,omitempty"`
		Remotes []Remote `yaml:"remotes,omitempty"`
		Default Default `yaml:"default,omitempty"`
		Server Server `yaml:"manifest-server,omitempty"`
		Projects []Project `yaml:"projects,omitempty"`
	} `yaml:"manifest"`
}

func (m Manifest) FindRemote(remoteName string) *Remote {
	for _, remote := range m.Manifest.Remotes {
		if remote.Name == remoteName {
			return &remote
		}
	}

	return nil
}


