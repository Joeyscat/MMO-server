package model

import (
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer"
	"github.com/NumberMan1/common/summer/timeunit"
	"sync"
)

var (
	singleCharacterManager = singleton.Singleton{}
)

// CharacterManager 统一管理全部的角色（创建，移除，获取）
type CharacterManager struct {
	//游戏里全部的角色 <ChrId,ChrObj>
	characters *sync.Map
}

func GetCharacterManagerInstance() *CharacterManager {
	result, _ := singleton.GetOrDo[*CharacterManager](&singleCharacterManager, func() (*CharacterManager, error) {
		return &CharacterManager{
			characters: &sync.Map{},
		}, nil
	})
	//每隔5秒保存Data到数据库
	summer.GetScheduleInstance().AddTask(result.save, timeunit.Seconds, 5, 0)
	return result
}

func (cm *CharacterManager) CreateCharacter(dbChr *database.DbCharacter) *Character {
	character := NewCharacter(dbChr)
	cm.characters.Store(character.Id(), character)
	GetEntityManagerInstance().AddEntity(dbChr.SpaceId, character)
	return character
}

func (cm *CharacterManager) RemoveCharacter(chrId int) {
	character, ok := cm.characters.Load(chrId)
	if ok {
		cm.characters.Delete(chrId)
		chr := character.(*Character)
		GetEntityManagerInstance().RemoveEntity(chr.Data.SpaceId, chr)
	}
}

func (cm *CharacterManager) GetCharacter(chrId int) *Character {
	c, ok := cm.characters.Load(chrId)
	if ok {
		return c.(*Character)
	} else {
		return nil
	}
}

func (cm *CharacterManager) Clear() {
	cm.characters = &sync.Map{}
}

func (cm *CharacterManager) save() {
	cm.characters.Range(func(key, value any) bool {
		chr := value.(*Character)
		chr.Data.X = int(chr.Position().X)
		chr.Data.Y = int(chr.Position().Y)
		chr.Data.Z = int(chr.Position().Z)
		database.OrmDb.Save(chr.Data)
		return true
	})
}
