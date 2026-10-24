

enum NotifyServerClassID {
	ServerRegID=1,
	AnnID=2,
	AttrChangeID=3,
	ServerStopID=4,
	RefreshVersionID=5,
	ServerStopWhiteIPListID=6,
    ServerHostAndBackUp = 7,
    AnnFishRegID = 8,               //]Բ~
    AnnFishID = 9,                  //~
    PayOrderID = 10,
    UpdateHallUserListID = 11,      //Ҫ@ȡÑϢUserID
    UserNotifyRetN2HID = 12,        //l֪ͨϢ Notify->Hall
    UserNotifyRetH2CID = 13,        //l֪ͨϢ Hall->Client
    UpdateUserNotifyIdxH2NID = 15,  //֪̖ͨ Hall->Notify
};