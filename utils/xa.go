package utils

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type ACTION int
const (
	_ ACTION = iota
	Init
	Commit
	Rollback
	RetryCommit
	RetryRollback
	Done
)

type XAParticipant struct {
	O []orm.Ormer
	Action []ACTION
	Prepare func() error
}

type XACoordinator struct {
	Action ACTION
	MaxRetry int
	Participants []*XAParticipant
}

func (c *XACoordinator) AddParticipant(p *XAParticipant) {
	c.Participants = append(c.Participants,p)
}

func (c *XACoordinator) PrepareTransaction() (err error) {
	c.Action = Rollback
	for _,p := range c.Participants {
		for i:=0; i<len(p.O);i++ {
			if err = p.O[i].Begin(); err != nil {
				return
			}
			if err = p.Prepare(); err != nil {
				p.Action[i] = Rollback
				return
			}
			p.Action[i] = Commit
		}
	}
	c.Action = Commit
	return
}

func (c *XACoordinator) FinishTransaction() (err error) {
	for i:=0;i<c.MaxRetry;i++{
		failCnt := 0
		if Commit == c.Action {
			for _,p := range c.Participants {
				for i:=0; i<len(p.O); i++ {
					if Commit == p.Action[i] || RetryCommit == p.Action[i] {
						if err = p.O[i].Commit(); err != nil {
							p.Action[i] = RetryCommit
							failCnt++
						}else {
							p.Action[i] = Done
						}
					}
				}
			}
		} else if Rollback == c.Action {
			for _,p := range c.Participants {
				for i:=0;i<len(p.O);i++ {
					if Rollback == p.Action[i] || RetryRollback == p.Action[i] {
						if err = p.O[i].Rollback(); err != nil {
							p.Action[i] = RetryRollback
							failCnt++
						} else {
							p.Action[i] = Done
						}
					}
				}
			}
		}
		if 0 == failCnt {
			return
		}
	}
	logs.Error(err.Error())
	err = errors.New("Retry max times still failed!")
	return
}