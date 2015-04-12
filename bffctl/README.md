# bffctl

bffctl is the BFF controller application.  It launches and kills daemonized instances of BFF.

## Usage

    bffctl {start|kill}

## Installation

The bffctl application requires the following directory structure in order to work.

    - bin/
        - bffctl
    - util/
        - bffcore
        - daemonize
    - var/
