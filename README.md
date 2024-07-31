
# PeerFund 

## Overview
PeerFund is a decentralized  P2P Lending Platform  built on the Ethereum blockchain that enables peer-to-peer lending and borrowing.It eliminates the need for traditional financial intermediaries, such as banks or lending institutions. By using blockchain technology, it enables direct interactions between lenders and borrowers, thus reducing costs and increasing efficiency.

## Architecture
#### The platform comprises several components:

- Frontend: A user-friendly interface built with React.js, allowing users to interact with the blockchain.
- Backend: A set of APIs and services powered by Node.js to handle off-chain data and user authentication.
- Ethereum Blockchain: The decentralized network where all transactions are recorded and users have wallets.
### Key Features
#### Decentralized Lending and Borrowing:

- Lenders can create offers specifying the loan amount, interest rate, and repayment terms.
- Borrowers can request loans by submitting an application with details about the amount needed and the repayment plan.
- Loan Contract: Manages loan agreements, including the principal amount, interest rate, duration, and repayment schedule.
- Repayment Contract: Automates the repayment process, ensuring timely payments and penalties for late repayments.
- Ethereum Wallet Integration:Support  wallets for authentication and transactions.
- Enables secure and transparent repayment through Ethereum wallets.
#### Transparency and Security:
All transactions and contract states are recorded on the blockchain, providing full transparency.
The platform uses secure cryptographic methods to protect user data and funds.
### Prerequisites
Development Environment:
    - html/css
    - Javascript
    - Etherum(ganachi)
    - Golang
### Installation
1. Install Dependencies
- [Etherum](https://archive.trufflesuite.com/ganache/)
- [Golang](https://go.dev/)
###  How it works
1. Clone the Repository
```bash
git clone https://github.com/onyango-granton/peerFund.git
```
2. cd in the folder
```bash
cd peerFund
```
3. Run Etherum(ganachi)
4. Run the Program
```bash
go run .
```
### Usage
Lendee User login:
- Users logins on the platform using their credentials.
- User data verified .
- Lendee user profile:The profile provides an overview of the lendee  loan activities, including outstanding balances, repayment status, and loan history.
- Lendee  can browse available loan offers and apply for one.

Repayment Process:
- Lendee repay the loan through scheduled installments via their Ethereum wallet.
- The Repayment Contract ensures that repayments are recorded and penalties are applied if necessary.

### Contributing
We welcome contributions from the community! Please follow our contributing guidelines for more information on how to get involved.

## license
[License](/home/nyagooh/Downloads/peerFund/license)



