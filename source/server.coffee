if process.argv.length isnt 4
  console.log 'Usage: node server.js <port> <directory>'
  process.exit 1

if isNaN port = parseInt process.argv[2]
  console.log 'invalid port: ' + process.argv[2]
  process.exit 1
imageDirectory = process.argv[3]

express = require 'express'
{MongoClient} = require 'mongodb'
http = require 'http'
multiparty = require 'multiparty'
fs = require 'fs'
path = require 'path'

connection = null
sequence = 0

main = ->
  app = express()
  app.use express.urlencoded()
  app.use express.json()
  app.use express.static __dirname + '/assets'
  app.post '/upload', handleUpload
  app.get '/nav/*', (req, res) ->
    return res.sendfile __dirname + '/assets/index.html'
  app.get '/error', (req, res) ->
    return res.sendfile __dirname + '/assets/index.html'
  
  app.param 'img', (req, res, next, id) ->
    return next() unless /^[0-9]+$/.exec(id)?
    query = sequence: parseInt id
    connection.collection('images').find(query).toArray (err, docs) ->
      return res.send 'internal error' if err?
      return res.redirect '/error' if docs.length isnt 1
      res.contentType docs[0].mime
      path = path.join imageDirectory, docs[0].sequence.toString()
      res.sendfile path
  app.get '/:img', (req, res) ->
    res.redirect '/error'
  app.get '*', (req, res) -> res.redirect '/error'

  server = http.createServer app
  server.listen port

handleUpload = (req, res) ->
  form = new multiparty.Form()
  res.sendJSON = (json) ->
    data = JSON.stringify json
    res.writeHead 200,
      'Content-Type': 'application/json'
      'Content-Length': data.length
    res.end data
  form.parse req, (err, fields, _files) ->
    __files = (x for _, x of _files)
    files = []
    files = files.concat x for x in __files
    # make sure they uploaded the right file count
    if files.length is 0
      return res.sendJSON error: 'no files uploaded'
    if files.length > 1
      for file in files
        fs.unlink file.path
      return res.sendJSON error: 'too many files uploaded'
    # make sure it's a supported file type
    file = files[0]
    mime = getMimeType path.extname file.originalFilename
    if not mime?
      return res.sendJSON error: 'invalid extension'
    uploadImageFromFile mime, file.path, (err, identifier) ->
      return res.sendJSON error: err.toString() if err?
      res.sendJSON identifier: identifier

uploadImageFromFile = (mime, tempPath, callback) ->
  collection = connection.collection 'images'
  doc = mime: mime, sequence: sequence++
  
  ident = doc.sequence.toString()
  fs.rename tempPath, path.join imageDirectory, ident
  
  # insert the record
  collection.insert doc, (err, docs) ->
    return callback err if err?
    return callback null, ident

getMimeType = (extension) ->
  list =
    png: 'image/png'
    jpeg: 'image/jpeg'
    jpg: 'image/jpeg'
    svg: 'image/svg+xml'
    gif: 'image/gif'
    bmp: 'image/bmp'
    ico: 'image/x-icon'
    jfif: 'image/jfif'
    tiff: 'image/tiff'
  return list[extension[1..].toLowerCase()]

MongoClient.connect 'mongodb://127.0.0.1:27017/1mage', (err, db) ->
  return console.log err if err?
  connection = db
  collection = db.collection 'images'
  collection.find({}).sort({sequence: -1}).nextObject (err, doc) ->
    return console.log err if err?
    if doc?
      sequence = doc.sequence + 1
    else sequence = 0
    console.log 'starting with sequence:', sequence
    main()

