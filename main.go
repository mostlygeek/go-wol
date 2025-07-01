package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/mostlygeek/go-wol/internal"
)

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
	t, err := template.New("index").Parse(internal.TEMPLATE)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

// sendWolHandler handles the POST request to send the WoL packet.
func sendWolHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Address string `json:"address"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	if request.Address == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "MAC address is required"})
		return
	}

	err = SendMagicPacket(request.Address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Failed to send WoL packet: %v", err)})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Magic packet sent successfully!"})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9080"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/wakeup", sendWolHandler)

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
