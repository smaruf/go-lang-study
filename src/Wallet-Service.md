**Features:**

	a) Core-Wallet (gPRC/API)
		1. Account Management (service: waccount: account-mng + pin-mng)
			1. Created
			2. Failed
			3. Closed
			4. Recovered
			5. Settled
			6. Locked
			7. Blacklist
			6. Type
		2. Transaction (service: wtrxn: transactions)
			1. Confirmed
			2. Failed
			3. Reverted
			4. Disbursed
			5. Locking (debit:y, credit: n)(?)
			6. Time latency
			7. Pending
		3. Account Queries (service: wquery: basic queries)
			1. Account balance
			2. Last Transaction
			3. Transaction History.. Only last 10
		4. PIN management (service: waccount)
			1. Transaction PIN
			2. Generated PIN
			3. User specific key for OTP
		4. Wallet Management (service: wmng: additional management)
		    1. Account recovery
		    2. Account Settlement
		    3. PIN reset (?)
		    4. Audit	
	b) OTP/PIN (service: wpin: pin-handle + otp-handle + link:waccount)
		1. Generator.. User specific
			1. User specific key based
			2. From system
			3. From Auth-app
		2. Verification
			1. Successful
			2. Failed
		3. Threshold
			1. Retry
			2. Blocking
	c) Messaging (service: wmsg: sms + email + push)
		1. SMS
		2. Email
		3. Push Notification
	d) Streaming (service: wqueue: kafka + db + cross-module-async)
		1. Logging
		2. Data Sync
		3. Message
		4. Auto Disburse
		5. Third party ledger (?)(tcash)
		6. Reporting
	e) Database
		1. In-memory DB (need to choose)
			1. Account opening
			2. Transaction processing
			3. Balance calculation
			4. Last transaction
			5. Quick query
		2. Stored (by streaming or by configuration)(?)
			1. Finalized Account
			2. Finalized Transaction
			3. Audit-trail
			4. Secure logging
			5. General Ledger
		3. Statistics (?)
			1. Reporting (service: wreport:...)
				1. Day End
				2. Weekly
				3. Monthly
				4. Quarterly
				5. Half yearly
				6. Yearly
			2. Analytics (service: wanalytic:...)
				1. Account opening/closing count
				2. Trx count
				3. Settlement count
				4. System monitoring
	f) Load management 
		1. Reverse-Proxy (service: wgateway: api-mng + grpc-mng + encript-decript? + auth)(?)
		2. Scaling (?)* [depends on platform]
		3. Availability (?)* [depends on CI/CD]
	g) Cache
		1. Account Balance
		2. Account Access or Locking (?)
		3. Last Trx
		
	* Load testing:
		1. Basic feature
		2. 10E10 loading
		3. Failed detection
		