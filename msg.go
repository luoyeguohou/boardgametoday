package main

const (
	// 创建房间  cs sc
	Create_Room = 30001
	// 加入房间  cs sc
	Join_Room = 30002
	// 房间信息更改 sc
	Room_State_Change = 30003
	// 更换座位 cs
	Change_Site = 30004
	// 离开房间 cs
	Exis_Room = 30005

	// 更换游戏
	Change_Game = 30006
	// 开始游戏
	Start_Game  = 30007
	Change_Name = 30008
	Heart_Beat  = 30103

	// 选择玩家 cs
	Toufu_Choose_Player = 30101
	// 回合结束 sc
	Toufu_Turn_End = 30102
	// 游戏结束 sc
	Toufu_Game_End = 30103
)
