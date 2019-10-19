const test = require('tape');
const getSandboxClient = require('./common').getSandboxClient;
const listParties = require('./common').listParties;
const allocateParties = require('./common').allocateParties;

const host = process.env.ENDPOINT_URL | "localhost";
const port = process.env.ENDPOINT_PORT | 6865;

var client
test('get client id should return an id associated with the ledger', async t =>{
    client = await getSandboxClient(host,port);
    t.equal(client.ledgerId.includes('sandbox'), true);
    t.end();
});


const party = `paul-${Date.now().toString()}`;
const displayName = `Paul ${Date.now().toString()}`;
test('allocating new parties should not throw error', async t => {
    try {
        const response = await allocateParties(client, party, displayName);
        const expected = { partyDetails: { party: party, isLocal: true, displayName: displayName } };
        t.deepEqual(expected, response);
    } catch (err){
        t.fail(`This should mot happen ${err}`);
    }
    t.end();
});

test('allocating same parties twice it should throw an Error number', async t => {
    try {
        await allocateParties(client, party, displayName);
    } catch (err){
        t.ok(err.toString().substring(7,8)==='3');
    }
    t.end();
});

test('get list of parties should return an array of parties details', async t => {
   const parties = await listParties(client);
    t.ok(Array.isArray(parties.partyDetails)===true);
    t.end();
});