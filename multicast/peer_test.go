package multicast

import (
	"fmt"
	"github.com/talostrading/sonic"
	"github.com/talostrading/sonic/net/ipv4"
	"net"
	"testing"
)

// Listing multicast group memberships: netstat -gsv

func TestUDPPeer_IPv4Addresses(t *testing.T) {
	ioc := sonic.MustIO()
	defer ioc.Close()

	{
		_, err := NewUDPPeer(ioc, "udp", net.IPv4zero.String())
		if err == nil {
			t.Fatal("should have received an error as the address is missing the port")
		}
	}
	{
		_, err := NewUDPPeer(ioc, "udp4", net.IPv4zero.String())
		if err == nil {
			t.Fatal("should have received an error as the address is missing the port")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp", "")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), net.IPv4zero.String(); given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp4", "")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), net.IPv4zero.String(); given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp", ":0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), net.IPv4zero.String(); given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp4", ":0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), net.IPv4zero.String(); given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), "127.0.0.1"; given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp4", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), "127.0.0.1"; given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp", "localhost:0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), "127.0.0.1"; given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
	{
		peer, err := NewUDPPeer(ioc, "udp4", "localhost:0")
		if err != nil {
			t.Fatal(err)
		}
		defer peer.Close()

		addr := peer.LocalAddr()
		if given, expected := addr.IP.String(), "127.0.0.1"; given != expected {
			t.Fatalf("given=%s expected=%s", given, expected)
		}
		if addr.Port == 0 {
			t.Fatal("port should not be 0")
		}

		if iface, _ := peer.Outbound(); iface != nil {
			t.Fatal("not explicit outbound interface should have been set")
		}
	}
}

func TestUDPPeer_JoinInvalidGroup(t *testing.T) {
	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	if err := peer.Join("0.0.0.0:4555"); err == nil {
		t.Fatal("should not have joined")
	}
}

func TestUDPPeer_Join(t *testing.T) {
	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	if err := peer.Join("224.0.0.0"); err != nil {
		t.Fatal(err)
	}

	addr, err := ipv4.GetMulticastInterface(peer.socket)
	if err != nil {
		t.Fatal(err)
	}
	if !addr.IsUnspecified() {
		t.Fatal("multicast address should be unspecified")
	}
}

func TestUDPPeer_SetLoop1(t *testing.T) {
	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	if peer.Loop() {
		t.Fatal("peer should not loop packets by default")
	}

	if err := peer.SetLoop(false); err != nil {
		t.Fatal(err)
	}
	if peer.Loop() {
		t.Fatal("peer should not loop packets")
	}

	if err := peer.SetLoop(true); err != nil {
		t.Fatal(err)
	}
	if !peer.Loop() {
		t.Fatal("peer should loop packets")
	}

	if err := peer.SetLoop(false); err != nil {
		t.Fatal(err)
	}
	if peer.Loop() {
		t.Fatal("peer should not loop packets")
	}
}

func TestUDPPeer_SetOutboundInterfaceOnUnspecifiedIPAndPort(t *testing.T) {
	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	for _, iff := range testInterfaces {
		fmt.Println("setting", iff.Name, "as outbound")

		if err := peer.SetOutboundIPv4(iff.Name); err != nil {
			t.Fatal(err)
		}

		outboundInterface, outboundIP := peer.Outbound()
		fmt.Println("outbound for", iff.Name, outboundInterface, outboundIP)

		{
			addr, err := ipv4.GetMulticastInterface(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%s GetMulticastInterface_Inet4Addr addr=%s\n", iff.Name, addr.String())
		}

		{
			interfaceAddr, multicastAddr, err := ipv4.GetMulticastInterface2(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf(
				"%s GetMulticastInterface_IPMreq4Addr interface_addr=%s multicast_addr=%s\n",
				iff.Name, interfaceAddr.String(), multicastAddr.String())
		}

		{
			interfaceIndex, err := ipv4.GetMulticastInterfaceIndex(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%s GetMulticastInterface_Index interface_index=%d\n", iff.Name, interfaceIndex)
		}
	}
}

func TestUDPPeer_SetOutboundInterfaceOnUnspecifiedPort(t *testing.T) {
	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()

	for _, iff := range testInterfaces {
		fmt.Println("setting", iff.Name, "as outbound")

		if err := peer.SetOutboundIPv4(iff.Name); err != nil {
			t.Fatal(err)
		}

		outboundInterface, outboundIP := peer.Outbound()
		fmt.Println("outbound for", iff.Name, outboundInterface, outboundIP)

		{
			addr, err := ipv4.GetMulticastInterface(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%s GetMulticastInterface_Inet4Addr addr=%s\n", iff.Name, addr.String())
		}

		{
			interfaceAddr, multicastAddr, err := ipv4.GetMulticastInterface2(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf(
				"%s GetMulticastInterface_IPMreq4Addr interface_addr=%s multicast_addr=%s\n",
				iff.Name, interfaceAddr.String(), multicastAddr.String())
		}

		{
			interfaceIndex, err := ipv4.GetMulticastInterfaceIndex(peer.NextLayer())
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%s GetMulticastInterface_Index interface_index=%d\n", iff.Name, interfaceIndex)
		}
	}
}

func TestUDPPeer_TTL(t *testing.T) {
	ioc := sonic.MustIO()
	defer ioc.Close()

	peer, err := NewUDPPeer(ioc, "udp", "")
	if err != nil {
		t.Fatal(err)
	}

	actualTTL, err := ipv4.GetMulticastTTL(peer.NextLayer())
	if err != nil {
		t.Fatal(err)
	}

	if actualTTL != peer.TTL() {
		t.Fatalf("wrong TTL expected=%d given=%d", actualTTL, peer.TTL())
	}

	if peer.TTL() != 1 {
		t.Fatalf("peer TTL should be 1 by default")
	}

	if err := peer.SetTTL(32); err != nil {
		t.Fatal(err)
	}

	if peer.TTL() != 32 {
		t.Fatalf("peer TTL should be 32")
	}
}
