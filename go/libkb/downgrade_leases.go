// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package libkb

import (
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
	"strings"
)

func kidsToString(kids []keybase1.KID) string {
	var tmp []string
	for _, k := range kids {
		tmp = append(tmp, string(k))
	}
	return strings.Join(tmp, ",")
}

func sigIDsToString(sigIDs []keybase1.SigID) string {
	var tmp []string
	for _, k := range sigIDs {
		tmp = append(tmp, string(k))
	}
	return strings.Join(tmp, ",")
}

func UidsToString(uids []keybase1.UID) string {
	s := make([]string, len(uids))
	for i, uid := range uids {
		s[i] = string(uid)
	}
	return strings.Join(s, ",")
}

type Lease struct {
	MerkleSeqno keybase1.Seqno    `json:"merkle_seqno"`
	LeaseID     keybase1.LeaseID  `json:"downgrade_lease_id"`
	HashMeta    keybase1.HashMeta `json:"hash_meta"`
}

type leaseReply struct {
	Lease
	Status AppStatus `json:"status"`
}

func (r *leaseReply) GetAppStatus() *AppStatus {
	return &r.Status
}

func RequestDowngradeLeaseByKID(m MetaContext, kids []keybase1.KID) (lease *Lease, mr *MerkleRoot, err error) {
	var res leaseReply
	err = m.G().API.PostDecode(m, APIArg{
		Endpoint:    "downgrade/key",
		SessionType: APISessionTypeREQUIRED,
		Args: HTTPArgs{
			"kids": S{kidsToString(kids)},
		},
	}, &res)
	if err != nil {
		return nil, nil, err
	}
	return leaseWithMerkleRoot(m, res)
}

func CancelDowngradeLease(m MetaContext, l keybase1.LeaseID) error {
	_, err := m.G().API.Post(m, APIArg{
		Endpoint:    "downgrade/cancel",
		SessionType: APISessionTypeREQUIRED,
		Args: HTTPArgs{
			"downgrade_lease_id": S{string(l)},
		},
	})
	return err
}

func RequestDowngradeLeaseBySigIDs(m MetaContext, sigIDs []keybase1.SigID) (lease *Lease, mr *MerkleRoot, err error) {
	var res leaseReply
	err = m.G().API.PostDecode(m, APIArg{
		Endpoint:    "downgrade/sig",
		SessionType: APISessionTypeREQUIRED,
		Args: HTTPArgs{
			"sig_ids": S{sigIDsToString(sigIDs)},
		},
	}, &res)
	if err != nil {
		return nil, nil, err
	}
	return leaseWithMerkleRoot(m, res)
}

func RequestDowngradeLeaseByTeam(m MetaContext, teamID keybase1.TeamID, uids []keybase1.UID) (lease *Lease, mr *MerkleRoot, err error) {
	var res leaseReply
	err = m.G().API.PostDecode(m, APIArg{
		Endpoint:    "downgrade/team",
		SessionType: APISessionTypeREQUIRED,
		Args: HTTPArgs{
			"team_id":     S{string(teamID)},
			"member_uids": S{UidsToString(uids)},
		},
	}, &res)
	if err != nil {
		return nil, nil, err
	}
	return leaseWithMerkleRoot(m, res)
}

func leaseWithMerkleRoot(m MetaContext, res leaseReply) (lease *Lease, mr *MerkleRoot, err error) {
	mr, err = m.G().MerkleClient.FetchRootFromServerBySeqno(m, res.Lease.MerkleSeqno)
	if err != nil {
		return nil, nil, err
	}
	return &res.Lease, mr, nil
}
