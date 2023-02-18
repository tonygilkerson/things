# LoRa - Hello World

This is my initial LoRa project. The goal is to program two LoRa devices with TinyGo and see if they can talk to each other.

## Erase Factory AT Firmware

As stated in the [Develop with STM32Cube MCU Package](https://wiki.seeedstudio.com/LoRa-E5_STM32WLE5JC_Module/#develop-with-stm32cube-mcu-package) doc:

>Factory AT Firmware is programmed with RDP(Read Protection) Level 1, developers need to remove RDP first with STM32Cube Programmer. Note that regression RDP to level 0 will cause a flash memory mass to erase and the Factory AT Firmware can't be restored again.

## Attach ST-Link V2

| ST-LINK Pin | LoRa-E5 Pin |
| ----------- | ----------- |
| 2 - SWCLK   | SWCLK       |
| 4 - SWDIO   | SWDIO       |
| 6 - GND     | GND         |
| 8 - 3.3V    | VCC         |

![wireup](wireup.drawio.png)
