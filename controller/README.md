# DIY Roller Shades Controller

This is the arduino-powered controller for the roller shades.

## Serial Commands

At startup the device will emit `Ready.`. Commands are case-insensitive and newline `\n` delimited.

- The basic syntax is `AT+` followed by a command.
- Commands that return values end in `?`.
- Return values are encoded as `+COMMAND_NAME:VALUE`.
- Commands that take an argument end with `=ARG_VALUE`.
- Responses will end with `OK` or an error.
- Error format: `ERR:CODE:NAME:DETAILS`.

Available Commands:

| Command | Arg/Return Type | Description |
| --- | --- | --- |
| `AT` | | NOOP operaion |
| `AT+STATUS?` | enum | Returns the current device status. |
| `AT+OPEN` | | Opens the shade. |
| `AT+CLOSE` | | Closes the shade. |
| `AT+STOP` | | Stops movement of the shade. |
| `AT+LOCK?` | 0 or 1 | Returns the lock state of the device. |
| `AT+LOCK=` | 0 or 1 | Sets the lock state of the device. |
| `AT+POS?` | int, 0 to 100 | Returns the current position of the shade. 100 is fully open. |
| `AT+POS=` | int, 0 to 100 | Sets the position of the shade, rounded to nearest 5. 100 is fully open. |

### AT

This command always returns `OK`.

```
SEND: AT
RECV: OK
```

## AT+STATUS?

Returns the current status of the device.

| Value | Description |
| --- | --- |
| `IDLE` | Device is not currently performing any action. |
| `CAL` | Needs calibration. |
| `HOME` | Performing calibration. |
| `UP` | Moving up (opening). |
| `DOWN` | Moving down (closing). |

Idle:
```
SEND: AT+STATUS?
RECV: +STATUS:IDLE
RECV: OK
```

### AT+OPEN

Commands the shade to move to the full-open position. This has the same effect as `AT+POS=100`.

```
SEND: AT+OPEN
RECV: OK
```

Locked or uncalibrated:
```
SEND: AT+OPEN
RECV: ERR:700:Forbidden:Device is locked.
```


### AT+CLOSE

Commands the shade to move to the full-closed position. This has the same effect as `AT+POS=0`.

```
SEND: AT+CLOSE
RECV: OK
```

Locked or uncalibrated:
```
SEND: AT+CLOSE
RECV: ERR:700:Forbidden:Device is locked.
```


### AT+STOP

Stops the shade movement at the current position.

```
SEND: AT+STOP
RECV: OK
```

Locked or uncalibrated:
```
SEND: AT+STOP
RECV: ERR:700:Forbidden:Device is locked.
```

### AT+LOCK?

Return the current lock state.

Locked:
```
SEND: AT+LOCK?
RECV: +LOCK:1
RECV: OK
```

Unlocked:
```
SEND: AT+LOCK?
RECV: +LOCK:0
RECV: OK
```

### AT+LOCK=1

Lock the device in position.

```
SEND: AT+LOCK=1
RECV: OK
```

### AT+LOCK=0

Unlock the device.

```
SEND: AT+LOCK=0
RECV: OK
```

Uncalibrated:
```
SEND: AT+LOCK=0
RECV: ERR:701:Not Allowed:Device must be calibrated first.
```

### AT+POS?

Return the current position of the shade.

```
SEND: AT+POS?
RECV: +POS:55
RECV: OK
```

Uncalibrated:
```
SEND: AT+POS?
RECV: ERR:501:Action Failed:Device must be calibrated first.
```

### AT+POS=

Set the new position. This command may return before the shade reaches the desired position.

```
SEND: AT+POS=55
RECV: OK
```

Out of range:
```
SEND: AT+POS=128
RECV: ERR:601:Out of Range:The argument value is not between 0 and 100 included.
```

Uncalibrated:
```
SEND: AT+POS=55
RECV: ERR:501:Action Failed:Device must be calibrated first.
```

Locked:
```
SEND: AT+POS=55
RECV: ERR:700:Forbidden:Device is locked.
```
