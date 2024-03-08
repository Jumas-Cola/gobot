let error = true;

let res = [
  db.ChatVotes.drop(),
  db.ChatVotes.insertOne({}),
  db.ChatVotes.deleteOne({}),
];

printjson(res);
