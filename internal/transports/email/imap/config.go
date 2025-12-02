package imap

// Config for IMAP/SMTP transport.
type Config struct {
	ID       string `yaml:"id" json:"id"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Folder   string `yaml:"folder" json:"folder"`
	Idle     bool   `yaml:"idle" json:"idle"`

	SMTPHost string `yaml:"smtp_host" json:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port" json:"smtp_port"`
	SMTPTLS  bool   `yaml:"smtp_tls" json:"smtp_tls"`

	AllowSenders []string `yaml:"allow_senders" json:"allow_senders"`
	MaxBytes     int      `yaml:"max_bytes" json:"max_bytes"`
}

func (c *Config) Defaults() {
	if c.Port == 0 {
		c.Port = 993
	}
	if c.Folder == "" {
		c.Folder = "INBOX"
	}
	if c.SMTPPort == 0 {
		c.SMTPPort = 587
	}
	if c.ID == "" {
		c.ID = "email-imap"
	}
	if c.SMTPHost == "" {
		c.SMTPHost = c.Host
	}
	if c.MaxBytes == 0 {
		c.MaxBytes = 262144 // 256 KiB
	}
}

func (c *Config) Validate() error {
	if c.Host == "" || c.Username == "" || c.Password == "" {
		return Err("host, username, password are required")
	}
	if len(c.AllowSenders) == 0 {
		c.AllowSenders = []string{c.Username}
	}
	return nil
}
