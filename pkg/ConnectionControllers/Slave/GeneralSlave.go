package Slave

import (
	"bufio"
	"crypto/tls"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"net"
	"time"
)

type General struct {
	MasterConnection       net.Conn
	MasterConnectionReader *bufio.Reader
	MasterConnectionWriter *bufio.Writer
	MasterHost             string
	MasterPort             string
	TLSConfiguration       *tls.Config
	ProxyProtocol          ConnectionControllers.ProxyProtocol
}

func (general *General) Serve() error {
	var finalError error
	for {
		timeoutSetError := general.MasterConnection.SetReadDeadline(time.Now().Add(20 * time.Second))
		if timeoutSetError != nil {
			log.Fatal(timeoutSetError)
		}
		NumberOfReceivedBytes, buffer, connectionError := Sockets.Receive(general.MasterConnectionReader, 1024)
		if connectionError != nil {
			if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				finalError = connectionError
				break
			}
		}
		if NumberOfReceivedBytes != 1 {
			continue
		}
		if buffer[0] != ConnectionControllers.NewConnection[0] {
			continue
		}
		clientConnection, connectionError := Sockets.TLSConnect(&general.MasterHost, &general.MasterPort, general.TLSConfiguration)

		if connectionError != nil {
			finalError = connectionError
			break
		}
		clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
		go general.ProxyProtocol.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
	return finalError
}
