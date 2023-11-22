package main

import "time"

type Service interface {
	OnPlyInit()
	OnPlyExit(ply *Player)
	OnInit()
	OnSec()
}

type ServiceMananger struct {
	Services     []Service
	ServiceRoom  ServiceRoom
	ServiceToufu ServiceToufu
}

var ServiceMGR ServiceMananger

func (s *ServiceMananger) InitAllService() {
	s.Services = append(s.Services, &s.ServiceRoom)
	s.Services = append(s.Services, &s.ServiceToufu)

	for _, service := range s.Services {
		service.OnInit()
	}

	TimerRoutine(time.Second, func() { s.OnSec() }, nil)
}

func (s *ServiceMananger) OnPlyExit(ply *Player) {
	for _, service := range s.Services {
		service.OnPlyExit(ply)
	}
}

func (s *ServiceMananger) OnSec() {
	for _, service := range s.Services {
		service.OnSec()
	}
}
