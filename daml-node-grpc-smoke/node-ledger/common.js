const ledger = require('@digitalasset/daml-ledger');
const fs = require('fs');

const getSandboxClient = async (host, port) => {
    try{
        const client = await ledger.DamlLedgerClient.connect({host,port});
        client.partyManagementClient
        return client;
    }catch(err){
        console.log(`Error: ${err}`);
        process.exit(1);
    }
}

const listParties = async (client) => {
    const parties = client.partyManagementClient.listKnownParties();
    return parties;
}

const allocateParties = async (client, party, displayName) => {
    const response = client.partyManagementClient.allocateParty({
        partyIdHint: party,
        displayName: displayName
    });

    return response;
}

const getPackages = async (client) => {
    const packages = client.packageClient.listPackages();
    return packages;
}

const extractModuleNames = (payload) => {
    return payload.getDamlLf1().getModulesList().map(a => 
        a.getName().getSegmentsList().reduce((prev, curr) => `${prev}.${curr}`)
    );
}

const extractPayloads = async (client, packages) => {
    return await packages.packageIds.map(async (id) => {
        const package = await client.packageClient.getPackage(id);
        const payload = await ledger.lf.ArchivePayload.deserializeBinary(package.archivePayload);
        return {
            packag_id: id,
            module_name: extractModuleNames(payload)
        };
    });
}

const uploadDar = async (client, file) => {
    const content = fs.readFileSync(file);
    const contentBase64 = content.toString('base64');
    client.packageManagementClient.uploadDarFile({
        darFile: contentBase64
    })
}

const listPackages = async (client) => {
    const packages = client.packageManagementClient.listKnownPackages();
    return packages;
}

module.exports = {
    getSandboxClient,
    listParties,
    allocateParties,
    extractPayloads,
    uploadDar,
    listPackages,
}