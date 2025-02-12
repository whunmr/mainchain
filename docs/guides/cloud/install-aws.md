# AWS 101

#### Contents

[[toc]]

## Introduction to installing `und` on AWS EC2 instances

::: danger
If you intend to become a `MainNet` Validator, it is **HIGHLY** recommended that you practice on `TestNet` first in
order to _fully familiarise_ yourself with the process.
:::

This guide introduces a _very simple_ "quick start" single AWS EC2 instance, using a VPC with a single public subnet
to connect to a Mainchain Public network. Validator node operators are highly encouraged to explore more sophisticated
architecture configurations to increase the security, reliability and availability of their Validator node - for
example, multi layered network with both private/public subnets, one or more "sentry" full (non-validator) nodes placed
in front of your main (hidden in a private subnet) Validator node to handle and relay communication from the outside
world, and reverse proxy for any RPC access via the sentries etc. in addition to implementing hardware KMS solutions to
protect validator private keys.

::: danger Important
This guide should not be considered the default, full, out of the box solution for a Validator node, but more an "AWS
101" guide to familiarise the reader with the core concepts involved in setting up the minimum AWS EC2 instance and
associated service requirements in order to operate a Validator node. It should be considered a starting point giving
you the initial building blocks from which to build a more sophisticated network/node architecture to support and
protect your Validator node.
:::

::: tip
any public IPs generated by AWS during this guide are not static - they will change if the EC2 instance is restarted. It
is recommended that the reader also investigates AWS's "elastic IPs" (static IPs) for their public-facing node(s).
:::

Where the guide prompts users to open a terminal, Windows 10 users should use PowerShell.

This guide assumes the reader already has an AWS account.

## Part 1: Create an SSH key pair

In order to configure your node(s) and install the necessary software, you will be connecting to the instance via SSH.
SSH connections to an EC2 instance must be authenticated using an SSH Key Pair - password authentication is (and should
be) disabled for security.

In the AWS Console, go to the EC2 dashboard and in the left menu, under "Network & Security", click "Key Pairs" followed
by the "Create Key Pair" button.

Enter a suitable name, for example `aws-ec2-und-validator-node`. Leave the file format as "pem", and click the "Crete
key pair" button.

This will download a private key called `aws-ec2-und-validator-node.pem` (or whatever you called it) which will be used
to log into your EC2 instance via SSH.

::: danger IMPORTANT
Do not lose this key - it is used to access your node's EC2 instance via SSH, and there is no way to recover it from the
AWS console if it's lost. Keep is safe, keep is secure, keep it secret!
:::

Open a terminal on your local PC, and check if you have a `$HOME/.ssh` directory:

```bash
ls -la $HOME/.ssh
```

If you do not have a `$HOME/.ssh` directory, create it:

```bash
$ mkdir $HOME/.ssh
$ chmod 700 $HOME/.ssh
```

Next, move the downloaded private key into the `$HOME/.ssh` directory, and tighten the file's permissions, replacing any
bold text appropriately:

```bash
$ mv /path/to/aws-ec2-und-validator-node.pem $HOME/.ssh
$ cd $HOME/.ssh
$ chmod 400 aws-ec2-und-validator-node.pem
```

## Part 2: Create a VPC (Network)

We need to create a network for our node. This will be a simple VPC with a single public subnet. Validator operators are
encouraged to explore more sophisticated configurations for their production node.

From the Services menu (top-left in the header) in AWS Console, scroll down to "Networking & Content Delivery" and
click "VPC", then click the "Launch VPC Wizard" button.

1. Select the "VPC with a single public subnet" and click "Select"
2. Leave the defaults and choose a suitable name in the "VPC Name" section.
3. Choose a more descriptive name for the Subnet, then click the "Create VPC" button.

::: tip
The names you choose can be anything, but should make the VPC and subnet easily identifiable to you. These names will
not be public and only visible to you through the AWS console.
:::

Next, we need to configure the VPC to automatically assign a public IP to instances in the network.

1. In the left menu of the VPC Dashboard, click "Subnets".
2. Select your new subnet from the list, and click the "Actions" dropdown, followed by "Modify auto-assign IP settings".
3. Tick the "Auto assign IPv4" box, and click "Save".

## Part 3: Create and launch an EC2 instance

The EC2 instance is the Virtual Machine where the node will be installed an run. We'll create and launch a single Linux
instance, and connect it to the network created in the previous part.

