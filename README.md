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
Version: v0.3.0
usage:
  -poll duration
        defines polling interval (default 1s)
  -restart
        restart process if terminated
  -version
        displays version information
  -wait duration
        time to wait before polling again after a restart (default 5s)

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
        "process.groups":
            {
                "firefox":
                    {
                        "match": [ "firefox$" ],
                        "cpu": 0.5,
                        "mem": 1
                    },
                "evernote":
                    {
                        "match": [ "Evernote$" ],
                        "cpu": 0.5,
                        "mem": 1
                    }
                }
    }
}
```
Each process group is identified by a name (like `firefox`) and all process bound to it are selected using one or more regular expression. The match is done using the executable path of the process (you can easly discover it with command `ps -A`). Defined constraints are expressed in percentage and are checked for all processes of a group.

### Notes
- Restarted processes will die, when `psguard` terminates in case it wasn't started in background (tested only on **macOS**).
- For now it's a [Go language](https://golang.org/) learning project and not much tests has been done on it. 