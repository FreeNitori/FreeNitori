package multiplexer

import (
	"encoding/json"
	"log"
	"net"
)

var ExitCode = make(chan int)

type IPCPacket struct {
	IssuerIdentifier   string   `json:"issuer_identifier"`
	ReceiverIdentifier string   `json:"receiver_identifier"`
	MessageIdentifier  string   `json:"message_identifier"`
	Response           bool     `json:"response"`
	Body               []string `json:"body"`
}

func WritePacket(connection net.Conn, outgoingPacket IPCPacket) (incomingPacket IPCPacket) {
	var err error
	var incoming IPCPacket

	jsonEncoder := json.NewEncoder(connection)
	jsonDecoder := json.NewDecoder(connection)

	// Encode the outgoing packet
	err = jsonEncoder.Encode(outgoingPacket)
	if err != nil {
		log.Printf("Failed to encode packet, %s", err)
		return incoming
	}

	// Decode and return the incoming packet
	err = jsonDecoder.Decode(&incoming)
	if err != nil {
		log.Printf("Failed to decode packet, %s", err)
		return incoming
	}
	return incoming
}

func (incomingPacket IPCPacket) SupervisorPacketHandler() (outgoingPacket IPCPacket) {
	switch incomingPacket.IssuerIdentifier {
	case "ChatBackendInitializer":

		// This should never talk to anything other than the supervisor
		if incomingPacket.ReceiverIdentifier != "Supervisor" ||
			incomingPacket.MessageIdentifier != "ChatBackendInitializationFinish" {
			log.Println("Invalid packet from Chat backend initializer.")
			ExitCode <- 1
		}

		// Print out the message otherwise
		log.Printf("User: %s | ID: %s | Prefix: %s",
			incomingPacket.Body[0],
			incomingPacket.Body[1],
			incomingPacket.Body[2])
		log.Printf("FreeNitori is now ready. Press Control-C to terminate.")
	default:
		return IPCPacket{
			IssuerIdentifier:   "Supervisor",
			ReceiverIdentifier: incomingPacket.IssuerIdentifier,
			MessageIdentifier:  "Error",
			Response:           true,
			Body:               []string{"Unknown issuer."},
		}
	}
	return IPCPacket{
		IssuerIdentifier:   "Supervisor",
		ReceiverIdentifier: incomingPacket.ReceiverIdentifier,
		MessageIdentifier:  "MessageAcknowledgement",
		Response:           true,
		Body:               []string{"Request has been received."},
	}
}
