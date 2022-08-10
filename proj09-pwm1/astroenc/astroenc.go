// This package contains the Struct and Methods for controlling an AMTT Encoder
// See this datasheet: https://www.cuidevices.com/product/resource/amt22.pdf
package astroenc

// In the case of odd parity, For a given set of bits, if the count of bits with a value of 1 is even,
// the parity bit value is set to 1 making the total count of 1s in the whole set (including the parity bit) an odd number. If the count of bits with a value of 1 is odd, the count is already odd so the parity bit's value is 0
/*


	Example 1:

		byte-1: 10101011 - 171
		byte-2: 10000000 - 128

					kk hhhhhh llllllll
					10 543210 76543210 (position in byte)

		full: 10 101011 10000000


		The odd bits
		h5 h3 h1 l7 l5 l3 l1
		1  1  1  1  0  0  0       even=1 == k1

		The even bits
		h4 h2 h0 l6 l4 l2 l0
		0  0  1  0  0  0  0        odd=0 == k0

	Example 2:

		   kk hhhhhh llllllll
			 10 543210 76543210

	full 01 100001 10101011
	14      100001 10101011

	The odd bits
	h5 h3 h1 l7 l5 l3 l1
	1  0  0  1  1  1  1        # of 1s is odd thus use 0 == k1

	The even bits
	h4 h2 h0 l6 l4 l2 l0
	0  0  1  0  0  0  1        # of 1s is even thus use 1 == k0

*/

import (
	"machine"
	"time"
)

// AMT22 constants
const AMT22_NOP byte = 0x00
const AMT22_RESET byte = 0x60
const AMT22_ZERO byte = 0x70
const RES12 int8 = 12
const RES14 int8 = 14

// Encoder
type RAEncoder struct {
	cs         machine.Pin
	resolution int8
	spi        *machine.SPI
}

func NewRA(spi *machine.SPI, cs machine.Pin, resolution int8) RAEncoder {

	return RAEncoder{
		spi:        spi,
		cs:         cs,
		resolution: resolution,
	}

}

// Configure RA encoder
func (ra *RAEncoder) Configure() {

	//
	// Channel select for encoder on the SPI bus
	// initialize high i.e. Not listening
	//
	ra.cs.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ra.cs.High()

}

// Zero the RA encoder
func (ra *RAEncoder) ZeroRA() {

	ra.WriteRead(AMT22_NOP, AMT22_ZERO)

	// allow time to reset
	time.Sleep(time.Millisecond * 240)

}

func (ra *RAEncoder) GetPositionRA() (r1, r2 byte) {

	r1, r2 = ra.WriteRead(AMT22_NOP, AMT22_NOP)

	return

}

// DEVTODO need to validate results and return as one int16
func (ra *RAEncoder) WriteRead(b1 byte, b2 byte) (r1, r2 byte) {

	// DEVTODO add mutex so that only one go routine can talk on a channel at a time
	//         we dont want to clash with the display which is shareing the SPI buss

	// Select RA channel
	ra.cs.Low()
	time.Sleep(time.Microsecond * 3) // wait min time see datasheet

	// byte 1
	r1, _ = ra.spi.Transfer(b1)
	time.Sleep(time.Microsecond * 3) // wait min time see datasheet

	// byte 2
	r2, _ = machine.SPI0.Transfer(b2)
	time.Sleep(time.Microsecond * 3) // wait min time see datasheet

	// de-select RA channel
	ra.cs.High()

	return

}

/*

work in progress


package main

import "fmt"

func main() {
	b1 := byte(0b10101011)
	b2 := byte(0b10000000)
	fmt.Printf("b1: %08b\nb2: %08b \n", b1, b2)

	r1 := uint16(b1)
	r1 = r1 << 8
	fmt.Printf("r1: %016b\n", r1)

	r2 := uint16(b2)
	fmt.Printf("r2: %016b\n", r2)

	r := r1 | r2
	fmt.Printf("r: %016b\n", r)

	var checkBitUpperIsSet bool = isKthBitSet(r,16)
	var checkBitLowerIsSet bool = isKthBitSet(r,15)


	fmt.Printf("checkBitUpperIsSet: %v,  checkBitLowerIsSet: %v", checkBitUpperIsSet, checkBitLowerIsSet)



  fmt.Printf("\n 15:%v 13:%v 11:%v 9:%v 7:%v 5:%v 3:%v 1:%v \n", isKthBitSet(r,15) ,isKthBitSet(r,13) ,isKthBitSet(r,11) ,isKthBitSet(r,9) ,isKthBitSet(r,7) ,isKthBitSet(r,5) ,isKthBitSet(r,3) , isKthBitSet(r,1) )
	odd := !(isKthBitSet(r,15) != isKthBitSet(r,13) != isKthBitSet(r,11) != isKthBitSet(r,9) != isKthBitSet(r,7) != isKthBitSet(r,5) != isKthBitSet(r,3) != isKthBitSet(r,1))

	fmt.Printf("odd check: %v\n", odd == checkBitLowerIsSet)

  fmt.Printf("\n 16:%v 14:%v 12:%v 10:%v 8:%v 6:%v 4:%v 2:%v \n", isKthBitSet(r,16) ,isKthBitSet(r,14) ,isKthBitSet(r,12) ,isKthBitSet(r,10) ,isKthBitSet(r,8) ,isKthBitSet(r,6) ,isKthBitSet(r,4) , isKthBitSet(r,2) )

  even := !(isKthBitSet(r,16) != isKthBitSet(r,14) != isKthBitSet(r,12) != isKthBitSet(r,10) != isKthBitSet(r,8) != isKthBitSet(r,6) != isKthBitSet(r,4) != isKthBitSet(r,2))

  fmt.Printf("even check: %v\n", even == checkBitUpperIsSet)


    //we got back a good position, so just mask away the checkbits
	r &= 0x3FFF
	fmt.Printf("\nPosition: %v", r)
}

func isKthBitSet(n, k uint16) bool {

	flag := n & (1 << (k - 1))

	if flag != 0 {
		return true
	} else {
		return false
	}

}


//For a given set of bits, if the count of bits with a value of 1 is even, the parity bit value is set to 1
func negitiveParity(n uint16, startBit int8) bool {

	var count int = 0

	for i := startBit; i < 16; i=i+2 {
		if isKthBitSet(n,i) {
			count++
		}
	}

	if(count%2==0){
      return true
    }else{
      return false
    }

}


*/
