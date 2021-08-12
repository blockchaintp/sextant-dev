const Minio = require('minio')

const run = async () => {
  const minioClient = new Minio.Client({
    endPoint: 'localhost',
    port: 8001,
    useSSL: false,
    accessKey: '',
    secretKey: ''
  })

  try {
    const stream = await minioClient.listObjectsV2('mybucket', '/hello/world', true, '')
    stream.on('data', function(obj) { console.log(obj) } )
    stream.on('error', function(err) { console.log(err) } )
    // console.log('--------------------------------------------')
    // console.dir(result)
  } catch(e) {
    console.error(e.toString())
    process.exit(1)
  }
  
}

run()