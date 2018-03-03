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

