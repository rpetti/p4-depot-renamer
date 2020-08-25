# Perforce Depot Renamer

## General

This tool can be used against a checkpoint file to perform a "deep rename" of a depot.

## Disclaimer

This code is provided as-is, with no warranty. The author(s) take no responsbility for any damage it may do. Always take a back up of your Perforce instance before running this, and test the results thoroughly before going live!

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

1. Do a full `p4 verify -q //...` and fix any issues.
2. **Take a full checkpoint and backup of your server.**
3. Stop p4d.
4. Remove db.* files.
5. If the latest checkpoint is compressed, uncompress it. `gunzip checkpoint.###.gz`
6. Run the tool against the checkpoint file with the necessary options. (outlined below)
7. Restore the checkpoint. `p4d -r . -jr checkpoint.renamed`
8. Run database verification tests. `p4d -r . -xv`
9. Move the lbr/versioned files for the renamed depot. For example, if you renamed `depot` to `myproduct`, then rename the `depot` directory in the P4ROOT to `myproduct`.
10. Bring the server back online. `p4d -r . -p <port>`
11. Run a full lbr verification. `p4 verify -q //...`

**Note**: It's likely that verify will encounter "BAD" files at this point. This is likely because of `$Id$` tags in the files messing up the calculated checksum since the path has changed. In order to fix these, you'll need to run `p4 verify -qv //newdepotname/...`.

## Usage

Example usage:

```
p4-depot-renamer -cp checkpoint.100 -depot mydepot -rename-to -mynewdepot -o checkpoint.100.renamed
```

- `-cp <checkpoint file>` - Which checkpoint file to process.
- `-depot <depot name>` - The name of the depot to be renamed.
- `-rename-to <depot name>` - The new name of the depot.
- `-o <new checkpoint file>` - The new checkpoint file to write.
- `-f` - Forces overwriting the new checkpoint file if it already exists.
- `-batch` - Batch mode - allows to batching multiple renames in one pass, also enables including/excluding certain transforms.

## Advanced usage

The batch argument allows to replace -depot/-rename-to with a YAML file that defines multiple transforms. Additionally, the batch argument allows
specifying either an inclusion list (only transforms in the list are applied) or an exclusion list (transforms in the list are not applied).

For example, the following YAML allows to move multiple depots under a single one:

```
- PathFrom: old_depot_first
  PathTo: new_depot/first
  ExcludedTransforms:
    - db.depot:0
    - db.depot:3
    - db.domain:0
- PathFrom: old_depot_second
  PathTo: new_depot/second
  ExcludedTransforms:
    - db.depot:0
    - db.depot:3
    - db.domain:0
```

Note that in this case we want to exclude changes to domains/depots - we should instead manually create the new depots prior to the rename.


## Running a demo example locally

```
go run . -batch demo_batch.yaml -cp demo_checkpoint -o new_checkpoint -f
```
