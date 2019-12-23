# psguard

Simple program to poll a process for resources consumption and existence. On request it can revive the process when dies. Built with [gopsutil](https://github.com/shirou/gopsutil) and [gjson](https://github.com/tidwall/gjson).

## Build

**NOTE**: Go 1.13 or higher is required.
```sh
# clone the repository
$ git clone https://github.com/gsscoder/psguard.git

# change the working directory
$ cd psguard

# build the executable
$ sh ./build.sh

# test if it works
$ ./artifacts/psguard -version
```

## Usage

```sh
$ ./psguard -help   
psguard: Polls processes for resources consumption and existence
Version: v0.2.0
usage:
  -poll duration
        defines polling interval (default 1s)
  -restart
        restart process if terminated
  -version
        displays version information
  -wait duration
        time to wait before polling again after a restart (default 5s)

$ ps -A | grep firefox
92697 ??         1:58.34 /Applications/Firefox.app/Contents/MacOS/firefox
$ ps -A | grep Evernote
92663 ??         0:06.16 /Applications/Evernote.app/Contents/MacOS/Evernote
$ jj -i psguard.json -o psguard.json -v 92697 constraints.firefox.pid
$ jj -i psguard.json -o psguard.json -v 92663 constraints.evernote.pid

$ ./psguard -restart 2>>psguard.log &
[2] 80757
$ tail -f psguard.log
2019/12/23 08:07:57 firefox: CPU constraint of 0.50% violated by +38.28%
2019/12/23 08:07:57 firefox: Memory constraint of 1.00% violated by +2.14%
2019/12/23 08:08:03 evernote: CPU constraint of 0.50% violated by +0.68%
2019/12/23 08:08:03 evernote: Memory constraint of 1.00% violated by +0.24%
...
$ sudo kill 80757
Password:
[2]  + terminated  ./psguard -restart 2>> psguard.log
```

## Configuration
**psguard.json**:
```json
{
    "constraints": {
        "firefox":
            {
                "pid": 56575,
                "cpu": 0.5,
                "mem": 1
            },
        "evernote":
            {
                "pid": 44208,
                "cpu": 0.5,
                "mem": 1
            }
    }
}
```

### Notes
- Command [JJ](https://github.com/tidwall/jj) used to edit `psguard.json` from terminal can be installed on **macOS** with `brew install tidwall/jj/jj`.