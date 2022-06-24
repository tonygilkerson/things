# proj07 - Astrostepper with OLED

In this project I am combining the astrostepper and OLED projectes into one.

Components used:

* [product page](https://www.amazon.com/STEPPERONLINE-Stepper-Bipolar-42x42x39mm-4-wires/dp/B00W98OYE4?pd_rd_w=A8Ylc&content-id=amzn1.sym.bc622850-a717-4d94-96c3-7cc183488298&pf_rd_p=bc622850-a717-4d94-96c3-7cc183488298&pf_rd_r=B5FRSWXQAHWXZYYD9T4A&pd_rd_wg=QT3en&pd_rd_r=b1a31da3-8fb8-4f7c-abe6-ccf39265831e&pd_rd_i=B00W98OYE4&psc=1&ref_=pd_bap_d_rp_1_t) - STEPPERONLINE 0.9deg Nema 17 Stepper Motor Bipolar 0.9A 36Ncm/50oz.in 42x42x39mm 4-wires DIY

> See also my [component reference](https://github.com/tonygilkerson/things#components)

## Project Demo

![setup](img/todo.drawio.png)

```bash
tinygo flash -target=esp32-coreboard-v2  -port=/dev/cu.usbserial-0001
picocom --baud 115200 /dev/cu.usbserial-0001
```

## Mount and Motor Specs

todo

## References

* todo

## Ref

I need:
GT2 12T 5mm bore 6mm belt pulley (DEC motor shaft)
GT2 12T 5mm bore 6mm belt pulley (RA motor shaft)
https://www.amazon.com/gp/product/B01IMPM44O/ref=ox_sc_act_title_1?smid=A12MRQC2NA7LMA&psc=1

GT2 48T 6mm bore 6mm belt pulley (The other end)
https://www.amazon.com/gp/product/B07SR78PKY/ref=ox_sc_act_title_1?smid=A1NQCH9MN8OPZG&th=1

Stepper Motor
https://www.amazon.com/gp/product/B00W98OYE4/ref=ox_sc_act_title_4?smid=AWQBCGWISS7BL&psc=1