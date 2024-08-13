# Transfer SOL

Simple example of transferring lamports (SOL).

### Config dev net

```shell
solana config set --url https://api.devnet.solana.com
```

### Creating the example keypairs:

```shell
solana-keygen new --no-bip39-passphrase -o transfer-sol/accounts/ringo.json
solana-keygen new --no-bip39-passphrase -o transfer-sol/accounts/paul.json
solana-keygen new --no-bip39-passphrase -o transfer-sol/accounts/john.json
solana-keygen new --no-bip39-passphrase -o transfer-sol/accounts/george.json
```

### Viewing their public keys:

```shell
solana-keygen pubkey transfer-sol/accounts/ringo.json
solana-keygen pubkey transfer-sol/accounts/paul.json
solana-keygen pubkey transfer-sol/accounts/john.json
solana-keygen pubkey transfer-sol/accounts/george.json
```

```shell
Ringo:      3c5di8sz3rkag4LBLpjMxiMo7fAnPDMuyRQB6o4L4G1r
George:     FmLBYcNq1PYHsrtBVK4byrwjeqYtTFtLL6GjvYy4fpCM
Paul:       2dw4Ff2P9NfNqL2eCCMKmh4wNunhsvzNTGxMZJEVHtSA
John:       2AXzcKA3cXr1SMmGESTkP8pqx232njXQBCJJPJCb9vfJ
```

### Airdropping:

```shell
solana airdrop --keypair transfer-sol/accounts/ringo.json 2
solana airdrop --keypair transfer-sol/accounts/paul.json 1
solana airdrop --keypair transfer-sol/accounts/john.json 0.2
solana airdrop --keypair transfer-sol/accounts/george.json 0.3
```

### Viewing their balances:

```shell
solana account <pubkey> 
```

## Run the example:

In one terminal:
```shell
npm run reset-and-build
npm run simulation
```

In another terminal:
```shell
solana logs | grep "<program id> invoke" -A 7
```