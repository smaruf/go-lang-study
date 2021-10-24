
### Dependencies:
 1. API gateway
 2. Kafka stream
 3. Python
 4. Golang: ??
 5. Impala + Kudu
 6. Postgres
 7. Redis: ??
 3. Proto: ??

### Integration:
 ```
 commchat-->wallet-bridge-->upay-api-->
                          |-->kafka-stream-->wallet-db (audit-tail, account-info, journal, ledger)
                                          |-->analytic-db (statistics, history, audit-tail-summary, reporting)
                                          |-->tcash-app (blockchain, audit-tail)
```
### Concerns:
 1. Time scheduling
 2. Data sync
 3. Fail management
 4. Report generation
