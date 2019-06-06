node-pinger
===========

A microservice that:

- Retrieves the latest block known by a BTC node
- Queries the BTC node about block information

**WIP (as of 7 Jun 2019):**
1. ~~Compatibility with SupervisorD~~ ✅
2. ~~Standardadisation with `go-module`.~~ ✅
3. Logging
4. Containerisation (Docker)
5. Front-End with Angular

Install
-------
***Preparation on the RPC Server side***<br />
Locate the `bitcoin.conf` file of your target BTC node, and therein
- Check that `server=1` so that the ir accepts JSON-RPC commands;
- Make sure you locate the values of the parameters `rpcuser`, and `rpcpassword`;
You'll need them later for `RPC_USER`, `RPC_PASS`. (**Never** use this as your wallet password.)
- Take note of the parameter `rpcport` –usually set at `8332`, please keep it that way–, and also of the IP address of the node.
You'll use them together later as, eg `RPC_HOST=192.168.1.100:8332`;
- Make sure that the system from which you'll execute this code is included in `rcpallowip` (to allow within the local network, try something like `192.168.1.1/24`, and please never set it to `*`).

**Compilation**<br />
While in `/path/to/node-pinger` execute:
```
go get ./btcinfo/...
go build -o btcinfo/bin/btcinfo btcinfo/cmd/main.go
```
This will get the dependencies, and create the binary of the `btcinfo` service to be executed by `Supervisord`.
<br /> 

**Configuration files, {supervisord,btcinfo}.conf**<br />
Locate the files `/path/to/node-pinger/btcinfo`.<br />
<br />
In the file `supervisord.example.conf`, adequate the following line:
```
files=/path/to/node-pinger/btcinfo/btcinfo.conf
```
and copy the file as `supervisor.conf` (don't worry, it's *.gitignore*d).<br />

In the file `btcinfo.example.conf`, adequate the following lines according to the values you got from the RPC Server (as mentioned above):
```
directory=/path/to/btcinfo
command=/path/to/btcinfo/bin/btcinfo
environment=RPC_HOST="192.168.1.100:8332",RPC_USER="YOUR_USER",RPC_PASS="YOUR_PASSWORD"
```
and copy the file as `btcinfo.conf` (don't worry, it's also *.gitignore*d).<br />

**Running with Supervisord**<br />
Assuming you have Supervisord installed (macOS: `brew install supervisord`), execute:
```
sudo supervisord -c btcinfo/supervisord.conf
```
and the service will be running on your system on port `8080`.

Usage
-----
**Endpoints**<br>
- `GET /` 
  - Will retrieve the latest block of the node, eg `{latestBlock: 579375}`
- `GET /blocks?hash=HASH_OF_THE_BLOCK_IN_QUESTION`
  - *TODO:* Move to or add endpoint `/blocks/{hash}`
  - Will retrieve information of the block with the hash you give as JSON.
  - (_See note on verbosity regarding the parameter `tx` of the response_)
 <br />
 
 **Note on verbosity of the RPC server**<br />
 The library used by this command, `github.com/btcsuite/btcd/btcjson` as a one-to-one implementation of the RPC commands. This module offers two functions to retrieve block information:
 `getTxsFromBlockHash` and `getTxsFromBlockHashTx`. <br />
 Although the second is more verbose about the parameter `tx` of the response, This command uses the first.
 The reason lies in a (only partially documented) change of signature of the corresponding RPC CLI command `getblock`.<br />

 **TL;DR**: If you need more verbose information, you need to replace the library with this [pull request](https://github.com/btcsuite/btcd/pull/1112) and follow the instructions there.
 