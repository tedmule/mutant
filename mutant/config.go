package mutant

type MutantConfig struct {
	// Context    string `env:"MUTANT_CONTEXT" envDefault:""`

	// Running in k8s cluser or out of k8s cluster, inclusert|outofcluster
	// K8S client authentication mode: <kubeconfig|sa|token>
	// 	- kubeconfig: use ~/.kube/config
	// 	- sa		: run in kubernetes, use service account
	// 	- token		: use token
	Mode       string `env:"MUTANT_MODE" envDefault:"outofcluster"`
	Listen     string `env:"MUTANT_LISTEN" envDefault:":8443"`
	Annotation string `env:"MUTANT_ANNOTATION" envDefault:"mutant"`
	// Length        int    `env:"MUTANT_JABBER_WORD" envDefault:"2"`
	// Name          string `env:"MUTANT_NAME" envDefault:"default"`
	// Count         int    `env:"MUTANT_COUNT" envDefault:"1"`
	LogLevel   string `env:"MUTANT_LOG_LEVEL" envDefault:"debug"`
	K8SAPI     string `env:"MUTANT_K8A_API" envDefault:"https://kubernetes.default.svc"`
	Production bool   `env:"MUTANT_PROD" envDefault:"false"`
	K8SToken   string `env:"MUTANT_K8S_TOKEN" envDefault:"eyJhbGciOiJSUzI1NiIsImtpZCI6InRBb1JyNzRaa3VYZmV6cmk4bHZybGJZcjVpOGN4cDhCSEtCdEJQMnp1RWMifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZXYiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoiZGV2LWNvbnRhaW5lci1zZWNyZXQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGV2LWNvbnRhaW5lci1zYSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjExMWY0OGRmLTFkYWEtNDljOS1hMzIzLTI0Nzc3ZWE0Y2U0ZCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZXY6ZGV2LWNvbnRhaW5lci1zYSJ9.o0ZKu_ziOO3-GJ_kzYDnNq3UslhjRkue0TJWFAC9wgAgndhQi37r6-HwtMx3syHnC8Q5sNdG_Df0vYAKSH5PjgA2RqbMIoOWUwRxEDIwNBHHZ9xJrOu4gCZoxWqHgBskmjsqE5zVw5D6ksltAEZKFke15t2NlYuiiaz1Mj9mcEdUk7ryo5Z18VGKe6lsdbqfu_6GkUvN5NvzvoZcSrnc6VTGxuBV_c1Mfhk0lJpIlzEZjjDCpi6w-V3aH1oIJE5xmBxSOo9i8GRCV1SmEMsOErF9Qsc2QRwIiuIe4R4ALS-xSxqbrDBEAnI95feZDlsJU8yrqMsm0zxpkpHWSHQ13Q"`
	//
	ReadTimtout    int    `env:"MUTANT_READ_TIMEOUT" envDefault:"10"`
	WriteTimtout   int    `env:"MUTANT_WRITE_TIMEOUT" envDefault:"10"`
	MaxHeaderBytes int    `env:"MUTANT_MAX_HEADER_BYTES" envDefault:"0"`
	CertFile       string `env:"MUTANT_CERT_FILE" envDefault:"./certs/server.crt"`
	KeyFile        string `env:"MUTANT_KEY_FILE" envDefault:"./certs/server.key"`
}