1. From the Services menu in AWS Console, click "EC2" under "Compute", followed by the "Launch Instance" button.
2. On the "1. Choose AMI" tab, use the search input to find the AMI ID "ami-0f2b4fc905b0bd1f1". Click "Community AMIs"
   in the results, find "CentOS Linux 7 x86_64 HVM EBS ENA" and click the "Select" button.
3. On the "2. Choose Instance Type" tab, tick `t2.small` or `t2.medium` and click "Next: Configure Instance Details".
4. In the "Network" section, select the VPC created in the previous part. Leave the rest as the defaults, and click
   the "6. Configure Security Group" tab at the top.
5. Give the security group a meaningful name and description - for example "und-validator-node". You will be able to
   find and edit this security group in the AWS EC2 console, under "Network & Security -> Security Groups" once it has
   been created.

Next, we need to configure some firewall rules for the instance. There is already a default rule for SSH (port 22)
access, but this needs tightening to restrict SSH access. We also need to add rules for P2P so the node can communicate
with other nodes, and RPC so that you can broadcast Txs to your own node.

1. In the "Source" column, for the SSH rule, click the dropdown that currently says "Custom" and select "My IP". This
   will ensure only your computer's IP can log in via SSH. Keep in mind, you will need to update this value if you do
   not have a static IP for your PC. This step is important, since the default value of `0.0.0.0/0` means that anyone,
   anywhere can attempt to access your EC2 instance via SSH.

::: tip Note
if your IP changes, you will need to update this value.
:::

2. Add a description, for example "SSH for my PC"

Next, we need to create a rule for the `P2P` port:

1. Click "Add Rule"
2. Leave the "type" as "Custom TCP Rule"
3. Leave the "Protocol" as "TCP
4. Set the Port Range as `26656`
5. Click the drop-down in the Source column and select "Anywhere"
6. Set the description as something like "UND Node P2P"

Next, we need to create a rule for the `RPC` port:

1. Click "Add Rule"
2. Leave the "type" as "Custom TCP Rule"
3. Leave the "Protocol" as "TCP
4. Set the Port Range as 26657
5. Click the drop-down in the Source column and select "My IP" (you want to restrict sending Txs to the node to your own
   IP)
6. Set the description as something like "UND Node RPC"

::: tip Note
port `26657` can be closed on your Validator node once you have registered your validator.
:::

**Next**, click on the **"4. Add Storage"** tab.

