# psguard

Simple program to poll a process for resources consumption and existence. On request it can revive the process when dies. Built with [gopsutil](github.com/shirou/gopsutil).

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
psguard: Polls a process for resources consumption and existence
Version: v0.1.0
usage:
  -cpu float
    	allowed CPU % usage (default 1)
  -mem float
    	allowed memory % usage (default 1)
  -pid int
    	pid of the process to monitor
  -poll duration
    	defines polling interval (default 1s)
  -restart
    	restart process if terminated
  -version
    	displays version information
  -wait duration
    	time to wait before polling again after a restart (default 5s)

$ ps -A | grep firefox
14195 ??         1:01.92 /Applications/Firefox.app/Contents/MacOS/firefox
$ ./psguard -pid 14195 -restart 2>>psguard.log &
[1] 30867
$ tail -f psguard.log
2019/12/18 20:52:49 CPU constraint of 1.00% violated by +9.63%
2019/12/18 20:52:49 Memory constraint of 1.00% violated by +2.84%
2019/12/18 20:52:51 CPU constraint of 1.00% violated by +9.60%
2019/12/18 20:52:51 Memory constraint of 1.00% violated by +2.80%
...
$ sudo kill 30867
Password:
[1]  + terminated  ./psguard -pid 14195 2>> psguard.log
```