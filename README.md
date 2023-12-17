## go telegram antispam

Bot for checking for spam in telegram groups

## Make hamspam database

[Gonzofilter Getting Started](https://github.com/gsauthof/gonzofilter#getting-started)

## Usage

1. Put hamspam.db in go-telegram-antispam directory

2. Run bot by docker compose

```bash
echo "TOKEN=<YOUR_TELEGRAM_BOT_TOKEN>" > .env
docker compose up -d
```

## Spam (RU) database for training classifier

[dbhub.io/Jumas-Cola/spam.db](https://dbhub.io/Jumas-Cola/spam.db)

### Make files for training from sqlite spam.db

```python
import sqlite3
import uuid
import os

con = sqlite3.connect("spam.db")
cur = con.cursor()
sel = cur.execute("SELECT * FROM messages")

learn_spam_dir = 'ex/learn/spam/'
test_spam_dir = 'ex/test/spam/'
learn_ham_dir = 'ex/learn/ham/'
test_ham_dir = 'ex/test/ham/'

for d in [learn_spam_dir, learn_ham_dir, test_spam_dir, test_ham_dir]:
    if not os.path.exists(d):
        os.makedirs(d)

# 80% for training, 20% for testing
learn_test_ratio = .8

rows = sel.fetchall()
for i, row in enumerate(rows):
    path = test_spam_dir if i > (len(rows) * learn_test_ratio) \
                         else learn_spam_dir
    with open(path + uuid.uuid4().urn[9:], 'w') as f:
        f.write(row[0])
```

Ham messages can be downloaded from telegram chat history.
