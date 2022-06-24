# proj02 - Serial Monitor

This project shows how to use command line tool to monitor or debug your app running on a ESP32 microcontroller. Here are the project components.

* **TinyGo** - A Go compiler for small places, see [tinygo.org](https://tinygo.org/)
* **picocom CLI** - minimal dumb-terminal emulation program, see [man page](https://linux.die.net/man/8/picocom)
* **ESP32-WROOM-32D** - ESP32 dev board, see [Data Sheet](https://www.espressif.com/sites/default/files/documentation/esp32-wroom-32d_esp32-wroom-32u_datasheet_en.pdf) 

## Project Demo

Clone [this repo](https://github.com/tonygilkerson/things) from GitHub and navigate to the this projec's subfolder.

```sh
git clone https://github.com/tonygilkerson/things.git
cd things/proj02
```

Use `tinygo` to flash the program onto the board.

```sh
tinygo flash -target=esp32-coreboard-v2 -port /dev/cu.usbserial-0001
# The opress the "flash button" on the board
```

Use `picocom` to monitor the serial output

```sh
picocom --baud 115200 /dev/cu.usbserial-0001
```

To exit `picocom` type `control+a` then `control+x`

Here are is the `picocom` help:

```text
*** Picocom commands (all prefixed by [C-a])

*** [C-x] : Exit picocom
*** [C-q] : Exit without reseting serial port
*** [C-b] : Set baudrate
*** [C-u] : Increase baudrate (baud-up)
*** [C-d] : Decrease baudrate (baud-down)
*** [C-i] : Change number of databits
*** [C-j] : Change number of stopbits
*** [C-f] : Change flow-control mode
*** [C-y] : Change parity mode
*** [C-p] : Pulse DTR
*** [C-t] : Toggle DTR
*** [C-g] : Toggle RTS
*** [C-|] : Send break
*** [C-c] : Toggle local echo
*** [C-w] : Write hex
*** [C-s] : Send file
*** [C-r] : Receive file
*** [C-v] : Show port settings
*** [C-h] : Show this message
```

