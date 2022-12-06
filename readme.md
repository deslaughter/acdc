# ACDC: Automated Campbell Diagram Capabilities

This application provides a web-based graphical user interface to generating Campbell Diagrams and visualizing mode shapes for OpenFAST turbine models.

Additional information can be found in [Aeroelastic Modeling for Distributed Wind Turbines](https://www.nrel.gov/docs/fy22osti/81724.pdf).

## Current State

This software is currently under development and the results it produces, if it does produce results, should not be considered reliable.

## Building

Building the `acdc` executable requires the [Go compiler](https://go.dev/dl/), version 1.19 or newer. Once the compiler is installed, run the following command from the terminal to download and build the program:

`go install github.com/deslaughter/acdc@latest`

The `acdc` executable will be created in `~/go/bin/`.