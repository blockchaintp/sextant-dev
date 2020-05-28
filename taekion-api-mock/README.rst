TFS API Endpoint Mockup Interface
---------------------------------

``GET /debug``

    ``/address``
        Displays the blockchain address for a structure.

        Query Parameters
            object
                String. Values accepted:

                ``datablock`` - Dumps the contents of a data block.

                ``directory`` - Returns a list of files in target directory.

                ``inode`` - Returns inode internal data.

                ``volume`` - Returns volume metadata.

            id
                String. Varies by object:

                **datablock** - Hash of the data block to query.

                Example Input:  ``2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea``

                **directory** - Directory UUID.

                Example Input: ``11111111-1111-1111-1111-111111111111``

                **inode** - inode UUID.

                Example Input: ``11111111-1111-1111-1111-111111111111``

                **volume** - Name of the volume.

                Example Input: ``testVolume``

        Usage
            |   ``http://localhost:8000/debug/address?object=datablock&id=2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea``
            |   ``http://localhost:8000/debug/address?object=directory&id=11111111-1111-1111-1111-111111111111``
            |   ``http://localhost:8000/debug/address?object=volume&id=testvolume``

    ``/datablock``
        Datablock returns the contents of a data block stored in TFS.

        Query Parameters
            hash
                String. Hash of the data block to query.

                Example Input: ``2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea``

            raw
                String. Returns raw bytes.

                Example Input: ``""``

        Usage
            ``http://localhost:8000/debug/datablock?hash=2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea``

    ``/directory``
        Returns the file names and UUIDs contained in a TFS directory.

        Query Parameters
            UUID
                String. Directory UUID.

                Example Input: ``11111111-1111-1111-1111-111111111111``

        Usage
            ``http://localhost:8000/debug/directory?id=11111111-1111-1111-1111-111111111111``

    ``/inode``
        Returns a TFS cache inode's metadata.

        Query Parameters
            id
                String. Inode UUID.

                Example Input: ``11111111-1111-1111-1111-111111111111``

        Usage
            ``http://localhost:8000/debug/inode?id=11111111-1111-1111-1111-111111111111``

    ``/volume``
        Displays metadata for a volume.

        Query Parameters
            name
                String. Name assigned to volume.

                Example Input: ``testVolume``

        Usage
            ``http://localhost:8000/debug/debugvolume?id=testvolume``

    ``/wrapper``
        Displays the object wrapper data for a structure.

        Query Parameters
            object
                String. Values accepted:

                ``datablock`` - Dumps the contents of a data block.

                ``directory`` - Returns a list of files in target directory.

                ``inode`` - Returns inode internal data.

            id
                String. Varies by object_type:

                **datablock** - Hash of the data block to query.

                Example Input: ``2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea``

                **directory** - Directory UUID.

                Example Input: ``11111111-1111-1111-1111-111111111111``

                **inode** - inode UUID.

                Example Input: ``11111111-1111-1111-1111-111111111111``

            dump
                String. Hex dump of wrapper data.

        Usage
            |   ``http://localhost:8000/debug/wrapper?object=datablock&id=2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea``
            |   ``http://localhost:8000/debug/wrapper?object=inode&id=11111111-1111-1111-1111-111111111111``

``GET /snapshot``
    Commands for manipulating TFS snapshots for a volume.

    Query Parameters
        create
            String. Name of a new snapshot on target volume.

            Example Input: ``testSnapshot``

        volume
            String. Name target volume.

            Example Input: ``testVolume``

        list
            String. List snapshots present on target volume.

            Example Input: ``""``

    Usage
        |   ``tfs-cli volume snapshot create <volume name> <snapshot name>``
        |   ``tfs-cli volume snapshot list <volume name>``

``GET /volume``
    Commands for manipulating and querying TCDP Volumes.

    Query Parameters
        create
            String. Name assigned to new volume.

            Example Input: ``testVolume``

        encryption
            String. Encryption type to use on new volume: AES_GCM/None.

            Example Input: ``AES_GCM``

        key
            String. Fingerprint of encryption key.

            Example Input: ``2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea``

        compression
            String. Compression type to use on new volume: LZ4/None.

            Example Input: ``LZ4``

        list
            String. List existing TFS volumes.

            Example Input: ``""``

        Usage
            | http://localhost:8000/volume?list=
            | http://localhost:8000/volume?create=testvol&encryption=aes_gcm&fingerprint=2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea&compression=none


Json Response
^^^^^^^^^^^^^

JSON responses share a common format:

``action`` - API action requested.

``object`` - Data structure of payload.

``payload`` - Data for action/object.

Example output::

    {
        "action": "snapshot",
        "object": "snapshot",
        "payload": {
            "Name": "example_Snapshot",
            "Volume": "testVolume"
        }
    }

Contents of ``payload`` will vary based on the request action and object. As an
example: If a request is sent to ``/GET /debug/address/``::

    {
        "action": "debug/address",
        "object": "volume",
        "payload": {
            "address": "26c66900ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0d",
            "id": "72adb771-40b1-4a1f-9472-8e684b691f99",
            "name": "test"
        }
    }

If a request is sent to ``/GET /debug/volume``::

    {
        "action": "debug/volume",
        "object": "volume",
        "payload": {
            "Compression Type": "LZ4",
            "Encryption Type": "AES_GCM",
            "Key Fingerprint": "9f86d081884c7d659a2feaa0c55ad015",
            "Last Hash": "26c66900ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0d",
            "Root Inode UUID": "1a03a5e5-648c-442d-9ac8-3c191e4627d2",
            "Volume Name": "test"
        }
    }

