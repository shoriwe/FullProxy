package ArgumentParser

import (
	"flag"
	"os"
)

func ParseSocks5Arguments() (*string, *string, []byte, []byte, *bool) {
	protocolFlagSet := flag.NewFlagSet("socks5", flag.ExitOnError)
	host := protocolFlagSet.String("host", "", "Host to listen on. When \"-slave\" flag is set, is the IP of master to connect")
	port := protocolFlagSet.String("port", "1080", "Port to listen on. When \"-slave\" flag is set, is the Port of the master to connect. I both modes the default port is 1080")
	username := protocolFlagSet.String("username", "", "Username of the running proxy, requires \"-password\". It will be ignored if is an empty string")
	password := protocolFlagSet.String("password", "", "Password of the running proxy, requires \"-username\". It will be ignored if is an empty string")
	slave := protocolFlagSet.Bool("slave", false, "Connect to a master, no bind proxying")
	_ = protocolFlagSet.Parse(os.Args[2:])
	if len(*host) == 0 {
		if *slave {
			*host = "127.0.0.1"
		} else {
			*host = "0.0.0.0"
		}
	}
	return host, port, []byte(*username), []byte(*password), slave
}