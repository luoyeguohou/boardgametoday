package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Room struct {
	Players [8]*Player
	RoomID  int
	Host    *Player
}

type ServiceRoom struct {
	Rooms  map[int]*Room
	Rand   *rand.Rand
	logCnt int
}

func (s *ServiceRoom) OnInit() {
	s.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	s.Rooms = make(map[int]*Room)
}

func (s *ServiceRoom) OnPlyInit() {

}

func (s *ServiceRoom) FindRoom(ply *Player) *Room {
	for _, room := range s.Rooms {
		if room == nil {
			continue
		}
		for _, rPly := range room.Players {
			if rPly == ply {
				return room
			}
		}
	}

	return nil
}

func (s *ServiceRoom) FindIdx(ply *Player) int {
	room := s.FindRoom(ply)
	for idx, ePly := range room.Players {
		if ePly == ply {
			return idx
		}
	}
	return -1
}

func (s *ServiceRoom) OnPlyExit(ply *Player) {
	room := s.FindRoom(ply)
	if room == nil {
		return
	}
	s.ExitRoom(room.RoomID, ply)
}

func (s *ServiceRoom) OnSec() {
	s.logCnt += 1
	if s.logCnt < 10 {
		return
	}
	s.logCnt = 0
	str := "start room log:"
	for _, room := range s.Rooms {
		if room == nil {
			continue
		}
		str += "[ room id: " + strconv.Itoa(room.RoomID) + " ply:"
		for _, player := range room.Players {
			if player == nil {
				continue
			}
			str += " " + player.NickName + " "
		}
		str += "]\n"
	}
	fmt.Print(str)
}

func (s *ServiceRoom) BroadcastRoom(roomID int, cmd int, param SImap) {
	for _, ply := range s.Rooms[roomID].Players {
		if ply != nil {
			ply.SendMsg(cmd, param)
		}
	}
}

func (s *ServiceRoom) RoomStateChanged(roomID int) int {
	if s.Rooms[roomID] == nil {
		return -1
	}

	for _, ply := range s.Rooms[roomID].Players {
		if ply != nil {
			ply.SendMsg(Room_State_Change, SImap{
				"state":   s.GetRoomState(roomID),
				"isHost":  s.Rooms[roomID].Host == ply,
				"roomID":  roomID,
				"selfIdx": s.FindIdx(ply),
			})
		}
	}
	return 0
}

func nextRoomID(id int) int {
	if id >= 9999 {
		return 1000
	}
	return id + 1
}

func (s *ServiceRoom) CreateRoom(ply *Player) int {
	var roomID = s.Rand.Intn(8000) + 1000
	for {
		if s.Rooms[roomID] == nil {
			break
		}
		roomID = nextRoomID(roomID)
	}

	s.Rooms[roomID] = &Room{}
	s.Rooms[roomID].Players[0] = ply
	s.Rooms[roomID].RoomID = roomID
	s.Rooms[roomID].Host = ply
	return roomID
}

func (s *ServiceRoom) GetRoomState(roomID int) SImap {
	if s.Rooms[roomID] == nil {
		return nil
	}

	var playerDatas []SImap
	for pos, ply := range s.Rooms[roomID].Players {
		if ply == nil {
			continue
		}
		playerData := ply.ToSImap()
		playerData["Pos"] = pos
		playerData["IsHost"] = s.Rooms[roomID].Host == ply
		playerDatas = append(playerDatas, playerData)
	}
	return SImap{
		"Player": playerDatas,
	}
}

func (s *ServiceRoom) isEmptyRoom(roomID int) bool {
	if s.Rooms[roomID] == nil {
		return true
	}
	var allEmpty = true
	for _, ply := range s.Rooms[roomID].Players {
		if ply != nil {
			allEmpty = false
			break
		}
	}
	return allEmpty
}

func (s *ServiceRoom) isFullRoom(roomID int) bool {
	if s.Rooms[roomID] == nil {
		return true
	}
	var allFull = true
	for _, ply := range s.Rooms[roomID].Players {
		if ply == nil {
			allFull = false
			break
		}
	}
	return allFull
}

func (s *ServiceRoom) ExitRoom(roomID int, ply *Player) int {
	// 无效房间
	if s.Rooms[roomID] == nil {
		return -1
	}

	for pos, ePly := range s.Rooms[roomID].Players {
		if ply == ePly {
			s.Rooms[roomID].Players[pos] = nil
			if s.Rooms[roomID].Host == ply {
				for _, ele := range s.Rooms[roomID].Players {
					if ele != nil {
						s.Rooms[roomID].Host = ele
						break
					}
				}
			}
			// 空房间释放
			if s.isEmptyRoom(roomID) {
				s.Rooms[roomID] = nil
				return 0
			}
			s.RoomStateChanged(roomID)
			return 0
		}
	}
	// 房间里没有该玩家
	return -2
}

func (s *ServiceRoom) ChangePos(roomID int, ply *Player, aimPos int) int {
	// 无效房间
	if s.Rooms[roomID] == nil {
		return -1
	}

	// 目标位置如果非空
	if s.Rooms[roomID].Players[aimPos] != nil {
		return -2
	}

	// 去掉之前的位置
	for idx, iPly := range s.Rooms[roomID].Players {
		if iPly == ply {
			s.Rooms[roomID].Players[idx] = nil
			break
		}
	}
	s.Rooms[roomID].Players[aimPos] = ply
	s.RoomStateChanged(roomID)
	return 0
}

func (s *ServiceRoom) JoinRoom(roomID int, ply *Player) int {
	// 无效房间
	if s.Rooms[roomID] == nil {
		return -1
	}
	// 满员
	if s.isFullRoom(roomID) {
		return -1
	}
	// 加入
	for idx, iPly := range s.Rooms[roomID].Players {
		if iPly == nil {
			s.Rooms[roomID].Players[idx] = ply
			break
		}
	}

	s.RoomStateChanged(roomID)

	return 0
}
