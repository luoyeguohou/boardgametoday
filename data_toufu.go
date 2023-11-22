package main

type ToufuCfg struct {
	id      int
	camp    int
	winCamp int
	winGold int
}

var ToufuCfgs map[int]ToufuCfg = map[int]ToufuCfg{
	1: {id: 1, camp: 1, winCamp: 1, winGold: 2},
	2: {id: 2, camp: 1, winCamp: 1, winGold: 2},
	3: {id: 3, camp: 3, winCamp: 1, winGold: 1},
	4: {id: 4, camp: 2, winCamp: 2, winGold: 1},
	5: {id: 5, camp: 3, winCamp: 2, winGold: 1},
	6: {id: 6, camp: 3, winCamp: 2, winGold: 1},
	7: {id: 7, camp: 3, winCamp: 3, winGold: 1},
	8: {id: 8, camp: 3, winCamp: 3, winGold: 1},
}
