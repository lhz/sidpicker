# sidpicker
Terminal interface for browsing [HVSC](http://www.hvsc.c64.org/) and playing SID tunes therein.

## Requirements

- High Voltage SID Collection (download [here](http://www.hvsc.c64.org/#download))
- sidplayerfp  (download [here](https://sourceforge.net/projects/sidplay-residfp/files/sidplayfp/1.4/))
- C-64 ROM dumps (basic, kernal and chargen) (download [here](https://sourceforge.net/p/vice-emu/code/HEAD/tree/trunk/vice/data/C64/))

### Setup

#### Environment variables

The `HVSC_BASE` environment variable must be set to point to the base directory of your extracted HVSC installation.

Examples:

Linux or MacOS:
```
export HVSC_BASE=~/Download/C64Music
```
Windows:
```
setx HVSC_BASE C:\Download\C64Music
```

#### sidplayfp

The configuration file `sidplayfp.ini` needs to be edited to point to the C-64 ROM dumps, since many
SID tunes contain code that rely on these to run properly.

Linux:
```
Kernal Rom = /usr/local/lib64/vice/C64/kernal
Basic Rom = /usr/local/lib64/vice/C64/basic
Chargen Rom = /usr/local/lib64/vice/C64/chargen
```
