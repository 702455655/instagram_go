package goinsta

import (
	"container/list"
	"makemoney/common/log"
	"sync"
	"time"
)

type ListWrap struct {
	list *list.List
	cur  *list.Element
}

type AccountPoolt struct {
	Accounts     []*Instagram
	notAvailable *list.List
	Available    map[string]*ListWrap
	Cooling      *list.List
	avalLock     sync.Mutex
	noAvalLock   sync.Mutex
	coolingLock  sync.Mutex
	checkTimer   *time.Ticker
}

var AccountPool *AccountPoolt

var CallBackCheckAccount func(inst *Instagram) bool = nil

func InitAccountPool(accounts []*Instagram) {
	AccountPool = &AccountPoolt{}
	AccountPool.Cooling = list.New()
	AccountPool.Available = make(map[string]*ListWrap)
	AccountPool.notAvailable = list.New()
	AccountPool.Accounts = make([]*Instagram, len(accounts))
	var accountsIndex = 0
	for idx := range accounts {
		if accounts[idx].IsLogin && accounts[idx].Status == "" {
			Available := AccountPool.Available[accounts[idx].Tags]
			if Available == nil {
				Available = &ListWrap{
					list: list.New(),
					cur:  nil,
				}
				AccountPool.Available[accounts[idx].Tags] = Available
			}
			Available.list.PushBack(accounts[idx])
			AccountPool.Accounts[accountsIndex] = accounts[idx]
			accountsIndex++
		} else {
			AccountPool.notAvailable.PushBack(accounts[idx])
		}
	}
	AccountPool.Accounts = AccountPool.Accounts[:accountsIndex]

	if CallBackCheckAccount != nil {
		AccountPool.checkTimer = time.NewTicker(time.Second * 10)
		go CheckAccount()
	}

	log.Info("init account pool available count: %d, bad account :%d",
		accountsIndex, AccountPool.notAvailable.Len())

	go AccountPool.checkBad()
}

func CheckAccount() {
	for range AccountPool.checkTimer.C {
		AccountPool.coolingLock.Lock()
		for item := AccountPool.Cooling.Front(); item != nil; item = item.Next() {
			inst := item.Value.(*Instagram)
			if CallBackCheckAccount(inst) {
				AccountPool.ReleaseOne(inst)
				next := item.Next()
				AccountPool.Cooling.Remove(item)
				item = next
				log.Info("check account: %s cooling finish", inst.User)
			}
		}
		AccountPool.coolingLock.Unlock()
	}
}

func (this *AccountPoolt) checkBad() {
	var check = func(inst *Instagram) (retErr error) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("account: %s doDevelopMeta panic error: %v", inst.User, err)
				retErr = err.(error)
			}
		}()
		err := inst.GetAccount().Sync()
		if err != nil {
			log.Warn("account pool check %s %v", inst.User, err)
		}
		return err
	}
	for true {
		ava := make([]*Instagram, 0)
		this.noAvalLock.Lock()
		for item := this.notAvailable.Front(); item != nil; item = item.Next() {
			inst := item.Value.(*Instagram)
			check(inst)
			if !inst.IsBad() {
				ava = append(ava, inst)
				next := item.Next()
				this.notAvailable.Remove(item)
				item = next
			}
		}
		this.noAvalLock.Unlock()
		if len(ava) > 0 {
			for _, item := range ava {
				this.ReleaseOne(item)
			}
		}
		time.Sleep(30 * time.Minute)
	}
}

func (this *AccountPoolt) GetOneBlock(OperName string, AccountTag string) *Instagram {
	for true {
		log.Warn("try require account tag: %s, for %s", AccountTag, OperName)
		inst := this.GetOneNoWait(OperName, AccountTag)
		if inst != nil {
			return inst
		}
		time.Sleep(time.Second * 10)
	}
	return nil
}

func (this *AccountPoolt) GetOneNoWait(OperName string, AccountTag string) *Instagram {
	this.avalLock.Lock()
	defer this.avalLock.Unlock()

	Available := this.Available[AccountTag]
	if Available == nil {
		log.Error("not find account tag: %s", AccountTag)
		return nil
	}

	if Available.cur == nil {
		Available.cur = Available.list.Front()
		if Available.cur == nil {
			return nil
		}
	}
	var oldCurl = Available.cur
	for true {
		inst := Available.cur.Value.(*Instagram)
		lastCur := Available.cur
		Available.cur = Available.cur.Next()

		if !inst.IsSpeedLimit(OperName) {
			Available.list.Remove(lastCur)
			return inst
		}
		if Available.cur == nil {
			Available.cur = Available.list.Front()
			if Available.cur == nil {
				return nil
			}
		}
		if Available.cur == oldCurl {
			break
		}
	}
	return nil
}

func (this *AccountPoolt) ReleaseOne(insta *Instagram) {

	if insta.IsBad() {
		this.noAvalLock.Lock()
		log.Error("add black account: %s ,status: %s", insta.User, insta.Status)
		this.notAvailable.PushBack(insta)
		this.noAvalLock.Unlock()
	} else {
		this.avalLock.Lock()
		Available := this.Available[insta.Tags]
		if Available == nil {
			log.Error("not find account tag: %s", insta.Tags)
		} else {
			Available.list.PushBack(insta)
			log.Info("add available account: %s", insta.User)
		}
		this.avalLock.Unlock()
	}
	SaveInstToDB(insta)
}

func (this *AccountPoolt) CoolingOne(inst *Instagram) {
	this.coolingLock.Lock()
	defer this.coolingLock.Unlock()

	this.Cooling.PushBack(inst)
}

func (this *AccountPoolt) BlackOne(insta *Instagram) {
	this.noAvalLock.Lock()
	defer this.noAvalLock.Unlock()

	this.notAvailable.PushBack(insta)
	SaveInstToDB(insta)
}