::: tip Note
As the state DB grows, the disk size requirements will grow. Check in our [Discord](https://discord.gg/SeB69w5)
for the latest values
:::

1. Change the size from 10 to **100** Gb.
2. Optionally, configure disk encryption.

Click "Review and Launch"

Review the details are correct, then click "Launch". This will prompt you to select a key pair to use on the instance.
Select the key you created in Part 1.

Once launched, make a note of the instance ID (top box titled "Your instances are now launching"), click your instance
ID link in the same box. This will take you to the Instances console, with your new instance already highlighted. Click
the pencil icon in the "Name" column, and give your instance a name.

Finally, make a note of the "IPv4 Public IP" value for your instance - you will need this to log in via SSH in the next
part.

## Part 4: Log in and configure instance via SSH

We now need to log in to the EC2 instance, install the prerequisites, then install and configure the UND node software.
This will all be done via SSH.

::: tip Note
any text in `[square_brackets]` (_including_ the square brackets) in the following commands should be replaced with your
own values. For example, `[aws_private_key]` should be replaced with the name of the file downloaded
in [Part 1](#part-1-create-an-ssh-key-pair), and `[vm_ip]` with the public IP address of your EC2 instance.
:::

The default username for our CentOS EC2 instance is `centos`.

Note for Windows users: Windows 10 should have an SSH client available in the PowerShell terminal. Older Windows
versions will require [PuTTY](https://www.chiark.greenend.org.uk/~sgtatham/putty/).

In a terminal on your local PC, run the following:

```bash
ssh -i $HOME/.ssh/[aws_private_key] centos@[vm_ip]
```

This will log you in to the EC2 instance via SSH.

### Part 4.1: Install the prerequisites

Once logged in to the VM via SSH, run:

```bash
sudo yum update -y
```

Install EPEL:

```bash
sudo yum install epel-release -y
```

Finally, install the following additional software:

```bash
sudo yum install nano jq wget -y
```

### Part 4.2: Install the und binary

See [https://github.com/unification-com/mainchain/releases](https://github.com/unification-com/mainchain/releases) for
the latest release, and replace the archive name with the latest `linux` version in the commands below.

```bash
wget https://github.com/unification-com/mainchain/releases/download/1.5.1/und_v1.5.1_linux_x86_64.tar.gz
tar -zxvf und_v1.5.1_linux_x86_64.tar.gz
sudo mv und /usr/local/bin/und
```

This should install the binary into `/usr/local/bin/und`. Verify the installation was successful:

```bash
which und
```

should output:

```bash
/usr/local/bin/und
```

and:

```bash
und version --log_level=""
```

should output something similar to:

```yaml
1.5.1
```

## Part 5: Initialising your full node

We'll now initialise and configure the `und` node itself. As previously, any text in `[square_brackets]` in the
following commands should be replaced with your own values accordingly. If you are not currently logged in to the EC2
instance via SSH, do so.

Once logged in, run:

```bash
und init [your_node_tag]
```

`[your_node_tag]` can be any ID you like, but is limited to ASCII characters (alphanumeric characters, hyphens and
underscores)

### Download the latest Genesis file.

The following command downloads the latest genesis for the respective network. Command is all one line:

#### TestNet

```bash
curl https://raw.githubusercontent.com/unification-com/testnet/master/latest/genesis.json > $HOME/.und_mainchain/config/genesis.json
```

#### MainNet

```bash
curl https://raw.githubusercontent.com/unification-com/mainnet/master/latest/genesis.json > $HOME/.und_mainchain/config/genesis.json
```

Get the current chain ID from genesis. Make a note of the output, it'll be required in commands later in the guide.
Command is all on one line:

```bash
$ jq --raw-output '.chain_id' $HOME/.und_mainchain/config/genesis.json
```

### Get seed nodes

::: danger IMPORTANT
Please ensure you get the correct seed node information for the network you would like to join! Remember to change the
directory if you are using something other than the default `$HOME/.und_mainchain` directory!
:::

Your node will need to know at least one seed node in order to join the network
and begin P2P communication with other nodes in the network. The latest seed information will always be available at
each network's respective Github repo:

#### TestNet: [https://github.com/unification-com/testnet/blob/master/latest/seed_nodes.md](https://github.com/unification-com/testnet/blob/master/latest/seed_nodes.md)

#### MainNet: [https://github.com/unification-com/mainnet/blob/master/latest/seed_nodes.md](https://github.com/unification-com/mainnet/blob/master/latest/seed_nodes.md)

Go to the repo for the network you are connecting to and copy one or more of the seed nodes (you only need
the `id@address:port`)

Edit your node configuration file using nano:

```bash
nano $HOME/.und_mainchain/config/config.toml
```

Hit <kbd>Ctrl</kbd>+<kbd>W</kbd>, type `[p2p]` (including the square brackets) and hit return - this will take you to
the `[p2p]` section of the config file, which begins with:

```toml
##### peer to peer configuration options #####
[p2p]
```

Find the `external_address = ""` variable about 9 lines below, and set it to your `vm_ip:26656` e.g.:

```toml
external_address = "11.22.33.44:26656"
```

Find the `seeds = ""` variable about 12 lines below, and add the seed node information between the double quotes (comma
separated, no spaces if more than one). For example:

```toml
seeds = "node_id@ip:port"
```

Next, hit <kbd>Ctrl</kbd>+<kbd>W</kbd>, type `[rpc]` (including the square brackets) and hit return - this will take you
to the `[rpc]` section of the config file, which begins with:

```toml
##### rpc server configuration options #####
[rpc]
```

About 3 lines under this, find:

```toml
laddr = "tcp://127.0.0.1:26657"
```

Change the value to:

```toml
laddr = "tcp://0.0.0.0:26657"
```

Hit <kbd>Ctrl</kbd>+<kbd>X</kbd> followed by `y` and then return to save the file and exit nano.

::: warning Note
you should revert the [rpc] configuration for port 26657 to:

```toml
laddr = "tcp://127.0.0.1:26657"
```

once you have run the `create-validator` command below. See [Part 8: Final cleanup](#part-8-final-cleanup) for further
details.
:::

### Gas Prices & Pruning

It is good practice to set the minimum-gas-prices value in `$HOME/.und_mainchain/config/app.toml`, in order to protect
your full node from spam transactions. This should be set as a decimal value in `nund`, and the recommended value is
currently `25.0nund`. This means your node will ignore any Txs with a gas price below this value. To do so, open
up `$HOME/.und_mainchain/config/app.toml` in a text editor, and set `minimum-gas-prices`

```bash
nano $HOME/.und_mainchain/config/app.toml
```

Change:

```toml
minimum-gas-prices = ""
```

To, for example:

```toml
minimum-gas-prices = "25.0nund"
```

:::tip Note
Validator nodes do not really need to keep a state history, so you should be able to safely set `pruning = "everything"`
:::

Hit <kbd>Ctrl</kbd>+<kbd>X</kbd> followed by `y` and then return to save the file and exit nano.

### State Syncing from Snapshots

The default method for syncing your node with the network is to start from `genesis`, and replay every block to the 
current block. As the chain grows, this can potentially take several days to complete. Thankfully, Cosmos SDK >= 0.42, 
which is used by the latest `und` software, can use State Syncing  from Snapshots to quickly sync your node from a safe 
checkpoint. This potentially reduces the sync time to no more than an hour or so, and in most cases mere minutes.

Setting this up requires a few more steps

1. Run the following command to get the latest block hash and height. For **TestNet**:

```bash
curl -s https://rest-testnet.unification.io/blocks/latest | jq '.|[.block_id.hash,.block.header.height]'
```

For **MainNet**:

```bash
curl -s https://rest.unification.io/blocks/latest | jq '.|[.block_id.hash,.block.header.height]'
```

Example output:

```json
[
  "820275B5EE63EDA2923886A01C0B1196A7CE1D96A89FA0D774942999C6698AAC",
  "1052423"
]
```

2. Using the output from the above command, configure `[statesync]` section in `.und_mainchain/config.toml`:

```toml
enable = true
rpc_servers = "TWO_RPC_NODES"
trust_height = 1052423
trust_hash = "820275B5EE63EDA2923886A01C0B1196A7CE1D96A89FA0D774942999C6698AAC"
trust_period = "168h0m0s"
discovery_time = "30s"
temp_dir = ""
chunk_request_timeout = "60s"
chunk_fetchers = "4"
```

The `rpc_servers` requires two RPC nodes for verification.

For **TestNet**, replace `TWO_RPC_NODES` with:

`sync1-testnet.unification.io:26657,sync2-testnet.unification.io:26657`

For **MainNet**:

`sync1.unification.io:26657,sync2.unification.io:26657`

e.g.:

```toml
rpc_servers = "sync1.unification.io:26657,sync2.unification.io:26657"
```

Or any RPC servers of your choice for the target network.

### Check connection

Finally, check that your node can connect to and sync with the network:

```bash
und start
```

Depending on the method used to sync your node, you should start seeing some output. If you are using `statesync`,
you should start seeing the following:

```bash
11:53AM INF Discovered new snapshot format=1 hash="V0���&�U1�J0�yP4A%�/���GŽ@\x05�<�j" height=1051600 module=statesync
```

After a few seconds (or at most, minutes), you should see your node start downloading the blocks:

```bash
11:56AM INF received proposal module=consensus proposal={"Type":32,"block_id":{"hash":"632E122ADDF385954FB8598FEE7D89EB09D7E93746FB36D2F12DECFEB7F07D9E","parts":{"hash":"E8246C504B9BC14275874A90C95E6AA035678302AD3BF9269B6F253B04C038BE","total":1}},"height":1052494,"pol_round":-1,"round":0,"signature":"HYJz0rV7o6bNm7za82sj1Az1rV25qVkLh9Y4s0K95nf86uVq+YmuDIf3LtIP7pDfFEYErxNVyeSplPGh7IVHDQ==","timestamp":"2022-05-19T10:56:03.273030584Z"}
11:56AM INF received complete proposal block hash=632E122ADDF385954FB8598FEE7D89EB09D7E93746FB36D2F12DECFEB7F07D9E height=1052494 module=consensus
11:56AM INF finalizing commit of block hash=632E122ADDF385954FB8598FEE7D89EB09D7E93746FB36D2F12DECFEB7F07D9E height=1052494 module=consensus num_txs=0 root=7EC77102840743503BD71FD89F60FD0B912DD0DE27575408B6AD67990CE4A6B8
11:56AM INF executed block height=1052494 module=state num_invalid_txs=0 num_valid_txs=0
11:56AM INF commit synced commit=436F6D6D697449447B5B3234312031333020323235203834203837203234372031303220323820323435203234382032372032303920313432203133372031303620353920343520323220323020313737203135342032303320323338203136352030203231322034392031383620313638203433203839203233395D3A3130304634457D
11:56AM INF committed state app_hash=F182E15457F7661CF5F81BD18E896A3B2D1614B19ACBEEA500D431BAA82B59EF height=1052494 module=state num_txs=0
11:56AM INF indexed block height=1052494 module=txindex
```

Hit <kbd>Ctrl</kbd>+<kbd>C</kbd> to stop the node - it will be configured as a background service next.

::: danger IMPORTANT
keep your `$HOME/.und_mainchain/config/node_key.json` and `$HOME/.und_mainchain/config/priv_validator_key.json` files
safe! These are required for your node to propose and sign blocks. If you ever migrate your node to a different host/VM
instance, you will need these.
:::

## Part 6: Running und as a daemon

Once you have initialised and tested the `und` node, it can be set up as a background daemon on the server
using `systemctl`. This means that you can easily start/stop/restart the service, and do not need to leave the SSH
session open while `und` is running.

If you're not logged in to your EC2 instance via SSH, log in. If you are still logged in, and have not stopped the `und`
node, hit <kbd>Ctrl</kbd>+<kbd>C</kbd> to stop the node.

We need to use the nano text editor to create the service configuration. Run:

```bash
sudo nano /etc/systemd/system/und.service
```

Add the following:

```
[Unit]
Description=Unification Mainchain Validator Node

[Service]
User=centos
Group=centos
WorkingDirectory=/home/centos
ExecStart=/usr/local/bin/und start --home=/home/centos/.und_mainchain
LimitNOFILE=4096

[Install]
WantedBy=default.target
```

Hit <kbd>Ctrl</kbd>+<kbd>X</kbd> followed by `y` and then return to save the file and exit nano.

Run:

```bash
sudo systemctl daemon-reload
sudo systemctl enable und
```

to update `systemctl`

You can now start and stop the und daemon in the background using:

```bash
$ sudo systemctl start und
$ sudo systemctl stop und
```

Finally, you can monitor the log output for the service by running:

```bash
sudo journalctl -u und --follow
```

and use <kbd>Ctrl</kbd>+<kbd>C</kbd> to exit the journalctl command. You can now log out of your SSH session and und
will continue running in the background.

## Part 7: Become a Validator

Important: ensure your node has fully synced with the network before continuing.
If not already logged in, log in to your EC2 instance via SSH, and run:

```bash
sudo journalctl -u und --follow
```

Once fully synced (check the downloaded height against the current clock in the block explorer), hit <kbd>Ctrl</kbd>
+<kbd>C</kbd> to quit `journalctl` ( `und` will continue to run in the background).

You will need your Validator node's Tendermint public key in order to register it on the Mainchain network as a
validator. Whilst still in the SSH session, run:

```bash
und tendermint show-validator
```

Make a note of the output, as you will need this in later commands, where it will be referred to
as `[you_validator_public_key]`.

You can now exit the SSH session - the rest of the commands will be run in a terminal on your local PC.

Go
to [https://github.com/unification-com/mainchain/releases/latest](https://github.com/unification-com/mainchain/releases/latest)
and download the latest `und` archive for your OS - for example, `und_v1.5.1_windows_x86_64.tar.gz`.

Open a terminal/PowerShell, and `cd` into the directory where you extracted the `und` executable:

```bash
cd path/to/extracted/und_directory
```

As previously, any text in `[square_brackets]` in the following commands should be replaced with your own values
accordingly.

If you do not already have a wallet/account, you can create one (on your local PC, not in the SSH session) by running:

```bash
./und keys add account_name
```

If you already have a wallet, you can import the account using:

```bash
./und keys add account_name --recover
```

in which case, you will be prompted for the mnemonic and a password to secure the wallet.

From here, it is assumed the reader has an account with sufficient FUND from which to self-delegate and create their
Validator node. The account you use to self-delegate will become the "owner" account of the Validator node.

On your local PC, run the following command, replacing any text in `[square_brackets]` accordingly with your own values:

```bash
./und tx staking create-validator \
--amount=[stake_in_nund] \
--pubkey=[your_validator_public_key] \
--moniker="[your_ev_moniker]" \
--website="[your_website]" \
--details="[description]" \
--security-contact="[security_email]" \
--chain-id=[chain_id] \
--commission-rate="[0.10]" \
--commission-max-rate="[0.20]" \
--commission-max-change-rate="[0.01]" \
--min-self-delegation="1" \
--gas="auto" \
--gas-prices="25.0nund" \
--gas-adjustment=1.5 \
--from=account_name \
--node=tcp://[vm_ip]:26657
```

`[stake_in_nund]` = (required) the amount of FUND in `nund` you are self-delegating. You can use
the `und convert 1000 fund nund` command to convert FUND to nund. E.g. `1000000000000nund`.

::: warning Note
do not enter more nund than you have in your wallet and ensure you have enough left over to pay for this and future Tx
fees!
:::

`[your_validator_public_key]` = (required) the public key output from the previous und tendermint show-validator
command.

`[your_ev_moniker]` = (required) a publicly visible ID/tag for your Validator node.

`[your_website]` = (optional) website promoting your node

`[description]` = (optional) short description of your node

`[security_email]` = (optional) security contact for your organisation

`[chain_id]` = the network (e.g. `FUND-TestNet-2`, or `FUND-MainNet-2`) you are creating a
validator on - this was obtained earlier in the guide via the `jq` command

`[account_name]` = the account self-delegating the FUND, previously created/imported with the `und keys add` command

`[vm_ip]` = the IP address of your EC2 instance running the full node - you can get this from your AWS EC2 Instances
console.

#### Commission Rates

Your commission rates can be set using the `--commission-rate`
, `--commission-max-change-rate and` `--commission-max-rate` flags.

`--commission-rate`: The % commission you will earn from delegators' rewards. Keeping this low can attract more
delegators to your node.

`--commission-max-rate`: The maximum you will ever increase your commission rate to - you cannot raise commission above
this value. Again, keeping this low can attract more delegators.

`--commission-max-change-rate`: The maximum you can increase the commission-rate by per day. For example, if your
maximum change rate is 0.01, you can only make changes in 0.01 increments, so from 0.10 (10%) to 0.11 (11%).

::: warning
The values for `--commission-max-change-rate` and `--commission-max-rate` flags cannot be changed after the
create-validator command has been run.
:::

Finally, the `--min-self-delegation` flag is the minimum amount of `nund` you are required to keep self-delegated to
your validtor, meaning you must always have at least this amount self-delegated to your node.

For example:

```bash
./und tx staking create-validator \
--amount=1000000000000nund \
--pubkey=undvalconspub1zcjduepq6yq7drzefkavsrxhxk69cy63tj3r... \
--moniker="MyAwesomeNode" \
--website="https://my-node-site.com" \
--details="My node is awesome" \
--security-contact="security@my-node-site.com" \
--chain-id=FUND-TestNet-2 \
--commission-rate="0.05" \
--commission-max-rate="0.10" \
--commission-max-change-rate="0.01" \
--min-self-delegation="1" \
--gas="auto" \
--gas-prices="25.0nund" \
--gas-adjustment=1.5 \
--from=my_new_wallet \
--node=tcp://33.44.55.66:26657
```

Your validator node should now be registered and contributing to the network. To verify, on you local PC, run the
following:

```bash
./und query staking validator \
$(und keys show [account_name] --bech=val -a) \
--chain-id=[chain_id] \
--node=tcp://[vm_ip]:26657
```

Replacing `[account_name]`, `[chain_id]` and `[vm_ip]` accordingly. Assuming you are going through this guide
on `TestNet`, you should also see your node listed
in [https://explorer-testnet.unification.io/validators](https://explorer-testnet.unification.io/validators)

## Part 8: Final cleanup

Finally, it's a good idea to close the **RPC** port (`26657`) on your validator node, leaving only the **P2P**
port (`26656`) open so that it can communicate with other nodes. This can be done by deleting the firewall rule for the
RPC port, and by reverting the RPC `laddr` configuration value in `config.toml` to the `127.0.0.1` IP address (which
will restrict the node's RPC access to `localhost`).

Further interaction with the network can be done by spinning up a separate non-validator full node (on your local PC for
example), and broadcasting transactions via that node instead.

::: danger Important
Do not alter the `P2P` port (`26656`) firewall rules. If you do, your validator node will not be able to communicate
with its peers.
:::

That's it - you should now have a full Validator node up and running on a very basic AWS EC2 instance. Once again -
Validator node operators are highly encouraged to explore more sophisticated architecture configurations to increase the
security, reliability and availability of their Validator node.
