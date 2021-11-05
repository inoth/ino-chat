package test

import (
	"ino-chat/cache"
	svc "ino-chat/src/service"
	"os"
	"testing"
	"time"
)

func TestSaveUser(t *testing.T) {
	var (
		uid = "test123"
	)
	cache.Init()
	t1 := time.Now()

	os.Setenv("GORUNEVN", "dev")
	res := svc.SaveUserInCache(uid)
	if !res {
		t.Error("存入失败")
	}
	t.Logf("ok; time: %v", time.Since(t1))
}

func TestNewRoom(t *testing.T) {
	var (
		rid   = "testroom"
		onwer = "test123"
		title = "testroom"
	)
	t1 := time.Now()
	cache.Init()

	os.Setenv("GORUNEVN", "dev")
	res := svc.NewRoom(rid, onwer, title)
	if !res {
		t.Error("存入失败")
	}
	t.Logf("ok; time: %v", time.Since(t1))
}

func TestGetRooms(t *testing.T) {
	t1 := time.Now()
	os.Setenv("GORUNEVN", "dev")
	cache.Init()

	rooms, err := svc.GetAllRoom()
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("ok; time: %v", time.Since(t1))
	t.Logf("%v", rooms)
}

func TestJoinRoom(t *testing.T) {
	var (
		rid = "testroom"
		uid = "inoth111"
	)
	t1 := time.Now()
	os.Setenv("GORUNEVN", "dev")
	cache.Init()
	if !svc.JoinRoom(rid, uid) {
		t.Error("加入失败")
	}

	t.Logf("ok; time: %v", time.Since(t1))
}

func TestExitRoom(t *testing.T) {
	var (
		rid = "testroom"
		uid = "inoth111"
	)
	t1 := time.Now()
	os.Setenv("GORUNEVN", "dev")
	cache.Init()
	if !svc.ExitRoom(rid, uid) {
		t.Error("退出失败")
	}

	t.Logf("ok; time: %v", time.Since(t1))
}

func TestGetUsersByRoom(t *testing.T) {
	t1 := time.Now()
	os.Setenv("GORUNEVN", "dev")
	cache.Init()
	users := svc.GetUsersByRoom("testroom")
	if len(users) <= 0 {
		t.Error("用户为空")
	}
	t.Logf("ok; time: %v", time.Since(t1))
	t.Logf("%v", users)
}
