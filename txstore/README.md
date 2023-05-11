This lib implements a tx storage backend for long-time archival purpose.

# Example log format

    {"tx_hash":"0x0000000000000000000000000000000000000000000000000000000000000000","monitor_id":"testing","timeSeen":1603154983.9674246,"peerInfo":"testing node","peerName":"testing node name","peerIP":"128.0.0.1","peerPort":1000,"txFrom":"0x0000000000000000000000000000000000000000","txTo":"0x0000000000000000000000000000000000000000","nonce":100,"gasPrice":"0","gasLimit":234234,"value":"0","payload":"7061796c6f6164","sigV":"0","sigR":"100","sigS":"0"}
    
# Notes

**DB based tx storage is outdated and being phased out.** It stores transactions in a slightly different format than that of structured log based storage.