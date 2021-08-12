const express = require('express')
const axios = require('axios')

const app = express()

app.use(async function(req, res) {

  console.log('--------------------------------------------')
  console.log('--------------------------------------------')
  console.log('--------------------------------------------')
  console.log(req.method)
  console.log(req.url)
  try {
    const result = await axios({
      method: req.method,
      headers: req.headers,
      url: `http://play.min.io${req.url}`,
    })
  
    
    console.dir(req.headers)
    console.dir(result.headers)
    console.log(result.status)
    console.log(result.data)
  
    res.status(result.status)
    res.set(result.headers)
    res.end(result.data)
  } catch(e) {
    if(e.response) {

      console.log('error')
      console.log(e.response.status)
      res.status(e.response.status)
      res.set(e.response.headers)
      res.end(e.response.data)
    }
    else {
      res.status(500)
      res.end(e.toString())
    }
  }
  
});
 
console.log("listening on port 5050")
app.listen(5050)