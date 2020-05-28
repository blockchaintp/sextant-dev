## Fixtures for taekion-api-mock

First we use Docker to build and run the mock api:

```bash
cd taekion-api-mock
docker build -t taekion-api-mock .
docker run -d -p 8000:8000 taekion-api-mock
```

Then we hit some endpoints and get the JSON back:

#### list volumes

```bash
curl http://localhost:8000/volume?list
```

```json
{
  "action": "volume",
  "object": "volume",
  "payload": {
    "volumes": {
      "exampleVol0": {
        "compression": "LZ4",
        "encryption": "AES-GCM",
        "fingerprint": "2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea"
      },
      "exampleVol1": {
        "compression": "none",
        "encryption": "none",
        "fingerprint": "none"
      }
    }
  }
}
```

#### create volume

Values for compression:

 * none
 * lz4

Values for encryption:

 * none
 * aes_gcm

```bash
curl http://localhost:8000/volume?create=apples&compression=none&encryption=none
```

```json
{
  "action": "volume",
  "object": "volume",
  "payload": {
    "compression": "none",
    "encryption": "none",
    "fingerprint": "",
    "name": "apples"
  }
}
```

#### list snapshots

```bash
curl http://localhost:8000/snapshot?list&volume=apples
```

```json
{
  "action": "snapshot",
  "object": "snapshot",
  "payload": {
    "Data": {
      "demoSnapshot": "02 Jan 06 15:04 MST",
      "testSnapshot": "2020-05-28 11:17:43.832029071 +0000 UTC m=+1226.288307827",
      "volume": "apples"
    }
  }
}
```

#### create snapshot

```bash
curl http://localhost:8000/snapshot?create=snapshot1&volume=apples
```

```json
{
  "action": "snapshot",
  "object": "snapshot",
  "payload": {
    "name": "snapshot1",
    "volume": "apples"
  }
}
```