## steps to reproduce
* run `/bin/bash reproduce.sh`
* read errors in logs: 
```bash
client_1   | time="2021-10-11T14:27:20Z" level=error msg="[Failed to create producer]" error="server error: PersistenceError: org.apache.bookkeeper.mledger.ManagedLedgerException: Error while recovering ledger" producerID=11 producer_name=pulsar-cluster-1-5 topic="persistent://public/default/my-topic5-partition-1"
```