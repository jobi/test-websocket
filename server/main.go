package main

// #include <dns_sd.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"code.google.com/p/go.net/websocket"
)

func main() {
	counterChannel := make(chan int)

	go func() {
		counter := 0

		for {
			<-time.After(1 * time.Second)
			fmt.Println("Counter at", counter)
			counterChannel <-counter
			counter++
		}
	} ()

	receivers := make([]chan int, 0)
	lock := new(sync.Mutex)

	go func() {
		for {
			counter := <-counterChannel

			lock.Lock()
			for _, receiver := range(receivers) {
				fmt.Println("Sending to receiver")
				receiver <- counter
			}
			lock.Unlock()
		}
	} ()

	websocketHandler := func(ws *websocket.Conn) {
		fmt.Println("Received websocket connection")

		receiver := make(chan int)
		lock.Lock()
		receivers = append(receivers, receiver)
		lock.Unlock()

		for {
			value := <-receiver
			websocket.JSON.Send(ws, value)
		}
	}

	http.Handle("/ws", websocket.Handler(websocketHandler))

	var listener net.Listener
	var err error

	if listener, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
		fmt.Println("Listen failed", err)
		return
	}

	var tcpAddr *net.TCPAddr
	if tcpAddr, err = net.ResolveTCPAddr("tcp", listener.Addr().String()); err != nil {
		fmt.Println("Couldn't resolve TCP address", err)
		return
	}

	fmt.Println("Listening on", fmt.Sprintf("http://localhost:%d", tcpAddr.Port))

	// network order
	port := (tcpAddr.Port&0xff)<<8 | (tcpAddr.Port&0xff00)>>8

	// Interface -1 means localhost
	var iface uint32 = 0
	iface -= 1

	serviceName := C.CString("woven")
	serviceType := C.CString("_woven._tcp")
	serviceDomain := C.CString("local.")

	var dnsRef C.DNSServiceRef
	if dnsErr := C.DNSServiceRegister(&dnsRef, 0,
		C.uint32_t(iface),
		serviceName,
		serviceType,
		serviceDomain, nil,
		C.uint16_t(port), 0, nil,
		nil, nil); dnsErr != 0 {
		fmt.Println("DNSServiceRegister failed", dnsErr)
		return
	}

	cleanup := func() {
		if dnsRef != nil {
			C.DNSServiceRefDeallocate(dnsRef)
			dnsRef = nil

			C.free(unsafe.Pointer(serviceDomain))
			C.free(unsafe.Pointer(serviceType))
			C.free(unsafe.Pointer(serviceName))

			dnsRef = nil
		}
	}

	defer cleanup()

	signalChannel:= make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGTERM)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		cleanup()
		os.Exit(1)
	} ()

	if err = http.Serve(listener, nil); err != nil {
		fmt.Println("Serve failed", err)
		return
	}
}
