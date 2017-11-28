package algo

import (
	"fmt"
	"strconv"
	"testing"

	"gopkg.in/mgo.v2"
)

func TestPrefixHash(t *testing.T) {
	if bytes, err := To16Bytes("u_github_012345678"); err != nil {
		t.Fatalf(err.Error())
	} else if l := len(bytes); l != 16 {
		t.Fatalf(fmt.Sprintf("hash length is %s", l))
	}

	m := make(map[[16]byte]bool)
	bs := [16]byte{}

	bytes, _ := To16Bytes("u_github_012345678")
	copy(bs[:], bytes[:])
	m[bs] = true
	bytes, _ = To16Bytes("p_gitlab_012345678")
	copy(bs[:], bytes[:])
	m[bs] = true
	bytes, _ = To16Bytes("u_source_012345678")
	copy(bs[:], bytes[:])
	m[bs] = true
	bytes, _ = To16Bytes("u_git_012345678")
	copy(bs[:], bytes[:])
	m[bs] = true
	bytes, _ = To16Bytes("u_bitbkt_012345678")
	copy(bs[:], bytes[:])
	m[bs] = true
	if len(m) < 5 {
		t.Fatalf(fmt.Sprintf("Conflict found. only %s keys found", len(m)))
	}
}

type doc struct {
	Id []byte `bson:"_id"`
	I  int
}

func TestIdHash(t *testing.T) {
	session, err := mgo.Dial("localhost")
	defer func() {
		if session != nil {
			session.Close()
		}
	}()

	if err != nil {
		t.Fatalf(err.Error())
	}

	h := session.DB("test").C("hash")

	var max int = 1 << 16 // 2 ** 16
	var i int = 0
	for ; i < max; i++ {
		if bytes, err := To16Bytes("u_github_" + strconv.Itoa(i)); err != nil {
			t.Fatalf(err.Error())
		} else {
			_, err := h.UpsertId(bytes, doc{bytes, i})
			if err != nil {
				t.Errorf(err.Error())
			}
		}
	}
	// echo db.hash.count()| mongo.
	// check conflicts if count < max manually
	//   db.hash.aggregate([{$group: {_id: '$i'} }, {$group: {_id: 1, count: {$sum: 1}}}]);

	// db.hash.find().limit(100) to check data distribution
}
