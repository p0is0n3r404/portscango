package ports

// TopPorts list of most commonly used ports
var TopPorts = []int{
	21,    // FTP
	22,    // SSH
	23,    // Telnet
	25,    // SMTP
	53,    // DNS
	80,    // HTTP
	110,   // POP3
	111,   // RPC
	135,   // MSRPC
	139,   // NetBIOS
	143,   // IMAP
	443,   // HTTPS
	445,   // SMB
	993,   // IMAPS
	995,   // POP3S
	1723,  // PPTP
	3306,  // MySQL
	3389,  // RDP
	5432,  // PostgreSQL
	5900,  // VNC
	6379,  // Redis
	8080,  // HTTP Proxy
	8443,  // HTTPS Alt
	27017, // MongoDB
}

// Top100Ports top 100 most common ports
var Top100Ports = []int{
	7, 9, 13, 21, 22, 23, 25, 26, 37, 53, 79, 80, 81, 82, 83, 84, 85, 88, 89, 90,
	99, 100, 106, 110, 111, 113, 119, 125, 135, 139, 143, 144, 146, 161, 163, 179,
	199, 211, 212, 222, 254, 255, 256, 259, 264, 280, 301, 306, 311, 340, 366, 389,
	406, 407, 416, 417, 425, 427, 443, 444, 445, 458, 464, 465, 481, 497, 500, 512,
	513, 514, 515, 524, 541, 543, 544, 545, 548, 554, 555, 563, 587, 593, 616, 617,
	625, 631, 636, 646, 648, 666, 667, 668, 683, 687, 691, 700, 705, 711, 714, 720,
}

// ServiceNames returns service names by port number
var ServiceNames = map[int]string{
	7:     "Echo",
	20:    "FTP-Data",
	21:    "FTP",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	67:    "DHCP",
	68:    "DHCP",
	69:    "TFTP",
	80:    "HTTP",
	110:   "POP3",
	111:   "RPC",
	119:   "NNTP",
	123:   "NTP",
	135:   "MSRPC",
	137:   "NetBIOS-NS",
	138:   "NetBIOS-DGM",
	139:   "NetBIOS-SSN",
	143:   "IMAP",
	161:   "SNMP",
	162:   "SNMP-Trap",
	179:   "BGP",
	194:   "IRC",
	389:   "LDAP",
	443:   "HTTPS",
	445:   "SMB",
	464:   "Kerberos",
	465:   "SMTPS",
	514:   "Syslog",
	515:   "LPD",
	520:   "RIP",
	521:   "RIPng",
	543:   "Klogin",
	544:   "Kshell",
	548:   "AFP",
	554:   "RTSP",
	563:   "NNTPS",
	587:   "Submission",
	593:   "HTTP-RPC",
	631:   "IPP",
	636:   "LDAPS",
	646:   "LDP",
	873:   "Rsync",
	902:   "VMware",
	989:   "FTPS-Data",
	990:   "FTPS",
	993:   "IMAPS",
	995:   "POP3S",
	1080:  "SOCKS",
	1194:  "OpenVPN",
	1433:  "MSSQL",
	1434:  "MSSQL-UDP",
	1521:  "Oracle",
	1723:  "PPTP",
	1883:  "MQTT",
	2049:  "NFS",
	2082:  "cPanel",
	2083:  "cPanel-SSL",
	2181:  "ZooKeeper",
	2222:  "SSH-Alt",
	2375:  "Docker",
	2376:  "Docker-SSL",
	3000:  "Node.js",
	3306:  "MySQL",
	3389:  "RDP",
	3690:  "SVN",
	4000:  "ICQ",
	4443:  "HTTPS-Alt",
	4444:  "Metasploit",
	5000:  "Flask",
	5432:  "PostgreSQL",
	5672:  "RabbitMQ",
	5900:  "VNC",
	5984:  "CouchDB",
	6379:  "Redis",
	6443:  "Kubernetes",
	6666:  "IRC-Alt",
	6667:  "IRC",
	7001:  "WebLogic",
	8000:  "HTTP-Alt",
	8008:  "HTTP-Alt",
	8080:  "HTTP-Proxy",
	8081:  "HTTP-Alt",
	8443:  "HTTPS-Alt",
	8888:  "HTTP-Alt",
	9000:  "PHP-FPM",
	9090:  "WebSM",
	9200:  "Elasticsearch",
	9300:  "Elasticsearch",
	9418:  "Git",
	9999:  "Urchin",
	10000: "Webmin",
	11211: "Memcached",
	27017: "MongoDB",
	27018: "MongoDB",
	28017: "MongoDB-Web",
	50000: "SAP",
}

// GetServiceName returns service name for given port
func GetServiceName(port int) string {
	if name, ok := ServiceNames[port]; ok {
		return name
	}
	return "Unknown"
}
