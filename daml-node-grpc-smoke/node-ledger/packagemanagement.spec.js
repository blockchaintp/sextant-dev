const test = require('tape');
const getSandboxClient = require('./common').getSandboxClient;
const uploadDar = require('./common').uploadDar;
const listPackages = require('./common').listPackages;

const host = process.env.ENDPOINT_URL | "localhost";
const port = process.env.ENDPOINT_PORT | 6865;

var client
test('get client and it should return an id', async t => {
    client = await getSandboxClient(host, port);
    t.equal(client.ledgerId.includes('sandbox'), true);
    t.end();
});

test('Upload Dar', async t => {
    const beforeUpload = await listPackages(client);
    console.log(`Before upload --> ${JSON.stringify(beforeUpload.packageDetailsList)}`);
    await uploadDar(client, '../dist/daml-node-ledger-api.dar');
    const afterUpload = await listPackages(client);
    console.log(`After upload ---> ${JSON.stringify(afterUpload.packageDetailsList)}`);
    t.end();
});