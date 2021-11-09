package main

import (
	"fmt"
	"net"
	"time"
	"bufio"
	//"strings"
	"encoding/hex"
	"strconv"
	)


func main() {
	var counter int 
	for{
		counter++
		fmt.Println("Connecting to server ... try ", counter , " times ")
		con, err := net.Dial("tcp", "localhost:5050")
		if err != nil {
			fmt.Println("Error to connect , Lets Try again ")
			time.Sleep(2 * time.Second)
			continue;
		}
		defer con.Close()
		fmt.Println("Connected to server successfully ... ")
		
		for {
			dataBytes := []byte{0x78, 0x78, 0x0D, 0x01,  0x01, 0x23, 0x45, 0x67, 0x89 ,0x01, 0x23, 0x45, 0x00, 0x01, 0x8C, 0xDD, 0x0D, 0x0A}
			_, err = con.Write(dataBytes)
			if err != nil {
				fmt.Println("Fail to send data to server ... ")
				break;
			}
			//fmt.Println("Sent to server .. ", string(dataBytes))
			
			message, err := bufio.NewReader(con).ReadBytes('\n')
		
			if err != nil {
				fmt.Println("Fail to get response from server ... ")
				break;
			}else{
				data := hex.EncodeToString(message[:len(message)])
				fmt.Println("DATA  :: ", data)
				strToCheckCrc := data[4:12]
				fmt.Println("str for crc check " , strToCheckCrc)
				CRC := data[12:16]
				fmt.Println("CRC CODE : ", CRC)
				status, err := CRCcheck(strToCheckCrc, CRC)
				if err != nil {
					fmt.Println("log: CRC error")
				}else{
					fmt.Println("log: status ", status)
				}
				
			}	
	 		time.Sleep(10 * time.Second)
		}
	}

}

type Params struct {
	Poly   uint16
	Init   uint16
	RefIn  bool
	RefOut bool
	XorOut uint16
	Check  uint16
	Name   string
}

// Predefined CRC-16 algorithms.
// List of algorithms with their parameters borrowed from here -  http://reveng.sourceforge.net/crc-catalogue/16.htm
//
// The variables can be used to create Table for the selected algorithm.
var (
	CRC16_ARC         = Params{0x8005, 0x0000, true, true, 0x0000, 0xBB3D, "CRC-16/ARC"}
	CRC16_AUG_CCITT   = Params{0x1021, 0x1D0F, false, false, 0x0000, 0xE5CC, "CRC-16/AUG-CCITT"}
	CRC16_BUYPASS     = Params{0x8005, 0x0000, false, false, 0x0000, 0xFEE8, "CRC-16/BUYPASS"}
	CRC16_CCITT_FALSE = Params{0x1021, 0xFFFF, false, false, 0x0000, 0x29B1, "CRC-16/CCITT-FALSE"}
	CRC16_CDMA2000    = Params{0xC867, 0xFFFF, false, false, 0x0000, 0x4C06, "CRC-16/CDMA2000"}
	CRC16_DDS_110     = Params{0x8005, 0x800D, false, false, 0x0000, 0x9ECF, "CRC-16/DDS-110"}
	CRC16_DECT_R      = Params{0x0589, 0x0000, false, false, 0x0001, 0x007E, "CRC-16/DECT-R"}
	CRC16_DECT_X      = Params{0x0589, 0x0000, false, false, 0x0000, 0x007F, "CRC-16/DECT-X"}
	CRC16_DNP         = Params{0x3D65, 0x0000, true, true, 0xFFFF, 0xEA82, "CRC-16/DNP"}
	CRC16_EN_13757    = Params{0x3D65, 0x0000, false, false, 0xFFFF, 0xC2B7, "CRC-16/EN-13757"}
	CRC16_GENIBUS     = Params{0x1021, 0xFFFF, false, false, 0xFFFF, 0xD64E, "CRC-16/GENIBUS"}
	CRC16_MAXIM       = Params{0x8005, 0x0000, true, true, 0xFFFF, 0x44C2, "CRC-16/MAXIM"}
	CRC16_MCRF4XX     = Params{0x1021, 0xFFFF, true, true, 0x0000, 0x6F91, "CRC-16/MCRF4XX"}
	CRC16_RIELLO      = Params{0x1021, 0xB2AA, true, true, 0x0000, 0x63D0, "CRC-16/RIELLO"}
	CRC16_T10_DIF     = Params{0x8BB7, 0x0000, false, false, 0x0000, 0xD0DB, "CRC-16/T10-DIF"}
	CRC16_TELEDISK    = Params{0xA097, 0x0000, false, false, 0x0000, 0x0FB3, "CRC-16/TELEDISK"}
	CRC16_TMS37157    = Params{0x1021, 0x89EC, true, true, 0x0000, 0x26B1, "CRC-16/TMS37157"}
	CRC16_USB         = Params{0x8005, 0xFFFF, true, true, 0xFFFF, 0xB4C8, "CRC-16/USB"}
	CRC16_CRC_A       = Params{0x1021, 0xC6C6, true, true, 0x0000, 0xBF05, "CRC-16/CRC-A"}
	CRC16_KERMIT      = Params{0x1021, 0x0000, true, true, 0x0000, 0x2189, "CRC-16/KERMIT"}
	CRC16_MODBUS      = Params{0x8005, 0xFFFF, true, true, 0x0000, 0x4B37, "CRC-16/MODBUS"}
	CRC16_X_25        = Params{0x1021, 0xFFFF, true, true, 0xFFFF, 0x906E, "CRC-16/X-25"}
	CRC16_XMODEM      = Params{0x1021, 0x0000, false, false, 0x0000, 0x31C3, "CRC-16/XMODEM"}
)

