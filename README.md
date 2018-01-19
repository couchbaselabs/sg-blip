CLI tool to connect to Sync Gateway over the BLIP protocol.

## Example

```
$ ./sg-blip subChanges http://localhost:4985/db
subChanges called.... args [http://localhost:4985/db]
2018/01/18 17:00:01 Got change: [[3,"TestMigrateMetadataRequiresImport","1-9f1dd0022bd0ffdf7c7f158ef08b6211"],[4,"TestMigrateMetadataRequiresImport2","1-9f1dd0022bd0ffdf7c7f158ef08b6211"],[5,"test","1-f97ffb79945badf2fc8f7708ddbf6667"],[6,"test2","1-f97ffb79945badf2fc8f7708ddbf6667"]]
2018/01/18 17:00:01 Got change: null
2018/01/18 17:00:12 Got change: [[7,"test3","1-f97ffb79945badf2fc8f7708ddbf6667"]]
^C
```