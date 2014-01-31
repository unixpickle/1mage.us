fs = require 'fs'
path = require 'path'

multiparty = require 'multiparty'

db = require './db'
args = require './args'


###
This is the express upload handler
###
module.exports = (req, res) ->
  form = new multiparty.Form()
  form.parse req, (err, fields, _files) ->
    # _files = {key: [file1, file2], key2: [...], ...}
    fileLists = x for _, x of _files
    files = []
    files = files.concat x for x in fileLists
    
    try
      info = getUploadedInfo files
      completeUpload info, res
    catch e
      fs.unlink x.path for x in files
      res.json error: e.toString()

###
Returns an object with a `mime` and `path` key or
throws an exception
###
getUploadedInfo = (files) ->
  # make sure they uploaded the right file count
  throw 'no files uploaded' if files.length is 0
  
  # if they uploaded too many files, we need to delete them
  throw 'too many files uploaded' if files.length > 1
  
  # get the mime type
  file = files[0]
  mime = getMimeType path.extname file.originalFilename
  throw 'invalid extension' if not mime?
  return mime: mime, path: file.path

###
Gets the MIME type for an extension or returns null.
###
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

completeUpload = (info, res) ->
  # rename tmp. file and then insert into DB
  seq = db.grabSequence()
  newPath = path.join args.directory, seq.toString()
  fs.rename info.path, newPath, (err) ->
    if err?
      fs.unlink info.path
      return res.json error: err.toString()
    db.insert info.mime, seq, (err) ->
      if err?
        fs.unlink newPath
        return res.json error: err.toString()
      res.json identifier: seq
