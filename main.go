package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
)

// HTML template for the web interface
var tpl = `
<!DOCTYPE html>
<html>
<head>
    <title>Wake-on-LAN</title>
</head>
<body>
    <h1>Wake-on-LAN</h1>
    <form method="POST" action="/send-wol">
        <button type="submit">Send Magic Packet</button>
    </form>
</body>
</html>
`

var macAddress string

// SendMagicPacket sends a Wake-on-LAN magic packet to the specified MAC address.
func SendMagicPacket(macAddr string) error {
	hwAddr, err := net.ParseMAC(macAddr)
	if err != nil {
		return err
	}

	if len(hwAddr) != 6 {
		return errors.New("invalid MAC address")
	}

	// Create the magic packet.
	packet := make([]byte, 102)
	// Add 6 bytes of 0xFF.
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	// Repeat the MAC address 16 times.
	for i := 1; i <= 16; i++ {
		copy(packet[i*6:], hwAddr)
	}

	// Send the packet using UDP.
	addr := net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 9,
	}
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	return err
}

// homeHandler serves the home page with the button to send the WoL packet.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index").Parse(tpl)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

// sendWolHandler handles the POST request to send the WoL packet.
func sendWolHandler(w http.ResponseWriter, r *http.Request) {
	err := SendMagicPacket(macAddress)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send WoL packet: %v", err), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9080"
	}

	macAddress = os.Getenv("MAC")
	if macAddress == "" {
		log.Fatal("MAC address environment variable MAC is required but not set")
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/send-wol", sendWolHandler)

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
