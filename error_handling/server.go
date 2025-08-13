package err 
type Server struct{
	Debug Debug
}

type Debug struct{
	Running string `yaml:"running"`
}