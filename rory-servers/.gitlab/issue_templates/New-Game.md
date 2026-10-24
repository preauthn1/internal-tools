* [ ]  Create `GameID` / `ServerID` 
  > `PrSys_CreateNewGame` 
  > * GameID: `` 
  > * ServerID: 
  >   * [ ] DEV `` ~ ``
  >   * [ ] QA `` ~ ``
  >   * [ ] Stage `` ~ ``
* [ ]  Determine `Port` Number
  > `` ~ ``
* [ ]  Game Logic Develop
* [ ]  Add Build Scripts
* [ ]  Add Platform Area on CMS & Database
``` SQL
INSERT INTO ServerInfo.ServerPlatformRelations (PlatformID, GameID, ServerID, Status)
SELECT {PlatformID}, GameID, ServerID, 1 FROM ServerInfo.GameServerList WHERE GameID = {GameID}
```
* [ ]  Environments Config Setting / Hall Server Setting
  * [ ] DEV
  * [ ] QA
  * [ ] Stage
* [ ]  Deploy
  * [ ] DEV
  * [ ] QA
  * [ ] Stage
* [ ]  Firewall
  * [ ] DEV
  * [ ] QA
  * [ ] Stage
* [ ]  TF Server Config
  * [ ] Stage