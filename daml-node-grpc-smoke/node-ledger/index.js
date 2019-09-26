const ledger = require('@digitalasset/daml-ledger');

let [, , host, port] = process.argv;
host = host || "localhost";
port = port || 6865;

const getSandboxClient = async (host, port) => {
    try{
        const client = await ledger.DamlLedgerClient.connect({host,port});
        return client;
    }catch(err){
        console.log(`Error: ${err}`);
        process.exit(1);
    }
}

const getPackages = async (client) => {
    try{
        const packages = client.packageClient.listPackages();
        return packages;
    }catch(err){
        console.log(`Error: ${err}`);
        process.exit(1);
    }
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

async function main(){
    const client = await getSandboxClient(host,port);
    console.log(`Sandbox id: ${client.ledgerId}\n`);
    const packages = await getPackages(client);
    console.log(`Packages : ${JSON.stringify(packages)}\n`);
    const payloads = await extractPayloads(client, packages);
    payloads.forEach(payload =>{
        payload.then(a => console.log(a));
    })
}

main()