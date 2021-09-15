## taekion-writer

 * `api/app_api_update.go`: `TFSUpdate` is what we add `mutations` to
 * `fs/handler/types_mutation.go` - the mutation types
 * `fuse/dir.go`: `140` - the file create fn - it makes a bundle
 * `fuse/file.go` - the read / write handlers
 * `fuse/main.go` - booting up with a journal
 * `api/app_api_sync.go` - the sync api - this writes the journal to sawtooth
 * `journal/sqllite/api_sync.go` - the journal implementation
 * `journal/sqllite/api_get.go` - get previous transaction
 * `client/api_update_mutation.go` - build mutation bundles from updates

basic flow of a write:

 * create the inode in an update and submit the bundle - get the bundle hash (so we can make more bundles)
 * create further updates for data blocks
 * loop over all the journal entry ids we get from the updates and block on them being confirmed
 * keep a queue of incoming writes to the same key
 * once all jounrnal entries are confirmed, return the current request and pop the next one