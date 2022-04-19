/*
We are performing following steps using this file
1. We are making connection with the network using connection profile.
2. We will create a new CA Client.
3. We will create a new wallet folder and create a new file there.
4. Check if admin user is already enrolled or not. If it is enrolled it 
    will throw an error message otherwise it will enroll admin user by creating
    all the required identities for the user like Certificate and the private
    key.
5. We save a all the certificate and private key to a wallet folder.
*/

'use strict';
const FabricCAServices = require('fabric-ca-client');
const { Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function main() {
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new CA client for interacting with the CA.
        // Here we will try to interact with the CA.
        // before we go and enroll any user we have to make new ca client 
        // which will connect with the organization1 Certificate Authority(CA)
        const caInfo = ccp.certificateAuthorities['ca.org1.example.com'];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        // this will create new fabric CA client connections
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        // Create a new file system based wallet for managing identities.
        // this will create user identities
        const walletPath = path.join(process.cwd(), 'wallet');
        // this is a function which will write a new file to our walletPath
        // because when ever we create a new identity for admin user, it has
        // to write one admin file into this wallet folder  
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const identity = await wallet.get('admin');
        if (identity) {
            console.log('An identity for the admin user "admin" already exists in the wallet');
            return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        // this will give us the secret id 
        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
        // we will use x509 algorithim to create the identities
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: 'Org1MSP',
            type: 'X.509',
        };
        // now this identity would be put inside the wallet folder
        await wallet.put('admin', x509Identity);
        console.log('Successfully enrolled admin user "admin" and imported it into the wallet');

    } catch (error) {
        console.error(`Failed to enroll admin user "admin": ${error}`);
        process.exit(1);
    }
}

main();
