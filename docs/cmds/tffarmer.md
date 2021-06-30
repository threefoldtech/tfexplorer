# tffarmer
### Explorer client for farmers to create and manage farms

## Setup
### Generate binary 
In the tfexplorer base dir:
```
make tffarmer
```
This will create the bin in the tfexplorer directory found as `bin/tffarmer`

## Usage

### Main options
- explorer url:  `--explorer=<EXPLORER_URL>`
- seed: `--seed=<PATH_TO_SEED_FILE>`

    NOTE: if seed is not given, a new identity is created and registered on explorer and saved in the location `~/.config/tffarmer.seed` for further usage
- debug enabling: `--debug`

### Farms
### Register new farm

`tffarmer farm register [command options] <farm_name>`

Command options:

- --addresses : wallet addresses. The format is `ASSET:ADDRESS`
- --email : email address of the farmer. It is used to send communication to the farmer and for the minting


### Update an existing farm
`tffarmer farm update [command options] [arguments...]`

Command options:
- --id : farm ID (default: 0) of farm to be updated
- --addresses : wallet address. the format is '`ASSET:ADDRESS`: e.g: 'TFT:GBUPOYJ7I4D4TYSFXPJNLSATHCCF2QDDQCIIIXBG7CV7S2U36UMAQENV'
- --email : email address of the farmer. It is used to send communication to the farmer and for the minting

### Add IP address to an existing farm
`tffarmer farm addip [command options] [arguments...]`

Command options:
- --id : farm ID (default: 0) to add IP to
- --address : IP address
- --gateway : gateway address

### Delete IP address from an existing farm
`tffarmer farm deleteip [command options] [arguments...]`

Command options:

- --id : farm ID (default: 0) to delete IP from
- --address : IP address

### List user Farms
`tffarmer farm list`