# DIY Motorized Roller Shades

This repo has the code and parts used for DIY roller shades.

## Rendering STLs for Printing

1. Adjust any values in `config.toml`
2. Install dependencies
    ```go get ./...```
3. Build the tool
    ```go build -o parts.bin ./parts```
4. Render all (to 0.1mm)
    ```./parts.bin -res 0.1```

## Assembly Instructions

To render the assembly diagram:
```./parts.bin -res 0.5```

Note: It may take a long time to render at high detail.


## Calibration

Upon power-up the device will need to be calibrated for it's start and end positions.

The device will be in a locked state until calibrated.

### Calibration Procedure

1. Pull the shade down to the desired "closed" position
2. Hold the programming button down for a minimum of **3 seconds**
3. Release the button (it will begin to go up)
4. The "open" position will be recorded when the shade is no longer able to move.
5. (optional) you may press the programming button again at any time while moving to signal the "open" position instead.

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
| `AT` || NOOP operaion |
| `AT+NAME?` | string | Returns the current name of the device. |
| `AT+NAME= `| string | Sets the device name to the provided value. 32 bytes max |
| `AT+OPEN` || Opens the shade. |
| `AT+CLOSE` || Closes the shade. |
| `AT+STOP` || Stops movement of the shade. |
| `AT+LOCK?` | 0 or 1 | Returns the lock state of the device. |
| `AT+LOCK=` | 0 or 1 | Sets the lock state of the device. |
| `AT+POS?` | int, 0 to 100 | Returns the current position of the shade. 100 is fully open. |
| `AT+POS=` | int, 0 to 100 | Sets the position of the shade, rounded to nearest 5. 100 is fully open. |

### AT

This command always returns `OK`.

```
> AT
< OK
```

### AT+NAME?

Returns the current name. Default is `DIY Roller Shades`.

```
> AT+NAME?
< +NAME:DIY Roller Shades
< OK
```

### AT+NAME=

Saves the new name to EEPROM. Must be between 1 and 32 bytes.

```
> AT+NAME=foobar
< OK
```

Too Short:
```
> AT+NAME=
< ERR:402:Invalid Args:New name must be specified.
```

Too Long:
```
> AT+NAME=
< ERR:402:Invalid Args:New name must be 32 characters or less.
```

### AT+OPEN

Commands the shade to move to the full-open position. This has the same effect as `AT+POS=100`.

```
> AT+OPEN
< OK
```

Locked or uncalibrated:
```
> AT+OPEN
< ERR:700:Forbidden:Device is locked.
```


### AT+CLOSE

Commands the shade to move to the full-closed position. This has the same effect as `AT+POS=0`.

```
> AT+CLOSE
< OK
```

Locked or uncalibrated:
```
> AT+CLOSE
< ERR:700:Forbidden:Device is locked.
```


### AT+STOP

Stops the shade movement at the current position.

```
> AT+STOP
< OK
```

Locked or uncalibrated:
```
> AT+STOP
< ERR:700:Forbidden:Device is locked.
```

### AT+LOCK?

Return the current lock state.

Locked:
```
> AT+LOCK?
< +LOCK:1
< OK
```

Unlocked:
```
> AT+LOCK?
< +LOCK:0
< OK
```

### AT+LOCK=1

Lock the device in position.

```
> AT+LOCK=1
< OK
```

### AT+LOCK=0

Unlock the device.

```
> AT+LOCK=0
< OK
```

Uncalibrated:
```
> AT+LOCK=0
< ERR:701:Not Allowed:Device must be calibrated first.
```

### AT+POS?

Return the current position of the shade.

```
> AT+POS?
< +POS:55
< OK
```

Uncalibrated:
```
> AT+POS?
< ERR:501:Action Failed:Device must be calibrated first.
```

### AT+POS=

Set the new position. This command may return before the shade reaches the desired position.

```
> AT+POS=55
< OK
```

Out of range:
```
> AT+POS=128
< ERR:601:Out of Range:The argument value is not between 0 and 100 included.
```

Uncalibrated:
```
> AT+POS=55
< ERR:501:Action Failed:Device must be calibrated first.
```

Locked:
```
> AT+POS=55
< ERR:700:Forbidden:Device is locked.
```
