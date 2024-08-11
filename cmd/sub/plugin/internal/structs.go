package internal

//// JenkinsDownload type = "jenkins".
//type JenkinsDownload struct {
//	Job      string `json:"job,omitempty"`
//	Artifact string `json:"artifact,omitempty"`
//}
//
//// URLDownload type = "url"
//type URLDownload struct {
//	URL string `json:"url,omitempty"`
//}
//
//type Download struct {
//	Type     string `json:"type"`
//	Filename string `json:"filename"`
//
//	JenkinsDownload
//	URLDownload
//}
//
//type Commands struct {
//	Windows []string `json:"windows"`
//	Unix    []string `json:"unix"`
//}
//
//type File struct {
//	Type string `json:"type"`
//	Path string `json:"path"`
//}
//
//type Install struct {
//	Type     string    `json:"type"`
//	Commands *Commands `json:"commands,omitempty"`
//}
//
//type Uninstall struct {
//	Files []File `json:"files"`
//}
//
//// DefinedLater Defined in pap, not in the json files themselves.
//type DefinedLater struct {
//	Path                 string `json:"path,omitempty"`
//	URL                  string `json:"url,omitempty"`
//	Source               string `json:"source,omitempty"`
//	IsDependency         bool   `json:"isDependency,omitempty"`
//	IsOptionalDependency bool   `json:"isOptionalDependency,omitempty"`
//}
//
//type AllDependencies struct {
//	Dependencies         []string `json:"dependencies,omitempty"`
//	OptionalDependencies []string `json:"optionalDependencies,omitempty"`
//}

type PluginInfo struct {
	URL  string
	Name string
	//Version string `json:"version"`
	//AllDependencies
	//
	//Downloads []Download `json:"downloads"`
	//DefinedLater
	//
	//Alias string `json:"alias,omitempty"`
}
