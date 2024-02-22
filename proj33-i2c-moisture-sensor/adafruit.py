  def read(self, reg_base, reg, buf, delay=0.008):
        """Read an arbitrary I2C register range on the device"""
        self.write(reg_base, reg)
        if self._drdy is not None:
            while self._drdy.value is False:
                pass
        else:
            time.sleep(delay)
        with self.i2c_device as i2c:
            i2c.readinto(buf)

  def write(self, reg_base, reg, buf=None):
      """Write an arbitrary I2C register range on the device"""
      full_buffer = bytearray([reg_base, reg])
      if buf is not None:
          full_buffer += buf

      if self._drdy is not None:
          while self._drdy.value is False:
              pass
      with self.i2c_device as i2c:
          i2c.write(full_buffer)            

_TOUCH_BASE           = const(0x0F)
_TOUCH_CHANNEL_OFFSET = const(0x10)

    def moisture_read(self):
        """Read the value of the moisture sensor"""
        buf = bytearray(2)

        self.read(_TOUCH_BASE, _TOUCH_CHANNEL_OFFSET, buf, 0.005)
        ret = struct.unpack(">H", buf)[0]
        time.sleep(0.001)

        # retry if reading was bad
        count = 0
        while ret > 4095:
            self.read(_TOUCH_BASE, _TOUCH_CHANNEL_OFFSET, buf, 0.005)
            ret = struct.unpack(">H", buf)[0]
            time.sleep(0.001)
            count += 1
            if count > 3:
                raise RuntimeError("Could not get a valid moisture reading.")

        return ret            