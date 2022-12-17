# Generic Example

```golang

func PublishMsg[M MsgInterface](m M, mb MsgBroker) {

	//
	// reflect to get message properties
	//
	msg := reflect.ValueOf(&m).Elem()
	msgKind := fmt.Sprintf(",%v", msg.Field(0).Interface())

	//
	// If it is a log message check to see if it is loggable
	//
	if msgKind == string(Log) {
		msgLogLevel := LogLevel(fmt.Sprintf(",%v", msg.Field(1).Interface()))
		if !mb.isLoggable(msgLogLevel) {
			fmt.Println("msg.PublishMsg Don't publish log message due to logging level")
			return
		}
	}

	//
	// Create msgStr
	//
	msgStr := fmt.Sprintf("^%v", msg.Field(0).Interface())
	for i := 1; i < msg.NumField(); i++ {
		msgStr = msgStr + fmt.Sprintf(",%v", msg.Field(i).Interface())
	}
	msgStr = msgStr + "~"

	//
	// Write to uart
	//
	if mb.uartUp != nil {
		mb.uartUp.Write([]byte(msgStr))
	}
	if mb.uartDn != nil {
		mb.uartDn.Write([]byte(msgStr))
	}

}
```