// Table is a 256-word table representing polinomial and algorithm settings for efficient processing.
type Table struct {
	params Params
	data   [256]uint16
}

func CRCcheck(data string, errorCode string) (status string, e error) {
	var dataType string
	/* check error code */
	incomingErrorHex, _ := hex.DecodeString(data)
	incomingDataCRC := Checksum(incomingErrorHex)               //Error code in uint16
	crcCheck := strconv.FormatUint(uint64(incomingDataCRC), 16) //Error code in string
	if len(crcCheck) == 3 {
		crcCheck = fmt.Sprint("0", crcCheck)
	}
	if errorCode != crcCheck { //consider as void data
		fmt.Println("** VOID Data")
		dataType = "V"
	} else {
		dataType = "A"
	}

	return dataType, nil
}

// MakeTable returns the Table constructed from the specified algorithm.
func MakeTable(params Params) *Table {
	table := new(Table)
	table.params = params
	for n := 0; n < 256; n++ {
		crc := uint16(n) << 8
		for i := 0; i < 8; i++ {
			bit := (crc & 0x8000) != 0
			crc <<= 1
			if bit {
				crc ^= params.Poly
			}
		}
		table.data[n] = crc
	}
	return table
}

// Init returns the initial value for CRC register corresponding to the specified algorithm.
func Init(table *Table) uint16 {
	return table.params.Init
}

// Update returns the result of adding the bytes in data to the crc.
func Update(crc uint16, data []byte, table *Table) uint16 {
	for _, d := range data {
		if table.params.RefIn {
			d = ReverseByte(d)
		}
		crc = crc<<8 ^ table.data[byte(crc>>8)^d]
	}
	return crc
}

// Complete returns the result of CRC calculation and post-calculation processing of the crc.
func Complete(crc uint16, table *Table) uint16 {
	if table.params.RefOut {
		return ReverseUint16(crc) ^ table.params.XorOut
	}
	return crc ^ table.params.XorOut
}

// Checksum returns CRC checksum of data usign scpecified algorithm represented by the Table.
func Checksum(data []byte) uint16 {
	table := MakeTable(CRC16_X_25)
	crc := Init(table)
	crc = Update(crc, data, table)
	return Complete(crc, table)
}



func Hex2Int(hexStr string) (int64, error) {
	intValue, err := strconv.ParseInt(hexStr, 16, 0)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func Bin2Int(binStr string) (int64, error) {
	intValue, err := strconv.ParseInt(binStr, 2, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func hex2Bin(hexStr string) (string, error) {
	ui, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return "", err
	}

	format := fmt.Sprintf("%%0%db", len(hexStr)*4)
	return fmt.Sprintf(format, ui), nil
}

func ReverseByte(val byte) byte {
	var rval byte = 0
	for i := uint(0); i < 8; i++ {
		if val&(1<<i) != 0 {
			rval |= 0x80 >> i
		}
	}
	return rval
}

func ReverseUint8(val uint8) uint8 {
	return ReverseByte(val)
}

func ReverseUint16(val uint16) uint16 {
	var rval uint16 = 0
	for i := uint(0); i < 16; i++ {
		if val&(uint16(1)<<i) != 0 {
			rval |= uint16(0x8000) >> i
		}
	}
	return rval
}















