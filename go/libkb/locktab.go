// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package libkb

import (
	"sync"
)

type NamedLock struct {
	sync.Mutex
	refs   int
	name   string
	parent *LockTable
}

func (l *NamedLock) incref() {
	l.refs++
}

func (l *NamedLock) decref() {
	l.refs--
}

func (l *NamedLock) Release(m MetaContext) {
	l.Unlock()
	l.parent.Lock()
	l.decref()
	if l.refs == 0 {
		delete(l.parent.locks, l.name)
	}
	l.parent.Unlock()
}

type LockTable struct {
	sync.Mutex
	locks map[string]*NamedLock
}

func (t *LockTable) init() {
	if t.locks == nil {
		t.locks = make(map[string]*NamedLock)
	}
}

func (t *LockTable) AcquireOnName(m MetaContext, s string) (ret *NamedLock) {
	t.Lock()
	t.init()
	if ret = t.locks[s]; ret == nil {
		ret = &NamedLock{refs: 0, name: s, parent: t}
		t.locks[s] = ret
	}
	ret.incref()
	t.Unlock()
	ret.Lock()
	return ret
}
