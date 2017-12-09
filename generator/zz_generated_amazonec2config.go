package generator

type Amazonec2Config struct {
    
    AccessKey string `json:"accessKey,omitempty"`
    
    Ami string `json:"ami,omitempty"`
    
    BlockDurationMinutes string `json:"blockDurationMinutes,omitempty"`
    
    DeviceName string `json:"deviceName,omitempty"`
    
    Endpoint string `json:"endpoint,omitempty"`
    
    IamInstanceProfile string `json:"iamInstanceProfile,omitempty"`
    
    InsecureTransport bool `json:"insecureTransport,omitempty"`
    
    InstanceType string `json:"instanceType,omitempty"`
    
    KeypairName string `json:"keypairName,omitempty"`
    
    Monitoring bool `json:"monitoring,omitempty"`
    
    OpenPort []string `json:"openPort,omitempty"`
    
    PrivateAddressOnly bool `json:"privateAddressOnly,omitempty"`
    
    Region string `json:"region,omitempty"`
    
    RequestSpotInstance bool `json:"requestSpotInstance,omitempty"`
    
    Retries string `json:"retries,omitempty"`
    
    RootSize string `json:"rootSize,omitempty"`
    
    SecretKey string `json:"secretKey,omitempty"`
    
    SecurityGroup []string `json:"securityGroup,omitempty"`
    
    SessionToken string `json:"sessionToken,omitempty"`
    
    SpotPrice string `json:"spotPrice,omitempty"`
    
    SshKeypath string `json:"sshKeypath,omitempty"`
    
    SshUser string `json:"sshUser,omitempty"`
    
    SubnetId string `json:"subnetId,omitempty"`
    
    Tags string `json:"tags,omitempty"`
    
    UseEbsOptimizedInstance bool `json:"useEbsOptimizedInstance,omitempty"`
    
    UsePrivateAddress bool `json:"usePrivateAddress,omitempty"`
    
    Userdata string `json:"userdata,omitempty"`
    
    VolumeType string `json:"volumeType,omitempty"`
    
    VpcId string `json:"vpcId,omitempty"`
    
    Zone string `json:"zone,omitempty"`
    
}
