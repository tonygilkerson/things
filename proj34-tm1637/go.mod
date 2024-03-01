module aeg

go 1.22.0

require tinygo.org/x/drivers/tm1637 v0.27.0


// see: https://github.com/tinygo-org/drivers/issues/554
// Create a directory tm1637) at the same level with my current project (cd .. ; mkdir tm1637)
// Create a Go module in that directory (cd tm1637; go mod init tm1637)
// Copy the files tm1637.go and registers.go from this commit: https://github.com/bxparks/tinygo-drivers/tree/3caaa7af6250b51d222108fd2c953e98f4e4b784/tm1637
// In your main project's directory change the go.mod file to include:


replace tinygo.org/x/drivers/tm1637 => ./tm1637