package generator

type DigitaloceanConfig struct {
    
    AccessToken string `json:"accessToken,omitempty"`
    
    Backups bool `json:"backups,omitempty"`
    
    Image string `json:"image,omitempty"`
    
    Ipv6 bool `json:"ipv6,omitempty"`
    
    PrivateNetworking bool `json:"privateNetworking,omitempty"`
    
    Region string `json:"region,omitempty"`
    
    Size string `json:"size,omitempty"`
    
    SshKeyFingerprint string `json:"sshKeyFingerprint,omitempty"`
    
    SshKeyPath string `json:"sshKeyPath,omitempty"`
    
    SshPort string `json:"sshPort,omitempty"`
    
    SshUser string `json:"sshUser,omitempty"`
    
    Tags string `json:"tags,omitempty"`
    
    Userdata string `json:"userdata,omitempty"`
    
}
