# Tempshelf

Tempshelf is archive file assets, upload specified s3 backet, and restore.


## usage

compose
```sh
./tempshelf compose <manifest path>
```

restore
```sh
./tempshelf pull <manifest path>
```


## manifest file
```json
{
    "meta" : {
        "storage" : "s3",
        "bucket" : "(bucket-name)",
        "region" : "(bucket-region)",
        "token" : "(IAM-token)",
        "secret" : "(IAM-secret)",
        "prefix" : "(prefix of object key)"
    },
    "files" : [
        "NOTE: this section auto generated"
    ]
}
```

## files structure
- assets
    - manifest.json
    - asset1.jpeg
    - asset2.png
    - asset3/
        - asset4.txt
        - asset5.zip
    - .tmp/
