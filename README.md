# Perforce Depot Renamer

## General

This tool can be used against a checkpoint file to perform a "deep rename" of a depot.

## Disclaimer

This code is provided as-is, with no warranty. The author(s) take no responsbility for any damage it may do. Always take a back up of your Perforce instance before running this!

## Limitations

- This was developed and tested against p4d 2018.1 checkpoints. YMMV for other versions.
- This tool has only been tested against standard depots. **Graph, Stream, and Remote depots will likely not work!**
- It has only been tested against a full checkpoint. **It may not work against journal data.**
- It has only been tested against a depot with the standard mapping. If your depot's Map is not the default ie. `<depotname>/...`, you may wish to change it before proceeding.

## Installation

```
go install github.com/rpetti/p4-depot-renamer
```

## How to Use

1. Take a full checkpoint and backup of your server, just in case!
2. Make sure the server (p4d) is not running.
3. Remove db.* files.
4. Un-gzip the latest checkpoint file.
5. Run the tool against the checkpoint file with the necessary options. (outlined below)
6. Restore the checkpoint. `p4d -r . -jr checkpoint.renamed`
7. Run database verification tests. `p4d -r . -xv`
8. Move the lbr/versioned files for the renamed depot. For example, if you renamed 'depot' to 'myproduct', then rename the 'depot' directory in the P4ROOT to 'myproduct'.
9. Bring the server back online. `p4d -r . -p <port>`
10. Run a full lbr verification. `p4 verify -q //...`

## Usage

Example usage:

```
p4-depot-renamer -cp checkpoint.100 -depot mydepot -rename-to -mynewdepot -o checkpoint.100.renamed
```
- `-cp <checkpoint file>` - Which checkpoint file to process.
- `-depot <depot name>` - The name of the depot to be renamed.
- `-rename-to <depot name>` - The new name of the depot.
- `-o <new checkpoint file>` - The new checkpoint file to write.