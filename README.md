# Deepin Sync

Deepin Sync is the sync backend of user person setting and data.

Sync will compare local file and server file version, then decide whow to sync.

if a client get a new setting to server, first create a .lock on server, so another client on can read setting.

how ever, all setting save as file system.

Default:

/deepin/setting/system/

/deepin/setting/theme/

type Sync interface {
    //Sync will call force sync
    Sync(id string)
}


#Sync Flow

1 Check Lock
  Lock |---> Wait ...
  NoLock |---> Push Lock
    |---> Recheck Lock
        |--->Wait

  CallBackup
     Push File

 Get Newest Version
    Check New 
        Decide Put
