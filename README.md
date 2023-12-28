# sfstypes

WIP Go library providing SmartFoxServer 2X data types.


## Installation

```
go get github.com/jannikdc/sfstypes
```

## Usage

```go
sfsobj := sfstypes.NewSFSObject() // Creates a new SFSObject
sfsarr := sfstypes.NewSFSArray()  // Creates a new SFSArray

// Put data into the SFSObject
sfsobj.Put("Test", "Test")
sfsobj.Put("Value", 0)

// Put data into the SFSArray
sfsarr.Add("Test")
sfsarr.Add(0)

// Print the SFSObject in json and hex format
fmt.Println(sfsobj.ToJson())
fmt.Println(sfsobj.GetHexDump())

// Print the SFSObject in json and hex format
fmt.Println(sfsarr.ToJson())
fmt.Println(sfsarr.GetHexDump())
```

## Disclaimer

All rights to the original code and protocol belong to their respective owner. This repository does not grant rights to the original code. If you are the owner of the original code and have concerns about its presence in this repository, please contact me, and I will promptly address the issue.

