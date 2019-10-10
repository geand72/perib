package main

import (
	// "encoding/binary"

	"fmt"
	"log"

	// "os"
	"time"

	"github.com/geand72/go-i2c"
)

// Write to M24M-01 EEPROM one page, page size 0xff
func WritePage(m24mAddr uint8, page byte, dataEEPage []byte) error {

	m24w, err := i2c.NewI2C(m24mAddr, 1)

	if err != nil {
		log.Fatalf("Can't open i2c bus for write: %v", err)
	}

	defer m24w.Close()

	_, err = m24w.WriteBytesEE(page, 0x00, dataEEPage)

	if err != nil {
		log.Fatalf("Can't write to eeprom: %v", err)
	}

	time.Sleep(time.Millisecond * 5)
	return err
}

// Read from M24M-01 EEPROM one page, page size 0xff
func ReadPage(m24mAddr uint8, page byte) (error, []byte) {

	m24r, err := i2c.NewI2C(m24mAddr, 1)
	if err != nil {
		log.Fatalf("Can't open i2c bus for write: %v", err)
	}

	defer m24r.Close()

	dataEEPage, _, err := m24r.ReadRegBytesEE(page, 0x00, 0x100)

	if err != nil {
		log.Fatalf("Can't read from eeprom: %v", err)
	}
	return err, dataEEPage
}

// Check slice byte with value
func isFillValue(b []byte, value byte) bool {
	if len(b) == 0 {
		return false
	}
	for _, v := range b {
		if v != value {
			return false
		}
	}
	return true
}

// Fill slice byte with value
func fillWithValue(b []byte, value byte) {
	for i := range b {
		b[i] = value
	}
}

// Test EEPROM M24M-01, 1 Mbit, 512 page, 256 byte per page.
// Address 0-255 page - 0x50, address 256-511 page - 0x51.
func testM24M() {
	var (
		ma    uint8   = 0x50
		pg, j byte    = 0x0, 0x00
		pr    float32 = 0x00
		dp    []byte  = make([]byte, 0x100)
	)

	for z := 0; z <= 1; z++ {
		fillWithValue(dp, j)
		// Write 512 page
		ma = 0x50
		pg = 0x00
		for ma <= 0x51 {
			for i := 0; i <= 0xff; i++ {
				if err := WritePage(ma, pg, dp); err != nil {
					log.Fatal("can,t write page", err)
				}
				err, dp := ReadPage(ma, pg)

				if err != nil {
					log.Fatalf("Can't read from eeprom: %v", err)
				}

				if !isFillValue(dp, j) {
					log.Fatalf("EEPROM error: SubPage %x,Page %x, Cell %x", ma-0x50, pg, i)
				}
				fmt.Printf("Write-Read value '%2x' to SubPage %1d Page %3d (%.0f%%)\r", j, ma-0x50, pg, pr*0.098)
				pg++
				pr++
			}
			fmt.Println()
			pg = 0x00
			ma++
		}
		j = +0xff
	}
	fmt.Printf("Test EEPROM M24M OK!\n")
}

// Reading Mac address from eeprom SAMA5D27 SOM1, address 0x51, bus 0
func readMac() []byte {
	i2cc, err := i2c.NewI2C(0x51, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer i2cc.Close()

	macaddr, _, err := i2cc.ReadRegBytes(0xfa, 0x6)
	if err != nil {
		log.Fatalf("Can't Read mac address%v", err)
	}
	return macaddr
}

// Test EEPROM SAMA5D27 SOM1
func testSomEE() {
	var (
		j, jj, jjj byte = 0x0, 0x0, 00
	)
	i2cc, err := i2c.NewI2C(0x51, 0)
	if err != nil {
		log.Fatalf("Can't open SOMeeprom:%v", err)
	}
	defer i2cc.Close()

	for z := 0; z <= 1; z++ {
		jj = 0
		jjj = jjj + 50
		for i := 0; i <= 0x7f; i++ {
			err = i2cc.WriteRegU8(jj, j)
			if err != nil {
				log.Fatalf("Can't Write SOMeeprom:%v", err)
			}
			time.Sleep(time.Millisecond * 5)

			mm, err := i2cc.ReadRegU8(jj)
			if err != nil {
				log.Fatalf("Can't Read SOMeeprom:%v", err)
			}
			if mm != j {
				fmt.Printf("error EEPROM SOM %x\n", jj)
			}
			jj++
		}
		fmt.Printf("Write-Read value '%2X' to byte %d (%d%%)\n", j, jj, jjj)
		j = 0xff
	}
	fmt.Printf("Test EEPROM SOM OK!\n")
}

//	Test mcp23008 ind, address 0x20, bus 1
func testMcpInd() {
	var (
		fl byte
	)

	mcp, err := i2c.NewI2C(0x20, 1)
	if err != nil {
		log.Fatalf("no open mcp ind%v", err)
	}

	defer mcp.Close()

	err = mcp.WriteRegU8(0x00, 0x03)
	if err != nil {
		log.Fatalf("no write reg 0x00 mcp ind%v", err)
	}
	err = mcp.WriteRegU8(0x01, 0x00)
	if err != nil {
		log.Fatalf("no write reg 0x01 mcp ind%v", err)
	}
	err = mcp.WriteRegU8(0x02, 0x03)
	if err != nil {
		log.Fatalf("no write reg 0x01 mcp ind%v", err)
	}
	err = mcp.WriteRegU8(0x03, 0x03)
	if err != nil {
		log.Fatalf("no write reg 0x01 mcp ind%v", err)
	}
	err = mcp.WriteRegU8(0x04, 0x03)
	if err != nil {
		log.Fatalf("no write reg 0x01 mcp ind%v", err)
	}

	fl = 0x28

	for i := 1; i < 2; i++ {

		err = mcp.WriteRegU8(0x09, fl)
		if err != nil {
			log.Fatal(err)
		}

		fl = fl ^ 0xd4

		time.Sleep(500 * time.Millisecond)
	}
}

func main() {

	testMcpInd()

	mcp, err := i2c.NewI2C(0x20, 1)
	if err != nil {
		log.Fatalf("no open mcp ind%v", err)
	}
	defer mcp.Close()

	for {
		intf, err := mcp.ReadRegU8(0x07)
		if err != nil {
			log.Fatalf("no open mcp ind%v", err)
		}

		switch intf {
		case 1:
			{
				intf, err := mcp.ReadRegU8(0x07)
				if err != nil {
					log.Fatalf("no open mcp ind%v", err)
				}
				fmt.Printf("intf %8b\r", intf)
			}
		case 2:
			{
				intf, err := mcp.ReadRegU8(0x07)
				if err != nil {
					log.Fatalf("no open mcp ind%v", err)
				}
				fmt.Printf("intf %8b\r", intf)
			}
		case 3:
			{
				intf, err := mcp.ReadRegU8(0x07)
				if err != nil {
					log.Fatalf("no open mcp ind%v", err)
				}
				fmt.Printf("intf %8b\r", intf)
			}
		}

	}

	// // Reading Mac address from eeprom SAMA5D27 SOM1, address 0x51, bus 0
	// macaddr := readMac()
	// fmt.Printf("Read mac address: %X\n", macaddr)

	// // Test M24M, first fill memory "0x00" and compare,
	// // then fill memory "0xFF" and compare.
	// testM24M()

	// // Test EEPROM SOM, first fill memory "0x00" and compare,
	// // then fill memory "0xFF" and compare.
	// testSomEE()
}
