# TODO

## rhash

- [ ] TODO(1): Add <https://github.com/spf13/viper> for configuration
- [ ] TODO(1a): Use XDG Base Directory Specification
- [ ] TODO(2) Add multiple inputs (Ex: --input $1 -i $2 -i $3)
- [ ] TODO(9): Check need of setting log level via flags (Ex: --log INFO, DEBUG, WARNING, ...)
- [ ] TODO(10): Recreate folder structure on destination Dir
- [ ] TODO(11): Create destinationPath if doesn't exist (maybe add a flag? force?)
- [ ] TODO(14): Check need of path validation or continue to use CustomFileInfo
- [ ] TODO(16): Check if has permission to move to destination
- [ ] TODO(19): Makefile; -ldflags "-s -w"
- [ ] TODO(21): Check WalkDir error/return
- [ ] TODO(24): Test simlinks, shorcuts, ...
- [X] TODO(3): Work on dry-run flag
- [X] TODO(4): Move some flags to rootCmd
- [X] TODO(5): Work on verbose flag
- [X] TODO(6): Work on uppercase flag
- [X] TODO(7): Work on recursive flag
- [X] TODO(8): Work on lenght/truncate flag
- [X] TODO(8a): Check MAX_PATH on windows
- [X] TODO(12): Check if go-git need git binary, if yes, drop module
- [X] TODO(13): Create a HashMachine interface, add Options
- [X] TODO(15): Add --hash fuzzy
- [X] TODO(17): Benchark reuse(.Reset()) vs recreate
- [X] TODO(18): Merge checkHash and getHashMachine?
- [X] TODO(20): Move fs related to mateusjdev/cmd/filesystem
- [X] TODO(22): Make CustomFileInfo private and create through ValidatePath()
- [X] TODO(23): Reuse enqueuePath for rename_helper_fuzzy and rename_helper_hash
