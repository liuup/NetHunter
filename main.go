package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	// "github.com/google/gopacket/pcap"
)

var (
	device string = "\\Device\\NPF_{7C8005B1-5678-4F83-8E55-D929A065FD68}"
	// device       string = "\\Device\\NPF_{BF415878-83C9-41A1-B07E-7963D76530AA}"
	snapshot_len int32 = 65535
	promiscuous  bool  = true
	err          error
	timeout      time.Duration = -1 * time.Second
	handle       pcap.Handle
)

func main() {
	fmt.Print()

	// find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Println(err)
	}
	for _, d := range devices {
		fmt.Println(d.Name)
		fmt.Println(d.Description)
		fmt.Println()
	}

	// 打开设备
	handle, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Println(err)
	}
	defer handle.Close()

	// filter := "tcp"
	// err = handle.SetBPFFilter(filter)
	// if err != nil {
	// 	log.Println(err)
	// }

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go func(*gopacket.PacketSource) {
		for p := range packetSource.Packets() {
			// fmt.Println(p.LinkLayer())
			fmt.Println("----- ----- ----- -----")
			fmt.Println(p)
			fmt.Println(p.NetworkLayer().LayerContents())
			fmt.Println(p.NetworkLayer().LayerPayload())
			fmt.Println("----- ----- ----- -----")
		}
	}(packetSource)

	select {}
}
