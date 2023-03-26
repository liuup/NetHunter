package main

import (
	"example/router"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/google/gopacket/pcap"
)

var (
	// 本地设备号
	device string = "\\Device\\NPF_{7C8005B1-5678-4F83-8E55-D929A065FD68}"
	// device       string = "\\Device\\NPF_{BF415878-83C9-41A1-B07E-7963D76530AA}"
	snapshot_len int32 = 65535
	promiscuous  bool  = true
	// err          error
	timeout time.Duration = -1 * time.Second
	// handle       *pcap.Handle
)

func main() {
	r := router.GetRouter()

	r.Run(":9876")

	// find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Println(err)
	}
	for _, d := range devices {
		fmt.Println(d.Name)
		fmt.Println(d.Description)
		fmt.Println("----- ----- ----- -----")
	}

	// 打开设备
	handle, err := pcap.OpenLive("\\Device\\NPF_{6859611A-A240-4DE0-8FC1-B7645ABB5A84}", snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Println(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	go func(*gopacket.PacketSource) {
		for p := range packetSource.Packets() {
			// fmt.Println(p.LinkLayer())
			fmt.Println("----- ----- ----- -----")
			// fmt.Println(p)
			// // fmt.Println(p.NetworkLayer().LayerContents())
			// // fmt.Println(p.NetworkLayer().LayerPayload())
			// fmt.Println(p.TransportLayer().LayerContents())
			// fmt.Println(p.TransportLayer().LayerPayload())

			GetPacketInfo(p)
			// fmt.Println("----- ----- ----- -----")
		}
	}(packetSource)

	select {}

}

func GetPacketInfo(p gopacket.Packet) {
	// 提取dns信息
	/*
		if dns := p.Layer(layers.LayerTypeDNS); dns != nil {
			fmt.Println("dns detected")
			dnsPacket, _ := dns.(*layers.DNS)
			// fmt.Println(dnsPacket.Questions)
			// fmt.Println(dnsPacket.Answers)
			// 获取域名
			for _, q := range dnsPacket.Questions {

				fmt.Println("question.Name:", string(q.Name))
			}
			// 获取域名和ip
			for _, a := range dnsPacket.Answers {
				fmt.Println("answer.Name:", string(a.Name))
				fmt.Println("answer.IP", a.IP)
			}

			// fmt.Println(dnsPacket)
			// fmt.Println(dnsPacket.CanDecode())
			// fmt.Println(dnsPacket.LayerType())
			// fmt.Println(dnsLayer)
			fmt.Println("packets:" + p.String())
		}
	*/

	// 提取icmp协议信息
	if icmp := p.Layer(layers.LayerTypeICMPv4); icmp != nil {
		fmt.Println("icmp detected")
		icmpPacket, _ := icmp.(*layers.ICMPv4)

		fmt.Println(icmpPacket)
	}
	if icmp := p.Layer(layers.LayerTypeICMPv6); icmp != nil {
		fmt.Println("icmp detected")
		icmpPacket, _ := icmp.(*layers.ICMPv6)

		fmt.Println(icmpPacket)
	}

	fmt.Println(p)
}
