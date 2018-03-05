# DIY Motorized Roller Shades

This repo has the code and parts used for DIY roller shades.

[Controller Serial Commands](controller/README.md#serial-commands)

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
